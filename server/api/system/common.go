package system

import (
	"encoding/json"
	"net/http"
	"server/model/common/code"
	"server/model/common/response"

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

func QueryMap(c *gin.Context) map[string]string {
	return queryMap(c)
}

func BindJSONMap(c *gin.Context) (map[string]interface{}, bool) {
	return bindJSONMap(c)
}

func SuccessOrFail(c *gin.Context, data interface{}, err error) {
	successOrFail(c, data, err)
}
