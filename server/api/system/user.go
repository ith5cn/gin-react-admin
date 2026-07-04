package system

import (
	"errors"
	"net/http"
	"server/model/common/code"
	"server/model/common/response"
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// Login 是系统模块的登录 HTTP 入口。
// api 层只负责参数绑定、调用 service、组织 HTTP 响应，不直接写业务规则。
func Login(c *gin.Context) {
	var loginReq systemRequest.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}

	// ClientIP 会正确处理反向代理头（X-Forwarded-For 等），UserAgent 用于登录日志。
	tokens, err := systemService.Login(loginReq.UserName, loginReq.Password, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		// 登录失败统一返回泛化错误，避免暴露用户是否存在。
		status := http.StatusInternalServerError
		message := "login failed"
		bizCode := code.SystemError
		if errors.Is(err, systemService.ErrLoginFailed) {
			status = http.StatusUnauthorized
			message = "user_name or password is incorrect"
			bizCode = code.LoginRequired
		}

		response.FailWithHTTP(c, status, bizCode, message)
		return
	}

	// data 中直接返回 TokenPair，字段名由 utils.TokenPair 的 json tag 控制。
	response.Success(c, tokens)
}

// CurrentUser 返回当前登录用户的前端初始化上下文。
// 该接口需要经过 JWT 中间件，user_id 由 middleware 写入 Gin context。
func CurrentUser(c *gin.Context) {
	userIDValue, ok := c.Get("user_id")
	if !ok {
		response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
		return
	}

	userContext, err := systemService.CurrentUserContext(userID)
	if err != nil {
		if errors.Is(err, systemService.ErrLoginFailed) {
			response.FailWithHTTP(c, http.StatusUnauthorized, code.LoginRequired)
			return
		}

		response.FailWithHTTP(c, http.StatusInternalServerError, code.SystemError)
		return
	}

	response.Success(c, userContext)
}
