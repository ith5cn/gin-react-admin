package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// attachmentRoutes 注册附件管理路由。
func attachmentRoutes(system *gin.RouterGroup) {
	system.GET("/attachment/index", middleware.Perm("system/attachment/index"), systemApi.AttachmentList)
	system.POST("/attachment/delete", middleware.Perm("system/attachment/destroy"), systemApi.DeleteAttachments)
}
