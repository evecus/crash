package models

type DNSConfig struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Nameserver     string `json:"nameserver"`      // 国内 DNS，逗号分隔
	Fallback       string `json:"fallback"`        // 国外 DNS，逗号分隔
	FakeIPFilter   string `json:"fake_ip_filter"`  // fake-ip 过滤，逗号/换行分隔
	FallbackFilter string `json:"fallback_filter"` // fallback 过滤条件
}
