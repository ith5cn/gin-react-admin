package system

import (
	"server/model/common/response"
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	data, err := systemService.UserList(queryMap(c))
	successOrFail(c, data, err)
}

func CreateUser(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.UserPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateUser(payload)
	successOrFail(c, result, err)
}

func UpdateUser(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.UserPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateUser(c.Param("id"), payload)
	successOrFail(c, result, err)
}

func DeleteUser(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteUser(c.Param("id")))
}

func RefreshUserCache(c *gin.Context) {
	response.Success(c, true)
}

func SetUserPassword(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.SetPasswordPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.SetUserPassword(c.Param("id"), payload.Password))
}

func BindUserRole(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.IDsPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.BindUserRolesStringID(c.Param("id"), payload.IDs))
}

func UserAuthList(c *gin.Context) {
	result, err := systemService.UserAuthList()
	successOrFail(c, result, err)
}
