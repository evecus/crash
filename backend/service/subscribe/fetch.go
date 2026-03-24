package subscribe

import (
	"crashpanel/models"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type FetchResult struct {
	NodeCount int
	FilePath  string
}

// Fetch 拉取订阅，经过 subconverter 转换后保存到核心工作目录
func Fetch(sub *models.Subscription) (*FetchResult, error) {
	convertedURL, err := buildSubConverterURL(sub)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", convertedURL, nil)
	if err != nil {
		return nil, err
	}
	if sub.UserAgent != "" {
		req.Header.Set("User-Agent", sub.UserAgent)
	} else {
		req.Header.Set("User-Agent", "clash.meta/mihomo")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20)) // 最大 10MB
	if err != nil {
		return nil, err
	}

	// 写入到工作目录
	filePath, err := saveSub(sub, body)
	if err != nil {
		return nil, err
	}

	nodeCount := countNodes(body, sub.Target)

	return &FetchResult{
		NodeCount: nodeCount,
		FilePath:  filePath,
	}, nil
}

func buildSubConverterURL(sub *models.Subscription) (string, error) {
	if sub.SubConverterURL == "" {
		// 直接使用原始订阅链接（不经过 subconverter）
		return sub.URL, nil
	}

	target := sub.Target
	if target == "" {
		target = "clash"
	}

	params := url.Values{}
	params.Set("target", target)
	params.Set("url", sub.URL)
	params.Set("insert", "true")
	params.Set("new_name", "true")
	params.Set("scv", "true")
	params.Set("udp", "true")
	if sub.Include != "" {
		params.Set("include", sub.Include)
	}
	if sub.Exclude != "" {
		params.Set("exclude", sub.Exclude)
	}
	if sub.ConfigURL != "" {
		params.Set("config", sub.ConfigURL)
	}

	base := strings.TrimRight(sub.SubConverterURL, "/")
	return fmt.Sprintf("%s/sub?%s", base, params.Encode()), nil
}

func saveSub(sub *models.Subscription, data []byte) (string, error) {
	// 保存到 /etc/crashpanel/core/providers/ 目录
	dir := "/etc/crashpanel/core/providers"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	ext := "yaml"
	if sub.Target == "singbox" {
		ext = "json"
	}
	filePath := fmt.Sprintf("%s/sub_%d.%s", dir, sub.ID, ext)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}
	return filePath, nil
}

func countNodes(data []byte, target string) int {
	content := string(data)
	count := 0
	if target == "singbox" {
		count = strings.Count(content, `"type": "vmess"`) +
			strings.Count(content, `"type": "shadowsocks"`) +
			strings.Count(content, `"type": "trojan"`) +
			strings.Count(content, `"type": "hysteria"`) +
			strings.Count(content, `"type": "vless"`)
	} else {
		// clash yaml：数 "- name:" 的行数作为近似
		for _, line := range strings.Split(content, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- name:") || strings.HasPrefix(trimmed, "-  name:") {
				count++
			}
		}
	}
	return count
}
