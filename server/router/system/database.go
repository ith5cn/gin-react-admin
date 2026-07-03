package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func databaseRoutes(PrivateGroup *gin.RouterGroup) {
	data := PrivateGroup.Group("/data")
	data.GET("/database/index", systemApi.DatabaseTableList)
	data.GET("/database/columns/:tableName", systemApi.DatabaseTableColumns)
	data.POST("/database/fragment", systemApi.DatabaseClearFragments)
	data.POST("/database/optimize", systemApi.DatabaseOptimizeTables)
	data.GET("/database/recycle", systemApi.DatabaseRecycleList)
	data.POST("/database/recycle/recover", systemApi.DatabaseRecycleRecover)
	data.POST("/database/recycle/destroy", systemApi.DatabaseRecycleDestroy)
}
