package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func DeptList(c *gin.Context) {
	result, err := systemService.DeptList(queryMap(c))
	successOrFail(c, result, err)
}

func CreateDept(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DeptPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateDept(payload)
	successOrFail(c, result, err)
}

func UpdateDept(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DeptPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateDept(c.Param("id"), payload)
	successOrFail(c, result, err)
}

func DeleteDept(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDept(c.Param("id")))
}

func DeptAccess(c *gin.Context) {
	result, err := systemService.DeptAccess(c.Query("tree") == "true")
	successOrFail(c, result, err)
}
