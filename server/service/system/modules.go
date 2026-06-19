package system

import (
	"errors"
	"fmt"
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func UserList(query map[string]string) (*commonResponse.PageResult, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	page := parsePage(query)
	base := softDelete(db.Model(&systemModel.AISystemUser{}))
	base = applyFilters(base, query,
		map[string]string{"username": "username", "nickname": "nickname", "phone": "phone", "email": "email"},
		map[string]string{"status": "status", "deptId": "dept_id"},
	)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}

	var users []systemModel.AISystemUser
	if err := base.Order("id ASC").Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(&users).Error; err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Roles, _ = RoleIDsByUserID(users[i].ID)
	}

	return &commonResponse.PageResult{List: users, Total: total}, nil
}

func CreateUser(data map[string]interface{}) (*systemModel.AISystemUser, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	roles := idsFromAny(data["roles"])
	payload := requestData(data, userColumns())
	if password, ok := data["password"].(string); ok && password != "" {
		hash, err := hashPassword(password)
		if err != nil {
			return nil, err
		}
		payload["password"] = hash
	}
	setDefaultTimes(payload, true)

	if payload["password"] == nil {
		return nil, errors.New("password is required")
	}

	if err := db.Model(&systemModel.AISystemUser{}).Create(payload).Error; err != nil {
		return nil, err
	}

	var user systemModel.AISystemUser
	if err := db.Where("username = ?", payload["username"]).First(&user).Error; err != nil {
		return nil, err
	}
	if len(roles) > 0 {
		if err := BindUserRoles(user.ID, roles); err != nil {
			return nil, err
		}
	}
	user.Roles = roles
	return &user, nil
}

func UpdateUser(id string, data map[string]interface{}) (*systemModel.AISystemUser, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	payload := requestData(data, userColumns())
	delete(payload, "password")
	setDefaultTimes(payload, false)
	if len(payload) > 0 {
		if err := db.Model(&systemModel.AISystemUser{}).Where("id = ?", id).Updates(payload).Error; err != nil {
			return nil, err
		}
	}
	if rolesValue, ok := data["roles"]; ok {
		if err := BindUserRolesStringID(id, idsFromAny(rolesValue)); err != nil {
			return nil, err
		}
	}

	var user systemModel.AISystemUser
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	user.Roles, _ = RoleIDsByUserID(user.ID)
	return &user, nil
}

func DeleteUser(id string) error {
	return deleteByID(&systemModel.AISystemUser{}, id)
}

func SetUserPassword(id string, data map[string]interface{}) error {
	password, _ := data["password"].(string)
	if password == "" {
		return errors.New("password is required")
	}
	hash, err := hashPassword(password)
	if err != nil {
		return err
	}
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Model(&systemModel.AISystemUser{}).Where("id = ?", id).Updates(map[string]interface{}{"password": hash, "update_time": gorm.Expr("NOW()")}).Error
}

func UserAuthList() ([]map[string]interface{}, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var users []systemModel.AISystemUser
	if err := softDelete(db).Where("status = ?", 1).Order("id ASC").Find(&users).Error; err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		label := user.Username
		if user.Nickname != nil && *user.Nickname != "" {
			label = *user.Nickname
		}
		result = append(result, map[string]interface{}{"label": label, "value": user.ID})
	}
	return result, nil
}

func BindUserRolesStringID(userID string, roleIDs []uint) error {
	id, err := parseUint(userID)
	if err != nil {
		return err
	}
	return BindUserRoles(id, roleIDs)
}

