package install

// MysqlConfig 是安装向导提交的 MySQL 连接信息。
type MysqlConfig struct {
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password"`
	Dbname   string `json:"dbname" binding:"required"`
}

// RedisConfig 是安装向导提交的 Redis 连接信息。
type RedisConfig struct {
	Mode     string `json:"mode"`
	Addr     string `json:"addr"`
	Addrs    string `json:"addrs"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// CheckRequest 是连通性检测接口入参。
type CheckRequest struct {
	Mysql MysqlConfig `json:"mysql" binding:"required"`
	Redis RedisConfig `json:"redis"`
}

// InstallRequest 是执行安装接口入参。
type InstallRequest struct {
	Mysql     MysqlConfig `json:"mysql" binding:"required"`
	Redis     RedisConfig `json:"redis"`
	SQLFiles  []string    `json:"sqlFiles"`
	JWTSecret string      `json:"jwtSecret"`
}

// StatusResponse 是安装状态接口响应。
type StatusResponse struct {
	Installed bool     `json:"installed"`
	SQLFiles  []string `json:"sqlFiles"`
}

// CheckResponse 是连通性检测接口响应。
type CheckResponse struct {
	MysqlOK bool `json:"mysqlOk"`
	RedisOK bool `json:"redisOk"`
}
