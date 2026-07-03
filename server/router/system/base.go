package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

// BaseRouter 注册系统模块全部路由。
// 公开路由（登录、刷新 token）挂 PublicGroup，其余挂 PrivateGroup；
// 各业务模块的路由注册函数按模块拆分在同目录的独立文件中。
func BaseRouter(PublicGroup, PrivateGroup *gin.RouterGroup) {
	base := PublicGroup.Group("/base")
	base.POST("/login", systemApi.Login)
	base.POST("/token/refresh", systemApi.RefreshToken)

	privateBase := PrivateGroup.Group("/base")
	privateBase.POST("/logout", systemApi.Logout)

	system := PrivateGroup.Group("/system")
	system.GET("/user", systemApi.CurrentUser)

	userRoutes(system)
	menuRoutes(system)
	roleRoutes(system)
	deptRoutes(system)
	postRoutes(system)
	dictRoutes(system)
	configRoutes(system)
	logRoutes(system)
	attachmentRoutes(system)
	codegenRoutes(system)
	databaseRoutes(PrivateGroup)
}
