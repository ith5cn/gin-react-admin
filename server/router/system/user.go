package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func userRoutes(system *gin.RouterGroup) {
	system.GET("/user/index", systemApi.UserList)
	system.GET("/user/auth-list", systemApi.UserAuthList)
	system.POST("/user", systemApi.CreateUser)
	system.PUT("/user/:id/refresh-cache", systemApi.RefreshUserCache)
	system.PUT("/user/:id/set-password", systemApi.SetUserPassword)
	system.POST("/user/:id/role", systemApi.BindUserRole)
	system.PUT("/user/:id", systemApi.UpdateUser)
	system.DELETE("/user/:id", systemApi.DeleteUser)
}
