package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// monitorRoutes 注册服务监控路由。
func monitorRoutes(system *gin.RouterGroup) {
	system.GET("/monitor/server", middleware.Perm("system/monitor/index"), systemApi.ServerMonitor)
}
