package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

// profileRoutes 注册个人中心路由。
// 操作对象永远是当前登录用户自己，所以只要求登录，不挂权限码。
func profileRoutes(system *gin.RouterGroup) {
	system.PUT("/user/profile", systemApi.UpdateProfile)
	system.PUT("/user/password", systemApi.ChangePassword)
}
