package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeptList(c *gin.Context) {
	result, err := systemService.DeptList(queryMap(c))
	successOrFail(c, result, err)
}

func CreateDept(c *gin.Context) {
	var payload systemRequest.DeptPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.CreateDept(data)
	successOrFail(c, result, err)
}

func UpdateDept(c *gin.Context) {
	var payload systemRequest.DeptPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.UpdateDept(c.Param("id"), data)
	successOrFail(c, result, err)
}

func DeleteDept(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDept(c.Param("id")))
}

func DeptAccess(c *gin.Context) {
	result, err := systemService.DeptAccess(c.Query("tree") == "true")
	successOrFail(c, result, err)
}

func PostList(c *gin.Context) {
	result, err := systemService.PostList(queryMap(c))
	successOrFail(c, result, err)
}
func PostAccess(c *gin.Context) {
	result, err := systemService.PostAccess()
	successOrFail(c, result, err)
}
func CreatePost(c *gin.Context) {
	var payload systemRequest.PostPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.CreatePost(data)
		successOrFail(c, result, err)
	}
}
func UpdatePost(c *gin.Context) {
	var payload systemRequest.PostPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.UpdatePost(c.Param("id"), data)
		successOrFail(c, result, err)
	}
}
func DeletePost(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeletePost(c.Param("id")))
}

func DictTypeList(c *gin.Context) {
	result, err := systemService.DictTypeList(queryMap(c))
	successOrFail(c, result, err)
}
func CreateDictType(c *gin.Context) {
	var payload systemRequest.DictTypePayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.CreateDictType(data)
		successOrFail(c, result, err)
	}
}
func UpdateDictType(c *gin.Context) {
	var payload systemRequest.DictTypePayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.UpdateDictType(c.Param("id"), data)
		successOrFail(c, result, err)
	}
}
func DeleteDictType(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDictType(c.Param("id")))
}

func DictDataList(c *gin.Context) {
	result, err := systemService.DictDataList(queryMap(c))
	successOrFail(c, result, err)
}
func DictAll(c *gin.Context) { result, err := systemService.DictAll(); successOrFail(c, result, err) }
func CreateDictData(c *gin.Context) {
	var payload systemRequest.DictDataPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.CreateDictData(data)
		successOrFail(c, result, err)
	}
}
func UpdateDictData(c *gin.Context) {
	var payload systemRequest.DictDataPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.UpdateDictData(c.Param("id"), data)
		successOrFail(c, result, err)
	}
}
func DeleteDictData(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteDictData(c.Param("id")))
}

func ConfigGroupList(c *gin.Context) {
	result, err := systemService.ConfigGroupList(queryMap(c))
	successOrFail(c, result, err)
}
func CreateConfigGroup(c *gin.Context) {
	var payload systemRequest.ConfigGroupPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.CreateConfigGroup(data)
		successOrFail(c, result, err)
	}
}
func UpdateConfigGroup(c *gin.Context) {
	var payload systemRequest.ConfigGroupPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.UpdateConfigGroup(c.Param("id"), data)
		successOrFail(c, result, err)
	}
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
	if ok {
		result, err := systemService.CreateConfig(data)
		successOrFail(c, result, err)
	}
}
func UpdateConfig(c *gin.Context) {
	var payload systemRequest.ConfigPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if ok {
		result, err := systemService.UpdateConfig(c.Param("id"), data)
		successOrFail(c, result, err)
	}
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

func LoginLogList(c *gin.Context) {
	result, err := systemService.LoginLogList(queryMap(c))
	successOrFail(c, result, err)
}
func OperLogList(c *gin.Context) {
	result, err := systemService.OperLogList(queryMap(c))
	successOrFail(c, result, err)
}

func systemServiceIDs(value interface{}) []uint {
	items, ok := value.([]interface{})
	if !ok {
		return []uint{}
	}
	result := make([]uint, 0, len(items))
	for _, item := range items {
		switch v := item.(type) {
		case float64:
			result = append(result, uint(v))
		case string:
			if id, err := strconv.ParseUint(v, 10, 64); err == nil {
				result = append(result, uint(id))
			}
		}
	}
	return result
}
