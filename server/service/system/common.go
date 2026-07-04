package system

import (
	"errors"
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

// systemDB 获取 ai_system 库连接。
// 所有 service 都从这里拿连接，未初始化时返回明确错误而不是 nil 指针。
func systemDB() (*gorm.DB, error) {
	if gormInit.Gorm.Databases == nil || gormInit.Gorm.Databases.AISystem == nil {
		return nil, errors.New("ai_system is not initialized")
	}
	return gormInit.Gorm.Databases.AISystem, nil
}

// parsePage 解析分页参数，page 从 1 开始；size 兼容 limit 参数名，默认 10。
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

// intFromQuery 从查询参数取整数，缺失或非法时用默认值。
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

// applyFilters 按两组映射拼查询条件：likes 生成 LIKE '%v%'，equals 生成 =。
// 列名来自代码里写死的映射（不是用户输入），配合 ? 占位符防止 SQL 注入。
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

// softDelete 统一追加软删除过滤条件。
// 本项目用业务列 delete_time 实现软删除，而不是 GORM 内置的 gorm.DeletedAt。
func softDelete(db *gorm.DB) *gorm.DB {
	return db.Where("delete_time IS NULL")
}

// hashPassword 用 bcrypt 加密密码。
// bcrypt 自带随机盐且计算慢，能有效对抗彩虹表和暴力破解。
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// requestData 按白名单过滤动态 map：只保留 allowed 里声明的字段并映射为列名，
// 防止前端多传字段直接写库（Mass Assignment 漏洞）。
func requestData(data map[string]interface{}, allowed map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, column := range allowed {
		if value, ok := data[key]; ok {
			result[column] = normalizeValue(value)
		}
	}
	return result
}

// normalizeValue 修正 JSON 反序列化的类型偏差：
// encoding/json 会把所有数字解析成 float64，整数值在这里还原成 int64。
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

// sortTreeChildren 递归排序树的每一层子节点，泛型写法对任意节点类型通用。
func sortTreeChildren[T any](nodes []*T, getChildren func(*T) []*T, less func(a *T, b *T) bool) {
	sort.SliceStable(nodes, func(i, j int) bool { return less(nodes[i], nodes[j]) })
	for _, node := range nodes {
		children := getChildren(node)
		if len(children) > 0 {
			sortTreeChildren(children, getChildren, less)
		}
	}
}

// isZeroParent 判断节点是否是根节点（无父节点或父 ID 为 0）。
func isZeroParent(parentID *uint) bool {
	return parentID == nil || *parentID == 0
}

// setDefaultTimes 写入创建/更新时间。
// 用 gorm.Expr("NOW()") 让数据库生成时间，避免应用服务器时钟不一致。
func setDefaultTimes(data map[string]interface{}, isCreate bool) map[string]interface{} {
	if isCreate {
		data["create_time"] = gorm.Expr("NOW()")
	}
	data["update_time"] = gorm.Expr("NOW()")
	return data
}

// deleteByID 按 ID 硬删除一行（模型没有 gorm.DeletedAt，Delete 就是真删）。
func deleteByID(model interface{}, id string) error {
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Delete(model, "id = ?", id).Error
}

// hasChildren 判断树形表某节点下是否还有未删除的子节点。
func hasChildren(table string, parentID string) (bool, error) {
	db, err := systemDB()
	if err != nil {
		return false, err
	}
	var count int64
	err = db.Table(table).Where("parent_id = ? AND delete_time IS NULL", parentID).Count(&count).Error
	return count > 0, err
}

// camelToSnake 把 camelCase 转成 snake_case，如 userName → user_name。
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
