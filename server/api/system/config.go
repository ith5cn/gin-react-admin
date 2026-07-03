package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ConfigGroupList(c *gin.Context) {
	result, err := systemService.ConfigGroupList(queryMap(c))
	successOrFail(c, result, err)
}

func CreateConfigGroup(c *gin.Context) {
	var payload systemRequest.ConfigGroupPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.CreateConfigGroup(data)
	successOrFail(c, result, err)
}

func UpdateConfigGroup(c *gin.Context) {
	var payload systemRequest.ConfigGroupPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.UpdateConfigGroup(c.Param("id"), data)
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
	var payload systemRequest.ConfigPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.CreateConfig(data)
	successOrFail(c, result, err)
}

func UpdateConfig(c *gin.Context) {
	var payload systemRequest.ConfigPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.UpdateConfig(c.Param("id"), data)
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
	var payload systemRequest.BatchUpdateConfigPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	groupID := uint64(0)
	switch value := data["group_id"].(type) {
	case float64:
		groupID = uint64(value)
	case string:
		groupID, _ = strconv.ParseUint(value, 10, 64)
	}
	items, _ := data["config"].([]interface{})
	configs := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			configs = append(configs, m)
		}
	}
	successOrFail(c, true, systemService.BatchUpdateConfig(uint(groupID), configs))
}
