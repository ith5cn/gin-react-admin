package system

import (
	"server/model/common/code"
	"server/model/common/response"
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func DatabaseTableList(c *gin.Context) {
	result, err := systemService.DatabaseTableList(queryMap(c))
	successOrFail(c, result, err)
}

func DatabaseTableColumns(c *gin.Context) {
	result, err := systemService.DatabaseTableColumns(c.Param("tableName"))
	successOrFail(c, result, err)
}

func DatabaseOptimizeTables(c *gin.Context) {
	var data systemRequest.DatabaseTablesPayload
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	result, err := systemService.DatabaseOptimizeTables(data)
	successOrFail(c, result, err)
}

func DatabaseClearFragments(c *gin.Context) {
	var data systemRequest.DatabaseTablesPayload
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	result, err := systemService.DatabaseClearFragments(data)
	successOrFail(c, result, err)
}

func DatabaseRecycleList(c *gin.Context) {
	result, err := systemService.DatabaseRecycleList(queryMap(c))
	successOrFail(c, result, err)
}

func DatabaseRecycleRecover(c *gin.Context) {
	var data systemRequest.DatabaseRecyclePayload
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	successOrFail(c, map[string]interface{}{}, systemService.DatabaseRecycleRecover(data))
}

func DatabaseRecycleDestroy(c *gin.Context) {
	var data systemRequest.DatabaseRecyclePayload
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	successOrFail(c, map[string]interface{}{}, systemService.DatabaseRecycleDestroy(data))
}
