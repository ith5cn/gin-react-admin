package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// DictTypeList 字典类型分页列表。
func DictTypeList(c *gin.Context) {
	result, err := systemService.DictTypeList(queryMap(c))
	successOrFail(c, result, err)
}

// CreateDictType 新增字典类型。
func CreateDictType(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DictTypePayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateDictType(payload)
	successOrFail(c, result, err)
}

// UpdateDictType 更新字典类型。
func UpdateDictType(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DictTypePayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateDictType(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteDictType 删除字典类型。
func DeleteDictType(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDictType(c.Param("id")))
}

// DictDataList 字典数据分页列表。
func DictDataList(c *gin.Context) {
	result, err := systemService.DictDataList(queryMap(c))
	successOrFail(c, result, err)
}

// DictAll 返回全部启用字典（前端启动时缓存）。
func DictAll(c *gin.Context) {
	result, err := systemService.DictAll()
	successOrFail(c, result, err)
}

// CreateDictData 新增字典数据。
func CreateDictData(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DictDataPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateDictData(payload)
	successOrFail(c, result, err)
}

// UpdateDictData 更新字典数据。
func UpdateDictData(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.DictDataPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateDictData(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteDictData 删除字典数据。
func DeleteDictData(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDictData(c.Param("id")))
}
