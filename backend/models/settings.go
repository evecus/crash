package models

type Settings struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// 核心
	CoreType    string `json:"core_type"`    // meta | singbox
	CorePath    string `json:"core_path"`    // 二进制路径
	CoreWorkDir string `json:"core_work_dir"` // 工作目录

	// 端口
	MixPort       int `json:"mix_port"`
	RedirPort     int `json:"redir_port"`
	TproxyPort    int `json:"tproxy_port"`
	DashboardPort int `json:"dashboard_port"`
	DNSPort       int `json:"dns_port"`

	// 代理模式
	RedirMod string `json:"redir_mod"` // Redir | Tproxy | Tun | Mix

	// 防火墙
	FirewallMod  string `json:"firewall_mod"`  // iptables | nftables
	FirewallArea int    `json:"firewall_area"` // 1=全局 2=绕过CN 3=仅局域网 4=纯净

	// DNS
	DNSMod string `json:"dns_mod"` // redir-host | fake-ip | mix

	// 功能开关
	CNIPRoute   bool `json:"cn_ip_route"`
	CommonPorts bool `json:"common_ports"`
}
