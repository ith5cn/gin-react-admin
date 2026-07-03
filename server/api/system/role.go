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
	payload, ok := bindJSON[systemRequest.RolePayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateRole(payload)
	successOrFail(c, result, err)
}

func UpdateRole(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.RolePayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateRole(c.Param("id"), payload)
	successOrFail(c, result, err)
}

func DeleteRole(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteRole(c.Param("id")))
}

func BindRoleMenu(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.IDsPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.BindRoleMenus(c.Param("id"), payload.IDs))
}

func RoleAccess(c *gin.Context) {
	result, err := systemService.RoleAccess()
	successOrFail(c, result, err)
}
