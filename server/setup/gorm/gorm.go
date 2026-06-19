package gormInit

import (
	"fmt"
	"server/config"

	"gorm.io/gorm"
)

var Gorm = new(_gorm)

type _gorm struct {
	// Databases 保存启动时初始化完成的 MySQL 连接，业务层后续从这里取连接。
	Databases *Databases
}

// Databases 按业务库名称保存 *gorm.DB。
// 当前只有 ai_system；保留结构体是为了后续新增数据库时扩展入口稳定。
type Databases struct {
	AISystem *gorm.DB
}

// InitializeAll 初始化项目启动所需的全部 MySQL 连接。
// 以后新增数据库连接时，只需要在这里扩展，main/initSetup 不需要跟着变复杂。
func (g *_gorm) InitializeAll() error {
	dbAISystem, err := g.Initialize()
	if err != nil {
		return fmt.Errorf("%s connect failed: %w", config.MysqlAISystem, err)
	}
	if err := ensureAISystemSchema(dbAISystem); err != nil {
		return fmt.Errorf("%s schema check failed: %w", config.MysqlAISystem, err)
	}

	g.Databases = &Databases{
		AISystem: dbAISystem,
	}

	return nil
}

// Initialize 打开默认 MySQL 连接。
// 当前默认库是 ai_system。
func (g *_gorm) Initialize() (*gorm.DB, error) {
	return g.InitializeByName(config.MysqlAISystem)
}

// InitializeByName 按连接名打开指定 MySQL。
// 支持的连接名定义在 config 中，例如 ai_system。
func (g *_gorm) InitializeByName(name string) (*gorm.DB, error) {
	return g.initializeMysqlByName(name)
}
