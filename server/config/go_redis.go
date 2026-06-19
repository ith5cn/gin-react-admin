package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	// RedisModeSingle 表示连接单实例 Redis，适合本地开发和普通缓存场景。
	RedisModeSingle = "single"
	// RedisModeCluster 表示连接 Redis Cluster，适合分片和高可用场景。
	RedisModeCluster = "cluster"
)

// Redis 描述 Redis 的运行模式和连接信息。
// 单例模式使用 Addr + DB；集群模式使用 Addrs，且不支持选择 DB。
type Redis struct {
	Mode     string
	Addr     string
	Addrs    []string
	Password string
	DB       int
}

// RedisConfig 从环境变量读取 Redis 配置。
// REDIS_MODE 默认为 single；除 cluster 之外的值都会按 single 处理。
func RedisConfig() Redis {
	mode := strings.ToLower(envOrDefault("REDIS_MODE", RedisModeSingle))
	if mode != RedisModeCluster {
		mode = RedisModeSingle
	}

	return Redis{
		Mode:     mode,
		Addr:     envOrDefault("REDIS_ADDR", "127.0.0.1:6379"),
		Addrs:    redisAddrs(),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB(),
	}
}

// redisAddrs 解析集群模式下的 REDIS_ADDRS。
// 示例：127.0.0.1:7001,127.0.0.1:7002,127.0.0.1:7003
func redisAddrs() []string {
	value := os.Getenv("REDIS_ADDRS")
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	addrs := make([]string, 0, len(parts))
	for _, part := range parts {
		addr := strings.TrimSpace(part)
		if addr != "" {
			addrs = append(addrs, addr)
		}
	}

	return addrs
}

// redisDB 只在单例模式下使用；Redis Cluster 不支持 SELECT DB。
func redisDB() int {
	value := os.Getenv("REDIS_DB")
	if value == "" {
		return 0
	}

	db, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return db
}
