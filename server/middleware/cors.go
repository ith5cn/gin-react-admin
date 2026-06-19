package middleware

import (
	"server/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 返回 Gin 跨域中间件。
// 具体允许的来源、方法、请求头等全部从 .env 读取，方便不同环境独立配置。
func CORS() gin.HandlerFunc {
	corsConfig := config.Cors()

	return cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowOrigins,
		AllowMethods:     corsConfig.AllowMethods,
		AllowHeaders:     corsConfig.AllowHeaders,
		ExposeHeaders:    corsConfig.ExposeHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           time.Duration(corsConfig.MaxAgeHour) * time.Hour,
	})
}
