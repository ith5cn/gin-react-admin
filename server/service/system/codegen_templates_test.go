package system

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	systemModel "server/model/system"
)

func ptrStr(v string) *string { return &v }

// testCodegenColumn 构造一个测试字段。
func testCodegenColumn(name, columnType, viewType, queryType, dictType string, required, isQuery, isSort int16) systemModel.ToolGenerateColumn {
	column := systemModel.ToolGenerateColumn{
		ColumnName:    name,
		ColumnComment: ptrStr(name + "描述"),
		ColumnType:    ptrStr(columnType),
		IsPK:          1,
		IsRequired:    required,
		IsInsert:      2,
		IsEdit:        2,
		IsList:        2,
		IsQuery:       isQuery,
		IsSort:        isSort,
		QueryType:     queryType,
		ViewType:      viewType,
	}
	if dictType != "" {
		column.DictType = ptrStr(dictType)
	}
	if viewType == "select" || viewType == "treeSelect" || viewType == "cascader" {
		column.OptionSource = ptrStr(optionSourceStatic)
		column.OptionConfig = ptrStr(`{"options":[{"label":"选项一","value":1,"children":[{"label":"子选项","value":2}]}]}`)
	}
	return column
}

// fullViewTypeContext 覆盖全部 20 种页面组件和主要查询操作符。
func fullViewTypeContext(generateModel int16, componentType int16, withDeleteTime bool) codegenContext {
	columns := []systemModel.ToolGenerateColumn{
		{ColumnName: "id", ColumnType: ptrStr("bigint(20) unsigned"), IsPK: 2, IsRequired: 2, IsList: 2, QueryType: "eq", ViewType: "input"},
		testCodegenColumn("title", "varchar(255)", "input", "like", "", 2, 2, 1),
		testCodegenColumn("secret", "varchar(255)", "password", "eq", "", 1, 1, 1),
		testCodegenColumn("summary", "text", "textarea", "like", "", 1, 1, 1),
		testCodegenColumn("price", "decimal(10,2)", "inputNumber", "gte", "", 1, 2, 1),
		testCodegenColumn("tags", "varchar(500)", "inputTag", "in", "", 1, 2, 1),
		testCodegenColumn("enabled", "tinyint(1)", "switch", "eq", "", 2, 2, 1),
		testCodegenColumn("progress", "int(11)", "slider", "between", "", 1, 2, 1),
		testCodegenColumn("category_id", "int(11)", "select", "eq", "", 1, 1, 1),
		testCodegenColumn("status", "tinyint(4)", "saSelect", "eq", "status", 2, 2, 1),
		testCodegenColumn("dept_id", "int(11)", "treeSelect", "eq", "", 1, 1, 1),
		testCodegenColumn("level", "tinyint(4)", "radio", "eq", "level", 1, 1, 1),
		testCodegenColumn("channels", "varchar(200)", "checkbox", "notin", "channel", 1, 1, 1),
		testCodegenColumn("publish_time", "datetime", "date", "gt", "", 1, 2, 1),
		testCodegenColumn("birthday", "date", "date", "eq", "", 1, 1, 1),
		testCodegenColumn("remind_at", "time", "time", "eq", "", 1, 1, 1),
		testCodegenColumn("region_path", "varchar(200)", "cascader", "eq", "", 1, 1, 1),
		testCodegenColumn("owner_id", "int(11)", "userSelect", "eq", "", 1, 1, 1),
		testCodegenColumn("city", "varchar(200)", "cityLinkage", "eq", "", 1, 1, 1),
		testCodegenColumn("cover", "varchar(500)", "uploadImage", "eq", "", 1, 1, 1),
		testCodegenColumn("attachment", "varchar(500)", "uploadFile", "eq", "", 1, 1, 1),
		testCodegenColumn("content", "longtext", "wangEditor", "eq", "", 1, 1, 1),
		testCodegenColumn("sort", "smallint(6)", "inputNumber", "eq", "", 1, 1, 2),
	}
	if withDeleteTime {
		columns = append(columns, systemModel.ToolGenerateColumn{ColumnName: "delete_time", ColumnType: ptrStr("datetime"), IsRequired: 1, QueryType: "eq", ViewType: "date"})
	}

	table := systemModel.ToolGenerateTable{
		ID:             1,
		TableNameValue: "demo_article",
		TableComment:   ptrStr("演示文章"),
		PackageName:    ptrStr("system"),
		BusinessName:   ptrStr("demo-article"),
		ClassName:      ptrStr("DemoArticle"),
		MenuName:       ptrStr("演示文章"),
		GenerateModel:  generateModel,
		ComponentType:  componentType,
		FormWidth:      800,
	}

	ctx := codegenContext{
		Table:              table,
		Columns:            columns,
		PackageName:        "system",
		BusinessName:       "demo-article",
		ClassName:          "DemoArticle",
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
	}
	return ctx
}

