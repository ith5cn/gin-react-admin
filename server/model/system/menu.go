package system

import "time"

// AISystemMenu 映射 ai_system_menu 菜单表。
// 菜单类型 type 中，M 表示菜单，B 表示按钮权限，L 表示外链，I 表示 iframe。
type AISystemMenu struct {
	ID          uint            `json:"id" gorm:"column:id;primaryKey"`         // 菜单ID。
	ParentID    *uint           `json:"parentId" gorm:"column:parent_id"`       // 父菜单ID，可为空。
	Level       *string         `json:"level" gorm:"column:level"`              // 层级路径，可为空。
	Name        *string         `json:"name" gorm:"column:name"`                // 菜单名称。
	Code        *string         `json:"code" gorm:"column:code"`                // 菜单或按钮权限标识。
	Icon        *string         `json:"icon" gorm:"column:icon"`                // 菜单图标。
	Route       *string         `json:"route" gorm:"column:route"`              // 前端路由地址。
	Component   *string         `json:"component" gorm:"column:component"`      // 前端组件路径。
	Redirect    *string         `json:"redirect" gorm:"column:redirect"`        // 跳转地址。
	IsHidden    int16           `json:"isHidden" gorm:"column:is_hidden"`       // 是否隐藏：1 是，2 否。
	IsLayout    uint8           `json:"isLayout" gorm:"column:is_layout"`       // 是否继承 layout：1 是，2 否。
	Type        string          `json:"type" gorm:"column:type"`                // 菜单类型：M/B/L/I。
	GenerateID  *int            `json:"generateId" gorm:"column:generate_id"`   // 代码生成ID。
	GenerateKey *string         `json:"generateKey" gorm:"column:generate_key"` // 代码生成标识。
	Status      int16           `json:"status" gorm:"column:status"`            // 状态：1 正常，2 停用。
	Sort        uint16          `json:"sort" gorm:"column:sort"`                // 排序值。
	Remark      *string         `json:"remark" gorm:"column:remark"`            // 备注。
	CreatedBy   *int            `json:"createdBy" gorm:"column:created_by"`     // 创建者ID。
	UpdatedBy   *int            `json:"updatedBy" gorm:"column:updated_by"`     // 更新者ID。
	CreateTime  *time.Time      `json:"createTime" gorm:"column:create_time"`   // 创建时间。
	UpdateTime  *time.Time      `json:"updateTime" gorm:"column:update_time"`   // 更新时间。
	DeleteTime  *time.Time      `json:"-" gorm:"column:delete_time"`            // 删除时间。
	Children    []*AISystemMenu `json:"children,omitempty" gorm:"-"`            // 树形子菜单。
}

func (AISystemMenu) TableName() string {
	return "ai_system_menu"
}

// AISystemRoleMenu 映射 ai_system_role_menu 角色菜单关联表。
type AISystemRoleMenu struct {
	ID     uint `gorm:"column:id;primaryKey"` // 关联ID。
	RoleID uint `gorm:"column:role_id"`       // 角色ID。
	MenuID uint `gorm:"column:menu_id"`       // 菜单ID。
}

func (AISystemRoleMenu) TableName() string {
	return "ai_system_role_menu"
}
