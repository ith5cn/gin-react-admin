package request

// LoginRequest 是系统登录接口入参。
// binding:"required" 会让 Gin 在字段缺失时返回绑定错误。
type LoginRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserPayload 是用户创建/更新入参。
// Password 只在创建时生效，更新走独立的 set-password 接口；
// Roles 为 nil 表示未提交，不改动角色绑定。
type UserPayload struct {
	Username       *string `json:"username"`
	UserType       *string `json:"userType"`
	Nickname       *string `json:"nickname"`
	Phone          *string `json:"phone"`
	Email          *string `json:"email"`
	Avatar         *string `json:"avatar"`
	Signed         *string `json:"signed"`
	Dashboard      *string `json:"dashboard"`
	DeptID         *uint   `json:"deptId"`
	Status         *int16  `json:"status"`
	Remark         *string `json:"remark"`
	BackendSetting *string `json:"backendSetting"`
	Password       *string `json:"password"`
	Roles          []uint  `json:"roles"`
}

// SetPasswordPayload 是管理员重置用户密码入参。
type SetPasswordPayload struct {
	Password string `json:"password" binding:"required"`
}
