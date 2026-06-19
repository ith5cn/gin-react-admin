package middleware

import (
	"net/http"
	"runtime/debug"
	"server/model/common/code"
	"server/model/common/response"
	loggerInit "server/setup/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 捕获 panic，记录堆栈日志，并返回统一 JSON。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				loggerInit.Logger.Get().Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.ByteString("stack", debug.Stack()),
				)

				response.FailWithAbort(c, http.StatusInternalServerError, code.SystemError)
			}
		}()

		c.Next()
	}
}
