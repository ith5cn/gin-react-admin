package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// ConfigGroupList 配置分组分页列表。
func ConfigGroupList(c *gin.Context) {
	result, err := systemService.ConfigGroupList(queryMap(c))
	successOrFail(c, result, err)
}

// CreateConfigGroup 新增配置分组。
func CreateConfigGroup(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.ConfigGroupPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateConfigGroup(payload)
	successOrFail(c, result, err)
}

// UpdateConfigGroup 更新配置分组。
func UpdateConfigGroup(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.ConfigGroupPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateConfigGroup(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteConfigGroup 删除配置分组。
func DeleteConfigGroup(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteConfigGroup(c.Param("id")))
}

// ConfigList 配置项分页列表。
func ConfigList(c *gin.Context) {
	result, err := systemService.ConfigList(queryMap(c))
	successOrFail(c, result, err)
}

// CreateConfig 新增配置项。
func CreateConfig(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.ConfigPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateConfig(payload)
	successOrFail(c, result, err)
}

// UpdateConfig 更新配置项。
func UpdateConfig(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.ConfigPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateConfig(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteConfig 删除配置项。
func DeleteConfig(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteConfig(c.Param("id")))
}

// ConfigInfo 按分组编码或 key 查询配置键值。
func ConfigInfo(c *gin.Context) {
	result, err := systemService.ConfigInfo(c.Query("code"))
	successOrFail(c, result, err)
}

// BatchUpdateConfig 批量保存配置项。
func BatchUpdateConfig(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.BatchUpdateConfigPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.BatchUpdateConfig(payload))
}
