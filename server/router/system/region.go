package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

// regionRoutes 注册省市区数据路由。
// 纯下拉数据，登录即可访问，不挂 Perm（同 dict-data/dictAll 等下拉接口的约定）。
func regionRoutes(system *gin.RouterGroup) {
	system.GET("/region/options", systemApi.RegionOptions)
}