func BindUserRoles(userID uint, roleIDs []uint) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&systemModel.AISystemUserRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			if err := tx.Create(&systemModel.AISystemUserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func RoleIDsByUserID(userID uint) ([]uint, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var rels []systemModel.AISystemUserRole
	if err := db.Where("user_id = ?", userID).Find(&rels).Error; err != nil {
		return nil, err
	}
	result := make([]uint, 0, len(rels))
	for _, rel := range rels {
		result = append(result, rel.RoleID)
	}
	return result, nil
}

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

func DeptList(query map[string]string) ([]*systemModel.AISystemDept, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var depts []systemModel.AISystemDept
	q := softDelete(db.Model(&systemModel.AISystemDept{}))
	q = applyFilters(q, query, map[string]string{"name": "name"}, map[string]string{"parentId": "parent_id", "status": "status"})
	if err := q.Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return BuildDeptTree(depts), nil
}

func CreateDept(data map[string]interface{}) (*systemModel.AISystemDept, error) {
	payload := requestData(data, deptColumns())
	return createWithLevel[systemModel.AISystemDept]("ai_system_dept", payload)
}

func UpdateDept(id string, data map[string]interface{}) (*systemModel.AISystemDept, error) {
	payload := requestData(data, deptColumns())
	return updateWithLevel[systemModel.AISystemDept]("ai_system_dept", id, payload)
}

func DeleteDept(id string) error {
	has, err := hasChildren("ai_system_dept", id)
	if err != nil {
		return err
	}
	if has {
		return errors.New("部门下存在子部门，无法删除")
	}
	return deleteByID(&systemModel.AISystemDept{}, id)
}

func DeptAccess(tree bool) (interface{}, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var depts []systemModel.AISystemDept
	if err := softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	if tree {
		return BuildDeptTree(depts), nil
	}
	return depts, nil
}

func PostList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemPost
	return pageList(query, &systemModel.AISystemPost{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"}, "sort ASC, id ASC")
}

func PostAccess() ([]systemModel.AISystemPost, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var posts []systemModel.AISystemPost
	err = softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&posts).Error
	return posts, err
}

func DictTypeList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemDictType
	return pageList(query, &systemModel.AISystemDictType{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{"status": "status"}, "id ASC")
}

func DictDataList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemDictData
	return pageList(query, &systemModel.AISystemDictData{}, &data, map[string]string{"label": "label", "code": "code"}, map[string]string{"status": "status", "typeId": "type_id"}, "sort ASC, id ASC")
}

func DictAll() (map[string][]systemModel.AISystemDictData, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var dicts []systemModel.AISystemDictData
	if err := softDelete(db).Where("status = ?", 1).Order("sort ASC, id ASC").Find(&dicts).Error; err != nil {
		return nil, err
	}
	result := make(map[string][]systemModel.AISystemDictData)
	for _, item := range dicts {
		if item.Code == nil || *item.Code == "" {
			continue
		}
		result[*item.Code] = append(result[*item.Code], item)
	}
	return result, nil
}

func ConfigGroupList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemConfigGroup
	return pageList(query, &systemModel.AISystemConfigGroup{}, &data, map[string]string{"name": "name", "code": "code"}, map[string]string{}, "sort ASC, id ASC")
}

func ConfigList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemConfig
	return pageList(query, &systemModel.AISystemConfig{}, &data, map[string]string{"name": "name", "key": "`key`"}, map[string]string{"groupId": "group_id"}, "sort ASC, id ASC")
}

func ConfigInfo(code string) (map[string]string, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	var configs []systemModel.AISystemConfig
	if err := db.Table("ai_system_config AS c").
		Select("c.*").
		Joins("JOIN ai_system_config_group AS g ON g.id = c.group_id").
		Where("c.delete_time IS NULL AND (g.code = ? OR c.`key` = ?)", code, code).
		Order("c.sort ASC, c.id ASC").
		Find(&configs).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, item := range configs {
		if item.Value == nil {
			result[item.Key] = ""
			continue
		}
		result[item.Key] = *item.Value
	}
	return result, nil
}

func BatchUpdateConfig(groupID uint, configs []map[string]interface{}) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, item := range configs {
			id := item["id"]
			payload := requestData(item, configColumns())
			payload["group_id"] = groupID
			setDefaultTimes(payload, false)
			if id == nil {
				setDefaultTimes(payload, true)
				if err := tx.Model(&systemModel.AISystemConfig{}).Create(payload).Error; err != nil {
					return err
				}
				continue
			}
			if err := tx.Model(&systemModel.AISystemConfig{}).Where("id = ?", id).Updates(payload).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func LoginLogList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemLoginLog
	return pageList(query, &systemModel.AISystemLoginLog{}, &data, map[string]string{"username": "username"}, map[string]string{"status": "status", "ip": "ip"}, "id DESC")
}

func OperLogList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.AISystemOperLog
	return pageList(query, &systemModel.AISystemOperLog{}, &data, map[string]string{"username": "username", "serviceName": "service_name", "router": "router", "ip": "ip"}, map[string]string{}, "id DESC")
}

func CreatePost(data map[string]interface{}) (*systemModel.AISystemPost, error) {
	return createSimple[systemModel.AISystemPost]("ai_system_post", data, simpleColumns())
}

func UpdatePost(id string, data map[string]interface{}) (*systemModel.AISystemPost, error) {
	return updateSimple[systemModel.AISystemPost]("ai_system_post", id, data, simpleColumns())
}

func DeletePost(id string) error {
	return deleteByID(&systemModel.AISystemPost{}, id)
}

