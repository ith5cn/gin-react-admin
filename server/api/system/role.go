package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// RoleList 角色树查询。
func RoleList(c *gin.Context) {
	result, err := systemService.RoleList(queryMap(c))
	successOrFail(c, result, err)
}

// CreateRole 新增角色。
func CreateRole(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.RolePayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateRole(payload)
	successOrFail(c, result, err)
}

// UpdateRole 更新角色。
func UpdateRole(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.RolePayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateRole(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteRole 删除角色。
func DeleteRole(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteRole(c.Param("id")))
}

// BindRoleMenu 为角色绑定菜单权限（整体替换）。
func BindRoleMenu(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.IDsPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.BindRoleMenus(c.Param("id"), payload.IDs))
}

// RoleAccess 角色下拉数据。
func RoleAccess(c *gin.Context) {
	result, err := systemService.RoleAccess()
	successOrFail(c, result, err)
}
