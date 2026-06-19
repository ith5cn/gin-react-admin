package system

import "time"

// AISystemConfigGroup 映射 ai_system_config_group 配置分组表。
type AISystemConfigGroup struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey"`
	Name       *string    `json:"name" gorm:"column:name"`
	Code       *string    `json:"code" gorm:"column:code"`
	Sort       uint16     `json:"sort" gorm:"column:sort"`
	Remark     *string    `json:"remark" gorm:"column:remark"`
	CreatedBy  *int       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy  *int       `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time"`
	DeleteTime *time.Time `json:"-" gorm:"column:delete_time"`
}

func (AISystemConfigGroup) TableName() string {
	return "ai_system_config_group"
}

// AISystemConfig 映射 ai_system_config 配置项表。
type AISystemConfig struct {
	ID               uint       `json:"id" gorm:"column:id;primaryKey"`
	GroupID          *uint      `json:"groupId" gorm:"column:group_id"`
	Key              string     `json:"key" gorm:"column:key"`
	Value            *string    `json:"value" gorm:"column:value"`
	Name             *string    `json:"name" gorm:"column:name"`
	InputType        *string    `json:"inputType" gorm:"column:input_type"`
	ConfigSelectData *string    `json:"configSelectData" gorm:"column:config_select_data"`
	Sort             uint16     `json:"sort" gorm:"column:sort"`
	Remark           *string    `json:"remark" gorm:"column:remark"`
	CreatedBy        *int       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy        *int       `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime       *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime       *time.Time `json:"updateTime" gorm:"column:update_time"`
	DeleteTime       *time.Time `json:"-" gorm:"column:delete_time"`
}

func (AISystemConfig) TableName() string {
	return "ai_system_config"
}
