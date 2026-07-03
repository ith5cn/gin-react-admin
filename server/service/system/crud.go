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

func PageList(query map[string]string, model interface{}, dest interface{}, likes map[string]string, equals map[string]string, order string) (*commonResponse.PageResult, error) {
	return pageList(query, model, dest, likes, equals, order)
}

func CreateRecord[T any](table string, data map[string]interface{}) (*T, error) {
	return createSimple[T](table, data, passthroughColumns(data))
}

func UpdateRecord[T any](table string, id string, data map[string]interface{}) (*T, error) {
	return updateSimple[T](table, id, data, passthroughColumns(data))
}

func DeleteRecord(model interface{}, id string) error {
	return deleteByID(model, id)
}

func passthroughColumns(data map[string]interface{}) map[string]string {
	columns := make(map[string]string, len(data))
	for key := range data {
		columns[key] = camelToSnake(key)
	}
	return columns
}

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

func createSimple[T any](table string, data map[string]interface{}, allowed map[string]string) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	payload := requestData(data, allowed)
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

func updateSimple[T any](table string, id string, data map[string]interface{}, allowed map[string]string) (*T, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	payload := requestData(data, allowed)
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

func parseUint(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	return uint(id), err
}

func idsFromAny(value interface{}) []uint {
	items, ok := value.([]interface{})
	if !ok {
		return []uint{}
	}
	result := make([]uint, 0, len(items))
	for _, item := range items {
		switch v := item.(type) {
		case float64:
			result = append(result, uint(v))
		case int:
			result = append(result, uint(v))
		case string:
			if id, err := parseUint(v); err == nil {
				result = append(result, id)
			}
		}
	}
	return result
}
