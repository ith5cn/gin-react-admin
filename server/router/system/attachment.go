package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

// attachmentRoutes 注册附件管理路由。
func attachmentRoutes(system *gin.RouterGroup) {
	system.GET("/attachment/index", systemApi.AttachmentList)
	system.POST("/attachment/delete", systemApi.DeleteAttachments)
}
