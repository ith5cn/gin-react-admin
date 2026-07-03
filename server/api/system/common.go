package system

import (
	"encoding/json"
	"net/http"
	"server/model/common/code"
	"server/model/common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func queryMap(c *gin.Context) map[string]string {
	result := make(map[string]string, len(c.Request.URL.Query()))
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

func bindJSONMap(c *gin.Context) (map[string]interface{}, bool) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return nil, false
	}
	return data, true
}

func bindJSONStructAsMap(c *gin.Context, payload interface{}) (map[string]interface{}, bool) {
	if err := c.ShouldBindJSON(payload); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return nil, false
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return nil, false
	}
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return nil, false
	}
	return data, true
}

func successOrFail(c *gin.Context, data interface{}, err error) {
	if err != nil {
		response.FailWithHTTP(c, http.StatusInternalServerError, code.SystemError, err.Error())
		return
	}
	response.Success(c, data)
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

// QueryMap 及以下导出包装是 codegen 生成代码（api/generated/）的稳定契约，
// 见 service/system/codegen_templates.go；改名或改签名需同步更新模板。
func QueryMap(c *gin.Context) map[string]string {
	return queryMap(c)
}

func BindJSONMap(c *gin.Context) (map[string]interface{}, bool) {
	return bindJSONMap(c)
}

func SuccessOrFail(c *gin.Context, data interface{}, err error) {
	successOrFail(c, data, err)
}
