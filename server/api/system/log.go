package system

import (
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func LoginLogList(c *gin.Context) {
	result, err := systemService.LoginLogList(queryMap(c))
	successOrFail(c, result, err)
}

func OperLogList(c *gin.Context) {
	result, err := systemService.OperLogList(queryMap(c))
	successOrFail(c, result, err)
}
