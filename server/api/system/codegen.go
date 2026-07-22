package system

import (
	systemRequest "server/model/system/request"
	systemService "server/service/system"

	"github.com/gin-gonic/gin"
)

// CodegenList 生成配置分页列表。
func CodegenList(c *gin.Context) {
	result, err := systemService.CodegenList(queryMap(c))
	successOrFail(c, result, err)
}

// CodegenDatasources 可选数据源列表。
func CodegenDatasources(c *gin.Context) {
	result, err := systemService.CodegenDatasources()
	successOrFail(c, result, err)
}

// CodegenDBTables 指定数据源下的数据表清单。
func CodegenDBTables(c *gin.Context) {
	result, err := systemService.CodegenDBTables(queryMap(c))
	successOrFail(c, result, err)
}

// CodegenImportTables 装载选中的数据表为生成配置。
func CodegenImportTables(c *gin.Context) {
	data, ok := bindJSON[systemRequest.CodegenImportPayload](c)
	if !ok {
		return
	}
	result, err := systemService.CodegenImportTables(data)
	successOrFail(c, result, err)
}

// CodegenDelete 删除生成配置。
func CodegenDelete(c *gin.Context) {
	payload, ok := bindJSON[systemRequest.IDsPayload](c)
	if !ok {
		return
	}
	successOrFail(c, map[string]interface{}{}, systemService.CodegenDelete(payload.IDs))
}

// CodegenDetail 生成配置详情。
func CodegenDetail(c *gin.Context) {
	result, err := systemService.CodegenDetail(c.Param("id"))
	successOrFail(c, result, err)
}

// CodegenUpdate 保存生成配置。
func CodegenUpdate(c *gin.Context) {
	data, ok := bindJSONMap(c)
	if !ok {
		return
	}
	result, err := systemService.CodegenUpdate(c.Param("id"), data)
	successOrFail(c, result, err)
}

// CodegenPreview 预览将要生成的代码。
func CodegenPreview(c *gin.Context) {
	result, err := systemService.CodegenPreview(c.Param("id"))
	successOrFail(c, result, err)
}

// CodegenGenerate 生成代码文件并同步菜单。
func CodegenGenerate(c *gin.Context) {
	result, err := systemService.CodegenGenerate(c.Param("id"))
	successOrFail(c, result, err)
}

// CodegenComponents 返回代码生成器支持的页面组件能力。
func CodegenComponents(c *gin.Context) {
	successOrFail(c, systemService.CodegenComponentCapabilities(), nil)
}

// CodegenOptionRoutes 返回可作为组件数据源的系统 GET 路由。
func CodegenOptionRoutes(c *gin.Context) {
	successOrFail(c, systemService.CodegenOptionRoutes(), nil)
}
