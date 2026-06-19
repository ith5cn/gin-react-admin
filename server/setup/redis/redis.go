package redisInit

import (
	"context"
	"fmt"
	"server/config"

	"github.com/redis/go-redis/v9"
)

var Redis = new(_redis)

type _redis struct {
	// UniversalClient 是 go-redis 提供的统一接口。
	// 不管底层是单例 Redis 还是 Redis Cluster，业务层都可以用同一套 API。
	Client redis.UniversalClient
}

// Initialize 根据 REDIS_MODE 创建 Redis 客户端，并用 Ping 验证连接是否可用。
// JWT 的 token 状态会写入 Redis，所以启动阶段 Redis 不通应直接报错。
func (r *_redis) Initialize() error {
	redisConfig := config.RedisConfig()

	client := newClient(redisConfig)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("redis %s connect failed: %w", redisConfig.Mode, err)
	}

	r.Client = client
	return nil
}

// Get 返回初始化后的 Redis 客户端。
// 调用方需要保证 Initialize 已经在启动阶段执行过，否则 Client 会是 nil。
func (r *_redis) Get() redis.UniversalClient {
	return r.Client
}

// newClient 屏蔽单例和集群 Redis 的创建差异。
// 注意：集群模式不配置 DB，因为 Redis Cluster 不支持 SELECT DB。
func newClient(redisConfig config.Redis) redis.UniversalClient {
	if redisConfig.Mode == config.RedisModeCluster {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisConfig.Addrs,
			Password: redisConfig.Password,
		})
	}

	return redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
}
