package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func MenuList(c *gin.Context) {
	data, err := systemService.MenuList(queryMap(c))
	successOrFail(c, data, err)
}

func CreateMenu(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.MenuPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateMenu(payload)
	successOrFail(c, result, err)
}

func UpdateMenu(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.MenuPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateMenu(c.Param("id"), payload)
	successOrFail(c, result, err)
}

func DeleteMenu(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteMenu(c.Param("id")))
}

func AccessMenu(c *gin.Context) {
	userID, _ := c.Get("user_id")
	result, err := systemService.AccessMenu(userID.(uint))
	successOrFail(c, result, err)
}

func MenuByRole(c *gin.Context) {
	result, err := systemService.MenuIDsByRoleID(c.Param("roleId"))
	successOrFail(c, result, err)
}
