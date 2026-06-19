package system

import "time"

type ToolGenerateTable struct {
	ID             uint       `json:"id" gorm:"column:id;primaryKey"`
	CreatedBy      *int       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy      *int       `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime     *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime     *time.Time `json:"updateTime" gorm:"column:update_time"`
	TableNameValue string     `json:"table_name" gorm:"column:table_name"`
	TableComment   *string    `json:"table_comment" gorm:"column:table_comment"`
	PackageName    *string    `json:"package_name" gorm:"column:package_name"`
	BusinessName   *string    `json:"business_name" gorm:"column:business_name"`
	ClassName      *string    `json:"class_name" gorm:"column:class_name"`
	MenuName       *string    `json:"menu_name" gorm:"column:menu_name"`
	BelongMenuID   *uint      `json:"belong_menu_id" gorm:"column:belong_menu_id"`
	TplCategory    *string    `json:"tpl_category" gorm:"column:tpl_category"`
	GeneratePath   *string    `json:"generate_path" gorm:"column:generate_path"`
	GenerateModel  int16      `json:"generate_model" gorm:"column:generate_model"`
	FormWidth      int        `json:"form_width" gorm:"column:form_width"`
	IsFull         int16      `json:"is_full" gorm:"column:is_full"`
	Remark         *string    `json:"remark" gorm:"column:remark"`
	Source         *string    `json:"source" gorm:"column:source"`
	ComponentType  int16      `json:"component_type" gorm:"column:component_type"`
	Sort           int16      `json:"sort" gorm:"column:sort"`
	DeleteTime     *time.Time `json:"-" gorm:"column:delete_time"`
}

func (ToolGenerateTable) TableName() string {
	return "nest_tool_generate_tables"
}

type ToolGenerateColumn struct {
	ID            uint       `json:"id" gorm:"column:id;primaryKey"`
	CreatedBy     *int       `json:"createdBy" gorm:"column:created_by"`
	UpdatedBy     *int       `json:"updatedBy" gorm:"column:updated_by"`
	CreateTime    *time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime    *time.Time `json:"updateTime" gorm:"column:update_time"`
	TableID       uint       `json:"table_id" gorm:"column:table_id"`
	ColumnName    string     `json:"column_name" gorm:"column:column_name"`
	ColumnComment *string    `json:"column_comment" gorm:"column:column_comment"`
	ColumnType    *string    `json:"column_type" gorm:"column:column_type"`
	DefaultValue  *string    `json:"default_value" gorm:"column:default_value"`
	IsPK          int16      `json:"is_pk" gorm:"column:is_pk"`
	IsRequired    int16      `json:"is_required" gorm:"column:is_required"`
	IsInsert      int16      `json:"is_insert" gorm:"column:is_insert"`
	IsEdit        int16      `json:"is_edit" gorm:"column:is_edit"`
	IsList        int16      `json:"is_list" gorm:"column:is_list"`
	IsQuery       int16      `json:"is_query" gorm:"column:is_query"`
	IsSort        int16      `json:"is_sort" gorm:"column:is_sort"`
	QueryType     string     `json:"query_type" gorm:"column:query_type"`
	ViewType      string     `json:"view_type" gorm:"column:view_type"`
	DictType      *string    `json:"dict_type" gorm:"column:dict_type"`
	AllowRoles    *string    `json:"allow_roles" gorm:"column:allow_roles"`
	Sort          uint8      `json:"sort" gorm:"column:sort"`
	Remark        *string    `json:"remark" gorm:"column:remark"`
	DeleteTime    *time.Time `json:"-" gorm:"column:delete_time"`
}

func (ToolGenerateColumn) TableName() string {
	return "nest_tool_generate_columns"
}
