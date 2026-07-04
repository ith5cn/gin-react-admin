package config

import (
	"strconv"
	"strings"
)

// CorsConfig 是跨域中间件所需的全部配置。
type CorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAgeHour       int
}

// ServerAddr 读取 HTTP 服务监听地址。
// 示例：SERVER_ADDR=:8080 表示监听本机 8080 端口。
func ServerAddr() string {
	return envOrDefault("SERVER_ADDR", ":8080")
}

// RouterPrefix 读取业务路由总前缀。
// 示例：ROUTER_PREFIX=/api/v1 时，登录接口会是 /api/v1/base/login。
func RouterPrefix() string {
	return envOrDefault("ROUTER_PREFIX", "/api/v1")
}

// Cors 读取跨域配置。
// 多值配置用英文逗号分隔，例如 CORS_ALLOW_METHODS=GET,POST,OPTIONS。
func Cors() CorsConfig {
	return CorsConfig{
		AllowOrigins: splitEnvOrDefault("CORS_ALLOW_ORIGINS", []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:5174",
			"http://127.0.0.1:5174",
			"http://localhost:5175",
			"http://127.0.0.1:5175",
		}),
		AllowMethods: splitEnvOrDefault("CORS_ALLOW_METHODS", []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		}),
		AllowHeaders: splitEnvOrDefault("CORS_ALLOW_HEADERS", []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		}),
		ExposeHeaders:    splitEnvOrDefault("CORS_EXPOSE_HEADERS", []string{"Content-Length"}),
		AllowCredentials: boolEnvOrDefault("CORS_ALLOW_CREDENTIALS", true),
		MaxAgeHour:       intEnvOrDefault("CORS_MAX_AGE_HOUR", 12),
	}
}

func splitEnvOrDefault(key string, defaultValue []string) []string {
	value := envOrDefault(key, "")
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			values = append(values, item)
		}
	}
	if len(values) == 0 {
		return defaultValue
	}
	return values
}

func boolEnvOrDefault(key string, defaultValue bool) bool {
	value := strings.ToLower(envOrDefault(key, ""))
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}

func intEnvOrDefault(key string, defaultValue int) int {
	value := envOrDefault(key, "")
	if value == "" {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil || number < 0 {
		return defaultValue
	}
	return number
}
