package system

import (
	"net/http"
	"time"

	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

const xlsxContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

// ExportUsers 导出用户列表 Excel（与列表接口同一套过滤条件和数据权限）。
func ExportUsers(c *gin.Context) {
	operatorID, ok := currentUserID(c)
	if !ok {
		return
	}
	data, err := systemService.ExportUsersExcel(operatorID, queryMap(c))
	if err != nil {
		successOrFail(c, nil, err)
		return
	}
	sendXlsx(c, "users_"+time.Now().Format("20060102150405")+".xlsx", data)
}

// UserImportTemplate 下载用户导入模板 Excel。
func UserImportTemplate(c *gin.Context) {
	data, err := systemService.UserImportTemplateExcel()
	if err != nil {
		successOrFail(c, nil, err)
		return
	}
	sendXlsx(c, "user_import_template.xlsx", data)
}

// ImportUsers 从上传的 Excel 导入用户。
func ImportUsers(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		successOrFail(c, nil, systemService.ErrImportNotExcel)
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		successOrFail(c, nil, err)
		return
	}
	defer file.Close()

	result, importErr := systemService.ImportUsersExcel(file)
	successOrFail(c, result, importErr)
}

// sendXlsx 以附件形式返回 Excel 文件流。
func sendXlsx(c *gin.Context, filename string, data []byte) {
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Data(http.StatusOK, xlsxContentType, data)
}
