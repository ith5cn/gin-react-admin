package install

import (
	installAPI "server/api/install"

	"github.com/gin-gonic/gin"
)

func Router(group *gin.RouterGroup) {
	install := group.Group("/install")
	install.GET("/status", installAPI.Status)
	install.POST("/check", installAPI.Check)
	install.POST("/run", installAPI.Run)
}
