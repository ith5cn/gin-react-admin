package system

import (
	systemModel "server/model/system"
	systemRequest "server/model/system/request"

	"golang.org/x/crypto/bcrypt"
)

// 本文件是"个人中心"业务：用户操作自己的资料和密码。
// 与 user_manage.go（管理员管理别人）的关键区别：
// 目标用户 ID 永远取自 JWT 里的当前登录人，绝不信任前端传来的 id——这是防越权的根本。

// UpdateProfile 更新当前登录用户自己的资料。
// 只开放安全字段（昵称/联系方式/头像/签名/偏好设置），
// 用户名、状态、部门、角色这类管理属性不在此接口范围内。
func UpdateProfile(userID uint, payload systemRequest.ProfilePayload) (*systemModel.AISystemUser, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	setColumn(data, "nickname", payload.Nickname)
	setColumn(data, "phone", payload.Phone)
	setColumn(data, "email", payload.Email)
	setColumn(data, "avatar", payload.Avatar)
	setColumn(data, "signed", payload.Signed)
	setColumn(data, "dashboard", payload.Dashboard)
	setColumn(data, "backend_setting", payload.BackendSetting)
	setDefaultTimes(data, false)

	if err := db.Model(&systemModel.AISystemUser{}).Where("id = ?", userID).Updates(data).Error; err != nil {
		return nil, err
	}

	var user systemModel.AISystemUser
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ChangePassword 用户修改自己的密码：必须先验证原密码，防止会话被劫持后直接改密锁死账号主人。
func ChangePassword(userID uint, oldPassword, newPassword string) error {
	if newPassword == "" {
		return ErrPasswordRequired
	}

	db, err := systemDB()
	if err != nil {
		return err
	}

	var user systemModel.AISystemUser
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)) != nil {
		return ErrOldPasswordWrong
	}

	return SetUserPassword(uintToString(userID), newPassword)
}

// uintToString 是 parseUint 的反向小工具。
func uintToString(value uint) string {
	return stringFromAny(value)
}
