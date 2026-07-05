---
description: 数据库 schema、GORM 约定、migration 规范
paths:
  - server/database/**
  - server/model/**
---

# 数据库规范

## Schema 管理

- 唯一 SQL 源文件：`server/database/ai_system.sql`（全量 DDL）
- 目前无 migration 工具；schema 变更需同步更新 `ai_system.sql`
- 建议后续引入 `golang-migrate` 或 `goose`，在此之前每次 schema 变更附上变更 SQL 注释

## GORM Model 约定

- model 文件放 `server/model/system/`，文件名对应表名（snake_case）
- 必须嵌入 `gorm.Model`（含 `id`, `created_at`, `updated_at`, `deleted_at` 软删除）或显式定义这些字段
- 字段 json tag 用 snake_case；gorm tag 显式指定 `column` 名称
- 禁止在 model 层写业务方法；只写 GORM hook（BeforeCreate 等）

## 查询规范

- 所有查询通过 GORM 方法（`db.Where().First()` 等），禁止字符串拼接 SQL
- 批量操作用 `db.CreateInBatches` 而非循环单条 insert
- 软删除：GORM 的 `deleted_at` 已自动处理，禁止手动 `WHERE deleted_at IS NULL`
- 分页：使用 `db.Offset().Limit()`，接口统一返回 `total` 和 `data` 列表

## 连接池

- 连接池参数在 `config/go_mysql.go` 的 `MysqlPoolConfig` 中配置
- 生产环境建议 `MaxOpenConns=50`，`MaxIdleConns=10`，`ConnMaxLifetime=1h`
