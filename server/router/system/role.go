package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// roleRoutes 注册角色管理路由。
func roleRoutes(system *gin.RouterGroup) {
	system.GET("/role/index", middleware.Perm("system/role/index"), systemApi.RoleList)
	system.GET("/role/access", systemApi.RoleAccess)
	system.POST("/role/create", middleware.Perm("system/role/create"), systemApi.CreateRole)
	system.POST("/role/:id/menu", middleware.Perm("system/role/menu-permission"), systemApi.BindRoleMenu)
	system.GET("/role/:id/dept", middleware.Perm("system/role/read"), systemApi.RoleDeptIDs)
	system.PUT("/role/:id", middleware.Perm("system/role/update"), systemApi.UpdateRole)
	system.DELETE("/role/:id", middleware.Perm("system/role/destroy"), systemApi.DeleteRole)
}
