package system

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"server/config"
	commonResponse "server/model/common/response"
	systemModel "server/model/system"
	systemRequest "server/model/system/request"
	gormInit "server/setup/gorm"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type codegenTableMeta struct {
	TableName      string  `gorm:"column:TABLE_NAME" json:"TABLE_NAME"`
	TableComment   *string `gorm:"column:TABLE_COMMENT" json:"TABLE_COMMENT"`
	Engine         *string `gorm:"column:ENGINE" json:"ENGINE"`
	TableCollation *string `gorm:"column:TABLE_COLLATION" json:"TABLE_COLLATION"`
	CreateTime     *string `gorm:"column:CREATE_TIME" json:"CREATE_TIME"`
}

type codegenColumnMeta struct {
	ColumnName      string  `gorm:"column:COLUMN_NAME"`
	ColumnComment   *string `gorm:"column:COLUMN_COMMENT"`
	ColumnType      *string `gorm:"column:COLUMN_TYPE"`
	ColumnDefault   *string `gorm:"column:COLUMN_DEFAULT"`
	IsNullable      string  `gorm:"column:IS_NULLABLE"`
	ColumnKey       *string `gorm:"column:COLUMN_KEY"`
	OrdinalPosition int     `gorm:"column:ORDINAL_POSITION"`
	DataType        *string `gorm:"column:DATA_TYPE"`
}

type codegenContext struct {
	Table              systemModel.ToolGenerateTable
	Columns            []systemModel.ToolGenerateColumn
	PackageName        string
	BusinessName       string
	BusinessApiName    string
	ClassName          string
	EntityVarName      string
	RoutePath          string
	BackendModelPath   string
	BackendAPIPath     string
	BackendRoutePath   string
	BackendServicePath string
	FrontendPageDir    string
	FrontendAPIPath    string
	QueryColumns       []systemModel.ToolGenerateColumn
	ListColumns        []systemModel.ToolGenerateColumn
	FormColumns        []systemModel.ToolGenerateColumn
	EditableColumns    []systemModel.ToolGenerateColumn
}

// CodegenList 分页查询已装载的代码生成表配置。
func CodegenList(query map[string]string) (*commonResponse.PageResult, error) {
	var data []systemModel.ToolGenerateTable
	return pageList(query, &systemModel.ToolGenerateTable{}, &data, map[string]string{"table_name": "table_name"}, map[string]string{"source": "source"}, "id ASC")
}

// CodegenDatasources 返回可选数据源列表：当前库 + information_schema 里可见的其他库。
func CodegenDatasources() ([]map[string]interface{}, error) {
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	dbName, err := databaseName(db)
	if err != nil {
		return nil, err
	}
	options := []map[string]interface{}{
		{"value": config.MysqlAISystem, "label": dbName, "databaseName": dbName},
	}
	if dbName != config.MysqlAISystem {
		options = append(options, map[string]interface{}{"value": dbName, "label": dbName, "databaseName": dbName})
	}
	var schemaNames []string
	if err := db.Raw(`SELECT SCHEMA_NAME FROM information_schema.SCHEMATA
WHERE SCHEMA_NAME NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')
ORDER BY SCHEMA_NAME ASC`).Scan(&schemaNames).Error; err == nil {
		seen := map[string]struct{}{config.MysqlAISystem: {}, dbName: {}}
		for _, name := range schemaNames {
			if _, ok := seen[name]; ok {
				continue
			}
			seen[name] = struct{}{}
			options = append(options, map[string]interface{}{"value": name, "label": name, "databaseName": name})
		}
	}
	return options, nil
}

