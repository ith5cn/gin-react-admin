package system

import (
	"errors"

	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"

	"gorm.io/gorm"
)

// CrontabList 分页查询定时任务。
func CrontabList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AIToolCrontab
	return pageList(query, &systemModel.AIToolCrontab{}, &data,
		map[string]string{"name": "name"},
		map[string]string{"status": "status", "type": "type"},
		"id ASC")
}

// CreateCrontab 新增定时任务：先校验 cron 表达式，写库成功后刷新调度器。
func CreateCrontab(payload systemRequest.CrontabPayload) (*systemModel.AIToolCrontab, error) {
	if payload.Rule != nil {
		if err := validateCrontabRule(*payload.Rule); err != nil {
			return nil, err
		}
	}
	task, err := createRow[systemModel.AIToolCrontab]("ai_tool_crontab", crontabPayloadData(payload))
	if err != nil {
		return nil, err
	}
	ReloadCrontabScheduler()
	return task, nil
}

// UpdateCrontab 更新定时任务，表达式或状态变化会同步到调度器。
func UpdateCrontab(id string, payload systemRequest.CrontabPayload) (*systemModel.AIToolCrontab, error) {
	if payload.Rule != nil {
		if err := validateCrontabRule(*payload.Rule); err != nil {
			return nil, err
		}
	}
	task, err := updateRow[systemModel.AIToolCrontab]("ai_tool_crontab", id, crontabPayloadData(payload))
	if err != nil {
		return nil, err
	}
	ReloadCrontabScheduler()
	return task, nil
}

// DeleteCrontab 删除定时任务并将其移出调度器。
func DeleteCrontab(id string) error {
	if err := deleteByID(&systemModel.AIToolCrontab{}, id); err != nil {
		return err
	}
	ReloadCrontabScheduler()
	return nil
}

// RunCrontabOnce 手动触发一次执行（异步），停用状态的任务也允许手动执行，方便调试。
func RunCrontabOnce(id string) error {
	task, err := crontabByID(id)
	if err != nil {
		return err
	}
	go executeCrontab(*task, false)
	return nil
}

// CrontabLogList 分页查询任务执行日志，最新的在前。
func CrontabLogList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AIToolCrontabLog
	return pageList(query, &systemModel.AIToolCrontabLog{}, &data,
		map[string]string{"name": "name"},
		map[string]string{"crontabId": "crontab_id", "status": "status"},
		"id DESC")
}

// crontabByID 按 ID 查询未删除的任务，不存在时返回业务错误。
func crontabByID(id string) (*systemModel.AIToolCrontab, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var task systemModel.AIToolCrontab
	if err := softDelete(db).Where("id = ?", id).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCrontabNotFound
		}
		return nil, err
	}
	return &task, nil
}

// crontabPayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过（部分更新）。
func crontabPayloadData(payload systemRequest.CrontabPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "name", payload.Name)
	setColumn(data, "type", payload.Type.Int16Ptr())
	setColumn(data, "target", payload.Target)
	setColumn(data, "parameter", payload.Parameter)
	setColumn(data, "task_style", payload.TaskStyle.Int16Ptr())
	setColumn(data, "rule", payload.Rule)
	setColumn(data, "singleton", payload.Singleton.Int16Ptr())
	setColumn(data, "status", payload.Status.Int16Ptr())
	setColumn(data, "remark", payload.Remark)
	return data
}
