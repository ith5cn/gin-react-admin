package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// configRoutes 注册配置分组与配置项路由。
func configRoutes(system *gin.RouterGroup) {
	system.GET("/config-group/index", middleware.Perm("system/config/index"), systemApi.ConfigGroupList)
	system.POST("/config-group", middleware.Perm("system/config/create"), systemApi.CreateConfigGroup)
	system.PUT("/config-group/:id", middleware.Perm("system/config/update"), systemApi.UpdateConfigGroup)
	system.DELETE("/config-group/:id", middleware.Perm("system/config/destroy"), systemApi.DeleteConfigGroup)

	system.GET("/config/index", middleware.Perm("system/config/index"), systemApi.ConfigList)
	system.GET("/config/get-config-info", systemApi.ConfigInfo)
	system.POST("/config", middleware.Perm("system/config/create"), systemApi.CreateConfig)
	system.POST("/config/batch-update", middleware.Perm("system/config/update"), systemApi.BatchUpdateConfig)
	system.PUT("/config/:id", middleware.Perm("system/config/update"), systemApi.UpdateConfig)
	system.DELETE("/config/:id", middleware.Perm("system/config/destroy"), systemApi.DeleteConfig)
}
