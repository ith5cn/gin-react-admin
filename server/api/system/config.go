package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func ConfigGroupList(c *gin.Context) {
	result, err := systemService.ConfigGroupList(queryMap(c))
	successOrFail(c, result, err)
}

func CreateConfigGroup(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.ConfigGroupPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateConfigGroup(payload)
	successOrFail(c, result, err)
}

func UpdateConfigGroup(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.ConfigGroupPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateConfigGroup(c.Param("id"), payload)
	successOrFail(c, result, err)
}

func DeleteConfigGroup(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteConfigGroup(c.Param("id")))
}

func ConfigList(c *gin.Context) {
	result, err := systemService.ConfigList(queryMap(c))
	successOrFail(c, result, err)
}

func CreateConfig(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.ConfigPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateConfig(payload)
	successOrFail(c, result, err)
}

func UpdateConfig(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.ConfigPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateConfig(c.Param("id"), payload)
	successOrFail(c, result, err)
}

func DeleteConfig(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteConfig(c.Param("id")))
}

func ConfigInfo(c *gin.Context) {
	result, err := systemService.ConfigInfo(c.Query("code"))
	successOrFail(c, result, err)
}

func BatchUpdateConfig(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.BatchUpdateConfigPayload](c)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.BatchUpdateConfig(payload))
}
