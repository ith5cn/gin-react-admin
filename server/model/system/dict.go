package system

import "time"

// AISystemDictType 映射 ai_system_dict_type 字典类型表。
type AISystemDictType struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey"`
	Name       *string    `json:"name" gorm:"column:name"`
	Code       *string    `json:"code" gorm:"column:code"`
	Status     int16      `json:"status" gorm:"column:status"`
	Remark     *string    `json:"remark" gorm:"column:remark"`
	CreatedBy  *int       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy  *int       `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime *time.Time `json:"updateTime" gorm:"column:update_time"`
	DeleteTime *time.Time `json:"-" gorm:"column:delete_time"`
}

func (AISystemDictType) TableName() string {
	return "ai_system_dict_type"
}

// AISystemDictData 映射 ai_system_dict_data 字典数据表。
type AISystemDictData struct {
	ID         uint       `json:"id" gorm:"column:id;primaryKey"`
	TypeID     *uint      `json:"typeId" gorm:"column:type_id"`
	Label      *string    `json:"label" gorm:"column:label"`
	Value      *string    `json:"value" gorm:"column:value"`
	Color      *string    `json:"color" gorm:"column:color"`
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

func (AISystemDictData) TableName() string {
	return "ai_system_dict_data"
}
