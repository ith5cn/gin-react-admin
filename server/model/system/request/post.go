package request

// PostPayload 是岗位创建/更新入参。
type PostPayload struct {
	Name   *string `json:"name"`
	Code   *string `json:"code"`
	Sort   *uint16 `json:"sort"`
	Status *int16  `json:"status"`
	Remark *string `json:"remark"`
}
