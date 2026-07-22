package system

import (
	"strings"
	"testing"

	systemModel "server/model/system"
)

func TestCodegenComponentCapabilities(t *testing.T) {
	capabilities := CodegenComponentCapabilities()
	if len(capabilities) != 20 {
		t.Fatalf("页面组件数量 = %d, want 20", len(capabilities))
	}
	seen := make(map[string]bool, len(capabilities))
	for _, capability := range capabilities {
		if seen[capability.Value] {
			t.Fatalf("页面组件重复: %s", capability.Value)
		}
		seen[capability.Value] = true
	}
	for _, viewType := range []string{"select", "treeSelect", "cascader", "saSelect", "radio", "checkbox"} {
		if !seen[viewType] {
			t.Fatalf("页面组件缺少 %s", viewType)
		}
	}
}

func TestSetCodegenOptionRoutePaths(t *testing.T) {
	SetCodegenOptionRoutePaths([]string{"/api/v1/system/dept/index", "/api/system/dept/index", "/api/v1/system/user/:id", "/api/v1/system/files/*path", "/api/v1/base/login"})
	routes := CodegenOptionRoutes()
	if len(routes) != 1 || routes[0].Path != "/system/dept/index" {
		t.Fatalf("路由目录过滤结果异常: %#v", routes)
	}
}

func TestValidateCodegenColumns(t *testing.T) {
	SetCodegenOptionRoutePaths([]string{"/api/v1/system/dept/index"})
	validStatic := systemModel.ToolGenerateColumn{ColumnName: "category_id", ColumnType: ptrStr("int"), IsInsert: 2, ViewType: "select", OptionSource: ptrStr("static"), OptionConfig: ptrStr("{\"options\":[{\"label\":\"分类一\",\"value\":1}]}")}
	validRoute := systemModel.ToolGenerateColumn{ColumnName: "dept_id", ColumnType: ptrStr("int"), IsQuery: 2, ViewType: "treeSelect", OptionSource: ptrStr("route"), OptionConfig: ptrStr("{\"path\":\"/system/dept/index\",\"dataPath\":\"list\",\"params\":{\"status\":1}}")}
	validDict := systemModel.ToolGenerateColumn{ColumnName: "status", IsEdit: 2, ViewType: "saSelect", DictType: ptrStr("status")}
	if err := validateCodegenColumns([]systemModel.ToolGenerateColumn{validStatic, validRoute, validDict}); err != nil {
		t.Fatalf("合法组件配置校验失败: %v", err)
	}

	cases := []struct {
		name   string
		column systemModel.ToolGenerateColumn
		want   string
	}{
		{"空字典", systemModel.ToolGenerateColumn{ColumnName: "status", IsEdit: 2, ViewType: "radio"}, "必须配置数据字典"},
		{"空选项", systemModel.ToolGenerateColumn{ColumnName: "kind", IsEdit: 2, ViewType: "select", OptionSource: ptrStr("static"), OptionConfig: ptrStr("{\"options\":[]}")}, "静态选项不能为空"},
		{"非法路由", systemModel.ToolGenerateColumn{ColumnName: "dept", IsEdit: 2, ViewType: "treeSelect", OptionSource: ptrStr("route"), OptionConfig: ptrStr("{\"path\":\"/system/dept/:id\"}")}, "有效的系统 GET 路由"},
		{"查询禁用", systemModel.ToolGenerateColumn{ColumnName: "content", IsQuery: 2, ViewType: "wangEditor"}, "不支持查询"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCodegenColumns([]systemModel.ToolGenerateColumn{tc.column})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("校验错误 = %v, want contain %q", err, tc.want)
			}
		})
	}
}

func TestInferViewTypeTimeAndYear(t *testing.T) {
	if got := inferViewType("remind_at", "time"); got != "time" {
		t.Fatalf("time 推断 = %s, want time", got)
	}
	if got := inferViewType("report_year", "year"); got != "date" {
		t.Fatalf("year 推断 = %s, want date", got)
	}
}

func TestValidateCodegenColumnsSkipsUnusedDataSource(t *testing.T) {
	unused := []systemModel.ToolGenerateColumn{
		{ColumnName: "views", IsList: 2, ViewType: "select"},
		{ColumnName: "status", IsList: 2, ViewType: "saSelect"},
	}
	if err := validateCodegenColumns(unused); err != nil {
		t.Fatalf("纯列表字段不应校验组件数据源: %v", err)
	}

	queryColumn := systemModel.ToolGenerateColumn{ColumnName: "views", IsQuery: 2, ViewType: "select"}
	if err := validateCodegenColumns([]systemModel.ToolGenerateColumn{queryColumn}); err == nil || !strings.Contains(err.Error(), "组件配置不能为空") {
		t.Fatalf("查询字段缺少数据源时应拒绝，实际错误: %v", err)
	}
}
