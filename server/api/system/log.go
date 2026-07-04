package system

import (
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// LoginLogList 登录日志分页列表。
func LoginLogList(c *gin.Context) {
	result, err := systemService.LoginLogList(queryMap(c))
	successOrFail(c, result, err)
}

// OperLogList 操作日志分页列表。
func OperLogList(c *gin.Context) {
	result, err := systemService.OperLogList(queryMap(c))
	successOrFail(c, result, err)
}
