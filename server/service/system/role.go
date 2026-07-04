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

// CreateRole 新增角色，level 层级路径由 createWithLevel 自动维护；
// 携带 deptIds 时同步写入自定义数据权限的部门授权。
func CreateRole(payload systemRequest.RolePayload) (*systemModel.AISystemRole, error) {
	role, err := createWithLevel[systemModel.AISystemRole]("ai_system_role", rolePayloadData(payload))
	if err != nil {
		return nil, err
	}
	if payload.DeptIDs != nil {
		if err := BindRoleDepts(role.ID, *payload.DeptIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

// UpdateRole 更新角色，父级变化时同步重算 level；deptIds 非 nil 时整体重设部门授权。
func UpdateRole(id string, payload systemRequest.RolePayload) (*systemModel.AISystemRole, error) {
	role, err := updateWithLevel[systemModel.AISystemRole]("ai_system_role", id, rolePayloadData(payload))
	if err != nil {
		return nil, err
	}
	if payload.DeptIDs != nil {
		if err := BindRoleDepts(role.ID, *payload.DeptIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

// DeleteRole 删除角色；存在子角色时拒绝删除，删除后清理部门授权关联。
func DeleteRole(id string) error {
	has, err := hasChildren("ai_system_role", id)
	if err != nil {
		return err
	}
	if has {
		return ErrRoleHasChildren
	}
	if err := deleteByID(&systemModel.AISystemRole{}, id); err != nil {
		return err
	}
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Where("role_id = ?", id).Delete(&systemModel.AISystemRoleDept{}).Error
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

// BindRoleDepts 重设角色的自定义数据权限部门集合（事务内先删后插，整体替换）。
func BindRoleDepts(roleID uint, deptIDs []uint) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&systemModel.AISystemRoleDept{}).Error; err != nil {
			return err
		}
		for _, deptID := range deptIDs {
			if err := tx.Create(&systemModel.AISystemRoleDept{RoleID: roleID, DeptID: deptID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// RoleDeptIDsByRoleID 查询角色已授权的部门 ID，前端编辑自定义数据权限时回显用。
func RoleDeptIDsByRoleID(roleID string) ([]uint, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	id, err := parseUint(roleID)
	if err != nil {
		return nil, err
	}
	return roleDeptIDs(db, id)
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
