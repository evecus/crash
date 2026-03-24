package models

import "time"

type Task struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Cron      string    `json:"cron"`    // cron 表达式，如 "0 3 * * *"
	Command   string    `json:"command"` // shell 命令
	Enabled   bool      `json:"enabled"`
	LastRun   time.Time `json:"last_run"`
	LastLog   string    `json:"last_log"`
	LastCode  int       `json:"last_code"` // 上次退出码
	CreatedAt time.Time `json:"created_at"`
}
