package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func roleRoutes(system *gin.RouterGroup) {
	system.GET("/role/index", systemApi.RoleList)
	system.GET("/role/access", systemApi.RoleAccess)
	system.POST("/role/create", systemApi.CreateRole)
	system.POST("/role/:id/menu", systemApi.BindRoleMenu)
	system.PUT("/role/:id", systemApi.UpdateRole)
	system.DELETE("/role/:id", systemApi.DeleteRole)
}
