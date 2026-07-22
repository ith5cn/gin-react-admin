package gormInit

import (
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// ensureAISystemSchema 在启动时做轻量 schema 自检/补齐：
// 给早期版本缺失的列补 ALTER，并确保代码生成器的两张配置表存在。
// 项目暂未引入 migration 工具，这里相当于最简版的向前兼容迁移。
func ensureAISystemSchema(db *gorm.DB) error {
	if !db.Migrator().HasColumn("ai_system_config_group", "sort") {
		if err := db.Exec("ALTER TABLE `ai_system_config_group` ADD COLUMN `sort` smallint unsigned NOT NULL DEFAULT 0 COMMENT '排序' AFTER `code`").Error; err != nil {
			return err
		}
	}

	if err := ensurePermissionCodes(db); err != nil {
		return err
	}

	if err := ensureLoginLogSchema(db); err != nil {
		return err
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
  option_source varchar(20) NULL DEFAULT NULL COMMENT '选项数据来源',
  option_config text NULL COMMENT '选项组件配置',
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
	if err := removeCodegenFullScreenColumns(db); err != nil {
		return err
	}
	if err := ensureCodegenOptionColumns(db); err != nil {
		return err
	}

	if err := ensureRoleDeptTable(db); err != nil {
		return err
	}

	// 定时任务菜单的 component 旧值指向不存在的前端路径，修正到实际页面位置（幂等）。
	if err := db.Exec("UPDATE `ai_system_menu` SET `component` = 'system/tool/crontab/index' WHERE `component` = 'tool/crontab/index'").Error; err != nil {
		return err
	}

	if err := ensureNoticeTable(db); err != nil {
		return err
	}

	if err := ensureUserImportExportMenu(db); err != nil {
		return err
	}

	return ensureOnlineUserMenu(db)
}

// ensureLoginLogSchema 修正旧版登录日志时间字段和菜单路径。
func ensureLoginLogSchema(db *gorm.DB) error {
	if db.Migrator().HasTable("ai_system_login_log") {
		var extra string
		if err := db.Raw("SELECT EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?", "ai_system_login_log", "login_time").Scan(&extra).Error; err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(extra), "on update") {
			if err := db.Exec("ALTER TABLE ai_system_login_log MODIFY COLUMN login_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间'").Error; err != nil {
				return err
			}
		}
		if err := db.Exec("UPDATE ai_system_login_log SET ip_location = '本机' WHERE (ip_location IS NULL OR ip_location = '') AND ip IN ('127.0.0.1', '::1')").Error; err != nil {
			return err
		}
	}
	if err := db.Exec("UPDATE ai_system_menu SET component = 'system/login-log/index' WHERE component = 'system/logs/loginLog'").Error; err != nil {
		return err
	}
	return db.Exec("UPDATE ai_system_menu SET level = '0,3000,3300,3400' WHERE code = 'system/login-log/destroy'").Error
}

// ensureRoleDeptTable 创建角色-部门关联表（自定义数据权限用），幂等。
func ensureRoleDeptTable(db *gorm.DB) error {
	return db.Exec(`
CREATE TABLE IF NOT EXISTS ai_system_role_dept (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  role_id int unsigned NOT NULL COMMENT '角色ID',
  dept_id int unsigned NOT NULL COMMENT '部门ID',
  PRIMARY KEY (id),
  KEY idx_role_id (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='角色部门关联表(自定义数据权限)';
`).Error
}

// ensureNoticeTable 创建通知公告表，幂等。
func ensureNoticeTable(db *gorm.DB) error {
	return db.Exec(`
CREATE TABLE IF NOT EXISTS ai_system_notice (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  title varchar(255) NOT NULL COMMENT '公告标题',
  type smallint NOT NULL DEFAULT 2 COMMENT '类型 (1通知 2公告)',
  content text NULL COMMENT '公告内容',
  status smallint NOT NULL DEFAULT 1 COMMENT '状态 (1正常 2停用)',
  remark varchar(255) NULL DEFAULT NULL COMMENT '备注',
  created_by int NULL DEFAULT NULL COMMENT '创建者',
  updated_by int NULL DEFAULT NULL COMMENT '更新者',
  create_time datetime NULL DEFAULT NULL COMMENT '创建时间',
  update_time datetime NULL DEFAULT NULL COMMENT '修改时间',
  delete_time datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='通知公告表';
`).Error
}

// ensureUserImportExportMenu 为老库补齐用户导入/导出按钮权限码，按 code 判重保证幂等。
// 挂到"用户管理"菜单（route='permission/user'）下；找不到父菜单就跳过（权限码没有归属没意义）。
func ensureUserImportExportMenu(db *gorm.DB) error {
	var parent struct {
		ID    int
		Level *string
	}
	if err := db.Table("ai_system_menu").Select("id, level").
		Where("route = ? AND type = 'M' AND delete_time IS NULL", "permission/user").
		Order("id ASC").Limit(1).Scan(&parent).Error; err != nil {
		return err
	}
	if parent.ID == 0 {
		return nil
	}
	childLevel := "0," + strconv.Itoa(parent.ID)
	if parent.Level != nil && *parent.Level != "" {
		childLevel = *parent.Level + "," + strconv.Itoa(parent.ID)
	}

	buttons := []struct{ name, code string }{
		{"用户导入", "system/user/import"},
		{"用户导出", "system/user/export"},
	}
	for _, button := range buttons {
		var count int64
		if err := db.Table("ai_system_menu").Where("code = ?", button.code).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		if err := db.Exec(`INSERT INTO ai_system_menu
  (parent_id, level, name, code, is_hidden, is_layout, type, status, sort, remark, create_time, update_time)
  VALUES (?, ?, ?, ?, 2, 1, 'B', 1, 0, '', NOW(), NOW())`,
			parent.ID, childLevel, button.name, button.code).Error; err != nil {
			return err
		}
	}
	return nil
}

// ensureOnlineUserMenu 为老库补齐"在线用户"菜单和按钮权限码，按 code 判重保证幂等。
// 不用固定菜单 ID：老库里用户自建菜单可能已占用任何 ID，这里让数据库自增分配。
func ensureOnlineUserMenu(db *gorm.DB) error {
	var count int64
	if err := db.Table("ai_system_menu").Where("code = ?", "system/online/index").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// 挂到"监控"目录（seed 里 route='monitor'）下；找不到就挂根目录。
	parentID := 0
	parentLevel := "0"
	var parent struct {
		ID    int
		Level *string
	}
	if err := db.Table("ai_system_menu").Select("id, level").
		Where("route = ? AND type = 'M' AND delete_time IS NULL", "monitor").
		Order("id ASC").Limit(1).Scan(&parent).Error; err == nil && parent.ID > 0 {
		parentID = parent.ID
		parentLevel = "0," + strconv.Itoa(parent.ID)
		if parent.Level != nil && *parent.Level != "" {
			parentLevel = *parent.Level + "," + strconv.Itoa(parent.ID)
		}
	}

	if err := db.Exec(`INSERT INTO ai_system_menu
  (parent_id, level, name, code, icon, route, component, redirect, is_hidden, is_layout, type, status, sort, remark, create_time, update_time)
  VALUES (?, ?, '在线用户', 'monitor/online', 'TeamOutlined', 'monitor/online', 'system/monitor/online', NULL, 2, 1, 'M', 1, 30, '', NOW(), NOW())`,
		parentID, parentLevel).Error; err != nil {
		return err
	}

	var menuID int
	if err := db.Table("ai_system_menu").Select("id").
		Where("code = ?", "monitor/online").Order("id DESC").Limit(1).Scan(&menuID).Error; err != nil {
		return err
	}
	childLevel := parentLevel + "," + strconv.Itoa(menuID)

	buttons := []struct{ name, code string }{
		{"在线用户列表", "system/online/index"},
		{"在线用户强退", "system/online/kick"},
	}
	for _, button := range buttons {
		if err := db.Exec(`INSERT INTO ai_system_menu
  (parent_id, level, name, code, is_hidden, is_layout, type, status, sort, remark, create_time, update_time)
  VALUES (?, ?, ?, ?, 2, 1, 'B', 1, 0, '', NOW(), NOW())`,
			menuID, childLevel, button.name, button.code).Error; err != nil {
			return err
		}
	}
	return nil
}

// legacyPermCodeMapping 是旧系统按钮权限码到本项目统一规范（system/<模块>/<动作>）的映射。
// 接口级权限中间件（middleware.Perm）按新码校验，老库启动时在这里一次性归一化。
// UPDATE 按 code 精确匹配，重复执行无副作用（幂等）。
var legacyPermCodeMapping = map[string]string{
	"/core/user/index":               "system/user/index",
	"/core/user/save":                "system/user/create",
	"/core/user/update":              "system/user/update",
	"/core/user/destroy":             "system/user/destroy",
	"/core/user/read":                "system/user/read",
	"/core/user/changeStatus":        "system/user/change-status",
	"/core/user/initUserPassword":    "system/user/set-password",
	"/core/user/clearCache":          "system/user/refresh-cache",
	"/core/user/setHomePage":         "system/user/set-home-page",
	"/core/menu/index":               "system/menu/index",
	"/core/menu/save":                "system/menu/create",
	"/core/menu/update":              "system/menu/update",
	"/core/menu/destroy":             "system/menu/destroy",
	"/core/menu/read":                "system/menu/read",
	"/core/dept/index":               "system/dept/index",
	"/core/dept/save":                "system/dept/create",
	"/core/dept/update":              "system/dept/update",
	"/core/dept/destroy":             "system/dept/destroy",
	"/core/dept/read":                "system/dept/read",
	"/core/dept/leaders":             "system/dept/leaders",
	"/core/role/index":               "system/role/index",
	"/core/role/save":                "system/role/create",
	"/core/role/update":              "system/role/update",
	"/core/role/destroy":             "system/role/destroy",
	"/core/role/read":                "system/role/read",
	"/core/role/menuPermission":      "system/role/menu-permission",
	"/core/post/index":               "system/post/index",
	"/core/post/save":                "system/post/create",
	"/core/post/update":              "system/post/update",
	"/core/post/destroy":             "system/post/destroy",
	"/core/post/read":                "system/post/read",
	"/core/post/changeStatus":        "system/post/change-status",
	"/core/post/import":              "system/post/import",
	"/core/post/export":              "system/post/export",
	"/core/dictType/index":           "system/dict-type/index",
	"/core/dictType/save":            "system/dict-type/create",
	"/core/dictType/update":          "system/dict-type/update",
	"/core/dictType/destroy":         "system/dict-type/destroy",
	"/core/dictType/read":            "system/dict-type/read",
	"/core/dictType/changeStatus":    "system/dict-type/change-status",
	"/core/attachment/index":         "system/attachment/index",
	"/core/attachment/destroy":       "system/attachment/destroy",
	"/core/database/index":           "system/database/index",
	"/core/database/detailed":        "system/database/columns",
	"/core/database/fragment":        "system/database/fragment",
	"/core/database/optimize":        "system/database/optimize",
	"/core/database/recycle":         "system/database/recycle",
	"/core/database/delete":          "system/database/destroy",
	"/core/database/recovery":        "system/database/recover",
	"/core/config/index":             "system/config/index",
	"/core/config/save":              "system/config/create",
	"/core/config/update":            "system/config/update",
	"/core/config/destroy":           "system/config/destroy",
	"/core/logs/getLoginLogPageList": "system/login-log/index",
	"/core/logs/getOperLogPageList":  "system/oper-log/index",
	"/core/notice/index":             "system/notice/index",
	"/core/notice/save":              "system/notice/create",
	"/core/notice/update":            "system/notice/update",
	"/core/notice/destroy":           "system/notice/destroy",
	"/core/notice/read":              "system/notice/read",
	"/core/email/index":              "system/email/index",
	"/core/email/destroy":            "system/email/destroy",
	"/core/system/monitor":           "system/monitor/index",
	"/core/system/getUserList":       "system/user/auth-list",
	"/core/system/getUserInfoByIds":  "system/user/info-by-ids",
	"/tool/code/index":               "system/codegen/index",
	"/tool/code/access":              "system/codegen/access",
	"/tool/crontab/index":            "system/crontab/index",
	"/tool/crontab/save":             "system/crontab/create",
	"/tool/crontab/update":           "system/crontab/update",
	"/tool/crontab/destroy":          "system/crontab/destroy",
	"/tool/crontab/read":             "system/crontab/read",
	"/tool/crontab/changeStatus":     "system/crontab/change-status",
	"/tool/crontab/run":              "system/crontab/run",
	"/tool/crontab/deleteLog":        "system/crontab/delete-log",
}

// ensurePermissionCodes 把菜单表里遗留的旧权限码归一化成新规范。
func ensurePermissionCodes(db *gorm.DB) error {
	for oldCode, newCode := range legacyPermCodeMapping {
		if err := db.Exec("UPDATE `ai_system_menu` SET `code` = ? WHERE `code` = ? AND `type` = 'B'", newCode, oldCode).Error; err != nil {
			return err
		}
	}
	// 种子数据里"登录日志删除/操作日志删除"两行共用了同一个旧码，按名称拆开。
	if err := db.Exec("UPDATE `ai_system_menu` SET `code` = 'system/login-log/destroy' WHERE `code` = '/core/logs/deleteOperLog' AND `name` = '登录日志删除'").Error; err != nil {
		return err
	}
	return db.Exec("UPDATE `ai_system_menu` SET `code` = 'system/oper-log/destroy' WHERE `code` = '/core/logs/deleteOperLog'").Error
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

// removeCodegenFullScreenColumns 清理已废弃的全屏配置列。
func removeCodegenFullScreenColumns(db *gorm.DB) error {
	tables := []string{"nest_tool_generate_tables", "ai_tool_generate_tables"}
	for _, table := range tables {
		if db.Migrator().HasTable(table) && db.Migrator().HasColumn(table, "is_full") {
			if err := db.Migrator().DropColumn(table, "is_full"); err != nil {
				return err
			}
		}
	}
	return nil
}

// ensureCodegenOptionColumns 为已有代码生成字段配置补齐组件数据源列。
func ensureCodegenOptionColumns(db *gorm.DB) error {
	const table = "nest_tool_generate_columns"
	if !db.Migrator().HasColumn(table, "option_source") {
		if err := db.Exec("ALTER TABLE `nest_tool_generate_columns` ADD COLUMN `option_source` varchar(20) NULL DEFAULT NULL COMMENT '选项数据来源' AFTER `dict_type`").Error; err != nil {
			return err
		}
	}
	if !db.Migrator().HasColumn(table, "option_config") {
		if err := db.Exec("ALTER TABLE `nest_tool_generate_columns` ADD COLUMN `option_config` text NULL COMMENT '选项组件配置' AFTER `option_source`").Error; err != nil {
			return err
		}
	}
	return nil
}
