package firewall

import (
	"crashpanel/models"
	"fmt"
	"strings"
)

func applyNftables(s *models.Settings, rules []models.FirewallRule) error {
	// 创建 table
	if err := runCmd("nft", "add", "table", "inet", "crashpanel"); err != nil {
		return err
	}

	tproxyPort := fmt.Sprintf("%d", s.TproxyPort)
	redirPort := fmt.Sprintf("%d", s.RedirPort)
	routingMark := fmt.Sprintf("0x%x", s.RedirPort+2)

	// 保留地址集合
	reserved := strings.Join(reservedIPv4, ", ")

	// 添加 prerouting chain
	chainType := "nat"
	priority := "-100"
	if s.RedirMod == "Tproxy" {
		chainType = "filter"
		priority = "-150"
	}

	if err := runCmd("nft", "add", "chain", "inet", "crashpanel", "prerouting",
		fmt.Sprintf("{ type %s hook prerouting priority %s; }", chainType, priority)); err != nil {
		return err
	}

	// 防回环
	if err := runCmd("nft", "add", "rule", "inet", "crashpanel", "prerouting",
		"meta", "mark", routingMark, "return"); err != nil {
		return err
	}

	// 跳过 DNS
	if err := runCmd("nft", "add", "rule", "inet", "crashpanel", "prerouting",
		"tcp", "dport", "53", "return"); err != nil {
		return err
	}
	if err := runCmd("nft", "add", "rule", "inet", "crashpanel", "prerouting",
		"udp", "dport", "53", "return"); err != nil {
		return err
	}

	// 跳过保留地址
	if err := runCmd("nft", "add", "rule", "inet", "crashpanel", "prerouting",
		"ip", "daddr", fmt.Sprintf("{%s}", reserved), "return"); err != nil {
		return err
	}

	// MAC/IP 过滤
	mode := filterMode(rules)
	if mode == "blacklist" {
		for _, mac := range macRules(rules) {
			if err := runCmd("nft", "add", "rule", "inet", "crashpanel", "prerouting",
				"ether", "saddr", mac, "return"); err != nil {
				return err
			}
		}
		if ips := ipRules(rules); len(ips) > 0 {
			ipSet := strings.Join(ips, ", ")
			if err := runCmd("nft", "add", "rule", "inet", "crashpanel", "prerouting",
				"ip", "saddr", fmt.Sprintf("{%s}", ipSet), "return"); err != nil {
				return err
			}
		}
	}

	// 根据模式设置跳转
	switch s.RedirMod {
	case "Tproxy":
		if err := runCmd("nft", "add", "rule", "inet", "crashpanel", "prerouting",
			"meta", "l4proto", "{tcp, udp}", "tproxy", "to", ":"+tproxyPort,
			"meta", "mark", "set", routingMark); err != nil {
			return err
		}
	default: // Redir
		if err := runCmd("nft", "add", "rule", "inet", "crashpanel", "prerouting",
			"tcp", "redirect", "to", ":"+redirPort); err != nil {
			return err
		}
	}

	return nil
}
