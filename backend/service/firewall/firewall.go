package firewall

import (
	"crashpanel/models"
	"fmt"
)

// Apply 根据设置下发防火墙规则
func Apply(s *models.Settings, rules []models.FirewallRule) error {
	// 先清理旧规则
	if err := Flush(s.FirewallMod); err != nil {
		return err
	}

	switch s.FirewallMod {
	case "nftables":
		return applyNftables(s, rules)
	default:
		return applyIptables(s, rules)
	}
}

// Flush 清理所有 crashpanel 防火墙规则
func Flush(mod string) error {
	switch mod {
	case "nftables":
		return runCmd("nft", "delete", "table", "inet", "crashpanel")
	default:
		return flushIptables()
	}
}

func flushIptables() error {
	chains := [][]string{
		{"iptables", "-t", "nat", "-F", "CRASHPANEL"},
		{"iptables", "-t", "nat", "-X", "CRASHPANEL"},
		{"iptables", "-t", "mangle", "-F", "CRASHPANEL"},
		{"iptables", "-t", "mangle", "-X", "CRASHPANEL"},
		{"ip6tables", "-t", "nat", "-F", "CRASHPANEL"},
		{"ip6tables", "-t", "nat", "-X", "CRASHPANEL"},
	}
	for _, args := range chains {
		// 忽略错误（链不存在时会报错）
		_ = runCmd(args[0], args[1:]...)
	}
	return nil
}

// reservedIPv4 私有/保留地址段
var reservedIPv4 = []string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"224.0.0.0/4",
	"240.0.0.0/4",
}

// filterMode 获取过滤规则中 whitelist/blacklist 统计
func filterMode(rules []models.FirewallRule) string {
	for _, r := range rules {
		return r.FilterMode
	}
	return "blacklist"
}

func macRules(rules []models.FirewallRule) []string {
	var macs []string
	for _, r := range rules {
		if r.Type == "mac" {
			macs = append(macs, r.Value)
		}
	}
	return macs
}

func ipRules(rules []models.FirewallRule) []string {
	var ips []string
	for _, r := range rules {
		if r.Type == "ip" {
			ips = append(ips, r.Value)
		}
	}
	return ips
}

func validateMod(mod string) error {
	if mod != "iptables" && mod != "nftables" {
		return fmt.Errorf("unknown firewall mod: %s", mod)
	}
	return nil
}
