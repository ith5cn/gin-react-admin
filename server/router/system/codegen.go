package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

// codegenRoutes 注册代码生成器路由。
func codegenRoutes(system *gin.RouterGroup) {
	system.GET("/codegen/index", systemApi.CodegenList)
	system.GET("/codegen/datasources", systemApi.CodegenDatasources)
	system.GET("/codegen/db-tables", systemApi.CodegenDBTables)
	system.POST("/codegen/importTables", systemApi.CodegenImportTables)
	system.POST("/codegen/delete", systemApi.CodegenDelete)
	system.GET("/codegen/detail/:id", systemApi.CodegenDetail)
	system.PUT("/codegen/:id", systemApi.CodegenUpdate)
	system.POST("/codegen/generate/:id", systemApi.CodegenGenerate)
	system.GET("/codegen/preview/:id", systemApi.CodegenPreview)
}
