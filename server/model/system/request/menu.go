package request

// MenuPayload 是菜单创建/更新入参，字段类型与 AISystemMenu 对齐。
type MenuPayload struct {
	ParentID  *uint   `json:"parentId"`
	Name      *string `json:"name"`
	Code      *string `json:"code"`
	Icon      *string `json:"icon"`
	Route     *string `json:"route"`
	Component *string `json:"component"`
	Redirect  *string `json:"redirect"`
	IsHidden  *int16  `json:"isHidden"`
	IsLayout  *uint8  `json:"isLayout"`
	Type      *string `json:"type"`
	Status    *int16  `json:"status"`
	Sort      *uint16 `json:"sort"`
	Remark    *string `json:"remark"`
}
