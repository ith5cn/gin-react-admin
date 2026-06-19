package response

// UserInfo 是前端登录后需要的当前用户资料。
type UserInfo struct {
	ID             uint    `json:"id"`
	Username       string  `json:"username"`
	UserType       string  `json:"user_type"`
	Nickname       *string `json:"nickname"`
	Phone          *string `json:"phone"`
	Email          *string `json:"email"`
	Avatar         *string `json:"avatar"`
	Signed         *string `json:"signed"`
	Dashboard      *string `json:"dashboard"`
	DeptID         *uint   `json:"deptId"`
	Status         int16   `json:"status"`
	BackendSetting *string `json:"backendSetting"`
}

// MenuNode 是前端 RawBackendMenuNode 对应的数据结构。
type MenuNode struct {
	ID        uint        `json:"id"`
	ParentID  *uint       `json:"parentId"`
	Name      string      `json:"name"`
	Code      string      `json:"code"`
	Icon      *string     `json:"icon"`
	Route     *string     `json:"route"`
	Component *string     `json:"component"`
	Redirect  *string     `json:"redirect"`
	IsHidden  int16       `json:"isHidden"`
	IsLayout  uint8       `json:"isLayout"`
	Type      string      `json:"type"`
	Status    int16       `json:"status"`
	Sort      uint16      `json:"sort"`
	Children  []*MenuNode `json:"children,omitempty"`
}

// UserContext 是前端初始化用户、菜单、权限上下文需要的完整响应。
type UserContext struct {
	User    UserInfo      `json:"user"`
	Roles   []interface{} `json:"roles"`
	Routers []*MenuNode   `json:"routers"`
	Codes   interface{}   `json:"codes"`
	Posts   []interface{} `json:"posts"`
	Depts   []interface{} `json:"depts"`
}
