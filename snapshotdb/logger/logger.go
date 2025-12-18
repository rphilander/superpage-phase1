package logger

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Logger struct {
	filePath string
	mu       sync.Mutex
}

type ErrorEntry struct {
	Timestamp int64  `json:"timestamp"`
	Error     string `json:"error"`
	Context   string `json:"context"`
}

func New(filePath string) *Logger {
	return &Logger{
		filePath: filePath,
	}
}

func (l *Logger) LogError(errMsg, context string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := ErrorEntry{
		Timestamp: time.Now().Unix(),
		Error:     errMsg,
		Context:   context,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(data, '\n'))
	return err
}
