package middleware

import (
	"errors"
	"net/http"
	"server/model/common/code"
	"server/model/common/response"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth 是业务接口使用的 access token 鉴权中间件。
// 它只校验 access token；refresh token 只允许交给刷新接口处理。
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := BearerToken(c.GetHeader("Authorization"))
		if !ok {
			response.FailWithAbort(c, http.StatusUnauthorized, code.LoginRequired)
			return
		}

		claims, err := utils.ValidateToken(tokenString, utils.TokenTypeAccess)
		if err != nil {
			if errors.Is(err, utils.ErrTokenExpired) {
				response.FailWithAbort(c, http.StatusUnauthorized, code.AccessTokenExpired)
				return
			}

			response.FailWithAbort(c, http.StatusUnauthorized, code.LoginRequired)
			return
		}

		// 后续 handler 可以通过 c.Get("user_id") / c.Get("username") 获取当前登录用户。
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// BearerToken 从 Authorization 请求头中解析 Bearer token。
func BearerToken(authorization string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
	return token, token != ""
}
