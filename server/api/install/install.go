package install

import (
	"net/http"
	"server/model/common/code"
	"server/model/common/response"
	installModel "server/model/install"
	installService "server/service/install"

	"github.com/gin-gonic/gin"
)

// Status 返回安装状态和可执行的 SQL 文件列表。
func Status(c *gin.Context) {
	result, err := installService.Status()
	successOrFail(c, result, err)
}

// Check 检测 MySQL / Redis 连接配置是否可用。
func Check(c *gin.Context) {
	var data installModel.CheckRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	result, err := installService.Check(data)
	successOrFail(c, result, err)
}

// Run 执行安装：写 .env、执行 SQL、生成安装锁文件。
func Run(c *gin.Context) {
	var data installModel.InstallRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	result, err := installService.Run(data)
	successOrFail(c, result, err)
}

// successOrFail 是 install 模块的简化版响应封装。
// 安装向导的使用者就是运维/管理员，错误细节（如数据库连不上）对排障是必要的，
// 所以这里保留 err.Error() 透传，与业务接口的收敛策略不同。
func successOrFail(c *gin.Context, data interface{}, err error) {
	if err != nil {
		response.FailWithHTTP(c, http.StatusInternalServerError, code.SystemError, err.Error())
		return
	}
	response.Success(c, data)
}
