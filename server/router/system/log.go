package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func logRoutes(system *gin.RouterGroup) {
	system.GET("/login-log/index", systemApi.LoginLogList)
	system.GET("/oper-log/index", systemApi.OperLogList)
}
