package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

func deptRoutes(system *gin.RouterGroup) {
	system.GET("/dept/index", systemApi.DeptList)
	system.GET("/dept/access", systemApi.DeptAccess)
	system.POST("/dept", systemApi.CreateDept)
	system.PUT("/dept/:id", systemApi.UpdateDept)
	system.DELETE("/dept/:id", systemApi.DeleteDept)
}
