package system

import "time"

// AISystemPost 映射 ai_system_post 岗位表。
type AISystemPost struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey"`
	Name       *string    `json:"name" gorm:"column:name"`
	Code       *string    `json:"code" gorm:"column:code"`
	Sort       uint16     `json:"sort" gorm:"column:sort"`
	Status     int16      `json:"status" gorm:"column:status"`
	Remark     *string    `json:"remark" gorm:"column:remark"`
	CreatedBy  *int       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy  *int       `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time"`
	DeleteTime *time.Time `json:"-" gorm:"column:delete_time"`
}

func (AISystemPost) TableName() string {
	return "ai_system_post"
}
