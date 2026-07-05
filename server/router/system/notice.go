package system

import (
	systemApi "server/api/system"
	"server/middleware"

	"github.com/gin-gonic/gin"
)

// noticeRoutes 注册通知公告路由。
func noticeRoutes(system *gin.RouterGroup) {
	system.GET("/notice/index", middleware.Perm("system/notice/index"), systemApi.NoticeList)
	system.POST("/notice", middleware.Perm("system/notice/create"), systemApi.CreateNotice)
	system.PUT("/notice/:id", middleware.Perm("system/notice/update"), systemApi.UpdateNotice)
	system.DELETE("/notice/:id", middleware.Perm("system/notice/destroy"), systemApi.DeleteNotice)
}
