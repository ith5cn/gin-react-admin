package code

const (
	Success            = 0
	ParamError         = 40001
	OperationFailed    = 40002
	AccessTokenExpired = 40101
	LoginRequired      = 40102
	SystemError        = 50001
)

var messages = map[int]string{
	Success:            "操作成功",
	ParamError:         "参数错误",
	OperationFailed:    "操作失败",
	AccessTokenExpired: "access token expired, please refresh token",
	LoginRequired:      "please login again",
	SystemError:        "系统异常",
}

func Message(code int) string {
	if message, ok := messages[code]; ok {
		return message
	}
	return "操作失败"
}
