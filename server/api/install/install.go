package install

import (
	"net/http"
	"server/model/common/code"
	"server/model/common/response"
	installModel "server/model/install"
	installService "server/service/install"

	"github.com/gin-gonic/gin"
)

func Status(c *gin.Context) {
	result, err := installService.Status()
	successOrFail(c, result, err)
}

func Check(c *gin.Context) {
	var data installModel.CheckRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	result, err := installService.Check(data)
	successOrFail(c, result, err)
}

func Run(c *gin.Context) {
	var data installModel.InstallRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Fail(c, code.ParamError, err.Error())
		return
	}
	result, err := installService.Run(data)
	successOrFail(c, result, err)
}

func successOrFail(c *gin.Context, data interface{}, err error) {
	if err != nil {
		response.FailWithHTTP(c, http.StatusInternalServerError, code.SystemError, err.Error())
		return
	}
	response.Success(c, data)
}
