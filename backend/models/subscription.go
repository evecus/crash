package models

import "time"

type Subscription struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Target     string    `json:"target"`      // clash | singbox
	UserAgent  string    `json:"user_agent"`
	AutoUpdate bool      `json:"auto_update"`
	Interval   int       `json:"interval"`    // 自动更新间隔（小时）
	Status     string    `json:"status"`      // ok | error | pending
	NodeCount  int       `json:"node_count"`
	ErrorMsg   string    `json:"error_msg"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`

	// subconverter 参数
	SubConverterURL string `json:"sub_converter_url"` // 外部 subconverter 地址
	Include         string `json:"include"`
	Exclude         string `json:"exclude"`
	ConfigURL       string `json:"config_url"` // 远程规则配置
}
