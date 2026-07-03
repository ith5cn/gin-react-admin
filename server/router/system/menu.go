package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func menuRoutes(system *gin.RouterGroup) {
	system.GET("/menu/index", systemApi.MenuList)
	system.GET("/menu/accessMenu", systemApi.AccessMenu)
	system.GET("/menu/getMenuByRole/:roleId", systemApi.MenuByRole)
	system.POST("/menu/create", systemApi.CreateMenu)
	system.PUT("/menu/:id", systemApi.UpdateMenu)
	system.DELETE("/menu/:id", systemApi.DeleteMenu)
}
