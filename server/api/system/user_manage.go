package system

import (
	"server/model/common/response"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	data, err := systemService.UserList(queryMap(c))
	successOrFail(c, data, err)
}

func CreateUser(c *gin.Context) {
	data, ok := bindJSONMap(c)
	if !ok {
		return
	}
	result, err := systemService.CreateUser(data)
	successOrFail(c, result, err)
}

func UpdateUser(c *gin.Context) {
	data, ok := bindJSONMap(c)
	if !ok {
		return
	}
	result, err := systemService.UpdateUser(c.Param("id"), data)
	successOrFail(c, result, err)
}

func DeleteUser(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteUser(c.Param("id")))
}

func RefreshUserCache(c *gin.Context) {
	response.Success(c, true)
}

func SetUserPassword(c *gin.Context) {
	data, ok := bindJSONMap(c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.SetUserPassword(c.Param("id"), data))
}

func BindUserRole(c *gin.Context) {
	data, ok := bindJSONMap(c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.BindUserRolesStringID(c.Param("id"), systemServiceIDs(data["ids"])))
}

func UserAuthList(c *gin.Context) {
	result, err := systemService.UserAuthList()
	successOrFail(c, result, err)
}
