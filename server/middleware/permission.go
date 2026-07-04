package middleware

import (
	"net/http"
	"server/model/common/code"
	"server/model/common/response"
	systemService "server/service/system"
	loggerInit "server/setup/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Perm 返回接口级权限校验中间件，必须挂在 JWTAuth 之后。
// 前端的权限码只控制按钮显隐（拦君子），这里才是真正的防线（拦小人）：
// 没有权限码的用户即使直接用 curl 调接口也会被 403 拒绝。
//
// 用法：system.DELETE("/user/:id", middleware.Perm("system/user/destroy"), api.DeleteUser)
func Perm(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, ok := c.Get("user_id")
		userID, castOK := userIDValue.(uint)
		if !ok || !castOK {
			response.FailWithAbort(c, http.StatusUnauthorized, code.LoginRequired)
			return
		}

		allowed, err := systemService.UserHasPermission(userID, permCode)
		if err != nil {
			// 查询失败时拒绝放行（fail closed）：宁可误伤也不能在故障时敞开权限。
			loggerInit.Logger.Get().Error("permission check failed",
				zap.Uint("user_id", userID),
				zap.String("perm", permCode),
				zap.Error(err))
			response.FailWithAbort(c, http.StatusInternalServerError, code.SystemError)
			return
		}
		if !allowed {
			response.FailWithAbort(c, http.StatusForbidden, code.PermissionDenied)
			return
		}

		c.Next()
	}
}
