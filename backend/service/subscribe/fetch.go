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

// Fetch 根据 provider 类型分发拉取逻辑
func Fetch(sub *models.Subscription) (*FetchResult, error) {
	switch sub.LinkType {
	case models.LinkTypeFile:
		return fetchFile(sub)
	case models.LinkTypeURI:
		return fetchURI(sub)
	default: // url
		return fetchURL(sub)
	}
}

// fetchURL 订阅链接，经 subconverter 转换后保存
func fetchURL(sub *models.Subscription) (*FetchResult, error) {
	fetchAddr, err := buildConvertURL(sub)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", fetchAddr, nil)
	if err != nil {
		return nil, err
	}

	ua := sub.UserAgent
	if ua == "" {
		ua = "clash.meta/mihomo"
	}
	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, err
	}

	filePath := providerPath(sub)
	if err := os.MkdirAll(providerDir(), 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filePath, body, 0644); err != nil {
		return nil, err
	}

	return &FetchResult{NodeCount: countNodes(body), FilePath: filePath}, nil
}

// fetchFile 本地文件，直接确认存在即可
func fetchFile(sub *models.Subscription) (*FetchResult, error) {
	if sub.FilePath == "" {
		return nil, fmt.Errorf("file path is empty")
	}
	data, err := os.ReadFile(sub.FilePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %v", err)
	}
	return &FetchResult{NodeCount: countNodes(data), FilePath: sub.FilePath}, nil
}

// fetchURI 单节点分享链接，写入 uri_group 文件
func fetchURI(sub *models.Subscription) (*FetchResult, error) {
	if sub.URL == "" {
		return nil, fmt.Errorf("uri is empty")
	}
	if err := os.MkdirAll(providerDir(), 0755); err != nil {
		return nil, err
	}
	// uri_group 文件：每行一个 URI
	filePath := providerDir() + "/uri_group"
	// 追加或创建
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 读取已有内容，避免重复写入
	existing, _ := io.ReadAll(f)
	line := fmt.Sprintf("%s %s\n", sub.Name, sub.URL)
	if !strings.Contains(string(existing), sub.URL) {
		if _, err := f.WriteString(line); err != nil {
			return nil, err
		}
	}
	return &FetchResult{NodeCount: 1, FilePath: filePath}, nil
}

// ProviderFilePath 返回给定 sub 的 provider 文件路径（供外部使用）
func ProviderFilePath(sub *models.Subscription) string {
	if sub.LinkType == models.LinkTypeFile {
		return sub.FilePath
	}
	if sub.LinkType == models.LinkTypeURI {
		return providerDir() + "/uri_group"
	}
	return providerPath(sub)
}

func providerDir() string {
	return "/etc/crashpanel/core/providers"
}

func providerPath(sub *models.Subscription) string {
	return fmt.Sprintf("%s/sub_%d.yaml", providerDir(), sub.ID)
}

func buildConvertURL(sub *models.Subscription) (string, error) {
	// 没有配置 subconverter，直接用原链接
	if sub.SubConverterURL == "" {
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

func countNodes(data []byte) int {
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "- name:") || strings.HasPrefix(t, "-  name:") {
			count++
		}
	}
	return count
}
