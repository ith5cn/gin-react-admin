package system

import (
	systemModel "server/model/system"
	systemRequest "server/model/system/request"
)

// MenuList 查询菜单并组装成树形结构返回（菜单管理页面用，不分页）。
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

// CreateMenu 新增菜单。createWithLevel 会根据 parent_id 自动维护 level 层级路径。
func CreateMenu(payload systemRequest.MenuPayload) (*systemModel.AISystemMenu, error) {
	return createWithLevel[systemModel.AISystemMenu]("ai_system_menu", menuPayloadData(payload))
}

// UpdateMenu 更新菜单，父级变化时同步重算 level。
func UpdateMenu(id string, payload systemRequest.MenuPayload) (*systemModel.AISystemMenu, error) {
	return updateWithLevel[systemModel.AISystemMenu]("ai_system_menu", id, menuPayloadData(payload))
}

// DeleteMenu 删除菜单；存在子菜单时拒绝删除，防止树上出现"孤儿节点"。
func DeleteMenu(id string) error {
	has, err := hasChildren("ai_system_menu", id)
	if err != nil {
		return err
	}
	if has {
		return ErrMenuHasChildren
	}
	return deleteByID(&systemModel.AISystemMenu{}, id)
}

// AccessMenu 返回指定用户有权访问的菜单树（RBAC 核心查询）：
// 用户 → user_role 中间表 → 角色 → role_menu 中间表 → 菜单。
// 超级管理员跳过角色过滤，直接返回全部启用菜单。
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

// MenuIDsByRoleID 查询角色已绑定的菜单 ID，用于权限弹窗回显勾选。
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

// BuildMenuTree 把扁平的菜单列表组装成树。
// 经典两步法：先把所有节点放进 map（父节点查找 O(1)），
// 再遍历一次把每个节点挂到父节点下；找不到父节点的提升为根，数据异常时也不丢节点。
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

// menuPayloadData 把类型化入参转成 GORM 更新 map，nil 字段跳过（部分更新语义）。
func menuPayloadData(payload systemRequest.MenuPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "parent_id", payload.ParentID)
	setColumn(data, "name", payload.Name)
	setColumn(data, "code", payload.Code)
	setColumn(data, "icon", payload.Icon)
	setColumn(data, "route", payload.Route)
	setColumn(data, "component", payload.Component)
	setColumn(data, "redirect", payload.Redirect)
	setColumn(data, "is_hidden", payload.IsHidden)
	setColumn(data, "is_layout", payload.IsLayout)
	setColumn(data, "type", payload.Type)
	setColumn(data, "status", payload.Status)
	setColumn(data, "sort", payload.Sort)
	setColumn(data, "remark", payload.Remark)
	return data
}
