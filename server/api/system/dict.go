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
	var payload systemRequest.DictTypePayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.CreateDictType(data)
	successOrFail(c, result, err)
}

func UpdateDictType(c *gin.Context) {
	var payload systemRequest.DictTypePayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.UpdateDictType(c.Param("id"), data)
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
	var payload systemRequest.DictDataPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.CreateDictData(data)
	successOrFail(c, result, err)
}

func UpdateDictData(c *gin.Context) {
	var payload systemRequest.DictDataPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.UpdateDictData(c.Param("id"), data)
	successOrFail(c, result, err)
}

func DeleteDictData(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDictData(c.Param("id")))
}
