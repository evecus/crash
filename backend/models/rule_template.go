package models

// RuleTemplate 规则模板（subconverter 用的 .ini 模板）
type RuleTemplate struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`        // 显示名称
	URL         string `json:"url"`         // 远程 URL 或空（本地）
	LocalPath   string `json:"local_path"`  // 本地缓存路径
	Description string `json:"description"` // 说明
	IsDefault   bool   `json:"is_default"`  // 是否为当前选中模板
}
