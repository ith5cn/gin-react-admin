package request

// CodegenImportTable 描述一张待装载的数据表。
type CodegenImportTable struct {
	TableName    string `json:"tableName"`
	TableComment string `json:"tableComment"`
}

// CodegenImportPayload 是"装载数据表"接口入参：从哪个数据源导入哪些表。
type CodegenImportPayload struct {
	Source string               `json:"source"`
	Tables []CodegenImportTable `json:"tables"`
}
