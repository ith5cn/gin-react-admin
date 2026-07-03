package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func dictRoutes(system *gin.RouterGroup) {
	system.GET("/dict-type/index", systemApi.DictTypeList)
	system.POST("/dict-type", systemApi.CreateDictType)
	system.PUT("/dict-type/:id", systemApi.UpdateDictType)
	system.DELETE("/dict-type/:id", systemApi.DeleteDictType)

	system.GET("/dict-data/index", systemApi.DictDataList)
	system.GET("/dict-data/dictAll", systemApi.DictAll)
	system.POST("/dict-data", systemApi.CreateDictData)
	system.PUT("/dict-data/:id", systemApi.UpdateDictData)
	system.DELETE("/dict-data/:id", systemApi.DeleteDictData)
}