func CreateDictType(data map[string]interface{}) (*systemModel.AISystemDictType, error) {
	return createSimple[systemModel.AISystemDictType]("ai_system_dict_type", data, dictTypeColumns())
}

func UpdateDictType(id string, data map[string]interface{}) (*systemModel.AISystemDictType, error) {
	return updateSimple[systemModel.AISystemDictType]("ai_system_dict_type", id, data, dictTypeColumns())
}

func DeleteDictType(id string) error {
	return deleteByID(&systemModel.AISystemDictType{}, id)
}

func CreateDictData(data map[string]interface{}) (*systemModel.AISystemDictData, error) {
	return createSimple[systemModel.AISystemDictData]("ai_system_dict_data", data, dictDataColumns())
}

func UpdateDictData(id string, data map[string]interface{}) (*systemModel.AISystemDictData, error) {
	return updateSimple[systemModel.AISystemDictData]("ai_system_dict_data", id, data, dictDataColumns())
}

func DeleteDictData(id string) error {
	return deleteByID(&systemModel.AISystemDictData{}, id)
}

func CreateConfigGroup(data map[string]interface{}) (*systemModel.AISystemConfigGroup, error) {
	return createSimple[systemModel.AISystemConfigGroup]("ai_system_config_group", data, configGroupColumns())
}

func UpdateConfigGroup(id string, data map[string]interface{}) (*systemModel.AISystemConfigGroup, error) {
	return updateSimple[systemModel.AISystemConfigGroup]("ai_system_config_group", id, data, configGroupColumns())
}

func DeleteConfigGroup(id string) error {
	return deleteByID(&systemModel.AISystemConfigGroup{}, id)
}

func CreateConfig(data map[string]interface{}) (*systemModel.AISystemConfig, error) {
	return createSimple[systemModel.AISystemConfig]("ai_system_config", data, configColumns())
}

func UpdateConfig(id string, data map[string]interface{}) (*systemModel.AISystemConfig, error) {
	return updateSimple[systemModel.AISystemConfig]("ai_system_config", id, data, configColumns())
}

func DeleteConfig(id string) error {
	return deleteByID(&systemModel.AISystemConfig{}, id)
}

func pageList(query map[string]string, model interface{}, dest interface{}, likes map[string]string, equals map[string]string, order string) (*commonResponse.PageResult, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	page := parsePage(query)
	base := softDelete(db.Model(model))
	base = applyFilters(base, query, likes, equals)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := base.Order(order).Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(dest).Error; err != nil {
		return nil, err
	}
	return &commonResponse.PageResult{List: dest, Total: total}, nil
}

func PageList(query map[string]string, model interface{}, dest interface{}, likes map[string]string, equals map[string]string, order string) (*commonResponse.PageResult, error) {
	return pageList(query, model, dest, likes, equals, order)
}

func CreateRecord[T any](table string, data map[string]interface{}) (*T, error) {
	return createSimple[T](table, data, passthroughColumns(data))
}

func UpdateRecord[T any](table string, id string, data map[string]interface{}) (*T, error) {
	return updateSimple[T](table, id, data, passthroughColumns(data))
}

func DeleteRecord(model interface{}, id string) error {
	return deleteByID(model, id)
}

func passthroughColumns(data map[string]interface{}) map[string]string {
	columns := make(map[string]string, len(data))
	for key := range data {
		columns[key] = camelToSnake(key)
	}
	return columns
}

