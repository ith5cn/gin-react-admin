package system

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCodegenPackageName(t *testing.T) {
	valid := []string{"system", "ai_tools", "module2"}
	for _, value := range valid {
		if err := validateCodegenPackageName(value); err != nil {
			t.Fatalf("合法包名 %q 被拒绝: %v", value, err)
		}
	}

	invalid := []string{"", " system ", "ai-tools", "../system", "system/path", "2module", "System"}
	for _, value := range invalid {
		if err := validateCodegenPackageName(value); !errors.Is(err, ErrInvalidPackageName) {
			t.Fatalf("非法包名 %q 未返回 ErrInvalidPackageName: %v", value, err)
		}
	}
}

func TestRenderGoFilesUseConfiguredPackage(t *testing.T) {
	systemCtx := fullViewTypeContext(1, 1, true)
	systemService := renderGoService(systemCtx)
	systemAPI := renderGoAPI(systemCtx)
	if strings.Contains(systemService, `"server/service/system"`) {
		t.Fatalf("system service 不应导入自身:\n%s", systemService)
	}
	if strings.Contains(systemAPI, `"server/api/system"`) {
		t.Fatalf("system api 不应导入自身:\n%s", systemAPI)
	}
	for name, content := range map[string]string{
		"model.go":   renderGoModel(systemCtx),
		"service.go": systemService,
		"api.go":     systemAPI,
		"router.go":  renderGoRouter(systemCtx),
	} {
		if !strings.Contains(content, "package system") {
			t.Fatalf("%s 未使用 system 包声明:\n%s", name, content)
		}
		mustParseGo(t, name, content)
	}

	customCtx := fullViewTypeContext(2, 2, false)
	customCtx.PackageName = "reports"
	customCtx.BackendModelPath = filepath.Join("model", "reports")
	customCtx.BackendServicePath = filepath.Join("service", "reports")
	customCtx.BackendAPIPath = filepath.Join("api", "reports")
	customCtx.BackendRoutePath = filepath.Join("router", "reports")
	customCtx.RoutePath = "reports/" + customCtx.BusinessName

	checks := map[string][]string{
		renderGoModel(customCtx):   {"package reports"},
		renderGoService(customCtx): {"package reports", `"server/model/reports"`, `"server/service/system"`},
		renderGoAPI(customCtx):     {"package reports", `"server/api/system"`, `"server/service/reports"`},
		renderGoRouter(customCtx):  {"package reports", `"server/api/reports"`},
	}
	for content, expected := range checks {
		for _, item := range expected {
			if !strings.Contains(content, item) {
				t.Fatalf("自定义包生成内容缺少 %q:\n%s", item, content)
			}
		}
		mustParseGo(t, "custom.go", content)
	}

	files := buildGoPreviewFiles(customCtx)
	for _, file := range files {
		if !strings.Contains(filepath.ToSlash(file.Path), "/reports/") {
			t.Fatalf("后端预览路径未使用包名: %s", file.Path)
		}
	}
}

func TestRefreshGeneratedRouteRegistryAcrossPackages(t *testing.T) {
	root := t.TempDir()
	for _, packageName := range []string{"generated", "reports", "system"} {
		if err := os.MkdirAll(filepath.Join(root, packageName), 0755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(path, content string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(root, "generated", "legacy.go"), codegenFileMarker+"\npackage generated\nfunc RegisterLegacyRoutes() {}\n")
	write(filepath.Join(root, "reports", "report.go"), codegenFileMarker+"\npackage reports\nfunc RegisterReportRoutes() {}\n")
	write(filepath.Join(root, "system", "article.go"), codegenFileMarker+"\npackage system\nfunc RegisterArticleRoutes() {}\n")
	write(filepath.Join(root, "system", "manual.go"), "package system\nfunc RegisterIgnoredRoutes() {}\n")

	if err := refreshGeneratedRouteRegistryAt(root); err != nil {
		t.Fatal(err)
	}
	contentBytes, err := os.ReadFile(filepath.Join(root, "generated", "register.go"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(contentBytes)
	for _, expected := range []string{
		`reportsRouter "server/router/reports"`,
		`systemRouter "server/router/system"`,
		"RegisterLegacyRoutes(group)",
		"reportsRouter.RegisterReportRoutes(group)",
		"systemRouter.RegisterArticleRoutes(group)",
	} {
		if !strings.Contains(content, expected) {
			t.Fatalf("路由注册表缺少 %q:\n%s", expected, content)
		}
	}
	if strings.Contains(content, "RegisterIgnoredRoutes") {
		t.Fatalf("手写路由不应进入生成注册表:\n%s", content)
	}
	if strings.Index(content, "RegisterLegacyRoutes(group)") > strings.Index(content, "reportsRouter.RegisterReportRoutes(group)") {
		t.Fatalf("路由注册顺序不稳定:\n%s", content)
	}
	mustParseGo(t, "register.go", content)
}

func TestCleanupLegacyGeneratedFilesOnlyRemovesGeneratedFiles(t *testing.T) {
	root := t.TempDir()
	ctx := fullViewTypeContext(1, 1, true)
	legacyPaths := []string{
		filepath.Join("model", "generated", ctx.BusinessName+".go"),
		filepath.Join("service", "generated", ctx.BusinessName+".go"),
		filepath.Join("api", "generated", ctx.BusinessName+".go"),
		filepath.Join("router", "generated", ctx.BusinessName+".go"),
	}
	for index, relativePath := range legacyPaths {
		path := filepath.Join(root, relativePath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		content := codegenFileMarker + "\npackage generated\n"
		if index == len(legacyPaths)-1 {
			content = "package generated\n"
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := cleanupLegacyGeneratedFiles(root, ctx); err != nil {
		t.Fatal(err)
	}
	for _, relativePath := range legacyPaths[:3] {
		if _, err := os.Stat(filepath.Join(root, relativePath)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("生成文件未删除: %s", relativePath)
		}
	}
	if _, err := os.Stat(filepath.Join(root, legacyPaths[3])); err != nil {
		t.Fatalf("手写同名文件不应删除: %v", err)
	}
}
