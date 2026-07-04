package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	systemModel "server/model/system"
	loggerInit "server/setup/logger"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// crontabParser 同时支持 5 段标准 cron（分 时 日 月 周）和 6 段带秒格式（秒 分 时 日 月 周）。
// 老数据（saiadmin 迁移来的种子）用的是 6 段格式，标准写法也要能用。
var crontabParser = cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

// crontabScheduler 持有唯一的 cron 实例和"任务ID → cron entry"映射。
// 所有变更（启动、重载）都在互斥锁内进行，避免并发注册出现重复 entry。
type crontabScheduler struct {
	mu      sync.Mutex
	cron    *cron.Cron
	entries map[uint]cron.EntryID
}

var scheduler = &crontabScheduler{entries: map[uint]cron.EntryID{}}

// crontabTaskRegistry 是"系统内部任务"注册表：target 填注册名即可被调度执行。
// 新增内部任务用 RegisterCrontabTask 注册，不要直接改这个 map。
var crontabTaskRegistry = map[string]func(parameter string) error{
	"system/clean-logs": cleanExpiredLogs,
}

// RegisterCrontabTask 注册一个系统内部任务。
// name 是任务在 ai_tool_crontab.target 里填写的标识，fn 的入参是任务配置的 parameter 字符串。
func RegisterCrontabTask(name string, fn func(parameter string) error) {
	crontabTaskRegistry[name] = fn
}

// StartCrontabScheduler 启动定时任务调度器并装载所有启用状态的任务。
// 在 main 的基础设施初始化完成后调用一次；重复调用只会触发一次启动。
func StartCrontabScheduler() error {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	if scheduler.cron == nil {
		scheduler.cron = cron.New(cron.WithParser(crontabParser))
		scheduler.cron.Start()
	}
	return scheduler.reloadLocked()
}

// ReloadCrontabScheduler 重新装载全部任务，任务的增删改后调用。
// 任务量小（后台系统通常几十个以内），全量重载比精细的增量维护更不容易出错。
// 调度器尚未启动时（如 install 模式）静默跳过。
func ReloadCrontabScheduler() {
	scheduler.mu.Lock()
	defer scheduler.mu.Unlock()
	if scheduler.cron == nil {
		return
	}
	if err := scheduler.reloadLocked(); err != nil {
		loggerInit.Logger.Get().Error("reload crontab scheduler failed", zap.Error(err))
	}
}

// reloadLocked 清空现有 entry 后按数据库最新状态重新注册，调用方必须已持有锁。
func (s *crontabScheduler) reloadLocked() error {
	for _, entryID := range s.entries {
		s.cron.Remove(entryID)
	}
	s.entries = map[uint]cron.EntryID{}

	db, err := systemDB()
	if err != nil {
		return err
	}
	var tasks []systemModel.AIToolCrontab
	if err := softDelete(db).Where("status = ?", 1).Find(&tasks).Error; err != nil {
		return err
	}

	for _, task := range tasks {
		if task.Rule == nil || *task.Rule == "" {
			continue
		}
		// 闭包捕获循环变量的经典坑：复制一份再捕获，否则所有 entry 都指向最后一个 task。
		taskCopy := task
		entryID, err := s.cron.AddFunc(*task.Rule, func() {
			executeCrontab(taskCopy, true)
		})
		if err != nil {
			// 单个任务表达式非法不应拖垮整个调度器，记日志后跳过。
			loggerInit.Logger.Get().Error("register crontab failed",
				zap.Uint("id", task.ID), zap.String("rule", *task.Rule), zap.Error(err))
			continue
		}
		s.entries[task.ID] = entryID
	}
	return nil
}

// validateCrontabRule 校验 cron 表达式，无效时返回业务错误（前端可直接展示）。
func validateCrontabRule(rule string) error {
	if strings.TrimSpace(rule) == "" {
		return ErrCrontabRuleInvalid
	}
	if _, err := crontabParser.Parse(rule); err != nil {
		return ErrCrontabRuleInvalid
	}
	return nil
}

