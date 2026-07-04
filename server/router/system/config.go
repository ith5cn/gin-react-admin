package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

// configRoutes 注册配置分组与配置项路由。
func configRoutes(system *gin.RouterGroup) {
	system.GET("/config-group/index", systemApi.ConfigGroupList)
	system.POST("/config-group", systemApi.CreateConfigGroup)
	system.PUT("/config-group/:id", systemApi.UpdateConfigGroup)
	system.DELETE("/config-group/:id", systemApi.DeleteConfigGroup)

	system.GET("/config/index", systemApi.ConfigList)
	system.GET("/config/get-config-info", systemApi.ConfigInfo)
	system.POST("/config", systemApi.CreateConfig)
	system.POST("/config/batch-update", systemApi.BatchUpdateConfig)
	system.PUT("/config/:id", systemApi.UpdateConfig)
	system.DELETE("/config/:id", systemApi.DeleteConfig)
}
