package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// menuRoutes 注册菜单管理路由；accessMenu 返回当前用户可见菜单，供前端动态路由使用。
func menuRoutes(system *gin.RouterGroup) {
	system.GET("/menu/index", middleware.Perm("system/menu/index"), systemApi.MenuList)
	system.GET("/menu/accessMenu", systemApi.AccessMenu)
	system.GET("/menu/getMenuByRole/:roleId", middleware.Perm("system/role/menu-permission"), systemApi.MenuByRole)
	system.POST("/menu/create", middleware.Perm("system/menu/create"), systemApi.CreateMenu)
	system.PUT("/menu/:id", middleware.Perm("system/menu/update"), systemApi.UpdateMenu)
	system.DELETE("/menu/:id", middleware.Perm("system/menu/destroy"), systemApi.DeleteMenu)
}
