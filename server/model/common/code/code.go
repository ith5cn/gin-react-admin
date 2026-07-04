package code

// 业务状态码（和 HTTP 状态码是两套体系）：
// 前端拦截器根据 code 决定行为——40101 触发自动刷新 token，
// 40102 跳登录页，其余非 0 值弹错误提示。
const (
	Success            = 0     // 成功。
	ParamError         = 40001 // 请求参数绑定或校验失败。
	OperationFailed    = 40002 // 业务规则拒绝（BizError），msg 可直接展示给用户。
	AccessTokenExpired = 40101 // access token 过期，前端应拿 refresh token 换新。
	LoginRequired      = 40102 // 未登录或 token 无效，需要重新登录。
	SystemError        = 50001 // 系统内部错误，细节只进服务端日志。
)

var messages = map[int]string{
	Success:            "操作成功",
	ParamError:         "参数错误",
	OperationFailed:    "操作失败",
	AccessTokenExpired: "access token expired, please refresh token",
	LoginRequired:      "please login again",
	SystemError:        "系统异常",
}

// Message 返回业务码的默认文案，未登记的码统一返回"操作失败"。
func Message(code int) string {
	if message, ok := messages[code]; ok {
		return message
	}
	return "操作失败"
}
