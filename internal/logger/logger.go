package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Error     string                 `json:"error"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// LogError logs an error to ~/.nahkoda/error.log
func LogError(err error, context map[string]interface{}) {
	if err == nil {
		return
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		// Can't log if we can't find home dir
		return
	}

	logDir := filepath.Join(home, ".nahkoda")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return
	}

	logPath := filepath.Join(logDir, "error.log")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	entry := LogEntry{
		Timestamp: time.Now(),
		Error:     err.Error(),
		Context:   context,
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(entry); err != nil {
		// Ignore encoding errors
		return
	}
}

// LogErrorf logs a formatted error
func LogErrorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	LogError(fmt.Errorf(msg), nil)
}
