package system

import (
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// OnlineUserList 在线用户分页列表（基于 Redis 会话记录）。
func OnlineUserList(c *gin.Context) {
	data, err := systemService.OnlineUserList(queryMap(c))
	successOrFail(c, data, err)
}

// KickOnlineUser 把指定会话踢下线，路径参数是该会话 access token 的 jti。
func KickOnlineUser(c *gin.Context) {
	successOrFail(c, true, systemService.KickOnlineUser(c.Param("jti")))
}
