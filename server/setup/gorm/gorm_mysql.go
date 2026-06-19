package gormInit

import (
	"server/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// initializeMysqlByName 是 MySQL 的具体连接实现。
// 上层只传连接名，DSN 拼接和配置读取都交给 config 包。
func (g *_gorm) initializeMysqlByName(name string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.DsnByName(name)), &gorm.Config{
		Logger: newGormLogger(config.MysqlLog()),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	poolConfig := config.MysqlPool()
	sqlDB.SetMaxIdleConns(poolConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(poolConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(poolConfig.ConnMaxLifetimeMinute) * time.Minute)

	return db, nil
}
