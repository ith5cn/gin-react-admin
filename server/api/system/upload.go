package system

import (
	"mime/multipart"
	"server/model/common/code"
	"server/model/common/response"
	systemModel "server/model/system"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// UploadImage 上传图片（头像、富文本插图、生成器的图片控件共用）。
// 前端以 multipart/form-data 提交，文件字段名固定为 file。
func UploadImage(c *gin.Context) {
	handleUpload(c, systemService.SaveUploadedImage)
}

// UploadFile 上传普通文件。
func UploadFile(c *gin.Context) {
	handleUpload(c, systemService.SaveUploadedFile)
}

func handleUpload(c *gin.Context, save func(*multipart.FileHeader, uint) (*systemModel.AISystemAttachment, error)) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, code.ParamError, "缺少上传文件字段 file")
		return
	}

	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)

	attachment, err := save(fileHeader, uid)
	if err != nil {
		successOrFail(c, nil, err)
		return
	}

	// 前端上传组件约定读取 path 和 url 两个字段。
	response.Success(c, map[string]interface{}{
		"id":   attachment.ID,
		"path": stringOrEmpty(attachment.StoragePath),
		"url":  stringOrEmpty(attachment.URL),
	})
}

func stringOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
