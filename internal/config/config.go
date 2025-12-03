package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents config for setting up KAFKA API URL
type Config struct {
	KafkaConnectURL string `json:"kafka_connect_url"`
}

// LoadConfig loads the app configuration from a file
func LoadConfig() (*Config, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".gokafkacon", "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config.json: %w", err)
	}
	defer file.Close() //nolint:errcheck

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config.json: %w", err)
	}
	return &cfg, nil
}
