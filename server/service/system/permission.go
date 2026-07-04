package system

import (
	systemModel "server/model/system"
)

// UserHasPermission 判断用户是否拥有指定权限码（接口级 RBAC 校验的核心查询）。
// 判定链路：用户 → 启用角色 → 角色绑定的按钮菜单（type='B'）→ 比对权限码。
// 超级管理员直接放行。
//
// 每次请求两条小查询（角色 + EXISTS），后台系统的量级下无需缓存；
// 好处是改完角色权限立即生效，不存在缓存失效窗口。
func UserHasPermission(userID uint, permCode string) (bool, error) {
	if permCode == "" {
		return true, nil
	}

	db, err := systemDB()
	if err != nil {
		return false, err
	}

	var roles []systemModel.AISystemRole
	if err := db.Table("ai_system_role AS r").
		Select("r.id, r.code").
		Joins("JOIN ai_system_user_role AS ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = ? AND r.delete_time IS NULL", userID, 1).
		Scan(&roles).Error; err != nil {
		return false, err
	}

	roleIDs := make([]uint, 0, len(roles))
	for _, role := range roles {
		if isSuperAdminRole(role) {
			return true, nil
		}
		roleIDs = append(roleIDs, role.ID)
	}
	if len(roleIDs) == 0 {
		return false, nil
	}

	var count int64
	err = db.Table("ai_system_menu AS m").
		Joins("JOIN ai_system_role_menu AS rm ON rm.menu_id = m.id").
		Where("rm.role_id IN ? AND m.code = ? AND m.type = ? AND m.status = ? AND m.delete_time IS NULL",
			roleIDs, permCode, "B", 1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
