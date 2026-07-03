package system

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	commonResponse "server/model/common/response"
	systemRequest "server/model/system/request"
	"strings"

	"gorm.io/gorm"
)

type databaseTableInfo struct {
	TableName    string  `json:"tableName" gorm:"column:tableName"`
	TableComment string  `json:"tableComment" gorm:"column:tableComment"`
	Engine       string  `json:"engine" gorm:"column:engine"`
	UpdateTime   *string `json:"updateTime" gorm:"column:updateTime"`
	TableRows    uint64  `json:"tableRows" gorm:"column:tableRows"`
	FragmentSize uint64  `json:"fragmentSize" gorm:"column:fragmentSize"`
	DataSize     uint64  `json:"dataSize" gorm:"column:dataSize"`
	IndexSize    uint64  `json:"indexSize" gorm:"column:indexSize"`
	TableCharset string  `json:"tableCharset" gorm:"column:tableCharset"`
	CreateTime   *string `json:"createTime" gorm:"column:createTime"`
}

type databaseColumnInfo struct {
	ColumnName    string  `json:"columnName" gorm:"column:columnName"`
	ColumnComment string  `json:"columnComment" gorm:"column:columnComment"`
	ColumnType    string  `json:"columnType" gorm:"column:columnType"`
	IsNullable    string  `json:"isNullable" gorm:"column:isNullable"`
	ColumnKey     string  `json:"columnKey" gorm:"column:columnKey"`
	ColumnDefault *string `json:"columnDefault" gorm:"column:columnDefault"`
	Extra         string  `json:"extra" gorm:"column:extra"`
	Ordinal       int     `json:"ordinal" gorm:"column:ordinal"`
}

type recycleRecord struct {
	ID         interface{}            `json:"id"`
	DeleteTime interface{}            `json:"deleteTime"`
	Content    map[string]interface{} `json:"content"`
}

var safeIdentifierPattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func DatabaseTableList(query map[string]string) (*commonResponse.PageResult, error) {
	db, schema, err := databaseContext()
	if err != nil {
		return nil, err
	}
	page := parsePage(query)

	baseSQL := `FROM information_schema.tables
WHERE table_schema = ? AND table_type = 'BASE TABLE'`
	args := []interface{}{schema}
	if keyword := strings.TrimSpace(query["tableName"]); keyword != "" {
		baseSQL += ` AND table_name LIKE ?`
		args = append(args, "%"+keyword+"%")
	}

	var total int64
	if err := db.Raw(`SELECT COUNT(*) `+baseSQL, args...).Scan(&total).Error; err != nil {
		return nil, err
	}

	var rows []databaseTableInfo
	err = db.Raw(`SELECT
  table_name AS tableName,
  COALESCE(table_comment, '') AS tableComment,
  COALESCE(engine, '') AS engine,
  DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s') AS updateTime,
  COALESCE(table_rows, 0) AS tableRows,
  COALESCE(data_free, 0) AS fragmentSize,
  COALESCE(data_length, 0) AS dataSize,
  COALESCE(index_length, 0) AS indexSize,
  COALESCE(SUBSTRING_INDEX(table_collation, '_', 1), '') AS tableCharset,
  DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s') AS createTime
`+baseSQL+` ORDER BY table_name ASC LIMIT ? OFFSET ?`,
		append(args, page.Size, (page.Page-1)*page.Size)...,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return &commonResponse.PageResult{List: rows, Total: total}, nil
}

func DatabaseTableColumns(tableName string) ([]databaseColumnInfo, error) {
	db, schema, err := databaseContext()
	if err != nil {
		return nil, err
	}
	if err := ensureTableExists(db, schema, tableName); err != nil {
		return nil, err
	}

	var rows []databaseColumnInfo
	err = db.Raw(`SELECT
  column_name AS columnName,
  COALESCE(column_comment, '') AS columnComment,
  column_type AS columnType,
  is_nullable AS isNullable,
  column_key AS columnKey,
  column_default AS columnDefault,
  extra AS extra,
  ordinal_position AS ordinal
FROM information_schema.columns
WHERE table_schema = ? AND table_name = ?
ORDER BY ordinal_position ASC`, schema, tableName).Scan(&rows).Error
	return rows, err
}

func DatabaseOptimizeTables(payload systemRequest.DatabaseTablesPayload) (map[string]interface{}, error) {
	return maintainTables(payload.Tables, "OPTIMIZE TABLE")
}

func DatabaseClearFragments(payload systemRequest.DatabaseTablesPayload) (map[string]interface{}, error) {
	return maintainTables(payload.Tables, "OPTIMIZE TABLE")
}

func DatabaseRecycleList(query map[string]string) (*commonResponse.PageResult, error) {
	db, schema, err := databaseContext()
	if err != nil {
		return nil, err
	}
	tableName := strings.TrimSpace(query["tableName"])
	if err := ensureRecyclableTable(db, schema, tableName); err != nil {
		return nil, err
	}
	page := parsePage(query)

	var total int64
	if err := db.Table(quoteIdentifier(tableName)).Where("delete_time IS NOT NULL").Count(&total).Error; err != nil {
		return nil, err
	}

	rows, err := queryRecycleRows(db, tableName, page.Size, (page.Page-1)*page.Size)
	if err != nil {
		return nil, err
	}
	return &commonResponse.PageResult{List: rows, Total: total}, nil
}

func DatabaseRecycleRecover(payload systemRequest.DatabaseRecyclePayload) error {
	return updateRecycleRows(payload, map[string]interface{}{"delete_time": nil})
}

func DatabaseRecycleDestroy(payload systemRequest.DatabaseRecyclePayload) error {
	db, schema, err := databaseContext()
	if err != nil {
		return err
	}
	tableName := strings.TrimSpace(payload.TableName)
	if err := ensureRecyclableTable(db, schema, tableName); err != nil {
		return err
	}
	if len(payload.IDs) == 0 {
		return ErrNoRowsSelected
	}
	return db.Exec("DELETE FROM "+quoteIdentifier(tableName)+" WHERE id IN ? AND delete_time IS NOT NULL", payload.IDs).Error
}

func maintainTables(tables []string, command string) (map[string]interface{}, error) {
	db, schema, err := databaseContext()
	if err != nil {
		return nil, err
	}
	if len(tables) == 0 {
		return nil, ErrNoTablesSelected
	}

	done := make([]string, 0, len(tables))
	for _, table := range tables {
		tableName := strings.TrimSpace(table)
		if err := ensureTableExists(db, schema, tableName); err != nil {
			return nil, err
		}
		if err := db.Exec(command + " " + quoteIdentifier(tableName)).Error; err != nil {
			return nil, err
		}
		done = append(done, tableName)
	}
	return map[string]interface{}{"tables": done}, nil
}

func updateRecycleRows(payload systemRequest.DatabaseRecyclePayload, values map[string]interface{}) error {
	db, schema, err := databaseContext()
	if err != nil {
		return err
	}
	tableName := strings.TrimSpace(payload.TableName)
	if err := ensureRecyclableTable(db, schema, tableName); err != nil {
		return err
	}
	if len(payload.IDs) == 0 {
		return ErrNoRowsSelected
	}
	return db.Table(quoteIdentifier(tableName)).Where("id IN ? AND delete_time IS NOT NULL", payload.IDs).Updates(values).Error
}

func queryRecycleRows(db *gorm.DB, tableName string, limit int, offset int) ([]recycleRecord, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE delete_time IS NOT NULL ORDER BY delete_time DESC LIMIT ? OFFSET ?", quoteIdentifier(tableName))
	rows, err := sqlDB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rawRows, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	result := make([]recycleRecord, 0, len(rawRows))
	for _, row := range rawRows {
		result = append(result, recycleRecord{
			ID:         row["id"],
			DeleteTime: row["delete_time"],
			Content:    row,
		})
	}
	return result, nil
}

func scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		item := make(map[string]interface{}, len(columns))
		for i, column := range columns {
			item[column] = normalizeSQLValue(values[i])
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func normalizeSQLValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []byte:
		var decoded interface{}
		text := string(v)
		if json.Unmarshal(v, &decoded) == nil {
			return decoded
		}
		return text
	default:
		return value
	}
}

func databaseContext() (*gorm.DB, string, error) {
	db, err := systemDB()
	if err != nil {
		return nil, "", err
	}
	schema, err := databaseName(db)
	if err != nil {
		return nil, "", err
	}
	return db, schema, nil
}

func ensureRecyclableTable(db *gorm.DB, schema string, tableName string) error {
	if err := ensureTableExists(db, schema, tableName); err != nil {
		return err
	}
	if !hasTableColumn(db, schema, tableName, "delete_time") {
		return ErrNoRecycleSupport
	}
	if !hasTableColumn(db, schema, tableName, "id") {
		return ErrNoIDColumn
	}
	return nil
}

func ensureTableExists(db *gorm.DB, schema string, tableName string) error {
	if !safeIdentifierPattern.MatchString(tableName) {
		return ErrInvalidTableName
	}
	var count int64
	err := db.Raw(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ? AND table_type = 'BASE TABLE'`, schema, tableName).Scan(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTableNotFound
	}
	return nil
}

func hasTableColumn(db *gorm.DB, schema string, tableName string, columnName string) bool {
	var count int64
	err := db.Raw(`SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = ? AND table_name = ? AND column_name = ?`, schema, tableName, columnName).Scan(&count).Error
	return err == nil && count > 0
}

func quoteIdentifier(value string) string {
	return "`" + strings.ReplaceAll(value, "`", "``") + "`"
}
