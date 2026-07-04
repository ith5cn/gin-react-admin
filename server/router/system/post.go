package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// postRoutes 注册岗位管理路由。
func postRoutes(system *gin.RouterGroup) {
	system.GET("/post/index", middleware.Perm("system/post/index"), systemApi.PostList)
	system.GET("/post/access", systemApi.PostAccess)
	system.POST("/post", middleware.Perm("system/post/create"), systemApi.CreatePost)
	system.PUT("/post/:id", middleware.Perm("system/post/update"), systemApi.UpdatePost)
	system.DELETE("/post/:id", middleware.Perm("system/post/destroy"), systemApi.DeletePost)
}
