package request

// DeptPayload 是部门创建/更新入参；指针字段为 nil 表示前端未提交该字段，更新时跳过。
type DeptPayload struct {
	ParentID *uint   `json:"parentId"`
	Name     *string `json:"name"`
	Status   *int16  `json:"status"`
	Sort     *uint16 `json:"sort"`
	Remark   *string `json:"remark"`
}
