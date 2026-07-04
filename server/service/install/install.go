package install

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	installModel "server/model/install"
	gormInit "server/setup/gorm"
	redisInit "server/setup/redis"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

const (
	lockFilePath = "runtime/install.lock"
	envFilePath  = ".env"
)

// Status 返回安装状态和可导入的 SQL 文件列表。
func Status() (installModel.StatusResponse, error) {
	files, err := SQLFiles()
	if err != nil {
		return installModel.StatusResponse{}, err
	}

	return installModel.StatusResponse{
		Installed: Installed(),
		SQLFiles:  files,
	}, nil
}

// Installed 通过锁文件判断系统是否已完成安装。
// 用文件而不是数据库做标记，是因为未安装时数据库本身还不可用。
func Installed() bool {
	_, err := os.Stat(lockFilePath)
	return err == nil
}

// SQLFiles 扫描 sql/ 和 database/ 目录下可执行的初始化 SQL 文件。
func SQLFiles() ([]string, error) {
	roots := []string{"sql", "database"}
	seen := map[string]bool{}
	files := make([]string, 0)

	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".sql") {
				continue
			}
			name := filepath.ToSlash(filepath.Join(root, entry.Name()))
			if !seen[name] {
				seen[name] = true
				files = append(files, name)
			}
		}
	}

	sort.Strings(files)
	return files, nil
}

// Check 分别探测 MySQL 和 Redis 连接是否可用，安装前的连通性自检。
func Check(req installModel.CheckRequest) (installModel.CheckResponse, error) {
	if err := pingMysql(req.Mysql, false); err != nil {
		return installModel.CheckResponse{}, err
	}
	if err := pingRedis(req.Redis); err != nil {
		return installModel.CheckResponse{}, err
	}
	return installModel.CheckResponse{MysqlOK: true, RedisOK: true}, nil
}

// Run 执行完整安装流程：建库 → 导入 SQL → 探测 Redis → 写 .env →
// 初始化连接 → 写安装锁。锁文件生成后 installGuard 中间件放行业务接口。
func Run(req installModel.InstallRequest) (map[string]interface{}, error) {
	if Installed() {
		return nil, errors.New("系统已安装，如需重新安装请先删除 runtime/install.lock")
	}

	if strings.TrimSpace(req.JWTSecret) == "" {
		req.JWTSecret = "gin-react-admin-change-me"
	}
	if len(req.SQLFiles) == 0 {
		files, err := SQLFiles()
		if err != nil {
			return nil, err
		}
		req.SQLFiles = files
	}

	if err := createDatabase(req.Mysql); err != nil {
		return nil, err
	}
	for _, file := range req.SQLFiles {
		if err := importSQLFile(req.Mysql, file); err != nil {
			return nil, err
		}
	}
	if err := pingRedis(req.Redis); err != nil {
		return nil, err
	}
	if err := writeEnv(req); err != nil {
		return nil, err
	}
	applyEnv(req)
	if err := gormInit.Gorm.InitializeAll(); err != nil {
		return nil, err
	}
	if err := redisInit.Redis.Initialize(); err != nil {
		return nil, err
	}
	if err := writeInstallLock(req); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"installed": true,
		"sqlFiles":  req.SQLFiles,
	}, nil
}

// pingMysql 用 5 秒超时探测 MySQL 连通性。
func pingMysql(cfg installModel.MysqlConfig, withDB bool) error {
	db, err := openMysql(cfg, withDB)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

// createDatabase 建目标库（若不存在）；库名经 escapeIdentifier 处理防注入。
func createDatabase(cfg installModel.MysqlConfig) error {
	db, err := openMysql(cfg, false)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `" + escapeIdentifier(cfg.Dbname) + "` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	return err
}

// importSQLFile 整文件执行初始化 SQL（连接开启了 MultiStatements）。
func importSQLFile(cfg installModel.MysqlConfig, name string) error {
	path, err := resolveSQLFile(name)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	db, err := openMysql(cfg, true)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("%s import failed: %w", filepath.ToSlash(path), err)
	}
	return nil
}

// openMysql 用 database/sql 直连 MySQL（安装阶段还没有 GORM 连接可用）。
func openMysql(cfg installModel.MysqlConfig, withDB bool) (*sql.DB, error) {
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = strings.TrimSpace(cfg.Host) + ":" + strings.TrimSpace(defaultString(cfg.Port, "3306"))
	mysqlCfg.User = strings.TrimSpace(cfg.User)
	mysqlCfg.Passwd = cfg.Password
	mysqlCfg.ParseTime = true
	mysqlCfg.Loc = time.Local
	mysqlCfg.Params = map[string]string{"charset": "utf8mb4"}
	mysqlCfg.MultiStatements = true
	if withDB {
		mysqlCfg.DBName = strings.TrimSpace(cfg.Dbname)
	}
	return sql.Open("mysql", mysqlCfg.FormatDSN())
}

