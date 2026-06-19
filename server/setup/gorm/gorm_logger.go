package gormInit

import (
	"context"
	"errors"
	"server/config"
	loggerInit "server/setup/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type zapGormLogger struct {
	level         logger.LogLevel
	slowThreshold time.Duration
}

func newGormLogger(logConfig config.MysqlLogConfig) logger.Interface {
	return &zapGormLogger{
		level:         parseGormLogLevel(logConfig.LogLevel),
		slowThreshold: time.Duration(logConfig.SlowThresholdMillisecond) * time.Millisecond,
	}
}

func (l *zapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &zapGormLogger{
		level:         level,
		slowThreshold: l.slowThreshold,
	}
}

func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Info {
		loggerInit.Logger.Get().Sugar().Infof(msg, data...)
	}
}

func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Warn {
		loggerInit.Logger.Get().Sugar().Warnf(msg, data...)
	}
}

func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= logger.Error {
		loggerInit.Logger.Get().Sugar().Errorf(msg, data...)
	}
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	switch {
	case err != nil && l.level >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		loggerInit.Logger.Get().Error("gorm sql error", append(fields, zap.Error(err))...)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold && l.level >= logger.Warn:
		loggerInit.Logger.Get().Warn("gorm slow sql", append(fields, zap.Duration("slow_threshold", l.slowThreshold))...)
	case l.level >= logger.Info:
		loggerInit.Logger.Get().Info("gorm sql", fields...)
	}
}

func parseGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	default:
		return logger.Warn
	}
}
