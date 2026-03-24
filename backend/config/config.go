package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port      int    `json:"port"`
	DBPath    string `json:"db_path"`
	JWTSecret string `json:"jwt_secret"`
	Debug     bool   `json:"debug"`
	// 面板登录密码（首次启动若无则用默认值）
	AdminPassword string `json:"admin_password"`
}

const defaultConfigPath = "/etc/crashpanel/config.json"

func Load() *Config {
	cfg := &Config{
		Port:          8080,
		DBPath:        "/etc/crashpanel/crashpanel.db",
		JWTSecret:     "crashpanel-default-secret-change-me",
		Debug:         false,
		AdminPassword: "admin",
	}

	// 优先读取配置文件
	path := os.Getenv("CRASHPANEL_CONFIG")
	if path == "" {
		path = defaultConfigPath
	}

	data, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(data, cfg)
	}

	// 环境变量覆盖
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWTSecret = secret
	}
	if pwd := os.Getenv("ADMIN_PASSWORD"); pwd != "" {
		cfg.AdminPassword = pwd
	}

	return cfg
}
