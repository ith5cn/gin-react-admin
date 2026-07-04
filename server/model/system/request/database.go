package request

// DatabaseTablesPayload 是表维护类接口（优化/清理碎片）的入参。
type DatabaseTablesPayload struct {
	Tables []string `json:"tables"`
}

// DatabaseRecyclePayload 是回收站恢复/销毁接口入参。
// IDs 用 interface{} 是因为不同表主键类型不一（数字或字符串）。
type DatabaseRecyclePayload struct {
	TableName string        `json:"tableName"`
	IDs       []interface{} `json:"ids"`
}
