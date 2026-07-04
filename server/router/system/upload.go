package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

// uploadRoutes 注册通用上传路由。
// 只要求登录（头像等场景所有用户都要用），不额外挂权限码。
func uploadRoutes(system *gin.RouterGroup) {
	system.POST("/uploadImage", systemApi.UploadImage)
	system.POST("/uploadFile", systemApi.UploadFile)
}