func createWithLevel[T any](table string, payload map[string]interface{}) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	normalizeParentAndLevel(db, table, payload)
	setDefaultTimes(payload, true)
	if err := db.Table(table).Create(payload).Error; err != nil {
		return nil, err
	}
	var result T
	if err := db.Table(table).Order("id DESC").First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func createSimple[T any](table string, data map[string]interface{}, allowed map[string]string) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	payload := requestData(data, allowed)
	setDefaultTimes(payload, true)
	if err := db.Table(table).Create(payload).Error; err != nil {
		return nil, err
	}
	var result T
	if err := db.Table(table).Order("id DESC").First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func updateSimple[T any](table string, id string, data map[string]interface{}, allowed map[string]string) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	payload := requestData(data, allowed)
	setDefaultTimes(payload, false)
	if err := db.Table(table).Where("id = ?", id).Updates(payload).Error; err != nil {
		return nil, err
	}
	var result T
	if err := db.Table(table).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func updateWithLevel[T any](table string, id string, payload map[string]interface{}) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	normalizeParentAndLevel(db, table, payload)
	setDefaultTimes(payload, false)
	if err := db.Table(table).Where("id = ?", id).Updates(payload).Error; err != nil {
		return nil, err
	}
	var result T
	if err := db.Table(table).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func normalizeParentAndLevel(db *gorm.DB, table string, payload map[string]interface{}) {
	parentID, ok := payload["parent_id"]
	if !ok || fmt.Sprint(parentID) == "" || fmt.Sprint(parentID) == "0" || fmt.Sprint(parentID) == "<nil>" {
		payload["parent_id"] = 0
		payload["level"] = "0"
		return
	}
	var parent struct {
		Level *string `gorm:"column:level"`
	}
	if err := db.Table(table).Select("level").Where("id = ?", parentID).First(&parent).Error; err == nil && parent.Level != nil && *parent.Level != "" {
		payload["level"] = strings.Trim(*parent.Level+","+fmt.Sprint(parentID), ",")
		return
	}
	payload["level"] = "0"
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

func BuildDeptTree(depts []systemModel.AISystemDept) []*systemModel.AISystemDept {
	nodeMap := make(map[uint]*systemModel.AISystemDept, len(depts))
	roots := make([]*systemModel.AISystemDept, 0)
	for i := range depts {
		dept := depts[i]
		dept.Children = []*systemModel.AISystemDept{}
		nodeMap[dept.ID] = &dept
	}
	for _, dept := range nodeMap {
		if isZeroParent(dept.ParentID) {
			roots = append(roots, dept)
			continue
		}
		if parent, ok := nodeMap[*dept.ParentID]; ok {
			parent.Children = append(parent.Children, dept)
		} else {
			roots = append(roots, dept)
		}
	}
	sortTreeChildren(roots, func(n *systemModel.AISystemDept) []*systemModel.AISystemDept { return n.Children }, func(a, b *systemModel.AISystemDept) bool {
		if a.Sort == b.Sort {
			return a.ID < b.ID
		}
		return a.Sort < b.Sort
	})
	return roots
}

func mustDB() *gorm.DB {
	db, _ := systemDB()
	return db
}

func parseUint(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	return uint(id), err
}

func idsFromAny(value interface{}) []uint {
	items, ok := value.([]interface{})
	if !ok {
		return []uint{}
	}
	result := make([]uint, 0, len(items))
	for _, item := range items {
		switch v := item.(type) {
		case float64:
			result = append(result, uint(v))
		case int:
			result = append(result, uint(v))
		case string:
			if id, err := parseUint(v); err == nil {
				result = append(result, id)
			}
		}
	}
	return result
}

func userColumns() map[string]string {
	return map[string]string{"username": "username", "user_type": "user_type", "userType": "user_type", "nickname": "nickname", "phone": "phone", "email": "email", "avatar": "avatar", "signed": "signed", "dashboard": "dashboard", "deptId": "dept_id", "status": "status", "remark": "remark", "backendSetting": "backend_setting"}
}

func menuColumns() map[string]string {
	return map[string]string{"parentId": "parent_id", "name": "name", "code": "code", "icon": "icon", "route": "route", "component": "component", "redirect": "redirect", "isHidden": "is_hidden", "isLayout": "is_layout", "type": "type", "status": "status", "sort": "sort", "remark": "remark"}
}

func roleColumns() map[string]string {
	return map[string]string{"parentId": "parent_id", "parent_id": "parent_id", "name": "name", "code": "code", "dataScope": "data_scope", "data_scope": "data_scope", "status": "status", "sort": "sort", "remark": "remark"}
}

func deptColumns() map[string]string {
	return map[string]string{"parentId": "parent_id", "parent_id": "parent_id", "name": "name", "status": "status", "sort": "sort", "remark": "remark"}
}

func configColumns() map[string]string {
	return map[string]string{"groupId": "group_id", "group_id": "group_id", "key": "key", "value": "value", "name": "name", "inputType": "input_type", "input_type": "input_type", "configSelectData": "config_select_data", "config_select_data": "config_select_data", "sort": "sort", "remark": "remark"}
}

func simpleColumns() map[string]string {
	return map[string]string{"name": "name", "code": "code", "sort": "sort", "status": "status", "remark": "remark"}
}

func dictTypeColumns() map[string]string {
	return map[string]string{"name": "name", "code": "code", "status": "status", "remark": "remark"}
}

func dictDataColumns() map[string]string {
	return map[string]string{"typeId": "type_id", "type_id": "type_id", "label": "label", "value": "value", "color": "color", "code": "code", "sort": "sort", "status": "status", "remark": "remark"}
}

func configGroupColumns() map[string]string {
	return map[string]string{"name": "name", "code": "code", "sort": "sort", "remark": "remark"}
}
