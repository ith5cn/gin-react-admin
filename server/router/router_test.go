package router

import (
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
