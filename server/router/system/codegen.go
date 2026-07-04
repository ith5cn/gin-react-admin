package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// codegenRoutes 注册代码生成器路由。
func codegenRoutes(system *gin.RouterGroup) {
	system.GET("/codegen/index", middleware.Perm("system/codegen/index"), systemApi.CodegenList)
	system.GET("/codegen/datasources", middleware.Perm("system/codegen/access"), systemApi.CodegenDatasources)
	system.GET("/codegen/db-tables", middleware.Perm("system/codegen/access"), systemApi.CodegenDBTables)
	system.POST("/codegen/importTables", middleware.Perm("system/codegen/access"), systemApi.CodegenImportTables)
	system.POST("/codegen/delete", middleware.Perm("system/codegen/access"), systemApi.CodegenDelete)
	system.GET("/codegen/detail/:id", middleware.Perm("system/codegen/access"), systemApi.CodegenDetail)
	system.PUT("/codegen/:id", middleware.Perm("system/codegen/access"), systemApi.CodegenUpdate)
	system.POST("/codegen/generate/:id", middleware.Perm("system/codegen/access"), systemApi.CodegenGenerate)
	system.GET("/codegen/preview/:id", middleware.Perm("system/codegen/access"), systemApi.CodegenPreview)
}
