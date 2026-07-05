package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// userRoutes 注册用户管理路由（列表/增删改/重置密码/绑角色）。
func userRoutes(system *gin.RouterGroup) {
	system.GET("/user/index", middleware.Perm("system/user/index"), systemApi.UserList)
	system.GET("/user/auth-list", systemApi.UserAuthList)
	system.GET("/user/export", middleware.Perm("system/user/export"), systemApi.ExportUsers)
	system.GET("/user/import-template", middleware.Perm("system/user/import"), systemApi.UserImportTemplate)
	system.POST("/user/import", middleware.Perm("system/user/import"), systemApi.ImportUsers)
	system.POST("/user", middleware.Perm("system/user/create"), systemApi.CreateUser)
	system.PUT("/user/:id/refresh-cache", middleware.Perm("system/user/refresh-cache"), systemApi.RefreshUserCache)
	system.PUT("/user/:id/set-password", middleware.Perm("system/user/set-password"), systemApi.SetUserPassword)
	system.POST("/user/:id/role", middleware.Perm("system/user/update"), systemApi.BindUserRole)
	system.PUT("/user/:id", middleware.Perm("system/user/update"), systemApi.UpdateUser)
	system.DELETE("/user/:id", middleware.Perm("system/user/destroy"), systemApi.DeleteUser)
}
