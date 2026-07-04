package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// dictRoutes 注册字典类型与字典数据路由。
func dictRoutes(system *gin.RouterGroup) {
	system.GET("/dict-type/index", middleware.Perm("system/dict-type/index"), systemApi.DictTypeList)
	system.POST("/dict-type", middleware.Perm("system/dict-type/create"), systemApi.CreateDictType)
	system.PUT("/dict-type/:id", middleware.Perm("system/dict-type/update"), systemApi.UpdateDictType)
	system.DELETE("/dict-type/:id", middleware.Perm("system/dict-type/destroy"), systemApi.DeleteDictType)

	system.GET("/dict-data/index", middleware.Perm("system/dict-type/index"), systemApi.DictDataList)
	system.GET("/dict-data/dictAll", systemApi.DictAll)
	system.POST("/dict-data", middleware.Perm("system/dict-type/create"), systemApi.CreateDictData)
	system.PUT("/dict-data/:id", middleware.Perm("system/dict-type/update"), systemApi.UpdateDictData)
	system.DELETE("/dict-data/:id", middleware.Perm("system/dict-type/destroy"), systemApi.DeleteDictData)
}
