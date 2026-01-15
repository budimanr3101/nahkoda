package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	KubectlPath       string        `json:"kubectl_path,omitempty"`
	DefaultNamespace  string        `json:"default_namespace,omitempty"`
	CacheTTL          time.Duration `json:"cache_ttl,omitempty"`
	Timeout           time.Duration `json:"timeout,omitempty"`
	EnableSuggestions bool          `json:"enable_suggestions"`
}

// Default returns a config with sensible defaults
func Default() *Config {
	return &Config{
		KubectlPath:       "", // auto-detect
		DefaultNamespace:  "default",
		CacheTTL:          30 * time.Second,
		Timeout:           30 * time.Second,
		EnableSuggestions: true,
	}
}

// Load reads config from ~/.nahkoda/config.json
// Returns default config if file doesn't exist
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".nahkoda", "config.json")
	file, err := os.Open(configPath)
	if os.IsNotExist(err) {
		// Return default config if not found
		return Default(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply defaults for missing fields
	if config.DefaultNamespace == "" {
		config.DefaultNamespace = "default"
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 30 * time.Second
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Validate before returning
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate checks if the config is valid
func (c *Config) Validate() error {
	// Validate kubectl path if specified
	if c.KubectlPath != "" {
		if _, err := os.Stat(c.KubectlPath); err != nil {
			return fmt.Errorf("kubectl path tidak valid: %w", err)
		}
	}

	// Validate timeout and cache TTL
	if c.Timeout < 0 {
		return fmt.Errorf("timeout tidak boleh negatif")
	}
	if c.CacheTTL < 0 {
		return fmt.Errorf("cache_ttl tidak boleh negatif")
	}

	return nil
}

// Save writes config to ~/.nahkoda/config.json
func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".nahkoda")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// InitDefault creates a default config file at ~/.nahkoda/config.json
func InitDefault() error {
	cfg := Default()
	return cfg.Save()
}
