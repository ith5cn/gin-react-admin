package request

// RolePayload 是角色创建/更新入参。
// dataScope 与 AISystemRole.DataScope 保持 *int16 类型一致。
type RolePayload struct {
	ParentID  *uint   `json:"parentId"`
	Name      *string `json:"name"`
	Code      *string `json:"code"`
	DataScope *int16  `json:"dataScope"`
	Status    *int16  `json:"status"`
	Sort      *uint16 `json:"sort"`
	Remark    *string `json:"remark"`
}
