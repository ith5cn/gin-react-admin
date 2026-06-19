package system

import "time"

// AISystemDept 映射 ai_system_dept 部门表。
type AISystemDept struct {
	ID         uint            `json:"id" gorm:"column:id;primaryKey"`
	ParentID   *uint           `json:"parentId" gorm:"column:parent_id"`
	Level      *string         `json:"level" gorm:"column:level"`
	Name       *string         `json:"name" gorm:"column:name"`
	Status     int16           `json:"status" gorm:"column:status"`
	Sort       uint16          `json:"sort" gorm:"column:sort"`
	Remark     *string         `json:"remark" gorm:"column:remark"`
	CreatedBy  *int            `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy  *int            `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime *time.Time      `json:"createTime" gorm:"column:create_time"`
	UpdateTime *time.Time      `json:"updateTime" gorm:"column:update_time"`
	DeleteTime *time.Time      `json:"-" gorm:"column:delete_time"`
	Children   []*AISystemDept `json:"children,omitempty" gorm:"-"`
}

func (AISystemDept) TableName() string {
	return "ai_system_dept"
}
