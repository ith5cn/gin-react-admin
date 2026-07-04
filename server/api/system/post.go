package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// PostList 岗位分页列表。
func PostList(c *gin.Context) {
	result, err := systemService.PostList(queryMap(c))
	successOrFail(c, result, err)
}

// PostAccess 岗位下拉数据。
func PostAccess(c *gin.Context) {
	result, err := systemService.PostAccess()
	successOrFail(c, result, err)
}

// CreatePost 新增岗位。
func CreatePost(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.PostPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CreatePost(payload)
	successOrFail(c, result, err)
}

// UpdatePost 更新岗位。
func UpdatePost(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.PostPayload](c)
	if !ok {
		return
	}
	result, err := systemService.UpdatePost(c.Param("id"), payload)
	successOrFail(c, result, err)
}

// DeletePost 删除岗位。
func DeletePost(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeletePost(c.Param("id")))
}
