package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// databaseRoutes 注册数据库维护路由，单独挂在 /data 分组下。
func databaseRoutes(PrivateGroup *gin.RouterGroup) {
	data := PrivateGroup.Group("/data")
	data.GET("/database/index", middleware.Perm("system/database/index"), systemApi.DatabaseTableList)
	data.GET("/database/columns/:tableName", middleware.Perm("system/database/columns"), systemApi.DatabaseTableColumns)
	data.POST("/database/fragment", middleware.Perm("system/database/fragment"), systemApi.DatabaseClearFragments)
	data.POST("/database/optimize", middleware.Perm("system/database/optimize"), systemApi.DatabaseOptimizeTables)
	data.GET("/database/recycle", middleware.Perm("system/database/recycle"), systemApi.DatabaseRecycleList)
	data.POST("/database/recycle/recover", middleware.Perm("system/database/recover"), systemApi.DatabaseRecycleRecover)
	data.POST("/database/recycle/destroy", middleware.Perm("system/database/destroy"), systemApi.DatabaseRecycleDestroy)
}
