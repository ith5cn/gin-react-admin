package system

import (
	systemApi "server/api/system"

	"github.com/gin-gonic/gin"
)

// BaseRouter 注册系统基础模块路由。
// 当前包含登录、刷新 token、退出登录；公开和私有路由分别挂到传入的路由组。
func BaseRouter(PublicGroup, PrivateGroup *gin.RouterGroup) {
	base := PublicGroup.Group("/base")
	base.POST("/login", systemApi.Login)
	base.POST("/token/refresh", systemApi.RefreshToken)

	privateBase := PrivateGroup.Group("/base")
	privateBase.POST("/logout", systemApi.Logout)

	system := PrivateGroup.Group("/system")
	system.GET("/user", systemApi.CurrentUser)

	system.GET("/user/index", systemApi.UserList)
	system.GET("/user/auth-list", systemApi.UserAuthList)
	system.POST("/user", systemApi.CreateUser)
	system.PUT("/user/:id/refresh-cache", systemApi.RefreshUserCache)
	system.PUT("/user/:id/set-password", systemApi.SetUserPassword)
	system.POST("/user/:id/role", systemApi.BindUserRole)
	system.PUT("/user/:id", systemApi.UpdateUser)
	system.DELETE("/user/:id", systemApi.DeleteUser)

	system.GET("/menu/index", systemApi.MenuList)
	system.GET("/menu/accessMenu", systemApi.AccessMenu)
	system.GET("/menu/getMenuByRole/:roleId", systemApi.MenuByRole)
	system.POST("/menu/create", systemApi.CreateMenu)
	system.PUT("/menu/:id", systemApi.UpdateMenu)
	system.DELETE("/menu/:id", systemApi.DeleteMenu)

	system.GET("/role/index", systemApi.RoleList)
	system.GET("/role/access", systemApi.RoleAccess)
	system.POST("/role/create", systemApi.CreateRole)
	system.POST("/role/:id/menu", systemApi.BindRoleMenu)
	system.PUT("/role/:id", systemApi.UpdateRole)
	system.DELETE("/role/:id", systemApi.DeleteRole)

	system.GET("/dept/index", systemApi.DeptList)
	system.GET("/dept/access", systemApi.DeptAccess)
	system.POST("/dept", systemApi.CreateDept)
	system.PUT("/dept/:id", systemApi.UpdateDept)
	system.DELETE("/dept/:id", systemApi.DeleteDept)

	system.GET("/post/index", systemApi.PostList)
	system.GET("/post/access", systemApi.PostAccess)
	system.POST("/post", systemApi.CreatePost)
	system.PUT("/post/:id", systemApi.UpdatePost)
	system.DELETE("/post/:id", systemApi.DeletePost)

	system.GET("/dict-type/index", systemApi.DictTypeList)
	system.POST("/dict-type", systemApi.CreateDictType)
	system.PUT("/dict-type/:id", systemApi.UpdateDictType)
	system.DELETE("/dict-type/:id", systemApi.DeleteDictType)

	system.GET("/dict-data/index", systemApi.DictDataList)
	system.GET("/dict-data/dictAll", systemApi.DictAll)
	system.POST("/dict-data", systemApi.CreateDictData)
	system.PUT("/dict-data/:id", systemApi.UpdateDictData)
	system.DELETE("/dict-data/:id", systemApi.DeleteDictData)

	system.GET("/config-group/index", systemApi.ConfigGroupList)
	system.POST("/config-group", systemApi.CreateConfigGroup)
	system.PUT("/config-group/:id", systemApi.UpdateConfigGroup)
	system.DELETE("/config-group/:id", systemApi.DeleteConfigGroup)

	system.GET("/config/index", systemApi.ConfigList)
	system.GET("/config/get-config-info", systemApi.ConfigInfo)
	system.POST("/config", systemApi.CreateConfig)
	system.POST("/config/batch-update", systemApi.BatchUpdateConfig)
	system.PUT("/config/:id", systemApi.UpdateConfig)
	system.DELETE("/config/:id", systemApi.DeleteConfig)

	system.GET("/login-log/index", systemApi.LoginLogList)
	system.GET("/oper-log/index", systemApi.OperLogList)

	system.GET("/attachment/index", systemApi.AttachmentList)
	system.POST("/attachment/delete", systemApi.DeleteAttachments)

	system.GET("/codegen/index", systemApi.CodegenList)
	system.GET("/codegen/datasources", systemApi.CodegenDatasources)
	system.GET("/codegen/db-tables", systemApi.CodegenDBTables)
	system.POST("/codegen/importTables", systemApi.CodegenImportTables)
	system.POST("/codegen/delete", systemApi.CodegenDelete)
	system.GET("/codegen/detail/:id", systemApi.CodegenDetail)
	system.PUT("/codegen/:id", systemApi.CodegenUpdate)
	system.POST("/codegen/generate/:id", systemApi.CodegenGenerate)
	system.GET("/codegen/preview/:id", systemApi.CodegenPreview)

	data := PrivateGroup.Group("/data")
	data.GET("/database/index", systemApi.DatabaseTableList)
	data.GET("/database/columns/:tableName", systemApi.DatabaseTableColumns)
	data.POST("/database/fragment", systemApi.DatabaseClearFragments)
	data.POST("/database/optimize", systemApi.DatabaseOptimizeTables)
	data.GET("/database/recycle", systemApi.DatabaseRecycleList)
	data.POST("/database/recycle/recover", systemApi.DatabaseRecycleRecover)
	data.POST("/database/recycle/destroy", systemApi.DatabaseRecycleDestroy)
}
