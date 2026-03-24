package task

import (
	"bytes"
	"crashpanel/database"
	"crashpanel/models"
	"os/exec"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	mu      sync.Mutex
	cron    *cron.Cron
	entries map[uint]cron.EntryID
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:    cron.New(),
		entries: make(map[uint]cron.EntryID),
	}
}

// Start 启动调度器并加载所有已启用任务
func (s *Scheduler) Start() {
	var tasks []models.Task
	database.DB.Where("enabled = ?", true).Find(&tasks)
	for i := range tasks {
		s.Add(&tasks[i])
	}
	s.cron.Start()
}

// Add 注册一个任务到 cron
func (s *Scheduler) Add(t *models.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果已存在则先移除
	if id, ok := s.entries[t.ID]; ok {
		s.cron.Remove(id)
	}

	taskID := t.ID
	command := t.Command

	entryID, err := s.cron.AddFunc(t.Cron, func() {
		log, code := runCommand(command)
		database.DB.Model(&models.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
			"last_run":  time.Now(),
			"last_log":  log,
			"last_code": code,
		})
	})

	if err == nil {
		s.entries[t.ID] = entryID
	}
}

// Remove 移除任务
func (s *Scheduler) Remove(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, ok := s.entries[id]; ok {
		s.cron.Remove(entryID)
		delete(s.entries, id)
	}
}

// RunNow 立即执行任务（同步）
func (s *Scheduler) RunNow(t *models.Task) (string, int) {
	log, code := runCommand(t.Command)
	database.DB.Model(t).Updates(map[string]interface{}{
		"last_run":  time.Now(),
		"last_log":  log,
		"last_code": code,
	})
	return log, code
}

func runCommand(command string) (string, int) {
	cmd := exec.Command("sh", "-c", command)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			code = -1
		}
	}

	output := buf.String()
	// 截断超长日志
	if len(output) > 4096 {
		output = output[:4096] + "\n... (truncated)"
	}
	return output, code
}
