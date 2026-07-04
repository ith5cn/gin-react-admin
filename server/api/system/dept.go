package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// DeptList 部门树查询。
func DeptList(c *gin.Context) {
	result, err := systemService.DeptList(queryMap(c))
	successOrFail(c, result, err)
}

// CreateDept 新增部门。
func CreateDept(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DeptPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateDept(payload)
	successOrFail(c, result, err)
}

// UpdateDept 更新部门。
func UpdateDept(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DeptPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateDept(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteDept 删除部门。
func DeleteDept(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDept(c.Param("id")))
}

// DeptAccess 部门下拉数据，?tree=true 时返回树形。
func DeptAccess(c *gin.Context) {
	result, err := systemService.DeptAccess(c.Query("tree") == "true")
	successOrFail(c, result, err)
}
