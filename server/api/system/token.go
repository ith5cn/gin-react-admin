package system

import (
	"net/http"
	"server/middleware"
	"server/model/common/code"
	"server/model/common/response"
	systemRequest "server/model/system/request"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// RefreshToken 是刷新 token 的 HTTP 入口。
// 它只接收 refresh token，不依赖 access token 中间件。
func RefreshToken(c *gin.Context) {
	var req systemRequest.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}

	tokens, err := utils.RefreshToken(req.RefreshToken)
	if err != nil {
		response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
		return
	}

	response.Success(c, tokens)
}

// Logout 是退出登录入口。
// 当前 access token 必须经过 JWT 中间件校验，请求体可选携带 refresh_token 一起撤销。
func Logout(c *gin.Context) {
	accessToken, ok := middleware.BearerToken(c.GetHeader("Authorization"))
	if !ok {
		response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
		return
	}

	var req systemRequest.LogoutRequest
	_ = c.ShouldBindJSON(&req)

	if err := utils.RevokeToken(accessToken); err != nil {
		response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
		return
	}

	if req.RefreshToken != "" {
		if err := utils.RevokeToken(req.RefreshToken); err != nil {
			response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
			return
		}
	}

	response.Success(c, map[string]interface{}{})
}
