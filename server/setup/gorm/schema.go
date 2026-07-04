package gormInit

import "gorm.io/gorm"

// ensureAISystemSchema 在启动时做轻量 schema 自检/补齐：
// 给早期版本缺失的列补 ALTER，并确保代码生成器的两张配置表存在。
// 项目暂未引入 migration 工具，这里相当于最简版的向前兼容迁移。
func ensureAISystemSchema(db *gorm.DB) error {
	if !db.Migrator().HasColumn("ai_system_config_group", "sort") {
		if err := db.Exec("ALTER TABLE `ai_system_config_group` ADD COLUMN `sort` smallint unsigned NOT NULL DEFAULT 0 COMMENT '排序' AFTER `code`").Error; err != nil {
			return err
		}
	}

	if err := db.Exec(`
CREATE TABLE IF NOT EXISTS nest_tool_generate_tables (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  created_by int NULL DEFAULT NULL COMMENT '创建者',
  updated_by int NULL DEFAULT NULL COMMENT '更新者',
  create_time datetime(6) NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  update_time datetime(6) NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  table_name varchar(200) NULL DEFAULT NULL COMMENT '表名称',
  table_comment varchar(500) NULL DEFAULT NULL COMMENT '表注释',
  package_name varchar(100) NULL DEFAULT NULL COMMENT '包名',
  business_name varchar(50) NULL DEFAULT NULL COMMENT '业务名称',
  class_name varchar(50) NULL DEFAULT NULL COMMENT '类名称',
  menu_name varchar(100) NULL DEFAULT NULL COMMENT '生成菜单名称',
  belong_menu_id int NULL DEFAULT NULL COMMENT '所属菜单ID',
  tpl_category varchar(100) NULL DEFAULT NULL COMMENT '生成模板类型',
  generate_path varchar(100) NOT NULL DEFAULT 'web' COMMENT '生成路径',
  generate_model smallint NOT NULL DEFAULT 1 COMMENT '生成模式, 1 软删除 2 非软删除',
  form_width int NOT NULL DEFAULT 600 COMMENT '表单宽度',
  is_full tinyint NOT NULL DEFAULT 1 COMMENT '是否全屏',
  remark varchar(255) NULL DEFAULT NULL COMMENT '备注',
  source varchar(255) NULL DEFAULT NULL COMMENT '数据源',
  component_type smallint NOT NULL DEFAULT 1 COMMENT '组件显示方式',
  sort tinyint NOT NULL DEFAULT 0 COMMENT '排序',
  delete_time datetime(6) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  KEY idx_source_table (source, table_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='代码生成表配置';
`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
CREATE TABLE IF NOT EXISTS nest_tool_generate_columns (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  created_by int NULL DEFAULT NULL COMMENT '创建者',
  updated_by int NULL DEFAULT NULL COMMENT '更新者',
  create_time datetime(6) NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  update_time datetime(6) NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  table_id int unsigned NULL DEFAULT NULL COMMENT '所属表ID',
  column_name varchar(200) NULL DEFAULT NULL COMMENT '字段名称',
  column_comment varchar(255) NULL DEFAULT NULL COMMENT '字段注释',
  column_type varchar(50) NULL DEFAULT NULL COMMENT '字段类型',
  default_value varchar(50) NULL DEFAULT NULL COMMENT '默认值',
  is_pk smallint NOT NULL DEFAULT 1 COMMENT '1 非主键 2 主键',
  is_required smallint NOT NULL DEFAULT 1 COMMENT '1 非必填 2 必填',
  is_insert smallint NOT NULL DEFAULT 1 COMMENT '1 非插入字段 2 插入字段',
  is_edit smallint NOT NULL DEFAULT 1 COMMENT '1 非编辑字段 2 编辑字段',
  is_list smallint NOT NULL DEFAULT 1 COMMENT '1 非列表显示字段 2 列表显示字段',
  is_query smallint NOT NULL DEFAULT 1 COMMENT '1 非查询字段 2 查询字段',
  is_sort smallint NOT NULL DEFAULT 1 COMMENT '1 非排序 2 排序',
  query_type varchar(100) NOT NULL DEFAULT 'eq' COMMENT '查询方式',
  view_type varchar(100) NOT NULL DEFAULT 'input' COMMENT '页面控件',
  dict_type varchar(200) NULL DEFAULT NULL COMMENT '字典类型',
  allow_roles varchar(255) NULL DEFAULT NULL COMMENT '允许查看该字段的角色',
  sort tinyint unsigned NOT NULL DEFAULT 0 COMMENT '排序',
  remark varchar(255) NULL DEFAULT NULL COMMENT '备注',
  delete_time datetime(6) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (id),
  KEY idx_table_id (table_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='代码生成字段配置';
`).Error; err != nil {
		return err
	}

	if err := ensureCodegenSoftDeleteColumns(db); err != nil {
		return err
	}

	return nil
}

// ensureCodegenSoftDeleteColumns 为 codegen 配置表补齐软删除列。
func ensureCodegenSoftDeleteColumns(db *gorm.DB) error {
	tables := []string{"nest_tool_generate_tables", "nest_tool_generate_columns"}
	for _, table := range tables {
		if !db.Migrator().HasColumn(table, "delete_time") {
			if err := db.Exec("ALTER TABLE `" + table + "` ADD COLUMN `delete_time` datetime(6) NULL DEFAULT NULL COMMENT '删除时间'").Error; err != nil {
				return err
			}
		}
	}
	return nil
}
