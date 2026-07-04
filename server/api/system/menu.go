package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// MenuList 菜单树查询（菜单管理页面用）。
func MenuList(c *gin.Context) {
	data, err := systemService.MenuList(queryMap(c))
	successOrFail(c, data, err)
}

// CreateMenu 新增菜单。
func CreateMenu(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.MenuPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateMenu(payload)
	successOrFail(c, result, err)
}

// UpdateMenu 更新菜单。
func UpdateMenu(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.MenuPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateMenu(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteMenu 删除菜单。
func DeleteMenu(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteMenu(c.Param("id")))
}

// AccessMenu 返回当前登录用户可访问的菜单树。
// user_id 由 JWT 中间件写入 Gin context，这里直接取用。
func AccessMenu(c *gin.Context) {
	userID, _ := c.Get("user_id")
	result, err := systemService.AccessMenu(userID.(uint))
	successOrFail(c, result, err)
}

// MenuByRole 查询指定角色绑定的菜单 ID 列表。
func MenuByRole(c *gin.Context) {
	result, err := systemService.MenuIDsByRoleID(c.Param("roleId"))
	successOrFail(c, result, err)
}
