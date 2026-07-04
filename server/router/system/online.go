package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// onlineRoutes 注册在线用户管理路由（查询/踢下线）。
func onlineRoutes(system *gin.RouterGroup) {
	system.GET("/online/index", middleware.Perm("system/online/index"), systemApi.OnlineUserList)
	system.DELETE("/online/kick/:jti", middleware.Perm("system/online/kick"), systemApi.KickOnlineUser)
}
