package system

import (
	"errors"
	commonResponse "server/model/common/response"
	systemModel "server/model/system"

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

func userColumns() map[string]string {
	return map[string]string{"username": "username", "user_type": "user_type", "userType": "user_type", "nickname": "nickname", "phone": "phone", "email": "email", "avatar": "avatar", "signed": "signed", "dashboard": "dashboard", "deptId": "dept_id", "status": "status", "remark": "remark", "backendSetting": "backend_setting"}
}
