package system

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Info struct {
	OS        string  `json:"os"`
	Arch      string  `json:"arch"`
	CPUUsage  float64 `json:"cpu_usage"`
	MemTotal  uint64  `json:"mem_total"`
	MemUsed   uint64  `json:"mem_used"`
	MemFree   uint64  `json:"mem_free"`
	DiskTotal uint64  `json:"disk_total"`
	DiskUsed  uint64  `json:"disk_used"`
	GoVersion string  `json:"go_version"`
}

func GetInfo() Info {
	info := Info{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
	}

	memInfo := readMemInfo()
	info.MemTotal = memInfo["MemTotal"] * 1024
	info.MemFree = memInfo["MemFree"] * 1024
	info.MemUsed = info.MemTotal - info.MemFree

	return info
}

func readMemInfo() map[string]uint64 {
	result := make(map[string]uint64)
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return result
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSuffix(parts[0], ":")
		val, err := strconv.ParseUint(parts[1], 10, 64)
		if err == nil {
			result[key] = val
		}
	}
	return result
}

// FormatBytes 格式化字节数为可读字符串
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
