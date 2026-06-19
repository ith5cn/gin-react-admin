package system

import "time"

// AISystemLoginLog 映射 ai_system_login_log 登录日志表。
type AISystemLoginLog struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey"`
	Username   *string    `json:"username" gorm:"column:username"`
	IP         *string    `json:"ip" gorm:"column:ip"`
	IPLocation *string    `json:"ipLocation" gorm:"column:ip_location"`
	OS         *string    `json:"os" gorm:"column:os"`
	Browser    *string    `json:"browser" gorm:"column:browser"`
	Status     int16      `json:"status" gorm:"column:status"`
	Message    *string    `json:"message" gorm:"column:message"`
	LoginTime  *time.Time `json:"loginTime" gorm:"column:login_time"`
	Remark     *string    `json:"remark" gorm:"column:remark"`
	CreatedBy  *int       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy  *int       `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time"`
	DeleteTime *time.Time `json:"-" gorm:"column:delete_time"`
}

func (AISystemLoginLog) TableName() string {
	return "ai_system_login_log"
}

// AISystemOperLog 映射 ai_system_oper_log 操作日志表。
type AISystemOperLog struct {
	ID          uint       `json:"id" gorm:"column:id;primaryKey"`
	App         *string    `json:"app" gorm:"column:app"`
	Method      *string    `json:"method" gorm:"column:method"`
	RequestData *string    `json:"requestData" gorm:"column:request_data"`
	Remark      *string    `json:"remark" gorm:"column:remark"`
	Username    *string    `json:"username" gorm:"column:username"`
	ServiceName *string    `json:"serviceName" gorm:"column:service_name"`
	Router      *string    `json:"router" gorm:"column:router"`
	IP          *string    `json:"ip" gorm:"column:ip"`
	IPLocation  *string    `json:"ipLocation" gorm:"column:ip_location"`
	CreatedBy   *int       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy   *int       `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime  *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime  *time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (AISystemOperLog) TableName() string {
	return "ai_system_oper_log"
}