// executeCrontab 执行一次任务并写执行日志。
// scheduled 区分调度触发和手动触发：单次任务（singleton=1）只在调度触发后自动停用。
// 执行过程的 panic 由 recover 兜底，绝不让单个任务崩掉整个进程。
func executeCrontab(task systemModel.AIToolCrontab, scheduled bool) {
	defer func() {
		if r := recover(); r != nil {
			loggerInit.Logger.Get().Error("crontab task panic",
				zap.Uint("id", task.ID), zap.Any("panic", r))
		}
	}()

	err := runCrontabTarget(task)

	message := "执行成功"
	status := int16(1)
	if err != nil {
		message = err.Error()
		status = 2
	}
	recordCrontabLog(task, status, message)

	if scheduled && task.Singleton != nil && *task.Singleton == 1 {
		disableCrontab(task.ID)
	}
}

// runCrontabTarget 按执行方式分发：1 内部任务查注册表，2 HTTP 任务发请求。
func runCrontabTarget(task systemModel.AIToolCrontab) error {
	target := ""
	if task.Target != nil {
		target = strings.TrimSpace(*task.Target)
	}
	if target == "" {
		return fmt.Errorf("任务未配置调用目标")
	}
	parameter := ""
	if task.Parameter != nil {
		parameter = *task.Parameter
	}

	if task.TaskStyle != nil && *task.TaskStyle == 2 {
		return runHTTPTask(target, parameter)
	}

	fn, ok := crontabTaskRegistry[target]
	if !ok {
		return fmt.Errorf("未注册的内部任务: %s", target)
	}
	return fn(parameter)
}

// runHTTPTask 执行 HTTP 类型任务：有参数则 POST JSON，无参数则 GET，2xx/3xx 视为成功。
func runHTTPTask(url string, parameter string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("HTTP 任务的调用目标必须是 http(s) URL")
	}

	client := &http.Client{Timeout: 15 * time.Second}
	var resp *http.Response
	var err error
	if strings.TrimSpace(parameter) != "" {
		resp, err = client.Post(url, "application/json", strings.NewReader(parameter))
	} else {
		resp, err = client.Get(url)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 读掉响应体让连接可以复用，内容本身不关心。
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP 状态码 %d", resp.StatusCode)
	}
	return nil
}

// recordCrontabLog 写一条执行日志，尽力而为：失败只记服务端日志。
func recordCrontabLog(task systemModel.AIToolCrontab, status int16, message string) {
	db, err := systemDB()
	if err != nil {
		loggerInit.Logger.Get().Error("record crontab log failed", zap.Error(err))
		return
	}
	now := time.Now()
	crontabID := task.ID
	entry := systemModel.AIToolCrontabLog{
		CrontabID:     &crontabID,
		Name:          task.Name,
		Target:        task.Target,
		Parameter:     task.Parameter,
		ExceptionInfo: ptrString(message),
		Status:        status,
		CreateTime:    &now,
	}
	if err := db.Create(&entry).Error; err != nil {
		loggerInit.Logger.Get().Error("record crontab log failed", zap.Error(err))
	}
}

// disableCrontab 把单次任务置为停用并移出调度器。
func disableCrontab(id uint) {
	db, err := systemDB()
	if err != nil {
		loggerInit.Logger.Get().Error("disable crontab failed", zap.Uint("id", id), zap.Error(err))
		return
	}
	if err := db.Model(&systemModel.AIToolCrontab{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": 2, "update_time": gorm.Expr("NOW()")}).Error; err != nil {
		loggerInit.Logger.Get().Error("disable crontab failed", zap.Uint("id", id), zap.Error(err))
		return
	}
	ReloadCrontabScheduler()
}

// cleanExpiredLogs 是内置的日志清理任务：删除 N 天前的登录/操作日志。
// parameter 支持 {"days": 30}，缺省保留 30 天。
func cleanExpiredLogs(parameter string) error {
	days := 30
	if strings.TrimSpace(parameter) != "" {
		var config struct {
			Days int `json:"days"`
		}
		if err := json.Unmarshal([]byte(parameter), &config); err != nil {
			return fmt.Errorf("参数格式错误，期望 {\"days\": 30}: %w", err)
		}
		if config.Days > 0 {
			days = config.Days
		}
	}

	db, err := systemDB()
	if err != nil {
		return err
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	if err := db.Where("create_time < ?", cutoff).Delete(&systemModel.AISystemOperLog{}).Error; err != nil {
		return err
	}
	return db.Where("login_time < ?", cutoff).Delete(&systemModel.AISystemLoginLog{}).Error
}
