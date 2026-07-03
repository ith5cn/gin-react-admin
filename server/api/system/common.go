package system

import (
	"errors"
	"net/http"
	"server/model/common/code"
	"server/model/common/response"
	systemService "server/service/system"
	loggerInit "server/setup/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// bindJSON 绑定类型化请求体，失败时统一返回参数错误响应。
func bindJSON[T any](c *gin.Context) (T, bool) {
	var payload T
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return payload, false
	}
	return payload, true
}

// successOrFail 统一收敛 service 层错误：
// BizError 是给用户看的业务错误，透传消息（HTTP 200 + OperationFailed）；
// 其余一律视为内部错误，记日志后只返回泛化 SystemError，不向客户端泄露 SQL 等细节。
func successOrFail(c *gin.Context, data interface{}, err error) {
	if err == nil {
		response.Success(c, data)
		return
	}

	var bizErr *systemService.BizError
	if errors.As(err, &bizErr) {
		response.Fail(c, code.OperationFailed, bizErr.Error())
		return
	}

	loggerInit.Logger.Get().Error("request failed",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Error(err))
	response.FailWithHTTP(c, http.StatusInternalServerError, code.SystemError)
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
