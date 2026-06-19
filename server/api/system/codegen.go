package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func CodegenList(c *gin.Context) {
	result, err := systemService.CodegenList(queryMap(c))
	successOrFail(c, result, err)
}

func CodegenDatasources(c *gin.Context) {
	result, err := systemService.CodegenDatasources()
	successOrFail(c, result, err)
}

func CodegenDBTables(c *gin.Context) {
	result, err := systemService.CodegenDBTables(queryMap(c))
	successOrFail(c, result, err)
}

func CodegenImportTables(c *gin.Context) {
	var data systemRequest.CodegenImportPayload
	if err := c.ShouldBindJSON(&data); err != nil {
		successOrFail(c, nil, err)
		return
	}
	result, err := systemService.CodegenImportTables(data)
	successOrFail(c, result, err)
}

func CodegenDelete(c *gin.Context) {
	data, ok := bindJSONMap(c)
	if !ok {
		return
	}
	successOrFail(c, map[string]interface{}{}, systemService.CodegenDelete(systemServiceIDs(data["ids"])))
}

func CodegenDetail(c *gin.Context) {
	result, err := systemService.CodegenDetail(c.Param("id"))
	successOrFail(c, result, err)
}

func CodegenUpdate(c *gin.Context) {
	data, ok := bindJSONMap(c)
	if !ok {
		return
	}
	result, err := systemService.CodegenUpdate(c.Param("id"), data)
	successOrFail(c, result, err)
}

func CodegenPreview(c *gin.Context) {
	result, err := systemService.CodegenPreview(c.Param("id"))
	successOrFail(c, result, err)
}

func CodegenGenerate(c *gin.Context) {
	result, err := systemService.CodegenGenerate(c.Param("id"))
	successOrFail(c, result, err)
}
