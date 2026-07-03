package request

// IDsPayload 是 { ids: [] } 形式的通用入参，用于角色绑菜单、用户绑角色等接口。
type IDsPayload struct {
	IDs []uint `json:"ids"`
}
