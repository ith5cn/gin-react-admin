package middleware

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	systemService "server/service/system"
	loggerInit "server/setup/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// maxOperLogBody 是操作日志记录请求体的上限；超出部分截断，避免大请求撑爆日志表。
const maxOperLogBody = 4 << 10

// passwordFieldPattern 匹配 JSON 里的密码类字段，入库前统一脱敏。
var passwordFieldPattern = regexp.MustCompile(`"(password|oldPassword|newPassword)"\s*:\s*"[^"]*"`)

// OperLog 是操作日志中间件，挂在 JWTAuth 之后：
// 只记录写操作（POST/PUT/DELETE），GET 查询量大且无副作用，不记。
// 日志异步写入，不阻塞请求响应。
func OperLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodGet {
			c.Next()
			return
		}

		// 请求体是一次性的流，读完必须"塞回去"，否则后面的 handler 拿不到参数。
		// 文件上传（multipart）的 body 是二进制大块，跳过采集。
		requestData := ""
		contentType := c.GetHeader("Content-Type")
		if c.Request.Body != nil && !strings.Contains(contentType, "multipart/form-data") {
			peeked, err := io.ReadAll(io.LimitReader(c.Request.Body, maxOperLogBody))
			if err == nil {
				// MultiReader 把已读部分和未读残余拼回一个完整的 body。
				c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(peeked), c.Request.Body))
				requestData = passwordFieldPattern.ReplaceAllString(string(peeked), `"$1":"***"`)
			}
		}

		c.Next()

		username := c.GetString("username")
		router := c.FullPath()
		if router == "" {
			router = c.Request.URL.Path
		}
		// HandlerName 形如 server/api/system.UpdateUser，取末段当作业务名。
		serviceName := c.HandlerName()
		if idx := strings.LastIndex(serviceName, "."); idx >= 0 {
			serviceName = serviceName[idx+1:]
		}
		ip := c.ClientIP()

		// 异步写库：日志不该拖慢用户请求；goroutine 里必须自己兜 panic。
		go func() {
			defer func() {
				if r := recover(); r != nil {
					loggerInit.Logger.Get().Error("oper log goroutine panic", zap.Any("panic", r))
				}
			}()
			systemService.RecordOperLog(username, method, router, serviceName, ip, requestData)
		}()
	}
}