// CodegenDBTables 从 information_schema 读取指定库的数据表清单，
// 支持关键字过滤，分页在内存中完成（表数量通常很小）。
func CodegenDBTables(query map[string]string) (*commonResponse.PageResult, error) {
	source := sourceOrDefault(query["source"])
	db, dbName, err := codegenSourceDB(source)
	if err != nil {
		return nil, err
	}
	var rows []codegenTableMeta
	if err := db.Raw(`SELECT TABLE_NAME, TABLE_COMMENT, ENGINE, TABLE_COLLATION, DATE_FORMAT(CREATE_TIME, '%Y-%m-%d %H:%i:%s') AS CREATE_TIME
FROM information_schema.tables
WHERE table_schema = ? AND table_type = 'BASE TABLE'
ORDER BY TABLE_NAME ASC`, dbName).Scan(&rows).Error; err != nil {
		return nil, err
	}

	keyword := strings.ToLower(strings.TrimSpace(query["keyword"]))
	filtered := make([]codegenTableMeta, 0, len(rows))
	for _, row := range rows {
		if keyword == "" || strings.Contains(strings.ToLower(row.TableName), keyword) || strings.Contains(strings.ToLower(stringValue(row.TableComment)), keyword) {
			filtered = append(filtered, row)
		}
	}

	page := parsePage(query)
	total := int64(len(filtered))
	start := (page.Page - 1) * page.Size
	if start >= len(filtered) {
		return &commonResponse.PageResult{List: []codegenTableMeta{}, Total: total}, nil
	}
	end := start + page.Size
	if end > len(filtered) {
		end = len(filtered)
	}
	return &commonResponse.PageResult{List: filtered[start:end], Total: total}, nil
}

// CodegenImportTables 把选中的数据表结构导入生成配置。
// 重复导入会保留表级配置、重建字段配置（以数据库最新结构为准）。
func CodegenImportTables(payload systemRequest.CodegenImportPayload) ([]map[string]interface{}, error) {
	source := sourceOrDefault(payload.Source)
	if len(payload.Tables) == 0 {
		return nil, ErrNoImportTables
	}
	sourceDB, dbName, err := codegenSourceDB(source)
	if err != nil {
		return nil, err
	}
	system, err := systemDB()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0, len(payload.Tables))
	err = system.Transaction(func(tx *gorm.DB) error {
		for _, item := range payload.Tables {
			result, err := importCodegenTable(tx, sourceDB, dbName, source, item)
			if err != nil {
				return err
			}
			results = append(results, result)
		}
		return nil
	})
	return results, err
}

// CodegenDelete 删除生成配置，同一事务里先删字段再删表配置。
func CodegenDelete(ids []uint) error {
	if len(ids) == 0 {
		return ErrNoDeleteTables
	}
	db, err := systemDB()
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("table_id IN ?", ids).Delete(&systemModel.ToolGenerateColumn{}).Error; err != nil {
			return err
		}
		return tx.Delete(&systemModel.ToolGenerateTable{}, "id IN ?", ids).Error
	})
}

// CodegenDetail 返回一张表的生成配置详情（表配置 + 字段配置）。
func CodegenDetail(id string) (map[string]interface{}, error) {
	table, columns, err := codegenDetail(id)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"table": table, "columns": columns}, nil
}

