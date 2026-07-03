package system

// BizError 是可以直接展示给客户端的业务错误。
// api 层的 successOrFail 用 errors.As 识别：BizError 透传消息（HTTP 200 + OperationFailed），
// 其余错误一律视为内部错误，记日志后只返回泛化的 SystemError，不向客户端泄露细节。
type BizError struct {
	message string
}

func (e *BizError) Error() string {
	return e.message
}

// NewBizError 构造业务错误；固定文案优先使用下方具名 sentinel，动态文案再用本函数。
func NewBizError(message string) *BizError {
	return &BizError{message: message}
}

var (
	ErrMenuHasChildren    = NewBizError("菜单下存在子菜单，无法删除")
	ErrRoleHasChildren    = NewBizError("角色下存在子角色，无法删除")
	ErrDeptHasChildren    = NewBizError("部门下存在子部门，无法删除")
	ErrPasswordRequired   = NewBizError("password is required")
	ErrAttachmentIDsEmpty = NewBizError("附件ID不能为空")
	ErrNoRowsSelected     = NewBizError("请选择要操作的数据")
	ErrNoTablesSelected   = NewBizError("请选择要操作的数据表")
	ErrNoRecycleSupport   = NewBizError("该数据表没有 delete_time 字段，无法查看回收站")
	ErrNoIDColumn         = NewBizError("该数据表没有 id 字段，无法恢复或永久删除")
	ErrInvalidTableName   = NewBizError("非法数据表名称")
	ErrTableNotFound      = NewBizError("数据表不存在")
	ErrNoImportTables     = NewBizError("请选择要导入的数据表")
	ErrNoDeleteTables     = NewBizError("请选择要删除的数据表")
	ErrEmptyTableName     = NewBizError("存在未填写表名的数据表")
)
