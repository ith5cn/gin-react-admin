package system

import (
	systemModel "server/model/system"
	systemRequest "server/model/system/request"

	"gorm.io/gorm"
)

// RoleList 查询角色并组装成树（角色支持父子层级）。
func RoleList(query map[string]string) ([]*systemModel.AISystemRole, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var roles []systemModel.AISystemRole
	q := softDelete(db.Model(&systemModel.AISystemRole{}))
	q = applyFilters(q, query, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"})
	if err := q.Order("sort ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return BuildRoleTree(roles), nil
}

// CreateRole 新增角色，level 层级路径由 createWithLevel 自动维护。
func CreateRole(payload systemRequest.RolePayload) (*systemModel.AISystemRole, error) {
	return createWithLevel[systemModel.AISystemRole]("ai_system_role", rolePayloadData(payload))
}

// UpdateRole 更新角色，父级变化时同步重算 level。
func UpdateRole(id string, payload systemRequest.RolePayload) (*systemModel.AISystemRole, error) {
	return updateWithLevel[systemModel.AISystemRole]("ai_system_role", id, rolePayloadData(payload))
}

// DeleteRole 删除角色；存在子角色时拒绝删除。
func DeleteRole(id string) error {
	has, err := hasChildren("ai_system_role", id)
	if err != nil {
		return err
	}
	if has {
		return ErrRoleHasChildren
	}
	return deleteByID(&systemModel.AISystemRole{}, id)
}

// BindRoleMenus 重设角色绑定的菜单集合。
// 与用户绑角色一样采用事务内"先删后插"，保证绑定关系的原子替换。
func BindRoleMenus(roleID string, menuIDs []uint) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&systemModel.AISystemRoleMenu{}).Error; err != nil {
			return err
		}
		id, err := parseUint(roleID)
		if err != nil {
			return err
		}
		for _, menuID := range menuIDs {
			if err := tx.Create(&systemModel.AISystemRoleMenu{RoleID: id, MenuID: menuID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// RoleAccess 返回启用状态的角色列表，用于下拉选择。
func RoleAccess() ([]systemModel.AISystemRole, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var roles []systemModel.AISystemRole
	err = softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&roles).Error
	return roles, err
}

// BuildRoleTree 把扁平角色列表组装成树，算法同 BuildMenuTree。
func BuildRoleTree(roles []systemModel.AISystemRole) []*systemModel.AISystemRole {
	nodeMap := make(map[uint]*systemModel.AISystemRole, len(roles))
	roots := make([]*systemModel.AISystemRole, 0)
	for i := range roles {
		role := roles[i]
		role.Children = []*systemModel.AISystemRole{}
		nodeMap[role.ID] = &role
	}
	for _, role := range nodeMap {
		if isZeroParent(role.ParentID) {
			roots = append(roots, role)
			continue
		}
		if parent, ok := nodeMap[*role.ParentID]; ok {
			parent.Children = append(parent.Children, role)
		} else {
			roots = append(roots, role)
		}
	}
	sortTreeChildren(roots, func(n *systemModel.AISystemRole) []*systemModel.AISystemRole { return n.Children }, func(a, b *systemModel.AISystemRole) bool {
		if a.Sort == b.Sort {
			return a.ID < b.ID
		}
		return a.Sort < b.Sort
	})
	return roots
}

// rolePayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过。
func rolePayloadData(payload systemRequest.RolePayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "parent_id", payload.ParentID)
	setColumn(data, "name", payload.Name)
	setColumn(data, "code", payload.Code)
	setColumn(data, "data_scope", payload.DataScope)
	setColumn(data, "status", payload.Status)
	setColumn(data, "sort", payload.Sort)
	setColumn(data, "remark", payload.Remark)
	return data
}
