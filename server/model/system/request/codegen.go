package request

type CodegenImportTable struct {
	TableName    string `json:"tableName"`
	TableComment string `json:"tableComment"`
}

type CodegenImportPayload struct {
	Source string               `json:"source"`
	Tables []CodegenImportTable `json:"tables"`
}
