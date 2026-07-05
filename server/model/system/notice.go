package system

import "time"

// AISystemNotice 映射 ai_system_notice 通知公告表。
type AISystemNotice struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey"`       // 公告ID。
	Title      *string    `json:"title" gorm:"column:title"`            // 标题。
	Type       *int16     `json:"type" gorm:"column:type"`              // 类型：1 通知，2 公告。
	Content    *string    `json:"content" gorm:"column:content"`        // 内容。
	Status     int16      `json:"status" gorm:"column:status"`          // 状态：1 正常，2 停用。
	Remark     *string    `json:"remark" gorm:"column:remark"`          // 备注。
	CreatedBy  *int       `json:"createdBy" gorm:"column:created_by"`   // 创建者ID。
	UpdatedBy  *int       `json:"updatedBy" gorm:"column:updated_by"`   // 更新者ID。
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间。
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time"` // 更新时间。
	DeleteTime *time.Time `json:"-" gorm:"column:delete_time"`          // 删除时间。
}

func (AISystemNotice) TableName() string {
	return "ai_system_notice"
}
