package core

import (
	"bufio"
	"crashpanel/database"
	"crashpanel/models"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Status struct {
	Running   bool      `json:"running"`
	PID       int       `json:"pid"`
	CoreType  string    `json:"core_type"`
	StartTime time.Time `json:"start_time"`
	Uptime    string    `json:"uptime"`
}

type TrafficStats struct {
	Upload   int64 `json:"upload"`
	Download int64 `json:"download"`
}

type Manager struct {
	mu        sync.Mutex
	cmd       *exec.Cmd
	running   bool
	pid       int
	startTime time.Time
	logBuf    []string
	maxLog    int
	upload    int64
	download  int64
}

func NewManager() *Manager {
	return &Manager{
		maxLog: 1000,
		logBuf: make([]string, 0, 1000),
	}
}

func (m *Manager) getSettings() (*models.Settings, error) {
	var s models.Settings
	if err := database.DB.First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (m *Manager) buildCmd(s *models.Settings) (*exec.Cmd, error) {
	if _, err := os.Stat(s.CorePath); err != nil {
		return nil, fmt.Errorf("core binary not found: %s", s.CorePath)
	}

	if err := os.MkdirAll(s.CoreWorkDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work dir: %v", err)
	}

	var cmd *exec.Cmd
	switch s.CoreType {
	case "singbox":
		cmd = exec.Command(s.CorePath, "run", "-D", s.CoreWorkDir, "-C", s.CoreWorkDir+"/jsons")
	default: // meta / mihomo
		cmd = exec.Command(s.CorePath, "-d", s.CoreWorkDir, "-f", s.CoreWorkDir+"/config.yaml")
	}

	return cmd, nil
}

func (m *Manager) appendLog(line string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ts := time.Now().Format("15:04:05")
	m.logBuf = append(m.logBuf, fmt.Sprintf("[%s] %s", ts, line))
	if len(m.logBuf) > m.maxLog {
		m.logBuf = m.logBuf[len(m.logBuf)-m.maxLog:]
	}
}

func (m *Manager) pipeOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		m.appendLog(scanner.Text())
	}
}

func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("core is already running")
	}

	s, err := m.getSettings()
	if err != nil {
		return err
	}

	cmd, err := m.buildCmd(s)
	if err != nil {
		return err
	}

	// 捕获 stdout + stderr
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start core: %v", err)
	}

	m.cmd = cmd
	m.running = true
	m.pid = cmd.Process.Pid
	m.startTime = time.Now()
	m.appendLog(fmt.Sprintf("core started (pid=%d, type=%s)", m.pid, s.CoreType))

	// 异步读取输出
	go m.pipeOutput(stdout)
	go m.pipeOutput(stderr)

	// 监控进程退出
	go func() {
		_ = cmd.Wait()
		m.mu.Lock()
		m.running = false
		m.pid = 0
		m.mu.Unlock()
		m.appendLog("core process exited")
	}()

	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return fmt.Errorf("core is not running")
	}

	if err := m.cmd.Process.Signal(os.Interrupt); err != nil {
		_ = m.cmd.Process.Kill()
	}

	m.running = false
	m.pid = 0
	m.appendLog("core stopped")
	return nil
}

func (m *Manager) Restart() error {
	if m.running {
		if err := m.Stop(); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return m.Start()
}

func (m *Manager) Status() Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, _ := m.getSettings()
	coreType := "meta"
	if s != nil {
		coreType = s.CoreType
	}

	uptime := ""
	if m.running {
		d := time.Since(m.startTime)
		uptime = fmt.Sprintf("%02d:%02d:%02d",
			int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
	}

	return Status{
		Running:   m.running,
		PID:       m.pid,
		CoreType:  coreType,
		StartTime: m.startTime,
		Uptime:    uptime,
	}
}

func (m *Manager) GetLog(n int) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n >= len(m.logBuf) {
		result := make([]string, len(m.logBuf))
		copy(result, m.logBuf)
		return result
	}
	result := make([]string, n)
	copy(result, m.logBuf[len(m.logBuf)-n:])
	return result
}

func (m *Manager) GetTraffic() TrafficStats {
	m.mu.Lock()
	defer m.mu.Unlock()
	return TrafficStats{
		Upload:   m.upload,
		Download: m.download,
	}
}
