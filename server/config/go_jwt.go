package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	// JwtLoginModeSingle 表示单一登录：同一用户新 token 会让旧 token 失效。
	JwtLoginModeSingle = "single"
	// JwtLoginModeMulti 表示多端登录：同一用户可以同时持有多个有效 token。
	JwtLoginModeMulti = "multi"
)

// Jwt 描述 JWT 签发和校验所需的配置。
type Jwt struct {
	// Secret 是 JWT 签名密钥，生产环境必须替换成足够长的随机字符串。
	Secret string `json:"secret" mapstructure:"secret"`
	// Issuer 是 token 签发者，用于区分 token 来源。
	Issuer string `json:"issuer" mapstructure:"issuer"`
	// AccessExpiresMinute 是 access token 的有效期，单位分钟。
	AccessExpiresMinute int `json:"access_expires_minute" mapstructure:"access_expires_minute"`
	// RefreshExpiresHour 是 refresh token 的有效期，单位小时。
	RefreshExpiresHour int `json:"refresh_expires_hour" mapstructure:"refresh_expires_hour"`
	// LoginMode 控制单一登录还是多端登录。
	LoginMode string `json:"login_mode" mapstructure:"login_mode"`
}

// JwtConfig 从环境变量读取 JWT 配置。
// JWT_LOGIN_MODE 默认为 multi；除 single 之外的值都会按 multi 处理。
func JwtConfig() Jwt {
	loginMode := strings.ToLower(envOrDefault("JWT_LOGIN_MODE", JwtLoginModeMulti))
	if loginMode != JwtLoginModeSingle {
		loginMode = JwtLoginModeMulti
	}

	return Jwt{
		Secret:              os.Getenv("JWT_SECRET"),
		Issuer:              envOrDefault("JWT_ISSUER", "gin-react-admin"),
		AccessExpiresMinute: jwtIntOrDefault("JWT_ACCESS_EXPIRES_MINUTE", 120),
		RefreshExpiresHour:  jwtIntOrDefault("JWT_REFRESH_EXPIRES_HOUR", 168),
		LoginMode:           loginMode,
	}
}

// jwtIntOrDefault 读取正整数配置，非法或空值时回退默认值。
func jwtIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil || number <= 0 {
		return defaultValue
	}

	return number
}
