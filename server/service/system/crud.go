package system

import (
	"fmt"
	commonResponse "server/model/common/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// 本文件是通用 CRUD 助手，被各业务模块和 codegen 生成的代码共用。
// PageList / CreateRecord / UpdateRecord / DeleteRecord 是 codegen 模板的
// 稳定契约（见 codegen_templates.go），改签名需同步更新模板和已生成代码。

func pageList(query map[string]string, model interface{}, dest interface{}, likes map[string]string, equals map[string]string, order string) (*commonResponse.PageResult, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	page := parsePage(query)
	base := softDelete(db.Model(model))
	base = applyFilters(base, query, likes, equals)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := base.Order(order).Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(dest).Error; err != nil {
		return nil, err
	}
	return &commonResponse.PageResult{List: dest, Total: total}, nil
}

// PageList 是老版分页契约：likes 全部按 LIKE、equals 全部按 = 处理。
// 仅为兼容早期 codegen 生成的代码保留，新生成的代码走 PageListFiltered。
func PageList(query map[string]string, model interface{}, dest interface{}, likes map[string]string, equals map[string]string, order string) (*commonResponse.PageResult, error) {
	return pageList(query, model, dest, likes, equals, order)
}

// CreateRecord 是 codegen 生成代码的通用创建入口：
// 把前端传来的 camelCase 字段名转成 snake_case 列名后插入。
func CreateRecord[T any](table string, data map[string]interface{}) (*T, error) {
	return createRow[T](table, requestData(data, passthroughColumns(data)))
}

// UpdateRecord 是 codegen 生成代码的通用更新入口，逻辑同 CreateRecord。
func UpdateRecord[T any](table string, id string, data map[string]interface{}) (*T, error) {
	return updateRow[T](table, id, requestData(data, passthroughColumns(data)))
}

// DeleteRecord 是 codegen 生成代码的通用硬删除入口。
func DeleteRecord(model interface{}, id string) error {
	return deleteByID(model, id)
}

// QueryFilter 描述一个查询条件：Param 是前端参数名，Column 是数据库列名，
// Op 是操作符（eq/neq/gt/gte/lt/lte/like/in/notin/between）。
type QueryFilter struct {
	Param  string
	Column string
	Op     string
}

// PageListFiltered 是 codegen 生成代码的分页查询契约（PageList 的增强版）：
// 支持全量查询操作符、orderBy/orderType 排序白名单和软删除开关。
func PageListFiltered(query map[string]string, model interface{}, dest interface{}, filters []QueryFilter, sortable map[string]string, defaultOrder string, useSoftDelete bool) (*commonResponse.PageResult, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	page := parsePage(query)
	base := db.Model(model)
	if useSoftDelete {
		base = softDelete(base)
	}
	base = applyQueryFilters(base, query, filters)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	order := resolveOrder(query, sortable, defaultOrder)
	if err := base.Order(order).Offset((page.Page - 1) * page.Size).Limit(page.Size).Find(dest).Error; err != nil {
		return nil, err
	}
	return &commonResponse.PageResult{List: dest, Total: total}, nil
}

// SoftDeleteRecord 通过 UPDATE delete_time 实现软删除，供 codegen 生成的软删模型使用。
func SoftDeleteRecord(table string, id string) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Table(table).Where("id = ?", id).
		Updates(map[string]interface{}{"delete_time": gorm.Expr("NOW()"), "update_time": gorm.Expr("NOW()")}).Error
}

func applyQueryFilters(db *gorm.DB, query map[string]string, filters []QueryFilter) *gorm.DB {
	for _, filter := range filters {
		value := query[filter.Param]
		if value == "" {
			continue
		}
		column := filter.Column
		switch filter.Op {
		case "neq":
			db = db.Where(column+" <> ?", value)
		case "gt":
			db = db.Where(column+" > ?", value)
		case "gte":
			db = db.Where(column+" >= ?", value)
		case "lt":
			db = db.Where(column+" < ?", value)
		case "lte":
			db = db.Where(column+" <= ?", value)
		case "like":
			db = db.Where(column+" LIKE ?", "%"+value+"%")
		case "in":
			db = db.Where(column+" IN ?", splitQueryValues(value))
		case "notin":
			db = db.Where(column+" NOT IN ?", splitQueryValues(value))
		case "between":
			parts := splitQueryValues(value)
			if len(parts) >= 2 {
				db = db.Where(column+" BETWEEN ? AND ?", parts[0], parts[1])
			}
		default:
			db = db.Where(column+" = ?", value)
		}
	}
	return db
}

