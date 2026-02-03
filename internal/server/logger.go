package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	logDir string
	mu     sync.RWMutex
	logs   []*RequestLog
	posts  []*PostRecord
}

func NewLogger(logDir string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	return &Logger{
		logDir: logDir,
		logs:   make([]*RequestLog, 0),
	}, nil
}

func (l *Logger) Log(req *RequestLog) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logs = append(l.logs, req)

	l.saveToFile()

	return nil
}

func (l *Logger) GetAllLogs() []*RequestLog {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.logs
}

func (l *Logger) LogPost(post PostRecord) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.posts = append(l.posts, &post)

	l.savePostsToFile()

	return nil
}

func (l *Logger) GetAllPosts() []*PostRecord {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.posts
}

func (l *Logger) GetSummary() map[string]interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	ips := make(map[string]bool)
	postCount := 0

	for _, log := range l.logs {
		ips[log.RemoteAddr] = true
		if log.Method == "POST" {
			postCount++
		}
	}

	return map[string]interface{}{
		"total_requests": len(l.logs),
		"unique_ips":     len(ips),
		"post_requests":  postCount,
		"last_request":   l.getLastRequestTime(),
	}
}

func (l *Logger) SaveSummary() error {
	summaryPath := filepath.Join(l.logDir, "summary.json")

	summary := l.GetSummary()

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal summary: %w", err)
	}

	if err := os.WriteFile(summaryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}

	return nil
}

func (l *Logger) saveToFile() {
	if l.logDir == "" {
		return
	}

	date := time.Now().UTC().Format("20060102")
	filePath := filepath.Join(l.logDir, fmt.Sprintf("requests_%s.json", date))

	data, err := json.MarshalIndent(l.logs, "", "  ")
	if err != nil {
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return
	}
}

func (l *Logger) getLastRequestTime() string {
	if len(l.logs) == 0 {
		return "-"
	}
	return l.logs[len(l.logs)-1].Timestamp.Format(time.RFC3339)
}

func (l *Logger) savePostsToFile() {
	if l.logDir == "" {
		return
	}

	filePath := filepath.Join(l.logDir, "posts.json")

	data, err := json.MarshalIndent(l.posts, "", "  ")
	if err != nil {
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return
	}
}
