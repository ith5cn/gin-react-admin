package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// NoticeList 通知公告分页列表。
func NoticeList(c *gin.Context) {
	data, err := systemService.NoticeList(queryMap(c))
	successOrFail(c, data, err)
}

// CreateNotice 新增通知公告。
func CreateNotice(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.NoticePayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreateNotice(payload)
	successOrFail(c, result, err)
}

// UpdateNotice 更新通知公告。
func UpdateNotice(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.NoticePayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdateNotice(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeleteNotice 删除通知公告。
func DeleteNotice(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeleteNotice(c.Param("id")))
}
