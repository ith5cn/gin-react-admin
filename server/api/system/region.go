package system

import (
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// RegionOptions 省市区级联选项，供前端 CityLinkage 组件使用。
func RegionOptions(c *gin.Context) {
	result, err := systemService.RegionOptions()
	successOrFail(c, result, err)
}
