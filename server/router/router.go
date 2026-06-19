package router

import (
	"net/http"
	"server/config"
	"server/middleware"
	"server/model/common/code"
	"server/model/common/response"
	generatedRouter "server/router/generated"
	installRouter "server/router/install"
	systemRouter "server/router/system"
	installService "server/service/install"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewRouter 创建 Gin 路由实例，并挂载所有业务模块路由。
// router 层只负责“URL -> api handler”的映射，不写业务逻辑。
func NewRouter() *gin.Engine {
	Router := gin.New()
	Router.Use(middleware.Recovery())
	Router.Use(middleware.RequestLogger())
	Router.Use(middleware.CORS())
	Router.Use(installGuard())

	// registerGroups 按指定前缀注册一套公开路由和私有路由。
	// 私有路由统一挂 JWT 中间件，业务模块只需要关心自己的 URL 分组。
	registerGroups := func(prefix string) {
		PublicGroup := Router.Group(prefix)
		PrivateGroup := Router.Group(prefix)

		installRouter.Router(PublicGroup)
		PrivateGroup.Use(middleware.JWTAuth())

		systemRouter.BaseRouter(PublicGroup, PrivateGroup) // 系统基础路由
		generatedRouter.RegisterRoutes(PrivateGroup)
	}

	// 根分组使用环境变量控制，默认是 /api/v1。
	routerPrefix := config.RouterPrefix()
	registerGroups(routerPrefix)

	// 兼容前端开发态常用的 /api 前缀。
	// Vite 代理会把 /api 重写成 /api/v1；如果直接访问 Go 后端，/api 也能命中。
	if strings.TrimRight(routerPrefix, "/") != "/api" {
		registerGroups("/api")
	}

	return Router
}

func installGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if installService.Installed() || strings.Contains(c.Request.URL.Path, "/install") {
			c.Next()
			return
		}
		response.FailWithAbort(c, http.StatusServiceUnavailable, code.SystemError, "系统尚未安装，请先访问 /install 完成初始化")
	}
}
