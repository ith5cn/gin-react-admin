package install

import (
	installAPI "server/api/install"

	"github.com/gin-gonic/gin"
)

// Router 注册安装向导路由。挂在公开组上：系统未安装时也必须能访问。
func Router(group *gin.RouterGroup) {
	install := group.Group("/install")
	install.GET("/status", installAPI.Status)
	install.POST("/check", installAPI.Check)
	install.POST("/run", installAPI.Run)
}
