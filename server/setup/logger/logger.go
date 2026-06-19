package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger = new(_logger)

// _logger 是项目统一日志入口。
// 业务代码统一通过 loggerInit.Logger.Get() 获取 zap logger。
type _logger struct {
	Zap *zap.Logger
}

// Initialize 根据环境变量初始化 zap。
// LOG_MODE=prod 使用 JSON 格式生产日志；其它值使用更适合本地阅读的开发日志。
func (l *_logger) Initialize() error {
	config := zap.NewDevelopmentConfig()
	if strings.EqualFold(os.Getenv("LOG_MODE"), "prod") {
		config = zap.NewProductionConfig()
	}

	level, err := parseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return err
	}
	config.Level = zap.NewAtomicLevelAt(level)

	zapLogger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return err
	}

	l.Zap = zapLogger
	return nil
}

// Get 返回初始化后的 zap logger。
// 如果还没初始化，则返回 no-op logger，避免空指针导致程序崩溃。
func (l *_logger) Get() *zap.Logger {
	if l.Zap == nil {
		return zap.NewNop()
	}
	return l.Zap
}

// parseLevel 解析 LOG_LEVEL。
// 未配置时默认 debug；配置错误时把错误返回给启动流程，让问题尽早暴露。
func parseLevel(level string) (zapcore.Level, error) {
	if level == "" {
		return zapcore.DebugLevel, nil
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(strings.ToLower(level))); err != nil {
		return zapcore.DebugLevel, err
	}

	return zapLevel, nil
}
