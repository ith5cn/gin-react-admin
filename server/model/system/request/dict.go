package request

// DictTypePayload 是字典类型创建/更新入参。
type DictTypePayload struct {
	Name   *string `json:"name"`
	Code   *string `json:"code"`
	Status *int16  `json:"status"`
	Remark *string `json:"remark"`
}

// DictDataPayload 是字典数据创建/更新入参。
type DictDataPayload struct {
	TypeID *uint   `json:"typeId"`
	Label  *string `json:"label"`
	Value  *string `json:"value"`
	Color  *string `json:"color"`
	Code   *string `json:"code"`
	Sort   *uint16 `json:"sort"`
	Status *int16  `json:"status"`
	Remark *string `json:"remark"`
}
