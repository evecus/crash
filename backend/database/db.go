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

	if err := DB.AutoMigrate(
		&models.Settings{},
		&models.Subscription{},
		&models.FirewallRule{},
		&models.Task{},
		&models.DNSConfig{},
		&models.RuleTemplate{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	seedDefaults()
}

func seedDefaults() {
	var count int64

	// Settings
	DB.Model(&models.Settings{}).Count(&count)
	if count == 0 {
		DB.Create(&models.Settings{
			CoreType:       "meta",
			MixPort:        7890,
			RedirPort:      7892,
			TproxyPort:     7893,
			DashboardPort:  9999,
			DNSPort:        1053,
			RedirMod:       "Redir",
			FirewallMod:    "nftables",
			FirewallArea:   2,
			DNSMod:         "redir-host",
			CNIPRoute:      true,
			CommonPorts:    false,
			CorePath:       "/usr/local/bin/mihomo",
			CoreWorkDir:    "/etc/crashpanel/core",
		})
	}

	// DNSConfig
	DB.Model(&models.DNSConfig{}).Count(&count)
	if count == 0 {
		DB.Create(&models.DNSConfig{
			Nameserver:     "223.5.5.5, 1.2.4.8",
			Fallback:       "1.1.1.1, 8.8.8.8",
			FakeIPFilter:   "*.lan, *.local, *.home.arpa",
			FallbackFilter: "geoip:cn",
		})
	}

	// RuleTemplate 内置模板列表
	DB.Model(&models.RuleTemplate{}).Count(&count)
	if count == 0 {
		templates := []models.RuleTemplate{
			{Name: "ACL4SSR 全能优化版（推荐）",      URL: "https://github.com/juewuy/ShellCrash/raw/dev/rules/ShellClash.ini",            IsDefault: true},
			{Name: "ACL4SSR 精简优化版（推荐）",      URL: "https://github.com/juewuy/ShellCrash/raw/dev/rules/ShellClash_Mini.ini"},
			{Name: "ACL4SSR 全能+去广告",            URL: "https://github.com/juewuy/ShellCrash/raw/dev/rules/ShellClash_Block.ini"},
			{Name: "ACL4SSR 极简版（适合自建节点）",  URL: "https://github.com/juewuy/ShellCrash/raw/dev/rules/ShellClash_Nano.ini"},
			{Name: "ACL4SSR 分流+游戏增强",          URL: "https://github.com/juewuy/ShellCrash/raw/dev/rules/ShellClash_Full.ini"},
			{Name: "ACL4SSR 多国精简",               URL: "https://github.com/juewuy/ShellCrash/raw/dev/rules/ACL4SSR_Online_Mini_MultiCountry.ini"},
			{Name: "ACL4SSR 回国专用",               URL: "https://github.com/juewuy/ShellCrash/raw/dev/rules/ACL4SSR_BackCN.ini"},
			{Name: "DustinWin 精简规则（推荐）",      URL: "https://raw.githubusercontent.com/DustinWin/ruleset_geodata/master/rule_templates/DustinWin_Lite.ini"},
			{Name: "DustinWin 全分组规则",            URL: "https://raw.githubusercontent.com/DustinWin/ruleset_geodata/master/rule_templates/DustinWin_Full.ini"},
		}
		DB.Create(&templates)
	}
}
