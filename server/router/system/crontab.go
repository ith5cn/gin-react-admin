package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// crontabRoutes 注册定时任务路由（列表/增删改/手动执行/执行日志）。
func crontabRoutes(system *gin.RouterGroup) {
	system.GET("/crontab/index", middleware.Perm("system/crontab/index"), systemApi.CrontabList)
	system.POST("/crontab", middleware.Perm("system/crontab/create"), systemApi.CreateCrontab)
	system.PUT("/crontab/:id", middleware.Perm("system/crontab/update"), systemApi.UpdateCrontab)
	system.DELETE("/crontab/:id", middleware.Perm("system/crontab/destroy"), systemApi.DeleteCrontab)
	system.POST("/crontab/run/:id", middleware.Perm("system/crontab/run"), systemApi.RunCrontab)
	system.GET("/crontab/log/index", middleware.Perm("system/crontab/index"), systemApi.CrontabLogList)
}
