package system

import (
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"

	gormInit "server/setup/gorm"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PageQuery struct {
	Page int
	Size int
}

func systemDB() (*gorm.DB, error) {
	if gormInit.Gorm.Databases == nil || gormInit.Gorm.Databases.AISystem == nil {
		return nil, errors.New("ai_system is not initialized")
	}
	return gormInit.Gorm.Databases.AISystem, nil
}

func parsePage(query map[string]string) PageQuery {
	page := intFromQuery(query, "page", 1)
	size := intFromQuery(query, "size", intFromQuery(query, "limit", 10))
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	return PageQuery{Page: page, Size: size}
}

func intFromQuery(query map[string]string, key string, fallback int) int {
	if query[key] == "" {
		return fallback
	}
	value, err := strconv.Atoi(query[key])
	if err != nil {
		return fallback
	}
	return value
}

func applyFilters(db *gorm.DB, query map[string]string, likes map[string]string, equals map[string]string) *gorm.DB {
	for param, column := range likes {
		if query[param] != "" {
			db = db.Where(column+" LIKE ?", "%"+query[param]+"%")
		}
	}
	for param, column := range equals {
		if query[param] != "" {
			db = db.Where(column+" = ?", query[param])
		}
	}
	return db
}

func queryMap(query map[string][]string) map[string]string {
	result := make(map[string]string, len(query))
	for key, values := range query {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

func softDelete(db *gorm.DB) *gorm.DB {
	return db.Where("delete_time IS NULL")
}

func ptrUint(v uint) *uint {
	return &v
}

func uintSlice(values []uint) []uint {
	if values == nil {
		return []uint{}
	}
	return values
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func requestData(data map[string]interface{}, allowed map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, column := range allowed {
		if value, ok := data[key]; ok {
			result[column] = normalizeValue(value)
		}
	}
	return result
}

func normalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case float64:
		if v == float64(int64(v)) {
			return int64(v)
		}
		return v
	default:
		return value
	}
}

func sortTreeChildren[T any](nodes []*T, getChildren func(*T) []*T, less func(a *T, b *T) bool) {
	sort.SliceStable(nodes, func(i, j int) bool { return less(nodes[i], nodes[j]) })
	for _, node := range nodes {
		children := getChildren(node)
		if len(children) > 0 {
			sortTreeChildren(children, getChildren, less)
		}
	}
}

func isZeroParent(parentID *uint) bool {
	return parentID == nil || *parentID == 0
}

func setDefaultTimes(data map[string]interface{}, isCreate bool) map[string]interface{} {
	if isCreate {
		data["create_time"] = gorm.Expr("NOW()")
	}
	data["update_time"] = gorm.Expr("NOW()")
	return data
}

func deleteByID(model interface{}, id string) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Delete(model, "id = ?", id).Error
}

func hasChildren(table string, parentID string) (bool, error) {
	db, err := systemDB()
	if err != nil {
		return false, err
	}
	var count int64
	err = db.Table(table).Where("parent_id = ? AND delete_time IS NULL", parentID).Count(&count).Error
	return count > 0, err
}

func hasColumnValue(v interface{}, field string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	return rv.FieldByName(field).IsValid()
}

func camelToSnake(value string) string {
	var out strings.Builder
	for index, r := range value {
		if unicode.IsUpper(r) {
			if index > 0 {
				out.WriteByte('_')
			}
			out.WriteRune(unicode.ToLower(r))
			continue
		}
		out.WriteRune(r)
	}
	return out.String()
}
