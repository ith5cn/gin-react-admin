package install

type MysqlConfig struct {
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	User     string `json:"user" binding:"required"`
	Password string `json:"password"`
	Dbname   string `json:"dbname" binding:"required"`
}

type RedisConfig struct {
	Mode     string `json:"mode"`
	Addr     string `json:"addr"`
	Addrs    string `json:"addrs"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type CheckRequest struct {
	Mysql MysqlConfig `json:"mysql" binding:"required"`
	Redis RedisConfig `json:"redis"`
}

type InstallRequest struct {
	Mysql     MysqlConfig `json:"mysql" binding:"required"`
	Redis     RedisConfig `json:"redis"`
	SQLFiles  []string    `json:"sqlFiles"`
	JWTSecret string      `json:"jwtSecret"`
}

type StatusResponse struct {
	Installed bool     `json:"installed"`
	SQLFiles  []string `json:"sqlFiles"`
}

type CheckResponse struct {
	MysqlOK bool `json:"mysqlOk"`
	RedisOK bool `json:"redisOk"`
}
