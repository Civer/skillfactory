// Package config handles persistent configuration for SkillFactory
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configDir  = ".skillfactory"
	configFile = "config.json"
)

// Config holds persistent user settings
type Config struct {
	SkillsFolder string `json:"skills_folder,omitempty"`
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

// Load reads the config from disk, returns empty config if not found
func Load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &Config{}, nil
	}

	return &cfg, nil
}

// Save writes the config to disk
func (c *Config) Save() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
