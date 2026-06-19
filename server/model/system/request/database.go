package request

type DatabaseTablesPayload struct {
	Tables []string `json:"tables"`
}

type DatabaseRecyclePayload struct {
	TableName string        `json:"tableName"`
	IDs       []interface{} `json:"ids"`
}
