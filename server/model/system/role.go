package system

import "time"

// AISystemRole 映射 ai_system_role 角色表。
// 当前登录闭环中用于查询当前用户拥有的角色。
type AISystemRole struct {
	ID         uint            `json:"id" gorm:"column:id;primaryKey"`       // 角色ID。
	ParentID   *uint           `json:"parentId" gorm:"column:parent_id"`     // 父角色ID，可为空。
	Level      *string         `json:"level" gorm:"column:level"`            // 层级路径，可为空。
	Name       *string         `json:"name" gorm:"column:name"`              // 角色名称。
	Code       *string         `json:"code" gorm:"column:code"`              // 角色编码。
	DataScope  *int16          `json:"dataScope" gorm:"column:data_scope"`   // 数据权限范围。
	Status     int16           `json:"status" gorm:"column:status"`          // 状态：1 正常，2 停用。
	Sort       uint16          `json:"sort" gorm:"column:sort"`              // 排序值。
	Remark     *string         `json:"remark" gorm:"column:remark"`          // 备注。
	CreatedBy  *int            `json:"createdBy" gorm:"column:created_by"`   // 创建者ID。
	UpdatedBy  *int            `json:"updatedBy" gorm:"column:updated_by"`   // 更新者ID。
	CreateTime *time.Time      `json:"createTime" gorm:"column:create_time"` // 创建时间。
	UpdateTime *time.Time      `json:"updateTime" gorm:"column:update_time"` // 更新时间。
	DeleteTime *time.Time      `json:"-" gorm:"column:delete_time"`          // 删除时间。
	Children   []*AISystemRole `json:"children,omitempty" gorm:"-"`          // 树形子角色。
}

func (AISystemRole) TableName() string {
	return "ai_system_role"
}

// AISystemUserRole 映射 ai_system_user_role 用户角色关联表。
type AISystemUserRole struct {
	ID     uint `gorm:"column:id;primaryKey"` // 关联ID。
	UserID uint `gorm:"column:user_id"`       // 用户ID。
	RoleID uint `gorm:"column:role_id"`       // 角色ID。
}

func (AISystemUserRole) TableName() string {
	return "ai_system_user_role"
}

// AISystemRoleDept 映射 ai_system_role_dept 角色部门关联表。
// 只在角色 data_scope=2（自定义数据权限）时使用，记录该角色可见的部门集合。
type AISystemRoleDept struct {
	ID     uint `gorm:"column:id;primaryKey"` // 关联ID。
	RoleID uint `gorm:"column:role_id"`       // 角色ID。
	DeptID uint `gorm:"column:dept_id"`       // 部门ID。
}

func (AISystemRoleDept) TableName() string {
	return "ai_system_role_dept"
}
