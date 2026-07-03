package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func postRoutes(system *gin.RouterGroup) {
	system.GET("/post/index", systemApi.PostList)
	system.GET("/post/access", systemApi.PostAccess)
	system.POST("/post", systemApi.CreatePost)
	system.PUT("/post/:id", systemApi.UpdatePost)
	system.DELETE("/post/:id", systemApi.DeletePost)
}
