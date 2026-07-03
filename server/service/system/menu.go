package system

import (
	"errors"
	systemModel "server/model/system"
)

func MenuList(query map[string]string) (map[string]interface{}, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var menus []systemModel.AISystemMenu
	q := softDelete(db.Model(&systemModel.AISystemMenu{}))
	q = applyFilters(q, query, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"})
	if err := q.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{"data": BuildMenuTree(menus)}, nil
}

func CreateMenu(data map[string]interface{}) (*systemModel.AISystemMenu, error) {
	payload := requestData(data, menuColumns())
	return createWithLevel[systemModel.AISystemMenu]("ai_system_menu", payload)
}

func UpdateMenu(id string, data map[string]interface{}) (*systemModel.AISystemMenu, error) {
	payload := requestData(data, menuColumns())
	return updateWithLevel[systemModel.AISystemMenu]("ai_system_menu", id, payload)
}

func DeleteMenu(id string) error {
	has, err := hasChildren("ai_system_menu", id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("菜单下存在子菜单，无法删除")
	}
	return deleteByID(&systemModel.AISystemMenu{}, id)
}

func AccessMenu(userID uint) ([]*systemModel.AISystemMenu, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	var roles []systemModel.AISystemRole
	if err := db.Table("ai_system_role AS r").
		Select("r.*").
		Joins("JOIN ai_system_user_role AS ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = ? AND r.delete_time IS NULL", userID, 1).
		Order("r.sort ASC, r.id ASC").
		Scan(&roles).Error; err != nil {
		return nil, err
	}

	roleIDs := make([]uint, 0, len(roles))
	isSuperAdmin := false
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
		if isSuperAdminRole(role) {
			isSuperAdmin = true
		}
	}

	var menus []systemModel.AISystemMenu
	if isSuperAdmin {
		menus, err = allEnabledMenus(db)
	} else {
		menus, err = menusByRoleIDs(roleIDs, db)
	}
	if err != nil {
		return nil, err
	}
	return BuildMenuTree(menus), nil
}

func MenuIDsByRoleID(roleID string) ([]uint, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var rels []systemModel.AISystemRoleMenu
	if err := db.Where("role_id = ?", roleID).Find(&rels).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rels))
	for _, rel := range rels {
		ids = append(ids, rel.MenuID)
	}
	return ids, nil
}

func BuildMenuTree(menus []systemModel.AISystemMenu) []*systemModel.AISystemMenu {
	nodeMap := make(map[uint]*systemModel.AISystemMenu, len(menus))
	roots := make([]*systemModel.AISystemMenu, 0)
	for i := range menus {
		menu := menus[i]
		menu.Children = []*systemModel.AISystemMenu{}
		nodeMap[menu.ID] = &menu
	}
	for _, menu := range nodeMap {
		if isZeroParent(menu.ParentID) {
			roots = append(roots, menu)
			continue
		}
		if parent, ok := nodeMap[*menu.ParentID]; ok {
			parent.Children = append(parent.Children, menu)
		} else {
			roots = append(roots, menu)
		}
	}
	sortTreeChildren(roots, func(n *systemModel.AISystemMenu) []*systemModel.AISystemMenu { return n.Children }, func(a, b *systemModel.AISystemMenu) bool {
		if a.Sort == b.Sort {
			return a.ID < b.ID
		}
		return a.Sort < b.Sort
	})
	return roots
}

func menuColumns() map[string]string {
	return map[string]string{"parentId": "parent_id", "name": "name", "code": "code", "icon": "icon", "route": "route", "component": "component", "redirect": "redirect", "isHidden": "is_hidden", "isLayout": "is_layout", "type": "type", "status": "status", "sort": "sort", "remark": "remark"}
}
