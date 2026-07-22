package system

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	systemModel "server/model/system"
)

const (
	optionSourceStatic = "static"
	optionSourceRoute  = "route"
)

type codegenComponentCapability struct {
	Label         string `json:"label"`
	Value         string `json:"value"`
	ConfigType    string `json:"configType,omitempty"`
	QueryDisabled bool   `json:"queryDisabled,omitempty"`
}

type codegenOptionNode struct {
	Label    string              `json:"label"`
	Value    interface{}         `json:"value"`
	Children []codegenOptionNode `json:"children,omitempty"`
}

type codegenOptionConfig struct {
	Options       []codegenOptionNode    `json:"options,omitempty"`
	Path          string                 `json:"path,omitempty"`
	DataPath      string                 `json:"dataPath,omitempty"`
	LabelField    string                 `json:"labelField,omitempty"`
	ValueField    string                 `json:"valueField,omitempty"`
	ChildrenField string                 `json:"childrenField,omitempty"`
	Params        map[string]interface{} `json:"params,omitempty"`
}

type CodegenOptionRoute struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

var codegenComponentCapabilities = []codegenComponentCapability{
	{Label: "输入框", Value: "input"},
	{Label: "密码框", Value: "password"},
	{Label: "文本域", Value: "textarea"},
	{Label: "数字输入框", Value: "inputNumber"},
	{Label: "标签输入框", Value: "inputTag"},
	{Label: "开关", Value: "switch"},
	{Label: "滑块", Value: "slider"},
	{Label: "数据下拉框", Value: "select", ConfigType: "options"},
	{Label: "字典下拉框", Value: "saSelect", ConfigType: "dict"},
	{Label: "树形下拉框", Value: "treeSelect", ConfigType: "options"},
	{Label: "单选框", Value: "radio", ConfigType: "dict"},
	{Label: "复选框", Value: "checkbox", ConfigType: "dict"},
	{Label: "日期选择器", Value: "date"},
	{Label: "时间选择器", Value: "time"},
	{Label: "级联选择器", Value: "cascader", ConfigType: "options"},
	{Label: "用户选择器", Value: "userSelect"},
	{Label: "省市区联动", Value: "cityLinkage"},
	{Label: "图片上传", Value: "uploadImage", QueryDisabled: true},
	{Label: "文件上传", Value: "uploadFile", QueryDisabled: true},
	{Label: "富文本控件", Value: "wangEditor", QueryDisabled: true},
}

var (
	codegenRoutesMu sync.RWMutex
	codegenRoutes   = map[string]CodegenOptionRoute{}
)

func CodegenComponentCapabilities() []codegenComponentCapability {
	result := make([]codegenComponentCapability, len(codegenComponentCapabilities))
	copy(result, codegenComponentCapabilities)
	return result
}

func SetCodegenOptionRoutePaths(paths []string) {
	next := make(map[string]CodegenOptionRoute)
	for _, rawPath := range paths {
		index := strings.Index(rawPath, "/system/")
		if index < 0 {
			continue
		}
		path := rawPath[index:]
		if strings.ContainsAny(path, ":*") {
			continue
		}
		next[path] = CodegenOptionRoute{Label: "GET " + path, Path: path}
	}
	codegenRoutesMu.Lock()
	codegenRoutes = next
	codegenRoutesMu.Unlock()
}

func CodegenOptionRoutes() []CodegenOptionRoute {
	codegenRoutesMu.RLock()
	result := make([]CodegenOptionRoute, 0, len(codegenRoutes))
	for _, route := range codegenRoutes {
		result = append(result, route)
	}
	codegenRoutesMu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].Path < result[j].Path })
	return result
}

func isKnownCodegenOptionRoute(path string) bool {
	codegenRoutesMu.RLock()
	_, ok := codegenRoutes[path]
	codegenRoutesMu.RUnlock()
	return ok
}

func componentCapability(viewType string) (codegenComponentCapability, bool) {
	for _, capability := range codegenComponentCapabilities {
		if capability.Value == viewType {
			return capability, true
		}
	}
	return codegenComponentCapability{}, false
}

func parseCodegenOptionConfig(column systemModel.ToolGenerateColumn) (codegenOptionConfig, error) {
	var config codegenOptionConfig
	raw := stringValue(column.OptionConfig)
	if raw == "" {
		return config, fmt.Errorf("组件配置不能为空")
	}
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return config, fmt.Errorf("组件配置不是合法 JSON")
	}
	return config, nil
}

func validateCodegenColumns(columns []systemModel.ToolGenerateColumn) error {
	for _, column := range columns {
		capability, ok := componentCapability(column.ViewType)
		if !ok {
			return NewBizError(fmt.Sprintf("字段 %s 使用了不支持的页面组件 %s", column.ColumnName, column.ViewType))
		}
		active := column.IsInsert == 2 || column.IsEdit == 2 || column.IsQuery == 2
		if !active {
			continue
		}
		if column.IsQuery == 2 && capability.QueryDisabled {
			return NewBizError(fmt.Sprintf("字段 %s 的页面组件不支持查询", column.ColumnName))
		}
		switch capability.ConfigType {
		case "dict":
			if stringValue(column.DictType) == "" {
				return NewBizError(fmt.Sprintf("字段 %s 必须配置数据字典", column.ColumnName))
			}
		case "options":
			if err := validateCodegenOptionColumn(column); err != nil {
				return NewBizError(fmt.Sprintf("字段 %s: %v", column.ColumnName, err))
			}
		}
	}
	return nil
}

func validateCodegenOptionColumn(column systemModel.ToolGenerateColumn) error {
	source := stringValue(column.OptionSource)
	config, err := parseCodegenOptionConfig(column)
	if err != nil {
		return err
	}
	switch source {
	case optionSourceStatic:
		if len(config.Options) == 0 {
			return fmt.Errorf("静态选项不能为空")
		}
		return validateCodegenOptionNodes(config.Options, column.ViewType != "select")
	case optionSourceRoute:
		if !strings.HasPrefix(config.Path, "/system/") || strings.ContainsAny(config.Path, ":*") {
			return fmt.Errorf("请选择有效的系统 GET 路由")
		}
		if !isKnownCodegenOptionRoute(config.Path) {
			return fmt.Errorf("所选路由不存在或不可作为选项数据源")
		}
		for key, value := range config.Params {
			if strings.TrimSpace(key) == "" || !isScalarOptionParam(value) {
				return fmt.Errorf("路由参数只能使用非空名称和字符串、数字或布尔值")
			}
		}
		return nil
	default:
		return fmt.Errorf("请选择静态选项或系统路由")
	}
}

func validateCodegenOptionNodes(nodes []codegenOptionNode, allowChildren bool) error {
	for _, node := range nodes {
		if strings.TrimSpace(node.Label) == "" || node.Value == nil || fmt.Sprint(node.Value) == "" {
			return fmt.Errorf("每个静态选项都必须填写标签和值")
		}
		switch node.Value.(type) {
		case string, float64, bool:
		default:
			return fmt.Errorf("静态选项值只能是字符串、数字或布尔值")
		}
		if !allowChildren && len(node.Children) > 0 {
			return fmt.Errorf("数据下拉框不能配置子节点")
		}
		if err := validateCodegenOptionNodes(node.Children, allowChildren); err != nil {
			return err
		}
	}
	return nil
}

func isScalarOptionParam(value interface{}) bool {
	switch value.(type) {
	case string, float64, bool:
		return true
	default:
		return false
	}
}
