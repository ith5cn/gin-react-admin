package system

import "time"

// AISystemUser 映射 ai_system_user 用户表。
// model 层只描述数据库字段和表名，不承载登录、校验、签发 token 等业务逻辑。
type AISystemUser struct {
	ID             uint       `json:"id" gorm:"column:id;primaryKey"`               // 用户ID，主键。
	Username       string     `json:"username" gorm:"column:username"`              // 登录用户名。
	Password       string     `json:"-" gorm:"column:password"`                     // bcrypt 加密后的密码。
	UserType       *string    `json:"user_type" gorm:"column:user_type"`            // 用户类型，例如 100 表示系统用户，可为空。
	Nickname       *string    `json:"nickname" gorm:"column:nickname"`              // 用户昵称，可为空。
	Phone          *string    `json:"phone" gorm:"column:phone"`                    // 手机号，可为空。
	Email          *string    `json:"email" gorm:"column:email"`                    // 邮箱，可为空。
	Avatar         *string    `json:"avatar" gorm:"column:avatar"`                  // 头像地址，可为空。
	Signed         *string    `json:"signed" gorm:"column:signed"`                  // 个人签名，可为空。
	Dashboard      *string    `json:"dashboard" gorm:"column:dashboard"`            // 后台首页类型，可为空。
	DeptID         *uint      `json:"deptId" gorm:"column:dept_id"`                 // 部门ID，可为空。
	Status         int16      `json:"status" gorm:"column:status"`                  // 状态：1 正常，2 停用。
	LoginIP        *string    `json:"loginIp" gorm:"column:login_ip"`               // 最后登录IP，可为空。
	LoginTime      *time.Time `json:"loginTime" gorm:"column:login_time"`           // 最后登录时间，可为空。
	BackendSetting *string    `json:"backendSetting" gorm:"column:backend_setting"` // 后台设置数据，可为空。
	Remark         *string    `json:"remark" gorm:"column:remark"`                  // 备注，可为空。
	CreatedBy      *int       `json:"createdBy" gorm:"column:created_by"`           // 创建者ID，可为空。
	UpdatedBy      *int       `json:"updatedBy" gorm:"column:updated_by"`           // 更新者ID，可为空。
	CreateTime     *time.Time `json:"createTime" gorm:"column:create_time"`         // 创建时间，可为空。
	UpdateTime     *time.Time `json:"updateTime" gorm:"column:update_time"`         // 更新时间，可为空。
	DeleteTime     *time.Time `json:"-" gorm:"column:delete_time"`                  // 删除时间，可为空。
	Roles          []uint     `json:"roles" gorm:"-"`                               // 前端用户列表需要的角色ID数组。
}

// TableName 明确指定 GORM 查询 ai_system_user 表。
func (AISystemUser) TableName() string {
	return "ai_system_user"
}
