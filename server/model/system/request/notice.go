package request

// NoticePayload 是通知公告创建/更新入参，指针字段为 nil 表示不改动（部分更新）。
// Type/Status 用 FlexInt16 兼容字典下拉提交的字符串数字。
type NoticePayload struct {
	Title   *string    `json:"title"`
	Type    *FlexInt16 `json:"type"`
	Content *string    `json:"content"`
	Status  *FlexInt16 `json:"status"`
	Remark  *string    `json:"remark"`
}
