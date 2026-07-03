package system

import (
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"

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

func CreateUser(payload systemRequest.UserPayload) (*systemModel.AISystemUser, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	if payload.Password == nil || *payload.Password == "" {
		return nil, ErrPasswordRequired
	}
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
	// Roles 为 nil 表示前端未提交角色字段，保持原绑定不动。
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

func DeleteUser(id string) error {
	return deleteByID(&systemModel.AISystemUser{}, id)
}

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