// mustParseGo 断言生成的 Go 源码语法合法。
func mustParseGo(t *testing.T, name, content string) {
	t.Helper()
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, name, content, parser.AllErrors); err != nil {
		t.Fatalf("生成的 %s 语法错误: %v\n---\n%s", name, err, content)
	}
}

func TestRenderGoFilesParse(t *testing.T) {
	for _, tc := range []struct {
		name string
		ctx  codegenContext
	}{
		{"软删除表", fullViewTypeContext(1, 1, true)},
		{"非软删除表", fullViewTypeContext(2, 1, false)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mustParseGo(t, "model.go", renderGoModel(tc.ctx))
			mustParseGo(t, "service.go", renderGoService(tc.ctx))
			mustParseGo(t, "api.go", renderGoAPI(tc.ctx))
			mustParseGo(t, "router.go", renderGoRouter(tc.ctx))
		})
	}
}

func TestRenderGoModelNoTimeImport(t *testing.T) {
	ctx := fullViewTypeContext(2, 1, false)
	// 只保留非日期字段，确认不再产出未使用的 time 导入。
	ctx.Columns = ctx.Columns[:2]
	content := renderGoModel(ctx)
	if strings.Contains(content, `import "time"`) {
		t.Fatalf("无日期字段时不应导入 time:\n%s", content)
	}
	mustParseGo(t, "model.go", content)
}

func TestRenderGoModelNullableTypes(t *testing.T) {
	ctx := fullViewTypeContext(1, 1, true)
	content := renderGoModel(ctx)
	for _, expect := range []string{
		"Price *float64", // 可空 decimal → 指针
		"Progress *int",  // 可空 int → 指针
		"Title string",   // 非空 varchar → 值类型
		"Enabled int",    // 非空 tinyint → 值类型
		"`json:\"-\" gorm:\"column:delete_time\"`", // 软删列不输出
	} {
		if !strings.Contains(content, expect) {
			t.Fatalf("model 缺少 %q:\n%s", expect, content)
		}
	}
}

func TestRenderGoServiceSoftDelete(t *testing.T) {
	soft := renderGoService(fullViewTypeContext(1, 1, true))
	if !strings.Contains(soft, `systemService.SoftDeleteRecord("demo_article", id)`) {
		t.Fatalf("软删除表应使用 SoftDeleteRecord:\n%s", soft)
	}
	if !strings.Contains(soft, ", true)") {
		t.Fatalf("软删除表 List 应开启软删过滤:\n%s", soft)
	}

	hard := renderGoService(fullViewTypeContext(2, 1, false))
	if !strings.Contains(hard, "systemService.DeleteRecord(&generatedModel.DemoArticle{}, id)") {
		t.Fatalf("非软删除表应使用 DeleteRecord:\n%s", hard)
	}
	if !strings.Contains(hard, ", false)") {
		t.Fatalf("非软删除表 List 不应过滤 delete_time:\n%s", hard)
	}

	// 选了软删除模型但表里没有 delete_time 列时必须回退硬删，避免 SQL 报错。
	mismatched := renderGoService(fullViewTypeContext(1, 1, false))
	if !strings.Contains(mismatched, ", false)") {
		t.Fatalf("无 delete_time 列时应回退非软删查询:\n%s", mismatched)
	}
}

func TestRenderGoServiceQueryOperators(t *testing.T) {
	content := renderGoService(fullViewTypeContext(1, 1, true))
	for _, expect := range []string{
		`{Param: "title", Column: "title", Op: "like"},`,
		`{Param: "price", Column: "price", Op: "gte"},`,
		`{Param: "progress", Column: "progress", Op: "between"},`,
		`{Param: "tags", Column: "tags", Op: "in"},`,
		`"sort": "sort"`,
	} {
		if !strings.Contains(content, expect) {
			t.Fatalf("service 缺少 %q:\n%s", expect, content)
		}
	}
}

func TestRenderFrontendEditComponents(t *testing.T) {
	content := renderFrontendEdit(fullViewTypeContext(1, 1, true))
	for _, expect := range []string{
		`<Input.Password`, `<Input.TextArea`, `<InputNumber`, `<Select mode="tags"`,
		`<Switch />`, `<Slider />`, `<Select allowClear`, `<Ith5Select dict="status"`, `<TreeSelect`,
		`<Ith5Radio dict="level"`, `<Ith5Checkbox dict="channel"`,
		`<DatePicker showTime`, `<TimePicker`, `<Cascader`, `<AuthSelect />`,
		`<CityLinkage />`, `<ImageUpload />`, `<FileUpload />`, `<WangEditor />`,
		`import Ith5Select from "@/components/ith5ui/ith5-select";`,
		`import dayjs from "dayjs";`,
		`{ name: "publishTime", format: "YYYY-MM-DD HH:mm:ss" },`,
		`{ name: "birthday", format: "YYYY-MM-DD" },`,
		`{ name: "remindAt", format: "HH:mm:ss" },`,
		`{ name: "channels", numeric: false },`,
		`rules={[{ required: true, message: "请输入title描述" }]}`,
		`getValueProps={(value) => ({ checked: value === 1 })}`,
		"catch (error) {",
	} {
		if !strings.Contains(content, expect) {
			t.Fatalf("edit 模板缺少 %q:\n%s", expect, content)
		}
	}
	if strings.Contains(content, ": any") {
		t.Fatalf("edit 模板不应出现 any:\n%s", content)
	}
}

