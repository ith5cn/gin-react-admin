package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"

	"gorm.io/gorm"
)

// 本文件是"用户管理"的业务逻辑（后台管理员对用户的增删改查），
// 与 user.go（登录/当前用户上下文，属于"认证"）职责不同。

// UserList 分页查询用户列表。
// 查询条件分两类：likes（模糊匹配，生成 LIKE '%xx%'）和 equals（精确匹配）；
// map 的 key 是前端参数名（camelCase），value 是数据库列名（snake_case）。
func UserList(query map[string]string) (*commonResponse.PageResult, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	page := parsePage(query)
	// softDelete 统一追加 delete_time IS NULL，过滤掉已软删除的数据。
	base := softDelete(db.Model(&systemModel.AISystemUser{}))
	base = applyFilters(base, query,
		map[string]string{"username": "username", "nickname": "nickname", "phone": "phone", "email": "email"},
		map[string]string{"status": "status", "deptId": "dept_id"},
	)

	// 分页的标准套路：先 Count 拿总数，再 Offset/Limit 取当前页数据。
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}

	var users []systemModel.AISystemUser
	if err := base.Order("id ASC").Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(&users).Error; err != nil {
		return nil, err
	}

	// 逐个用户补查角色 ID（前端列表要显示角色）。
	// 注意：这是典型的 N+1 查询写法，数据量大时应改为一次 IN 查询后按 user_id 分组。
	for i := range users {
		users[i].Roles, _ = RoleIDsByUserID(users[i].ID)
	}

	return &commonResponse.PageResult{List: users, Total: total}, nil
}

// CreateUser 创建用户：密码必填并做 bcrypt 加密，可同时绑定角色。
func CreateUser(payload systemRequest.UserPayload) (*systemModel.AISystemUser, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	if payload.Password == nil || *payload.Password == "" {
		return nil, ErrPasswordRequired
	}
	// 数据库永远只存 bcrypt 哈希，绝不存明文密码。
	hash, err := hashPassword(*payload.Password)
	if err != nil {
		return nil, err
	}

	data := userPayloadData(payload)
	data["password"] = hash
	setDefaultTimes(data, true)

	if err := db.Model(&systemModel.AISystemUser{}).Create(data).Error; err != nil {
		return nil, err
	}

	// 用 map 方式 Create 拿不到自增 ID，这里按唯一的 username 回查一次。
	var user systemModel.AISystemUser
	if err := db.Where("username = ?", data["username"]).First(&user).Error; err != nil {
		return nil, err
	}
	roles := payload.Roles
	if len(roles) > 0 {
		if err := BindUserRoles(user.ID, roles); err != nil {
			return nil, err
		}
	}
	user.Roles = roles
	return &user, nil
}

// UpdateUser 更新用户资料并按需重绑角色（部分更新语义：nil 字段不改动）。
func UpdateUser(id string, payload systemRequest.UserPayload) (*systemModel.AISystemUser, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	// 更新接口不改密码，密码走独立的 set-password 接口。
	data := userPayloadData(payload)
	setDefaultTimes(data, false)
	if len(data) > 0 {
		if err := db.Model(&systemModel.AISystemUser{}).Where("id = ?", id).Updates(data).Error; err != nil {
			return nil, err
		}
	}
	// Roles 为 nil 表示前端未提交角色字段，保持原绑定不动；
	// 空数组 [] 则表示"清空全部角色"，两者语义不同。
	if payload.Roles != nil {
		if err := BindUserRolesStringID(id, payload.Roles); err != nil {
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

// DeleteUser 按 ID 删除用户。
func DeleteUser(id string) error {
	return deleteByID(&systemModel.AISystemUser{}, id)
}

// SetUserPassword 管理员重置指定用户的密码。
func SetUserPassword(id string, password string) error {
	if password == "" {
		return ErrPasswordRequired
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

// UserAuthList 返回启用用户的 {label, value} 下拉选项，label 优先取昵称。
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
		// Nickname 是 *string（数据库可空列），解引用前必须判 nil，否则会 panic。
		if user.Nickname != nil && *user.Nickname != "" {
			label = *user.Nickname
		}
		result = append(result, map[string]interface{}{"label": label, "value": user.ID})
	}
	return result, nil
}

// BindUserRolesStringID 是 BindUserRoles 的字符串 ID 版本，方便 handler 直接传路径参数。
func BindUserRolesStringID(userID string, roleIDs []uint) error {
	id, err := parseUint(userID)
	if err != nil {
		return err
	}
	return BindUserRoles(id, roleIDs)
}

// BindUserRoles 重设用户的角色绑定。
// "先删后插"必须包在一个事务里：任何一步失败整体回滚，
// 否则可能出现旧绑定删了、新绑定没插上的中间状态。
// db.Transaction 会在回调返回 error 时自动 Rollback，返回 nil 时自动 Commit。
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

// RoleIDsByUserID 查询用户当前绑定的角色 ID 列表（走 user-role 中间表）。
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

// userPayloadData 把类型化入参转成 GORM 更新 map，只收集非 nil 字段（部分更新）。
func userPayloadData(payload systemRequest.UserPayload) map[string]interface{} {
	data := map[string]interface{}{}
	setColumn(data, "username", payload.Username)
	setColumn(data, "user_type", payload.UserType)
	setColumn(data, "nickname", payload.Nickname)
	setColumn(data, "phone", payload.Phone)
	setColumn(data, "email", payload.Email)
	setColumn(data, "avatar", payload.Avatar)
	setColumn(data, "signed", payload.Signed)
	setColumn(data, "dashboard", payload.Dashboard)
	setColumn(data, "dept_id", payload.DeptID)
	setColumn(data, "status", payload.Status)
	setColumn(data, "remark", payload.Remark)
	setColumn(data, "backend_setting", payload.BackendSetting)
	return data
}
