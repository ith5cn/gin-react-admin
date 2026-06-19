package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// MysqlAISystem 是当前后台系统库连接名。
	// 现阶段项目只连接 ai_system，一个库里承载登录用户等后台管理数据。
	MysqlAISystem = "ai_system"

	// defaultMysqlConfig 是 MySQL DSN 的默认附加参数。
	// parseTime=True 让 datetime/time 能正确解析成 Go 的 time.Time。
	defaultMysqlConfig = "charset=utf8mb4&parseTime=True&loc=Local"
)

// Mysql 描述一个 MySQL 连接所需的最小配置。
type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	Config   string
}

type MysqlPoolConfig struct {
	MaxIdleConns          int
	MaxOpenConns          int
	ConnMaxLifetimeMinute int
}

type MysqlLogConfig struct {
	SlowThresholdMillisecond int
	LogLevel                 string
}

// Dsn 返回默认 MySQL 连接的 DSN。
// 当前默认库是 ai_system；如果业务要指定库，使用 DsnByName。
func Dsn() string {
	return DsnByName(MysqlAISystem)
}

// MysqlPool 读取 MySQL 连接池配置，所有 MySQL 连接共用这套参数。
func MysqlPool() MysqlPoolConfig {
	return MysqlPoolConfig{
		MaxIdleConns:          mysqlIntOrDefault("MYSQL_MAX_IDLE_CONNS", 10),
		MaxOpenConns:          mysqlIntOrDefault("MYSQL_MAX_OPEN_CONNS", 100),
		ConnMaxLifetimeMinute: mysqlIntOrDefault("MYSQL_CONN_MAX_LIFETIME_MINUTE", 60),
	}
}

// MysqlLog 读取 GORM SQL 日志配置。
func MysqlLog() MysqlLogConfig {
	level := strings.ToLower(envOrDefault("MYSQL_LOG_LEVEL", "warn"))
	switch level {
	case "silent", "error", "warn", "info":
	default:
		level = "warn"
	}

	return MysqlLogConfig{
		SlowThresholdMillisecond: mysqlIntOrDefault("MYSQL_SLOW_THRESHOLD_MILLISECOND", 200),
		LogLevel:                 level,
	}
}

// DsnByName 根据连接名生成 DSN。
func DsnByName(name string) string {
	mysql := MysqlByName(name)
	return mysql.Dsn()
}

// MysqlByName 根据连接名读取对应环境变量。
// 未知连接名会回退到 ai_system，避免因为传错名称直接拿到空配置。
func MysqlByName(name string) Mysql {
	switch name {
	case MysqlAISystem:
		return mysqlFromEnv("AI_SYSTEM_MYSQL", MysqlAISystem)
	default:
		mysql := mysqlFromEnv("AI_SYSTEM_MYSQL", MysqlAISystem)
		if strings.TrimSpace(name) != "" {
			mysql.Dbname = strings.TrimSpace(name)
		}
		return mysql
	}
}

// Dsn 拼接 GORM MySQL driver 需要的 DSN 字符串。
func (m Mysql) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		m.User,
		m.Password,
		m.Host,
		m.Port,
		m.Dbname,
		m.Config,
	)
}

// mysqlFromEnv 按统一前缀读取一组 MySQL 配置。
// 例如 prefix=AI_SYSTEM_MYSQL 时，会读取 AI_SYSTEM_MYSQL_HOST 等变量。
func mysqlFromEnv(prefix string, defaultDbname string) Mysql {
	return Mysql{
		Host:     envOrDefault(prefix+"_HOST", "127.0.0.1"),
		Port:     envOrDefault(prefix+"_PORT", "3306"),
		User:     envOrDefault(prefix+"_USER", "root"),
		Password: envOrDefault(prefix+"_PASSWORD", ""),
		Dbname:   envOrDefault(prefix+"_DB", defaultDbname),
		Config:   envOrDefault(prefix+"_CONFIG", defaultMysqlConfig),
	}
}

// envOrDefault 是配置层通用的小工具：环境变量为空时使用默认值。
func envOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func mysqlIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil || number < 0 {
		return defaultValue
	}

	return number
}
