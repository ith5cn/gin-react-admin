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
	// DeptIDs 只在 dataScope=2（自定义数据权限）时有意义：nil 表示不改动，空数组表示清空授权。
	DeptIDs *[]uint `json:"deptIds"`
}