// resolveOrder 只接受排序白名单里的列，避免 orderBy 参数注入任意 SQL。
func resolveOrder(query map[string]string, sortable map[string]string, defaultOrder string) string {
	column, ok := sortable[query["orderBy"]]
	if !ok || column == "" {
		return defaultOrder
	}
	direction := "ASC"
	if strings.EqualFold(query["orderType"], "desc") {
		direction = "DESC"
	}
	return column + " " + direction
}

func splitQueryValues(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// passthroughColumns 为动态数据生成"参数名 → 列名"映射（camelCase → snake_case）。
func passthroughColumns(data map[string]interface{}) map[string]string {
	columns := make(map[string]string, len(data))
	for key := range data {
		columns[key] = camelToSnake(key)
	}
	return columns
}

// createWithLevel 在插入前先维护树形表的 parent_id/level 字段，
// 用于菜单、角色、部门这类有层级关系的表。
func createWithLevel[T any](table string, payload map[string]interface{}) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	normalizeParentAndLevel(db, table, payload)
	setDefaultTimes(payload, true)
	if err := db.Table(table).Create(payload).Error; err != nil {
		return nil, err
	}
	var result T
	if err := db.Table(table).Order("id DESC").First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// setColumn 在 value 非 nil 时把解引用后的值写入 payload，用于类型化入参转 GORM 更新 map。
func setColumn[T any](payload map[string]interface{}, column string, value *T) {
	if value != nil {
		payload[column] = *value
	}
}

// createRow 通用插入：写入创建/更新时间后 INSERT，再回查一条最新记录返回。
// 注意：用 map 插入拿不到自增 ID，这里按 id DESC 回查，高并发下可能取错行。
func createRow[T any](table string, payload map[string]interface{}) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	setDefaultTimes(payload, true)
	if err := db.Table(table).Create(payload).Error; err != nil {
		return nil, err
	}
	var result T
	if err := db.Table(table).Order("id DESC").First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// updateRow 通用更新：只更新 payload 里出现的列，然后回查该行返回。
func updateRow[T any](table string, id string, payload map[string]interface{}) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	setDefaultTimes(payload, false)
	if err := db.Table(table).Where("id = ?", id).Updates(payload).Error; err != nil {
		return nil, err
	}
	var result T
	if err := db.Table(table).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// updateWithLevel 更新树形表数据，父级变化时同步重算 level。
func updateWithLevel[T any](table string, id string, payload map[string]interface{}) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	normalizeParentAndLevel(db, table, payload)
	setDefaultTimes(payload, false)
	if err := db.Table(table).Where("id = ?", id).Updates(payload).Error; err != nil {
		return nil, err
	}
	var result T
	if err := db.Table(table).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// normalizeParentAndLevel 维护树形表的层级路径：
// 无父节点时 parent_id=0、level="0"；有父节点时 level = 父节点 level + "," + 父ID，
// 这样查某节点的所有祖先只需要按逗号拆 level 字符串。
func normalizeParentAndLevel(db *gorm.DB, table string, payload map[string]interface{}) {
	parentID, ok := payload["parent_id"]
	if !ok || fmt.Sprint(parentID) == "" || fmt.Sprint(parentID) == "0" || fmt.Sprint(parentID) == "<nil>" {
		payload["parent_id"] = 0
		payload["level"] = "0"
		return
	}
	var parent struct {
		Level *string `gorm:"column:level"`
	}
	if err := db.Table(table).Select("level").Where("id = ?", parentID).First(&parent).Error; err == nil && parent.Level != nil && *parent.Level != "" {
		payload["level"] = strings.Trim(*parent.Level+","+fmt.Sprint(parentID), ",")
		return
	}
	payload["level"] = "0"
}

// parseUint 把字符串 ID 转成 uint，路径参数转数字的通用小工具。
func parseUint(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	return uint(id), err
}