// CodegenUpdate 保存生成配置。字段配置采用整体替换（先删后插），
// 因为前端每次都会提交完整字段列表。
func CodegenUpdate(id string, payload map[string]interface{}) (map[string]interface{}, error) {
	table, currentColumns, err := codegenDetail(id)
	if err != nil {
		return nil, err
	}
	db, err := systemDB()
	if err != nil {
		return nil, err
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		tablePatch := mapValue(payload["table"])
		tableData := requestData(tablePatch, codegenTableColumns())
		setDefaultTimes(tableData, false)
		if len(tableData) > 0 {
			if err := tx.Model(&systemModel.ToolGenerateTable{}).Where("id = ?", id).Updates(tableData).Error; err != nil {
				return err
			}
		}

		nextColumns := columnsFromAny(payload["columns"])
		if nextColumns == nil {
			nextColumns = currentColumns
		}
		if err := tx.Where("table_id = ?", id).Delete(&systemModel.ToolGenerateColumn{}).Error; err != nil {
			return err
		}
		for index := range nextColumns {
			nextColumns[index].ID = 0
			nextColumns[index].TableID = table.ID
			if nextColumns[index].QueryType == "" {
				nextColumns[index].QueryType = "eq"
			}
			if nextColumns[index].ViewType == "" {
				nextColumns[index].ViewType = "input"
			}
			if err := tx.Create(&nextColumns[index]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return CodegenDetail(id)
}

// CodegenPreview 渲染全部模板但不落盘，供前端预览生成结果。
func CodegenPreview(id string) (map[string]interface{}, error) {
	ctx, err := buildCodegenContext(id)
	if err != nil {
		return nil, err
	}
	files := append(buildGoPreviewFiles(ctx), buildFrontendPreviewFiles(ctx)...)
	return map[string]interface{}{"files": files}, nil
}

// CodegenGenerate 渲染模板并写入磁盘：后端文件写到 server/ 相应目录，
// 前端文件写到 generate_path 指定的项目，最后刷新路由注册表、同步菜单。
// 注意：新写入的 Go 文件需要重新编译并重启服务才会生效。
func CodegenGenerate(id string) (map[string]interface{}, error) {
	ctx, err := buildCodegenContext(id)
	if err != nil {
		return nil, err
	}
	files := append(buildGoPreviewFiles(ctx), buildFrontendPreviewFiles(ctx)...)
	written := make([]string, 0, len(files))
	for _, file := range files {
		root := "."
		if file.Group == "frontend" {
			root = resolveFrontendRoot(stringValue(ctx.Table.GeneratePath))
		}
		abs := filepath.Join(root, file.Path)
		if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(abs, []byte(file.Content), 0644); err != nil {
			return nil, err
		}
		written = append(written, abs)
	}
	if err := refreshGeneratedRouteRegistry(); err != nil {
		return nil, err
	}
	menuIDs, _ := syncCodegenMenus(ctx)
	return map[string]interface{}{"generated": true, "files": written, "menuIds": menuIDs}, nil
}

// refreshGeneratedRouteRegistry 扫描 router/generated 下所有 RegisterXxxRoutes 函数，
// 重新生成 register.go 汇总注册表，让新生成的路由随下次编译自动挂载。
func refreshGeneratedRouteRegistry() error {
	const routerDir = "router/generated"
	entries, err := os.ReadDir(routerDir)
	if err != nil {
		return err
	}

	pattern := regexp.MustCompile(`func\s+(Register\w+Routes)\s*\(`)
	functions := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "register.go" || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		content, err := os.ReadFile(filepath.Join(routerDir, entry.Name()))
		if err != nil {
			return err
		}
		matches := pattern.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			if len(match) == 2 {
				functions = append(functions, match[1])
			}
		}
	}
	sort.Strings(functions)

	var builder strings.Builder
	builder.WriteString("// Code generated by gin-react-admin codegen. DO NOT EDIT.\n")
	builder.WriteString("package generated\n\n")
	builder.WriteString("import \"github.com/gin-gonic/gin\"\n\n")
	builder.WriteString("func RegisterRoutes(group *gin.RouterGroup) {\n")
	for _, name := range functions {
		builder.WriteString("\t")
		builder.WriteString(name)
		builder.WriteString("(group)\n")
	}
	builder.WriteString("}\n")

	return os.WriteFile(filepath.Join(routerDir, "register.go"), []byte(builder.String()), 0644)
}

// importCodegenTable 导入单张表：读表注释和字段元数据，落库为生成配置。
func importCodegenTable(tx *gorm.DB, sourceDB *gorm.DB, dbName, source string, item systemRequest.CodegenImportTable) (map[string]interface{}, error) {
	tableName := strings.TrimSpace(item.TableName)
	if tableName == "" {
		return nil, ErrEmptyTableName
	}
	var tableMeta codegenTableMeta
	if err := sourceDB.Raw(`SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.tables WHERE table_schema = ? AND table_name = ? AND table_type = 'BASE TABLE' LIMIT 1`, dbName, tableName).Scan(&tableMeta).Error; err != nil {
		return nil, err
	}
	if tableMeta.TableName == "" {
		return nil, NewBizError(fmt.Sprintf("数据表 %s 不存在", tableName))
	}
	var columns []codegenColumnMeta
	if err := sourceDB.Raw(`SELECT COLUMN_NAME, COLUMN_COMMENT, COLUMN_TYPE, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_KEY, ORDINAL_POSITION, DATA_TYPE
FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ORDINAL_POSITION ASC`, dbName, tableName).Scan(&columns).Error; err != nil {
		return nil, err
	}
	if len(columns) == 0 {
		return nil, NewBizError(fmt.Sprintf("数据表 %s 未查询到字段信息", tableName))
	}

	var existing systemModel.ToolGenerateTable
	err := tx.Where("source = ? AND table_name = ?", source, tableName).First(&existing).Error
	action := "created"
	comment := coalesceString(tableMeta.TableComment, item.TableComment, tableName)
	tableRow := systemModel.ToolGenerateTable{
		TableNameValue: tableName,
		TableComment:   ptrString(comment),
		PackageName:    ptrString(defaultString(stringValue(existing.PackageName), "system")),
		BusinessName:   ptrString(defaultString(stringValue(existing.BusinessName), toKebabCase(tableName))),
		ClassName:      ptrString(defaultString(stringValue(existing.ClassName), toPascalCase(tableName))),
		MenuName:       ptrString(defaultString(stringValue(existing.MenuName), comment)),
		GeneratePath:   ptrString(defaultString(stringValue(existing.GeneratePath), "web")),
		GenerateModel:  defaultInt16(existing.GenerateModel, 1),
		FormWidth:      defaultInt(existing.FormWidth, 600),
		IsFull:         defaultInt16(existing.IsFull, 1),
		Source:         ptrString(source),
		ComponentType:  defaultInt16(existing.ComponentType, 1),
		Sort:           existing.Sort,
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := tx.Create(&tableRow).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		action = "updated"
		tableRow.ID = existing.ID
		if err := tx.Model(&systemModel.ToolGenerateTable{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
			"table_comment": comment,
			"update_time":   gorm.Expr("NOW()"),
		}).Error; err != nil {
			return nil, err
		}
	}

	if err := tx.Where("table_id = ?", tableRow.ID).Delete(&systemModel.ToolGenerateColumn{}).Error; err != nil {
		return nil, err
	}
	for index, column := range columns {
		row := mapCodegenColumn(tableRow.ID, column, index, len(columns))
		if err := tx.Create(&row).Error; err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{"tableName": tableName, "action": action, "tableId": tableRow.ID, "columnCount": len(columns)}, nil
}

// codegenDetail 读取表配置和字段配置；字段按 sort 倒序保持原始列顺序。
func codegenDetail(id string) (systemModel.ToolGenerateTable, []systemModel.ToolGenerateColumn, error) {
	db, err := systemDB()
	if err != nil {
		return systemModel.ToolGenerateTable{}, nil, err
	}
	var table systemModel.ToolGenerateTable
	if err := db.Where("id = ?", id).First(&table).Error; err != nil {
		return table, nil, err
	}
	var columns []systemModel.ToolGenerateColumn
	err = db.Where("table_id = ?", id).Order("sort DESC, id ASC").Find(&columns).Error
	return table, columns, err
}

// codegenSourceDB 按数据源名称拿连接：默认当前库，其他名称动态建连接。
func codegenSourceDB(source string) (*gorm.DB, string, error) {
	if source == "" || source == config.MysqlAISystem {
		db, err := systemDB()
		if err != nil {
			return nil, "", err
		}
		name, err := databaseName(db)
		return db, name, err
	}
	db, err := gormInit.Gorm.InitializeByName(source)
	if err != nil {
		return nil, "", err
	}
	return db, config.MysqlByName(source).Dbname, nil
}

// databaseName 查询当前连接实际使用的数据库名。
func databaseName(db *gorm.DB) (string, error) {
	var name string
	if err := db.Raw("SELECT DATABASE()").Scan(&name).Error; err != nil {
		return "", err
	}
	return name, nil
}

// mapCodegenColumn 把 information_schema 的列元数据换算成生成配置的默认值：
// 主键/系统列不进表单，常见列名（name/status/sort 等）默认勾选列表与查询。
func mapCodegenColumn(tableID uint, column codegenColumnMeta, index int, count int) systemModel.ToolGenerateColumn {
	name := strings.ToLower(column.ColumnName)
	isPK := int16(1)
	if stringValue(column.ColumnKey) == "PRI" {
		isPK = 2
	}
	isRequired := int16(1)
	if column.IsNullable == "NO" {
		isRequired = 2
	}
	isInsert := int16(2)
	if isSystemColumnName(name) {
		isInsert = 1
	}
	isEdit := int16(2)
	if isPK == 2 || isSystemColumnName(name) {
		isEdit = 1
	}
	return systemModel.ToolGenerateColumn{
		TableID:       tableID,
		ColumnName:    column.ColumnName,
		ColumnComment: column.ColumnComment,
		ColumnType:    column.ColumnType,
		DefaultValue:  column.ColumnDefault,
		IsPK:          isPK,
		IsRequired:    isRequired,
		IsInsert:      isInsert,
		IsEdit:        isEdit,
		IsList:        boolFlag(name == "name" || name == "title" || name == "code" || name == "status" || name == "sort"),
		IsQuery:       boolFlag(name == "name" || name == "title" || name == "code" || name == "status"),
		IsSort:        boolFlag(name == "sort"),
		QueryType:     inferQueryType(stringValue(column.DataType)),
		ViewType:      inferViewType(name, stringValue(column.DataType)),
		Sort:          uint8(maxInt(count-column.OrdinalPosition+1, 1)),
	}
}

// buildCodegenContext 把表配置和字段配置整理成模板渲染上下文，
// 预先算好类名、路由路径、各类字段分组（查询/列表/表单）。
func buildCodegenContext(id string) (codegenContext, error) {
	table, columns, err := codegenDetail(id)
	if err != nil {
		return codegenContext{}, err
	}
	ctx := codegenContext{
		Table:              table,
		Columns:            columns,
		PackageName:        toKebabCase(defaultString(stringValue(table.PackageName), "system")),
		BusinessName:       toKebabCase(defaultString(stringValue(table.BusinessName), table.TableNameValue)),
		ClassName:          toPascalCase(defaultString(stringValue(table.ClassName), table.TableNameValue)),
		BackendModelPath:   filepath.Join("model", "generated"),
		BackendAPIPath:     filepath.Join("api", "generated"),
		BackendRoutePath:   filepath.Join("router", "generated"),
		BackendServicePath: filepath.Join("service", "generated"),
	}
	ctx.BusinessApiName = toCamelCase(ctx.BusinessName)
	ctx.EntityVarName = toCamelCase(ctx.ClassName)
	ctx.RoutePath = ctx.PackageName + "/" + ctx.BusinessName
	ctx.FrontendPageDir = filepath.Join("src", "pages", ctx.PackageName, ctx.BusinessName)
	ctx.FrontendAPIPath = filepath.Join("src", "api", ctx.PackageName, ctx.BusinessName+".ts")
	for _, column := range columns {
		if column.IsQuery == 2 && !isRichOrUpload(column.ViewType) {
			ctx.QueryColumns = append(ctx.QueryColumns, column)
		}
		if column.IsList == 2 {
			ctx.ListColumns = append(ctx.ListColumns, column)
		}
		if column.IsInsert == 2 || column.IsEdit == 2 {
			ctx.FormColumns = append(ctx.FormColumns, column)
		}
		if !isSystemColumnName(column.ColumnName) {
			ctx.EditableColumns = append(ctx.EditableColumns, column)
		}
	}
	sort.SliceStable(ctx.ListColumns, func(i, j int) bool { return ctx.ListColumns[i].Sort > ctx.ListColumns[j].Sort })
	sort.SliceStable(ctx.FormColumns, func(i, j int) bool { return ctx.FormColumns[i].Sort > ctx.FormColumns[j].Sort })
	return ctx, nil
}

// columnsFromAny 把前端提交的动态字段配置转回强类型结构。
func columnsFromAny(value interface{}) []systemModel.ToolGenerateColumn {
	items, ok := value.([]interface{})
	if !ok {
		return nil
	}
	rows := make([]systemModel.ToolGenerateColumn, 0, len(items))
	for _, item := range items {
		m := mapValue(item)
		rows = append(rows, systemModel.ToolGenerateColumn{
			ID:            uintFromAny(m["id"]),
			ColumnName:    stringFromAny(m["column_name"]),
			ColumnComment: optionalString(m["column_comment"]),
			ColumnType:    optionalString(m["column_type"]),
			DefaultValue:  optionalString(m["default_value"]),
			IsPK:          int16FromAny(m["is_pk"], 1),
			IsRequired:    int16FromAny(m["is_required"], 1),
			IsInsert:      int16FromAny(m["is_insert"], 2),
			IsEdit:        int16FromAny(m["is_edit"], 2),
			IsList:        int16FromAny(m["is_list"], 1),
			IsQuery:       int16FromAny(m["is_query"], 1),
			IsSort:        int16FromAny(m["is_sort"], 1),
			QueryType:     defaultString(stringFromAny(m["query_type"]), "eq"),
			ViewType:      defaultString(stringFromAny(m["view_type"]), "input"),
			DictType:      optionalString(m["dict_type"]),
			AllowRoles:    optionalString(m["allow_roles"]),
			Sort:          uint8(uintFromAny(m["sort"])),
			Remark:        optionalString(m["remark"]),
		})
	}
	return rows
}

func codegenTableColumns() map[string]string {
	return map[string]string{"table_comment": "table_comment", "package_name": "package_name", "business_name": "business_name", "class_name": "class_name", "menu_name": "menu_name", "belong_menu_id": "belong_menu_id", "generate_path": "generate_path", "generate_model": "generate_model", "component_type": "component_type", "sort": "sort", "form_width": "form_width", "is_full": "is_full", "remark": "remark"}
}

// sourceOrDefault 数据源为空时回退到当前系统库。
func sourceOrDefault(source string) string {
	source = strings.TrimSpace(source)
	if source == "" {
		return config.MysqlAISystem
	}
	return source
}

// inferQueryType 按列类型推默认查询方式：文本类默认 LIKE，其余默认 =。
func inferQueryType(dataType string) string {
	switch strings.ToLower(dataType) {
	case "varchar", "char", "text", "tinytext", "mediumtext", "longtext":
		return "like"
	default:
		return "eq"
	}
}

// inferViewType 按列名和类型推默认页面组件：
// 长文本→textarea、日期→date、status/_type 结尾→字典下拉、数字→数字输入框。
func inferViewType(columnName, dataType string) string {
	switch strings.ToLower(dataType) {
	case "text", "tinytext", "mediumtext", "longtext":
		return "textarea"
	case "date", "datetime", "timestamp", "time", "year":
		return "date"
	}
	if strings.HasSuffix(columnName, "status") || strings.HasSuffix(columnName, "_type") {
		return "saSelect"
	}
	if isNumericBaseType(dataType) {
		return "inputNumber"
	}
	return "input"
}

func mapValue(value interface{}) map[string]interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		return m
	}
	return map[string]interface{}{}
}