func TestRenderFrontendEditPreservesPercentText(t *testing.T) {
	ctx := fullViewTypeContext(1, 1, true)
	ctx.FormColumns[0].ColumnComment = ptrStr("完成率 100%")
	content := renderFrontendEdit(ctx)
	if !strings.Contains(content, "完成率 100%") {
		t.Fatalf("百分号文本未被完整保留:\n%s", content)
	}
	for _, invalid := range []string{"%!", "(string="} {
		if strings.Contains(content, invalid) {
			t.Fatalf("生成内容包含格式化异常 %q:\n%s", invalid, content)
		}
	}
}
func TestRenderFrontendEditDrawer(t *testing.T) {
	content := renderFrontendEdit(fullViewTypeContext(1, 2, true))
	for _, expect := range []string{"<Drawer", "onClose={close}", "width={800}"} {
		if !strings.Contains(content, expect) {
			t.Fatalf("抽屉模式缺少 %q:\n%s", expect, content)
		}
	}

	modal := renderFrontendEdit(fullViewTypeContext(1, 1, true))
	if !strings.Contains(modal, "<Modal") || !strings.Contains(modal, "width={800}") {
		t.Fatalf("模态框模式渲染异常:\n%s", modal)
	}
}

func TestRenderFrontendIndexDictColumns(t *testing.T) {
	content := renderFrontendIndex(fullViewTypeContext(1, 1, true))
	for _, expect := range []string{
		`type: "dict", dict: "status"`,
		`<Ith5Select dict="status"`,
		`(record: { id: number })`,
	} {
		if !strings.Contains(content, expect) {
			t.Fatalf("index 模板缺少 %q:\n%s", expect, content)
		}
	}
	if strings.Contains(content, ": any") {
		t.Fatalf("index 模板不应出现 any:\n%s", content)
	}
}

// TestWriteFrontendSample 把生成的前端文件写到临时目录，
// 供 CI/本地用真实 tsc 编译验证（TSX 无法在 Go 里做语法校验）。
func TestWriteFrontendSample(t *testing.T) {
	outDir := os.Getenv("CODEGEN_SAMPLE_DIR")
	if outDir == "" {
		t.Skip("设置 CODEGEN_SAMPLE_DIR 后写出前端样例")
	}
	ctx := fullViewTypeContext(1, 1, true)
	files := map[string]string{
		"index.tsx": renderFrontendIndex(ctx),
		"edit.tsx":  renderFrontendEdit(ctx),
		"api.ts":    renderFrontendAPI(ctx),
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(outDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRenderFrontendSearchComponents(t *testing.T) {
	ctx := fullViewTypeContext(1, 1, true)
	ctx.QueryColumns = nil
	for _, column := range ctx.Columns {
		switch column.ViewType {
		case "select":
			column.QueryType = "in"
		case "treeSelect", "cascader", "saSelect", "userSelect", "cityLinkage":
			column.QueryType = "eq"
		case "date", "time", "inputNumber":
			column.QueryType = "between"
		default:
			continue
		}
		ctx.QueryColumns = append(ctx.QueryColumns, column)
	}
	content := renderFrontendIndex(ctx)
	for _, expect := range []string{
		"<Select allowClear mode=\"multiple\" options=",
		"<TreeSelect allowClear",
		"<Cascader allowClear options=",
		"<Ith5Select dict=\"status\" valueType=\"number\"",
		"<DatePicker.RangePicker",
		"<TimePicker.RangePicker",
		"<QueryRange numeric />",
		"<AuthSelect />",
		"<CityLinkage />",
		"const toSearchParams =",
		".join(\",\")",
	} {
		if !strings.Contains(content, expect) {
			t.Fatalf("search 模板缺少 %q:\n%s", expect, content)
		}
	}
	for _, invalid := range []string{"options={[]}", "treeData={[]}", "%!", "(string="} {
		if strings.Contains(content, invalid) {
			t.Fatalf("search 模板包含非法占位 %q:\n%s", invalid, content)
		}
	}
}
