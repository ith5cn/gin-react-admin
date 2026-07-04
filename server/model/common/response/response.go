package response

import (
	"net/http"
	"server/model/common/code"

	"github.com/gin-gonic/gin"
)

// Response 是全项目统一的响应外壳：code 业务码、data 数据、msg 提示。
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// PageResult 是分页接口统一的数据结构：list 当前页数据 + total 总条数。
type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

// Result 是最底层的响应写出函数，其余快捷方法都建立在它之上。
func Result(c *gin.Context, httpCode int, bizCode int, data interface{}, msg string) {
	c.JSON(httpCode, Response{
		bizCode,
		data,
		msg,
	})
}

// Success 返回 HTTP 200 + 业务码 0。
func Success(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, code.Success, data, code.Message(code.Success))
}

// Fail 返回 HTTP 200 + 指定业务码；messages 可覆盖默认文案。
func Fail(c *gin.Context, bizCode int, messages ...string) {
	Result(c, http.StatusOK, bizCode, map[string]interface{}{}, messageOf(bizCode, messages...))
}

// FailWithHTTP 在 Fail 基础上允许自定义 HTTP 状态码（如 401、500）。
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
