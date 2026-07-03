package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func DictTypeList(c *gin.Context) {
	result, err := systemService.DictTypeList(queryMap(c))
	successOrFail(c, result, err)
}

func CreateDictType(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DictTypePayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateDictType(payload)
	successOrFail(c, result, err)
}

func UpdateDictType(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DictTypePayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateDictType(c.Param("id"), payload)
	successOrFail(c, result, err)
}

func DeleteDictType(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDictType(c.Param("id")))
}

func DictDataList(c *gin.Context) {
	result, err := systemService.DictDataList(queryMap(c))
	successOrFail(c, result, err)
}

func DictAll(c *gin.Context) {
	result, err := systemService.DictAll()
	successOrFail(c, result, err)
}

func CreateDictData(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DictDataPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateDictData(payload)
	successOrFail(c, result, err)
}

func UpdateDictData(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DictDataPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateDictData(c.Param("id"), payload)
	successOrFail(c, result, err)
}

func DeleteDictData(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDictData(c.Param("id")))
}
