package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.DefaultNamespace != "default" {
		t.Errorf("Expected default namespace 'default', got %s", cfg.DefaultNamespace)
	}
	if cfg.CacheTTL != 30*time.Second {
		t.Errorf("Expected cache TTL 30s, got %v", cfg.CacheTTL)
	}
	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", cfg.Timeout)
	}
	if !cfg.EnableSuggestions {
		t.Error("Expected suggestions enabled by default")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	// Create temp home
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() should return default config when file missing: %v", err)
	}

	if cfg.DefaultNamespace != "default" {
		t.Errorf("Expected default config, got %+v", cfg)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	// Create temp home with config
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	configDir := filepath.Join(tmpHome, ".nahkoda")
	os.MkdirAll(configDir, 0755)

	testConfig := &Config{
		DefaultNamespace:  "production",
		CacheTTL:          60 * time.Second,
		Timeout:           45 * time.Second,
		EnableSuggestions: false,
	}

	configPath := filepath.Join(configDir, "config.json")
	file, _ := os.Create(configPath)
	json.NewEncoder(file).Encode(testConfig)
	file.Close()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.DefaultNamespace != "production" {
		t.Errorf("Expected namespace 'production', got %s", cfg.DefaultNamespace)
	}
	if cfg.CacheTTL != 60*time.Second {
		t.Errorf("Expected cache TTL 60s, got %v", cfg.CacheTTL)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	configDir := filepath.Join(tmpHome, ".nahkoda")
	os.MkdirAll(configDir, 0755)

	configPath := filepath.Join(configDir, "config.json")
	os.WriteFile(configPath, []byte("invalid json"), 0644)

	_, err := Load()
	if err == nil {
		t.Error("Load() should fail on invalid JSON")
	}
}

func TestValidate_InvalidKubectlPath(t *testing.T) {
	cfg := &Config{
		KubectlPath: "/path/that/does/not/exist",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() should fail for invalid kubectl path")
	}
}

func TestValidate_NegativeTimeout(t *testing.T) {
	cfg := &Config{
		Timeout: -5 * time.Second,
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() should fail for negative timeout")
	}
}

func TestValidate_NegativeCacheTTL(t *testing.T) {
	cfg := &Config{
		CacheTTL: -10 * time.Second,
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() should fail for negative cache TTL")
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := Default()

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Validate() should pass for default config: %v", err)
	}
}

func TestSave(t *testing.T) {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	cfg := &Config{
		DefaultNamespace:  "staging",
		CacheTTL:          45 * time.Second,
		Timeout:           20 * time.Second,
		EnableSuggestions: true,
	}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// Verify file exists
	configPath := filepath.Join(tmpHome, ".nahkoda", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Verify content
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() after Save() failed: %v", err)
	}

	if loaded.DefaultNamespace != "staging" {
		t.Errorf("Expected namespace 'staging', got %s", loaded.DefaultNamespace)
	}
}

func TestInitDefault(t *testing.T) {
	tmpHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	err := InitDefault()
	if err != nil {
		t.Fatalf("InitDefault() failed: %v", err)
	}

	// Verify file exists
	configPath := filepath.Join(tmpHome, ".nahkoda", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Default config file was not created")
	}

	// Verify it's default config
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() after InitDefault() failed: %v", err)
	}

	if loaded.DefaultNamespace != "default" {
		t.Errorf("Expected default namespace, got %s", loaded.DefaultNamespace)
	}
}
