package system

import "time"

// AIToolCrontab 映射 ai_tool_crontab 定时任务表。
// rule 是 cron 表达式（支持 5 段标准格式和 6 段带秒格式）；
// task_style 区分执行方式：1 系统内部任务（target 填注册名），2 HTTP 请求任务（target 填 URL）。
type AIToolCrontab struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey"`       // 任务ID。
	Name       *string    `json:"name" gorm:"column:name"`              // 任务名称。
	Type       *int16     `json:"type" gorm:"column:type"`              // 任务类型（字典 crontab_type，仅展示分类用）。
	Target     *string    `json:"target" gorm:"column:target"`          // 调用目标：内部任务名或 HTTP URL。
	Parameter  *string    `json:"parameter" gorm:"column:parameter"`    // 调用参数（JSON 字符串）。
	TaskStyle  *int16     `json:"taskStyle" gorm:"column:task_style"`   // 执行方式：1 内部任务，2 HTTP 请求。
	Rule       *string    `json:"rule" gorm:"column:rule"`              // cron 表达式。
	Singleton  *int16     `json:"singleton" gorm:"column:singleton"`    // 是否单次执行：1 是（跑一次后自动停用），2 否。
	Status     int16      `json:"status" gorm:"column:status"`          // 状态：1 正常（参与调度），2 停用。
	Remark     *string    `json:"remark" gorm:"column:remark"`          // 备注。
	CreatedBy  *int       `json:"createdBy" gorm:"column:created_by"`   // 创建者ID。
	UpdatedBy  *int       `json:"updatedBy" gorm:"column:updated_by"`   // 更新者ID。
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间。
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time"` // 更新时间。
	DeleteTime *time.Time `json:"-" gorm:"column:delete_time"`          // 删除时间。
}

func (AIToolCrontab) TableName() string {
	return "ai_tool_crontab"
}

// AIToolCrontabLog 映射 ai_tool_crontab_log 定时任务执行日志表。
// 每次执行（含手动执行一次）写一条，status 1 成功 2 失败，失败原因记在 exception_info。
type AIToolCrontabLog struct {
	ID            uint       `json:"id" gorm:"column:id;primaryKey"`             // 日志ID。
	CrontabID     *uint      `json:"crontabId" gorm:"column:crontab_id"`         // 所属任务ID。
	Name          *string    `json:"name" gorm:"column:name"`                    // 任务名称快照。
	Target        *string    `json:"target" gorm:"column:target"`                // 调用目标快照。
	Parameter     *string    `json:"parameter" gorm:"column:parameter"`          // 调用参数快照。
	ExceptionInfo *string    `json:"exceptionInfo" gorm:"column:exception_info"` // 执行结果/异常信息。
	Status        int16      `json:"status" gorm:"column:status"`                // 执行状态：1 成功，2 失败。
	CreateTime    *time.Time `json:"createTime" gorm:"column:create_time"`       // 执行时间。
	UpdateTime    *time.Time `json:"updateTime" gorm:"column:update_time"`       // 更新时间。
	DeleteTime    *time.Time `json:"-" gorm:"column:delete_time"`                // 删除时间。
}

func (AIToolCrontabLog) TableName() string {
	return "ai_tool_crontab_log"
}
