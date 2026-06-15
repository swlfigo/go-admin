package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Captcha  CaptchaConfig  `mapstructure:"captcha"`
	Login    LoginConfig    `mapstructure:"login"`
	Admin    AdminConfig    `mapstructure:"admin"`
}

// AdminConfig 初始超级管理员账号（仅首次建库时写入；之后改密走个人中心，不会被覆盖）。
type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	AccessMinutes int    `mapstructure:"accessMinutes"`
	RefreshDays   int    `mapstructure:"refreshDays"`
}

type CaptchaConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type LoginConfig struct {
	MaxFailCount       int `mapstructure:"maxFailCount"`
	LockMinutes        int `mapstructure:"lockMinutes"`
	RateLimitPerMinute int `mapstructure:"rateLimitPerMinute"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetDefault("admin.username", "admin")
	v.SetDefault("admin.password", "admin123")
	// 环境变量覆盖（容器/部署用）：如 GOADMIN_DATABASE_HOST 覆盖 database.host、
	// GOADMIN_JWT_SECRET 覆盖 jwt.secret，便于不改 config.yaml 注入敏感值。
	for _, k := range []string{
		"server.port",
		"database.host", "database.port", "database.user", "database.password", "database.name", "database.sslmode",
		"jwt.secret", "jwt.accessMinutes", "jwt.refreshDays",
		"captcha.enabled",
		"login.maxFailCount", "login.lockMinutes", "login.rateLimitPerMinute",
		"admin.username", "admin.password",
	} {
		_ = v.BindEnv(k, "GOADMIN_"+strings.ToUpper(strings.ReplaceAll(k, ".", "_")))
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
