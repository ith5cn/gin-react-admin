package system

import (
	"errors"
	systemModel "server/model/system"
	systemResponse "server/model/system/response"
	gormInit "server/setup/gorm"
	"server/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ErrLoginFailed 是登录失败的统一错误。
// 不区分用户不存在、密码错误、账号禁用，避免接口泄露过多账号状态信息。
var ErrLoginFailed = errors.New("user_name or password is incorrect")

// Login 负责完整登录业务：
// 1. 从 ai_system_user 表按 username 查询用户。
// 2. 校验账号状态和 bcrypt 密码。
// 3. 校验成功后签发 JWT access token 和 refresh token。
func Login(userName string, password string) (*utils.TokenPair, error) {
	if gormInit.Gorm.Databases == nil || gormInit.Gorm.Databases.AISystem == nil {
		return nil, errors.New("ai_system is not initialized")
	}

	db := gormInit.Gorm.Databases.AISystem

	// 接口入参仍叫 user_name，数据库字段是 username，这里做一次语义映射。
	var user systemModel.AISystemUser
	if err := db.Where("username = ?", userName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLoginFailed
		}
		return nil, err
	}

	if user.Status != 1 {
		return nil, ErrLoginFailed
	}
	if !checkPassword(user.Password, password) {
		return nil, ErrLoginFailed
	}

	return utils.GenerateToken(user.ID, user.Username)
}

// checkPassword 使用 bcrypt 校验明文密码和数据库中的哈希密码。
func checkPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// CurrentUserContext 查询当前登录用户的前端初始化上下文。
// 返回内容包括用户资料、角色、动态菜单树、按钮权限码等。
func CurrentUserContext(userID uint) (*systemResponse.UserContext, error) {
	if gormInit.Gorm.Databases == nil || gormInit.Gorm.Databases.AISystem == nil {
		return nil, errors.New("ai_system is not initialized")
	}

	db := gormInit.Gorm.Databases.AISystem

	var user systemModel.AISystemUser
	if err := db.Where("id = ? AND status = ? AND delete_time IS NULL", userID, 1).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLoginFailed
		}
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
	roleValues := make([]interface{}, 0, len(roles))
	isSuperAdmin := false
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
		if isSuperAdminRole(role) {
			isSuperAdmin = true
		}
		if role.Code != nil && *role.Code != "" {
			roleValues = append(roleValues, *role.Code)
			continue
		}
		roleValues = append(roleValues, role.ID)
	}

	var menus []systemModel.AISystemMenu
	var err error
	if isSuperAdmin {
		menus, err = allEnabledMenus(db)
	} else {
		menus, err = menusByRoleIDs(roleIDs, db)
	}
	if err != nil {
		return nil, err
	}

	codes := interface{}(permissionCodes(menus))
	if isSuperAdmin {
		codes = "*"
	}

	return &systemResponse.UserContext{
		User: systemResponse.UserInfo{
			ID:             user.ID,
			Username:       user.Username,
			UserType:       stringValue(user.UserType),
			Nickname:       user.Nickname,
			Phone:          user.Phone,
			Email:          user.Email,
			Avatar:         user.Avatar,
			Signed:         user.Signed,
			Dashboard:      user.Dashboard,
			DeptID:         user.DeptID,
			Status:         user.Status,
			BackendSetting: user.BackendSetting,
		},
		Roles:   roleValues,
		Routers: buildMenuTree(menus),
		Codes:   codes,
		Posts:   []interface{}{},
		Depts:   []interface{}{},
	}, nil
}

func isSuperAdminRole(role systemModel.AISystemRole) bool {
	if role.ID == 1 {
		return true
	}
	if role.Code == nil {
		return false
	}
	return *role.Code == "admin" || *role.Code == "super_admin"
}

func allEnabledMenus(db *gorm.DB) ([]systemModel.AISystemMenu, error) {
	var menus []systemModel.AISystemMenu
	err := db.Where("status = ? AND delete_time IS NULL", 1).
		Order("sort ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

func menusByRoleIDs(roleIDs []uint, db *gorm.DB) ([]systemModel.AISystemMenu, error) {
	if len(roleIDs) == 0 {
		return []systemModel.AISystemMenu{}, nil
	}

	var menus []systemModel.AISystemMenu
	err := db.Table("ai_system_menu AS m").
		Select("DISTINCT m.*").
		Joins("JOIN ai_system_role_menu AS rm ON rm.menu_id = m.id").
		Where("rm.role_id IN ? AND m.status = ? AND m.delete_time IS NULL", roleIDs, 1).
		Order("m.sort ASC, m.id ASC").
		Scan(&menus).Error
	return menus, err
}

func buildMenuTree(menus []systemModel.AISystemMenu) []*systemResponse.MenuNode {
	nodeMap := make(map[uint]*systemResponse.MenuNode, len(menus))
	roots := make([]*systemResponse.MenuNode, 0)

	for _, menu := range menus {
		if menu.Type == "B" {
			continue
		}
		nodeMap[menu.ID] = menuNode(menu)
	}

	for _, menu := range menus {
		if menu.Type == "B" {
			continue
		}

		node := nodeMap[menu.ID]
		if menu.ParentID == nil || *menu.ParentID == 0 {
			roots = append(roots, node)
			continue
		}

		parent, ok := nodeMap[*menu.ParentID]
		if !ok {
			roots = append(roots, node)
			continue
		}

		parent.Children = append(parent.Children, node)
	}

	return roots
}

func menuNode(menu systemModel.AISystemMenu) *systemResponse.MenuNode {
	return &systemResponse.MenuNode{
		ID:        menu.ID,
		ParentID:  menu.ParentID,
		Name:      stringValue(menu.Name),
		Code:      stringValue(menu.Code),
		Icon:      menu.Icon,
		Route:     menu.Route,
		Component: menu.Component,
		Redirect:  menu.Redirect,
		IsHidden:  menu.IsHidden,
		IsLayout:  menu.IsLayout,
		Type:      menu.Type,
		Status:    menu.Status,
		Sort:      menu.Sort,
		Children:  []*systemResponse.MenuNode{},
	}
}

func permissionCodes(menus []systemModel.AISystemMenu) []string {
	codes := make([]string, 0)
	seen := make(map[string]struct{})

	for _, menu := range menus {
		if menu.Type != "B" || menu.Code == nil || *menu.Code == "" {
			continue
		}
		if _, ok := seen[*menu.Code]; ok {
			continue
		}

		seen[*menu.Code] = struct{}{}
		codes = append(codes, *menu.Code)
	}

	return codes
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
