package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// logRoutes 注册登录日志与操作日志路由。
func logRoutes(system *gin.RouterGroup) {
	system.GET("/login-log/index", middleware.Perm("system/login-log/index"), systemApi.LoginLogList)
	system.DELETE("/login-log/:id", middleware.Perm("system/login-log/destroy"), systemApi.LoginLogDelete)
	system.GET("/oper-log/index", middleware.Perm("system/oper-log/index"), systemApi.OperLogList)
}
