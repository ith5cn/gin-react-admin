package request

// RefreshTokenRequest 是刷新 token 接口入参。
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest 是退出登录接口入参。
// refresh_token 可选，传入时会和当前 access token 一起撤销。
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}
