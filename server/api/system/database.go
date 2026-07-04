package system

import (
	"server/model/common/code"
	"server/model/common/response"
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// DatabaseTableList 数据表信息分页列表。
func DatabaseTableList(c *gin.Context) {
	result, err := systemService.DatabaseTableList(queryMap(c))
	successOrFail(c, result, err)
}

// DatabaseTableColumns 指定表的字段结构。
func DatabaseTableColumns(c *gin.Context) {
	result, err := systemService.DatabaseTableColumns(c.Param("tableName"))
	successOrFail(c, result, err)
}

// DatabaseOptimizeTables 优化选中的数据表。
func DatabaseOptimizeTables(c *gin.Context) {
	var data systemRequest.DatabaseTablesPayload
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	result, err := systemService.DatabaseOptimizeTables(data)
	successOrFail(c, result, err)
}

// DatabaseClearFragments 清理表碎片。
func DatabaseClearFragments(c *gin.Context) {
	var data systemRequest.DatabaseTablesPayload
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	result, err := systemService.DatabaseClearFragments(data)
	successOrFail(c, result, err)
}

// DatabaseRecycleList 回收站数据列表。
func DatabaseRecycleList(c *gin.Context) {
	result, err := systemService.DatabaseRecycleList(queryMap(c))
	successOrFail(c, result, err)
}

// DatabaseRecycleRecover 恢复软删除数据。
func DatabaseRecycleRecover(c *gin.Context) {
	var data systemRequest.DatabaseRecyclePayload
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	successOrFail(c, map[string]interface{}{}, systemService.DatabaseRecycleRecover(data))
}

// DatabaseRecycleDestroy 彻底删除回收站数据。
func DatabaseRecycleDestroy(c *gin.Context) {
	var data systemRequest.DatabaseRecyclePayload
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	successOrFail(c, map[string]interface{}{}, systemService.DatabaseRecycleDestroy(data))
}