func ptrString(value string) *string { return &value }
func optionalString(value interface{}) *string {
	s := stringFromAny(value)
	if s == "" {
		return nil
	}
	return &s
}
func stringFromAny(value interface{}) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}
func uintFromAny(value interface{}) uint {
	switch v := value.(type) {
	case float64:
		return uint(v)
	case int:
		return uint(v)
	case uint:
		return v
	case string:
		id, _ := parseUint(v)
		return id
	default:
		return 0
	}
}
func int16FromAny(value interface{}, fallback int16) int16 {
	if value == nil {
		return fallback
	}
	return int16(uintFromAny(value))
}
func boolFlag(ok bool) int16 {
	if ok {
		return 2
	}
	return 1
}
func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
func coalesceString(values ...interface{}) string {
	for _, value := range values {
		switch v := value.(type) {
		case *string:
			if v != nil && strings.TrimSpace(*v) != "" {
				return *v
			}
		case string:
			if strings.TrimSpace(v) != "" {
				return v
			}
		}
	}
	return ""
}
func defaultInt(value int, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}
func defaultInt16(value int16, fallback int16) int16 {
	if value == 0 {
		return fallback
	}
	return value
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func isRichOrUpload(viewType string) bool {
	return viewType == "wangEditor" || viewType == "uploadImage" || viewType == "uploadFile"
}
func isSystemColumnName(columnName string) bool {
	switch strings.ToLower(columnName) {
	case "id", "created_by", "updated_by", "create_time", "update_time", "delete_time":
		return true
	default:
		return false
	}
}
func isNumericBaseType(baseType string) bool {
	return regexp.MustCompile(`^(tinyint|smallint|mediumint|int|bigint|decimal|float|double)`).MatchString(strings.ToLower(baseType))
}
