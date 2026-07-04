package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// deptRoutes 注册部门管理路由。
func deptRoutes(system *gin.RouterGroup) {
	system.GET("/dept/index", middleware.Perm("system/dept/index"), systemApi.DeptList)
	system.GET("/dept/access", systemApi.DeptAccess)
	system.POST("/dept", middleware.Perm("system/dept/create"), systemApi.CreateDept)
	system.PUT("/dept/:id", middleware.Perm("system/dept/update"), systemApi.UpdateDept)
	system.DELETE("/dept/:id", middleware.Perm("system/dept/destroy"), systemApi.DeleteDept)
}
