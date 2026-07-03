package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

func PostList(c *gin.Context) {
	result, err := systemService.PostList(queryMap(c))
	successOrFail(c, result, err)
}

func PostAccess(c *gin.Context) {
	result, err := systemService.PostAccess()
	successOrFail(c, result, err)
}

func CreatePost(c *gin.Context) {
	var payload systemRequest.PostPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.CreatePost(data)
	successOrFail(c, result, err)
}

func UpdatePost(c *gin.Context) {
	var payload systemRequest.PostPayload
	data, ok := bindJSONStructAsMap(c, &payload)
	if !ok {
		return
	}
	result, err := systemService.UpdatePost(c.Param("id"), data)
	successOrFail(c, result, err)
}

func DeletePost(c *gin.Context) {
	successOrFail(c, map[string]interface{}{}, systemService.DeletePost(c.Param("id")))
}
