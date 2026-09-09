package logger

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLogErrorWritesPrivateJSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	LogError(errors.New("cluster unavailable"), map[string]interface{}{"command": "kubectl get pods"})

	logPath := filepath.Join(home, ".nahkoda", "error.log")
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("stat log: %v", err)
	}
	if got := info.Mode().Perm(); got != 0600 {
		t.Errorf("log permission = %o, want 600", got)
	}

	file, err := os.Open(logPath)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	defer file.Close()

	var entry LogEntry
	if err := json.NewDecoder(file).Decode(&entry); err != nil {
		t.Fatalf("decode log: %v", err)
	}
	if entry.Error != "cluster unavailable" {
		t.Errorf("entry error = %q", entry.Error)
	}
	if entry.Context["command"] != "kubectl get pods" {
		t.Errorf("entry context = %#v", entry.Context)
	}
}
