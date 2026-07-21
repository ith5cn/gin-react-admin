package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/model/common/code"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestNewRouterRegistersWithoutConflict 验证全部路由能注册成功。
// Gin 在路由冲突（如同一层级静态段与参数段不兼容）时会在注册阶段直接 panic，
// 这类错误编译期发现不了，靠这个测试兜底。
func TestNewRouterRegistersWithoutConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration panicked: %v", r)
		}
	}()
	if engine := NewRouter(); engine == nil {
		t.Fatal("NewRouter returned nil")
	}
}

func TestInstallGuardRedirectSignal(t *testing.T) {
	t.Chdir(t.TempDir())
	gin.SetMode(gin.TestMode)
	engine := NewRouter()

	t.Run("business API returns not-installed code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/system/user", nil)
		engine.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
		}

		var payload struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Code != code.SystemNotInstalled {
			t.Fatalf("business code = %d, want %d", payload.Code, code.SystemNotInstalled)
		}
		if payload.Msg != code.Message(code.SystemNotInstalled) {
			t.Fatalf("message = %q, want %q", payload.Msg, code.Message(code.SystemNotInstalled))
		}
	})

	t.Run("install status remains public", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/install/status", nil)
		engine.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}

		var payload struct {
			Code int `json:"code"`
			Data struct {
				Installed bool `json:"installed"`
			} `json:"data"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if payload.Code != code.Success {
			t.Fatalf("business code = %d, want %d", payload.Code, code.Success)
		}
		if payload.Data.Installed {
			t.Fatal("installed = true, want false")
		}
	})
}
