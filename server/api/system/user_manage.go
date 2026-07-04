package system

import (
	"server/model/common/response"
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// UserList 用户分页列表，按当前操作者的数据权限过滤可见范围。
func UserList(c *gin.Context) {
	operatorID, ok := currentUserID(c)
	if !ok {
		return
	}
	data, err := systemService.UserList(operatorID, queryMap(c))
	successOrFail(c, data, err)
}

// CreateUser 新增用户。
func CreateUser(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.UserPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateUser(payload)
	successOrFail(c, result, err)
}

// UpdateUser 更新用户资料。
func UpdateUser(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.UserPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateUser(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteUser 删除用户。
func DeleteUser(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteUser(c.Param("id")))
}

// RefreshUserCache 预留接口：当前项目未做用户缓存，直接返回成功，
// 保留路由是为了兼容前端已有的按钮。
func RefreshUserCache(c *gin.Context) {
	response.Success(c, true)
}

// SetUserPassword 管理员重置用户密码。
func SetUserPassword(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.SetPasswordPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.SetUserPassword(c.Param("id"), payload.Password))
}

// BindUserRole 为用户绑定角色（整体替换）。
func BindUserRole(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.IDsPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.BindUserRolesStringID(c.Param("id"), payload.IDs))
}

// UserAuthList 用户下拉数据。
func UserAuthList(c *gin.Context) {
	result, err := systemService.UserAuthList()
	successOrFail(c, result, err)
}
