package request

// ConfigGroupPayload 是配置分组创建/更新入参。
type ConfigGroupPayload struct {
	Name   *string `json:"name"`
	Code   *string `json:"code"`
	Sort   *uint16 `json:"sort"`
	Remark *string `json:"remark"`
}

// ConfigPayload 是配置项创建/更新入参。
// Value 保持 interface{}：radio/select 组件的值可能是数字或字符串，由 service 层归一化。
type ConfigPayload struct {
	ID               *uint       `json:"id"`
	GroupID          *uint       `json:"groupId"`
	Key              *string     `json:"key"`
	Value            interface{} `json:"value"`
	Name             *string     `json:"name"`
	InputType        *string     `json:"inputType"`
	ConfigSelectData *string     `json:"configSelectData"`
	Sort             *uint16     `json:"sort"`
	Remark           *string     `json:"remark"`
}

// BatchUpdateConfigPayload 是配置项批量保存入参。
type BatchUpdateConfigPayload struct {
	GroupID uint            `json:"groupId"`
	Config  []ConfigPayload `json:"config"`
}
