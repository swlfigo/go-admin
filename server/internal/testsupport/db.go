package testsupport

import (
	"os"

	"go-admin/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConfig 返回测试用数据库连接（独立于开发库 go_admin，默认 go_admin_test，
// 可用环境变量 GOADMIN_TEST_DB 覆盖）。
func DBConfig() config.DatabaseConfig {
	name := os.Getenv("GOADMIN_TEST_DB")
	if name == "" {
		name = "go_admin_test"
	}
	return config.DatabaseConfig{
		Host: "127.0.0.1", Port: 5432, User: "postgres",
		Password: "postgres", Name: name, SSLMode: "disable",
	}
}

func dsn(c config.DatabaseConfig) string {
	return "host=" + c.Host +
		" port=" + itoa(c.Port) +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.Name +
		" sslmode=" + c.SSLMode
}

func itoa(n int) string {
	// 端口数字转字符串，避免引入额外依赖
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

// ensureDatabase 连到维护库 postgres，若测试库不存在则创建。
func ensureDatabase(c config.DatabaseConfig) error {
	admin := c
	admin.Name = "postgres"
	db, err := gorm.Open(postgres.Open(dsn(admin)), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	var exists bool
	if err := db.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", c.Name).Scan(&exists).Error; err != nil {
		return err
	}
	if !exists {
		// CREATE DATABASE 不支持参数占位；c.Name 来自受控常量，非用户输入，安全。
		_ = db.Exec("CREATE DATABASE " + c.Name).Error // 并发/已存在忽略
	}
	return nil
}

// Connect 连接测试库（不存在则自动创建），返回 *gorm.DB。
func Connect() (*gorm.DB, error) {
	c := DBConfig()
	if err := ensureDatabase(c); err != nil {
		return nil, err
	}
	return gorm.Open(postgres.Open(dsn(c)), &gorm.Config{})
}

// EnsureDB 确保测试库存在（不存在则创建）。供需要直接调用 database.Connect 的测试预先调用，
// 这样单独运行某个包（如 go test ./internal/database/）也能自建库，不依赖其它包先跑。
func EnsureDB() error {
	return ensureDatabase(DBConfig())
}
