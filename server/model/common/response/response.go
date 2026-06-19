package response

import (
	"net/http"
	"server/model/common/code"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

func Result(c *gin.Context, httpCode int, bizCode int, data interface{}, msg string) {
	c.JSON(httpCode, Response{
		bizCode,
		data,
		msg,
	})
}

func Success(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, code.Success, data, code.Message(code.Success))
}

func Fail(c *gin.Context, bizCode int, messages ...string) {
	Result(c, http.StatusOK, bizCode, map[string]interface{}{}, messageOf(bizCode, messages...))
}

func FailWithHTTP(c *gin.Context, httpCode int, bizCode int, messages ...string) {
	Result(c, httpCode, bizCode, map[string]interface{}{}, messageOf(bizCode, messages...))
}

// FailWithAbort 发送失败响应并中断请求链（用于 middleware）。
// httpCode 为 HTTP 状态码（如 401），bizCode 为业务状态码。
func FailWithAbort(c *gin.Context, httpCode int, bizCode int, messages ...string) {
	c.AbortWithStatusJSON(httpCode, Response{
		bizCode,
		map[string]interface{}{},
		messageOf(bizCode, messages...),
	})
}

func messageOf(bizCode int, messages ...string) string {
	if len(messages) > 0 && messages[0] != "" {
		return messages[0]
	}
	return code.Message(bizCode)
}
