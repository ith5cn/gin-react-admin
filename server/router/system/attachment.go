package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func attachmentRoutes(system *gin.RouterGroup) {
	system.GET("/attachment/index", systemApi.AttachmentList)
	system.POST("/attachment/delete", systemApi.DeleteAttachments)
}
