package firewall

import (
	"crashpanel/models"
	"fmt"
	"os/exec"
)

func runCmd(name string, args ...string) error {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %v: %s", name, args, string(out))
	}
	return nil
}

func applyIptables(s *models.Settings, rules []models.FirewallRule) error {
	tproxyPort := fmt.Sprintf("%d", s.TproxyPort)
	redirPort := fmt.Sprintf("%d", s.RedirPort)
	routingMark := fmt.Sprintf("%d", s.RedirPort+2)

	// 创建 CRASHPANEL 链（nat 表，用于 TCP Redir）
	if err := runCmd("iptables", "-t", "nat", "-N", "CRASHPANEL"); err != nil {
		return err
	}

	// 防回环：跳过已打标记的流量
	if err := runCmd("iptables", "-t", "nat", "-A", "CRASHPANEL",
		"-m", "mark", "--mark", routingMark, "-j", "RETURN"); err != nil {
		return err
	}

	// 跳过 DNS
	for _, proto := range []string{"tcp", "udp"} {
		if err := runCmd("iptables", "-t", "nat", "-A", "CRASHPANEL",
			"-p", proto, "--dport", "53", "-j", "RETURN"); err != nil {
			return err
		}
	}

	// 跳过保留地址
	for _, ip := range reservedIPv4 {
		if err := runCmd("iptables", "-t", "nat", "-A", "CRASHPANEL",
			"-d", ip, "-j", "RETURN"); err != nil {
			return err
		}
	}

	// MAC 黑名单
	mode := filterMode(rules)
	if mode == "blacklist" {
		for _, mac := range macRules(rules) {
			if err := runCmd("iptables", "-t", "nat", "-A", "CRASHPANEL",
				"-m", "mac", "--mac-source", mac, "-j", "RETURN"); err != nil {
				return err
			}
		}
		for _, ip := range ipRules(rules) {
			if err := runCmd("iptables", "-t", "nat", "-A", "CRASHPANEL",
				"-s", ip, "-j", "RETURN"); err != nil {
				return err
			}
		}
	}

	// 根据 redir_mod 决定跳转目标
	switch s.RedirMod {
	case "Tproxy":
		// mangle 表 TPROXY
		if err := runCmd("iptables", "-t", "mangle", "-N", "CRASHPANEL"); err != nil {
			return err
		}
		if err := runCmd("iptables", "-t", "mangle", "-A", "CRASHPANEL",
			"-p", "tcp", "-j", "TPROXY",
			"--on-port", tproxyPort, "--tproxy-mark", routingMark); err != nil {
			return err
		}
		if err := runCmd("iptables", "-t", "mangle", "-A", "CRASHPANEL",
			"-p", "udp", "-j", "TPROXY",
			"--on-port", tproxyPort, "--tproxy-mark", routingMark); err != nil {
			return err
		}
		if err := runCmd("iptables", "-t", "mangle", "-A", "PREROUTING",
			"-j", "CRASHPANEL"); err != nil {
			return err
		}
	default: // Redir
		if err := runCmd("iptables", "-t", "nat", "-A", "CRASHPANEL",
			"-p", "tcp", "-j", "REDIRECT", "--to-ports", redirPort); err != nil {
			return err
		}
		if err := runCmd("iptables", "-t", "nat", "-A", "PREROUTING",
			"-j", "CRASHPANEL"); err != nil {
			return err
		}
		if err := runCmd("iptables", "-t", "nat", "-A", "OUTPUT",
			"-j", "CRASHPANEL"); err != nil {
			return err
		}
	}

	return nil
}
