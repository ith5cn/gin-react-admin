package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// AttachmentList 返回附件分页列表。
// api 层只接收查询参数，具体筛选和分页逻辑交给 service。
func AttachmentList(c *gin.Context) {
	result, err := systemService.AttachmentList(queryMap(c))
	successOrFail(c, result, err)
}

// DeleteAttachments 批量软删除附件记录。
// 请求体格式兼容前端：{ ids: [], removeSource?: boolean }。
func DeleteAttachments(c *gin.Context) {
	var payload systemRequest.AttachmentDeletePayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	successOrFail(c, true, systemService.DeleteAttachments(data))
}
