package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	loggerInit "server/setup/logger"
	"strings"
	"time"

	"go.uber.org/zap"
)

// LoginLogList 分页查询登录日志，按 id 倒序（最新的在前）。
func LoginLogList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemLoginLog
	return pageList(query, &systemModel.AISystemLoginLog{}, &data, map[string]string{"username": "username"}, map[string]string{"status": "status", "ip": "ip"}, "id DESC")
}

// OperLogList 分页查询操作日志，按 id 倒序。
func OperLogList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemOperLog
	return pageList(query, &systemModel.AISystemOperLog{}, &data, map[string]string{"username": "username", "serviceName": "service_name", "router": "router", "ip": "ip"}, map[string]string{}, "id DESC")
}

// RecordLoginLog 写一条登录日志。
// 日志属于"尽力而为"：写失败只记服务端日志，绝不让日志问题影响登录主流程，
// 所以本函数不向调用方返回 error。
func RecordLoginLog(username, ip, userAgent string, success bool, message string) {
	db, err := systemDB()
	if err != nil {
		loggerInit.Logger.Get().Error("record login log failed", zap.Error(err))
		return
	}

	status := int16(2)
	if success {
		status = 1
	}
	now := time.Now()
	osName, browser := parseUserAgent(userAgent)
	entry := systemModel.AISystemLoginLog{
		Username:  ptrString(username),
		IP:        ptrString(ip),
		OS:        ptrString(osName),
		Browser:   ptrString(browser),
		Status:    status,
		Message:   ptrString(message),
		LoginTime: &now,
	}
	if err := db.Create(&entry).Error; err != nil {
		loggerInit.Logger.Get().Error("record login log failed", zap.Error(err))
	}
}

// RecordOperLog 写一条操作日志，由操作日志中间件在写类请求完成后调用。
// 与登录日志一样尽力而为，不返回 error。
func RecordOperLog(username, method, router, serviceName, ip, requestData string) {
	db, err := systemDB()
	if err != nil {
		loggerInit.Logger.Get().Error("record oper log failed", zap.Error(err))
		return
	}

	entry := systemModel.AISystemOperLog{
		App:         ptrString("backend"),
		Method:      ptrString(method),
		Router:      ptrString(router),
		ServiceName: ptrString(serviceName),
		Username:    ptrString(username),
		IP:          ptrString(ip),
		RequestData: ptrString(requestData),
	}
	if err := db.Create(&entry).Error; err != nil {
		loggerInit.Logger.Get().Error("record oper log failed", zap.Error(err))
	}
}

// parseUserAgent 从 User-Agent 里粗略识别操作系统和浏览器。
// 只做展示用途，不追求精确，避免为此引入第三方 UA 解析库。
func parseUserAgent(userAgent string) (osName string, browser string) {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "windows"):
		osName = "Windows"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macintosh"):
		osName = "macOS"
	case strings.Contains(ua, "android"):
		osName = "Android"
	case strings.Contains(ua, "iphone"), strings.Contains(ua, "ipad"):
		osName = "iOS"
	case strings.Contains(ua, "linux"):
		osName = "Linux"
	default:
		osName = "Unknown"
	}

	// 判断顺序有讲究：Edge/Chrome 的 UA 都带 safari 字样，要先判特征更明显的。
	switch {
	case strings.Contains(ua, "edg/"):
		browser = "Edge"
	case strings.Contains(ua, "chrome/"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox/"):
		browser = "Firefox"
	case strings.Contains(ua, "safari/"):
		browser = "Safari"
	default:
		browser = "Unknown"
	}
	return osName, browser
}
