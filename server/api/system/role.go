package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func RoleList(c *gin.Context) {
	result, err := systemService.RoleList(queryMap(c))
	successOrFail(c, result, err)
}

func CreateRole(c *gin.Context) {
	var payload systemRequest.RolePayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.CreateRole(data)
	successOrFail(c, result, err)
}

func UpdateRole(c *gin.Context) {
	var payload systemRequest.RolePayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.UpdateRole(c.Param("id"), data)
	successOrFail(c, result, err)
}

func DeleteRole(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteRole(c.Param("id")))
}

func BindRoleMenu(c *gin.Context) {
	var payload systemRequest.RoleMenuPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		successOrFail(c, nil, err)
		return
	}
	successOrFail(c, true, systemService.BindRoleMenus(c.Param("id"), payload.IDs))
}

func RoleAccess(c *gin.Context) {
	result, err := systemService.RoleAccess()
	successOrFail(c, result, err)
}
