package system

import (
	"errors"
	systemModel "server/model/system"

	"gorm.io/gorm"
)

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

func CreateRole(data map[string]interface{}) (*systemModel.AISystemRole, error) {
	payload := requestData(data, roleColumns())
	return createWithLevel[systemModel.AISystemRole]("ai_system_role", payload)
}

func UpdateRole(id string, data map[string]interface{}) (*systemModel.AISystemRole, error) {
	payload := requestData(data, roleColumns())
	return updateWithLevel[systemModel.AISystemRole]("ai_system_role", id, payload)
}

func DeleteRole(id string) error {
	has, err := hasChildren("ai_system_role", id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("角色下存在子角色，无法删除")
	}
	return deleteByID(&systemModel.AISystemRole{}, id)
}

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

func RoleAccess() ([]systemModel.AISystemRole, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var roles []systemModel.AISystemRole
	err = softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&roles).Error
	return roles, err
}

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

func roleColumns() map[string]string {
	return map[string]string{"parentId": "parent_id", "parent_id": "parent_id", "name": "name", "code": "code", "dataScope": "data_scope", "data_scope": "data_scope", "status": "status", "sort": "sort", "remark": "remark"}
}
