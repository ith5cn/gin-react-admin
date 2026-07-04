package main

import (
	"server/config"
	"server/router"
	installService "server/service/install"
	systemService "server/service/system"
	gormInit "server/setup/gorm"
	loggerInit "server/setup/logger"
	redisInit "server/setup/redis"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// 先加载 .env，后面的日志、数据库、Redis、JWT 配置都会从环境变量读取。
	loadEnv()

	// 日志要尽早初始化，这样启动阶段的错误也能用统一格式输出。
	if err := loggerInit.Logger.Initialize(); err != nil {
		panic(err)
	}
	defer loggerInit.Logger.Get().Sync()

	// 初始化基础设施：MySQL 多连接、Redis 等。
	if err := initSetup(); err != nil {
		loggerInit.Logger.Get().Fatal("setup initialize failed", zap.Error(err))
	}

	// 所有基础设施准备完成后再启动 HTTP 服务。
	if err := router.NewRouter().Run(config.ServerAddr()); err != nil {
		loggerInit.Logger.Get().Fatal("server run failed", zap.Error(err))
	}
}

func loadEnv() {
	// godotenv.Load 默认读取当前工作目录下的 .env。
	// 这里忽略错误，是为了兼容生产环境直接使用系统环境变量的情况。
	_ = godotenv.Load()
}

// initSetup 是项目统一的初始化入口。
// 后续如果接入 Gin 中间件、定时任务、消息队列等，也可以在这里继续扩展。
func initSetup() error {
	if !installService.Installed() {
		loggerInit.Logger.Get().Info("install lock not found, skip business setup and enable install wizard")
		return nil
	}

	if err := gormInit.Gorm.InitializeAll(); err != nil {
		return err
	}

	if err := redisInit.Redis.Initialize(); err != nil {
		return err
	}

	// 定时任务调度器依赖数据库（读任务配置），必须在 gorm 初始化之后启动。
	return systemService.StartCrontabScheduler()
}
