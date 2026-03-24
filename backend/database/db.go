package database

import (
	"crashpanel/models"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(dbPath string) {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create db dir: %v", err)
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// 自动迁移所有模型
	if err := DB.AutoMigrate(
		&models.Settings{},
		&models.Subscription{},
		&models.FirewallRule{},
		&models.Task{},
		&models.DNSConfig{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 初始化默认数据
	seedDefaults()
}

func seedDefaults() {
	// Settings 默认值（只在表为空时插入）
	var count int64
	DB.Model(&models.Settings{}).Count(&count)
	if count == 0 {
		DB.Create(&models.Settings{
			CoreType:      "meta",
			MixPort:       7890,
			RedirPort:     7892,
			TproxyPort:    7893,
			DashboardPort: 9999,
			DNSPort:       1053,
			RedirMod:      "Redir",
			FirewallMod:   "nftables",
			FirewallArea:  2,
			DNSMod:        "redir-host",
			CNIPRoute:     true,
			CommonPorts:   false,
			CorePath:      "/usr/local/bin/mihomo",
			CoreWorkDir:   "/etc/crashpanel/core",
		})
	}

	// DNSConfig 默认值
	DB.Model(&models.DNSConfig{}).Count(&count)
	if count == 0 {
		DB.Create(&models.DNSConfig{
			Nameserver:     "223.5.5.5, 1.2.4.8",
			Fallback:       "1.1.1.1, 8.8.8.8",
			FakeIPFilter:   "*.lan, *.local, *.home.arpa",
			FallbackFilter: "geoip:cn",
		})
	}
}
