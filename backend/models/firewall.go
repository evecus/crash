package models

import "time"

type FirewallRule struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Type       string    `json:"type"`        // mac | ip
	Value      string    `json:"value"`
	FilterMode string    `json:"filter_mode"` // whitelist | blacklist
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}