// pingRedis 用 5 秒超时探测 Redis 连通性。
func pingRedis(cfg installModel.RedisConfig) error {
	client := redisClient(cfg)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}

func redisClient(cfg installModel.RedisConfig) redis.UniversalClient {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode == "cluster" {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    splitCSV(cfg.Addrs),
			Password: cfg.Password,
		})
	}
	return redis.NewClient(&redis.Options{
		Addr:     defaultString(cfg.Addr, "127.0.0.1:6379"),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func resolveSQLFile(name string) (string, error) {
	clean := filepath.Clean(name)
	if filepath.IsAbs(clean) || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("非法 SQL 文件路径: %s", name)
	}
	if strings.EqualFold(filepath.Ext(clean), ".sql") {
		if _, err := os.Stat(clean); err == nil {
			return clean, nil
		}
	}
	for _, root := range []string{"sql", "database"} {
		path := filepath.Join(root, filepath.Base(clean))
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("SQL 文件不存在: %s", name)
}

func writeEnv(req installModel.InstallRequest) error {
	values := map[string]string{
		"AI_SYSTEM_MYSQL_HOST":     req.Mysql.Host,
		"AI_SYSTEM_MYSQL_PORT":     defaultString(req.Mysql.Port, "3306"),
		"AI_SYSTEM_MYSQL_USER":     req.Mysql.User,
		"AI_SYSTEM_MYSQL_PASSWORD": req.Mysql.Password,
		"AI_SYSTEM_MYSQL_DB":       req.Mysql.Dbname,
		"AI_SYSTEM_MYSQL_CONFIG":   "charset=utf8mb4&parseTime=True&loc=Local",
		"REDIS_MODE":               defaultString(req.Redis.Mode, "single"),
		"REDIS_ADDR":               defaultString(req.Redis.Addr, "127.0.0.1:6379"),
		"REDIS_ADDRS":              req.Redis.Addrs,
		"REDIS_PASSWORD":           req.Redis.Password,
		"REDIS_DB":                 strconv.Itoa(req.Redis.DB),
		"JWT_SECRET":               req.JWTSecret,
	}
	content, _ := os.ReadFile(envFilePath)
	lines := mergeEnvLines(string(content), values)
	return os.WriteFile(envFilePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func mergeEnvLines(content string, values map[string]string) []string {
	used := map[string]bool{}
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	result := make([]string, 0, len(lines)+len(values))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		key, _, ok := strings.Cut(trimmed, "=")
		if ok {
			if value, exists := values[key]; exists {
				result = append(result, key+"="+value)
				used[key] = true
				continue
			}
		}
		result = append(result, line)
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		if !used[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		result = append(result, key+"="+values[key])
	}
	return result
}

func applyEnv(req installModel.InstallRequest) {
	_ = os.Setenv("AI_SYSTEM_MYSQL_HOST", req.Mysql.Host)
	_ = os.Setenv("AI_SYSTEM_MYSQL_PORT", defaultString(req.Mysql.Port, "3306"))
	_ = os.Setenv("AI_SYSTEM_MYSQL_USER", req.Mysql.User)
	_ = os.Setenv("AI_SYSTEM_MYSQL_PASSWORD", req.Mysql.Password)
	_ = os.Setenv("AI_SYSTEM_MYSQL_DB", req.Mysql.Dbname)
	_ = os.Setenv("AI_SYSTEM_MYSQL_CONFIG", "charset=utf8mb4&parseTime=True&loc=Local")
	_ = os.Setenv("REDIS_MODE", defaultString(req.Redis.Mode, "single"))
	_ = os.Setenv("REDIS_ADDR", defaultString(req.Redis.Addr, "127.0.0.1:6379"))
	_ = os.Setenv("REDIS_ADDRS", req.Redis.Addrs)
	_ = os.Setenv("REDIS_PASSWORD", req.Redis.Password)
	_ = os.Setenv("REDIS_DB", strconv.Itoa(req.Redis.DB))
	_ = os.Setenv("JWT_SECRET", req.JWTSecret)
}

func writeInstallLock(req installModel.InstallRequest) error {
	if err := os.MkdirAll(filepath.Dir(lockFilePath), 0755); err != nil {
		return err
	}
	return os.WriteFile(lockFilePath, []byte("installed_at="+time.Now().Format(time.RFC3339)+"\ndatabase="+req.Mysql.Dbname+"\n"), 0644)
}

func escapeIdentifier(value string) string {
	return strings.ReplaceAll(strings.TrimSpace(value), "`", "``")
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}
