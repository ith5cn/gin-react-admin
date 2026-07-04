package system

import (
	systemModel "server/model/system"

	"gorm.io/gorm"
)

// 数据权限（dataScope）：控制"能看到哪些行"，与接口权限（能不能调这个接口）互补。
// 角色的 data_scope 取值：1 全部 2 自定义（按 role_dept 表）3 本部门 4 本部门及以下 5 仅本人。
// 用户有多个角色时取并集（最宽松的生效）；历史数据 data_scope 为 NULL 时按 1（全部）处理，
// 保证升级后现有角色的可见范围不变。

// DataScope 描述一个操作者的数据可见范围，由 UserDataScope 计算得出。
type DataScope struct {
	All      bool   // 全部数据可见，后续条件全部跳过。
	SelfOnly bool   // 命中"仅本人"且没有任何部门授权。
	DeptIDs  []uint // 可见部门 ID 集合（已去重）。
}

// UserDataScope 计算用户的数据可见范围。
// 超级管理员直接可见全部；普通用户按其全部启用角色的 data_scope 求并集。
func UserDataScope(userID uint) (*DataScope, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	var roles []systemModel.AISystemRole
	if err := db.Table("ai_system_role AS r").
		Select("r.id, r.code, r.data_scope").
		Joins("JOIN ai_system_user_role AS ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.status = ? AND r.delete_time IS NULL", userID, 1).
		Scan(&roles).Error; err != nil {
		return nil, err
	}

	deptSet := map[uint]struct{}{}
	hitSelfOnly := false
	for _, role := range roles {
		if isSuperAdminRole(role) {
			return &DataScope{All: true}, nil
		}
		// NULL 视为全部：数据权限功能上线前建的角色不应突然缩小可见范围。
		scope := int16(1)
		if role.DataScope != nil {
			scope = *role.DataScope
		}
		switch scope {
		case 1:
			return &DataScope{All: true}, nil
		case 2:
			ids, err := roleDeptIDs(db, role.ID)
			if err != nil {
				return nil, err
			}
			for _, id := range ids {
				deptSet[id] = struct{}{}
			}
		case 3:
			if deptID, err := userDeptID(db, userID); err != nil {
				return nil, err
			} else if deptID != 0 {
				deptSet[deptID] = struct{}{}
			}
		case 4:
			deptID, err := userDeptID(db, userID)
			if err != nil {
				return nil, err
			}
			if deptID != 0 {
				ids, err := deptWithDescendants(db, deptID)
				if err != nil {
					return nil, err
				}
				for _, id := range ids {
					deptSet[id] = struct{}{}
				}
			}
		case 5:
			hitSelfOnly = true
		}
	}

	result := &DataScope{}
	for id := range deptSet {
		result.DeptIDs = append(result.DeptIDs, id)
	}
	// 无任何角色或角色未授出任何部门时收敛为"仅本人"，宁可少看不可多看（fail-closed）。
	result.SelfOnly = hitSelfOnly || len(result.DeptIDs) == 0
	return result, nil
}

// applyUserDataScope 把数据范围套到用户表查询上：
// 可见部门内的用户 + 自己本人（自己的数据永远可见）。
func applyUserDataScope(db *gorm.DB, scope *DataScope, selfUserID uint) *gorm.DB {
	if scope == nil || scope.All {
		return db
	}
	if len(scope.DeptIDs) == 0 {
		return db.Where("id = ?", selfUserID)
	}
	return db.Where("dept_id IN ? OR id = ?", scope.DeptIDs, selfUserID)
}

// userDeptID 查用户所属部门 ID，未分配部门返回 0。
func userDeptID(db *gorm.DB, userID uint) (uint, error) {
	var user systemModel.AISystemUser
	if err := db.Select("dept_id").Where("id = ?", userID).First(&user).Error; err != nil {
		return 0, err
	}
	if user.DeptID == nil {
		return 0, nil
	}
	return *user.DeptID, nil
}

// roleDeptIDs 查角色在 role_dept 表里授权的部门 ID 列表（自定义数据权限）。
func roleDeptIDs(db *gorm.DB, roleID uint) ([]uint, error) {
	var rels []systemModel.AISystemRoleDept
	if err := db.Where("role_id = ?", roleID).Find(&rels).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rels))
	for _, rel := range rels {
		ids = append(ids, rel.DeptID)
	}
	return ids, nil
}

// deptWithDescendants 返回指定部门及其所有下级部门的 ID。
// 部门总量很小，一次查全表后在内存里按 parent_id 递归展开，避免依赖数据库方言的层级查询。
func deptWithDescendants(db *gorm.DB, deptID uint) ([]uint, error) {
	var depts []systemModel.AISystemDept
	if err := softDelete(db.Model(&systemModel.AISystemDept{})).Select("id, parent_id").Find(&depts).Error; err != nil {
		return nil, err
	}

	childrenOf := map[uint][]uint{}
	for _, dept := range depts {
		if dept.ParentID == nil {
			continue
		}
		childrenOf[*dept.ParentID] = append(childrenOf[*dept.ParentID], dept.ID)
	}

	result := []uint{}
	queue := []uint{deptID}
	seen := map[uint]struct{}{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if _, ok := seen[current]; ok {
			continue
		}
		seen[current] = struct{}{}
		result = append(result, current)
		queue = append(queue, childrenOf[current]...)
	}
	return result, nil
}
