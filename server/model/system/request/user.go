package request

// LoginRequest 是系统登录接口入参。
// binding:"required" 会让 Gin 在字段缺失时返回绑定错误。
type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
