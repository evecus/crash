package models

import "time"

// LinkType 提供者链接类型
// url  = 订阅链接（http/https，经 subconverter 转换）
// file = 本地文件（已上传到 providers/ 目录）
// uri  = 单节点分享链接（vmess:// ss:// trojan:// 等）
type LinkType = string

const (
	LinkTypeURL  LinkType = "url"
	LinkTypeFile LinkType = "file"
	LinkTypeURI  LinkType = "uri"
)

type Subscription struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	LinkType  string    `json:"link_type"` // url | file | uri
	URL       string    `json:"url"`       // 订阅链接 或 uri分享链接
	FilePath  string    `json:"file_path"` // 本地文件路径（link_type=file 时使用）
	Target    string    `json:"target"`    // clash | singbox

	// subconverter 参数（link_type=url 时使用）
	SubConverterURL string `json:"sub_converter_url"`
	UserAgent       string `json:"user_agent"`
	Include         string `json:"include"`
	Exclude         string `json:"exclude"`
	ConfigURL       string `json:"config_url"` // 远程规则模板

	// 更新策略
	AutoUpdate bool `json:"auto_update"`
	Interval   int  `json:"interval"` // 自动更新间隔（小时）

	// 健康检查间隔
	HealthInterval int `json:"health_interval"` // 分钟，默认 3
	UpdateInterval int `json:"update_interval"` // 小时，默认 12

	// 状态
	Status    string    `json:"status"`     // ok | error | pending
	NodeCount int       `json:"node_count"`
	ErrorMsg  string    `json:"error_msg"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
