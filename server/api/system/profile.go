package system

import (
	"net/http"
	"server/model/common/code"
	"server/model/common/response"
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// currentUserID 从 Gin context 取 JWT 中间件写入的当前用户 ID。
// 个人中心的所有接口都以它为准，绝不接收前端传的用户 id。
func currentUserID(c *gin.Context) (uint, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
		return 0, false
	}
	userID, ok := value.(uint)
	if !ok {
		response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
		return 0, false
	}
	return userID, true
}

// UpdateProfile 当前用户更新自己的资料。
func UpdateProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	payload, ok := bindJSON[systemRequest.ProfilePayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateProfile(userID, payload)
	successOrFail(c, result, err)
}

// ChangePassword 当前用户修改自己的密码。
func ChangePassword(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	payload, ok := bindJSON[systemRequest.ChangePasswordPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.ChangePassword(userID, payload.OldPassword, payload.NewPassword))
}
