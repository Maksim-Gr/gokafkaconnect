package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// KeysFromMap extracts and returns a slice of keys from the given map.
func KeysFromMap(m map[string]string) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func ToPrettyJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gokafkaconnect", "config.yaml"), nil
}

func ValidateURL(input string) error {
	if input == "" {
		return errors.New("URL cannot be empty")
	}

	// Reject explicit non-http/https schemes (e.g. ftp://)
	if strings.Contains(input, "://") {
		if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			return errors.New("URL scheme must be http or https")
		}
	}

	testURL := input
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		testURL = "http://" + input
	}

	parsed, err := url.ParseRequestURI(testURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	hostname := parsed.Hostname()
	if hostname == "" {
		return errors.New("URL must contain a host (e.g. localhost:8083 or example.com)")
	}

	// Reject ambiguous single-character bare hostnames (e.g. "d")
	// while still allowing service names like "kafkaconnect:8083"
	if len(hostname) <= 1 && !strings.Contains(hostname, ".") && hostname != "localhost" {
		return fmt.Errorf("invalid host: %q", hostname)
	}

	return nil
}

// SaveConfig saves config to file
func SaveConfig(cfg RestAPIConfig, configPath string) error {
	// Create a directory if not exists
	err := os.MkdirAll(filepath.Dir(configPath), 0o700)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o600)
}

func LoadConfig() (RestAPIConfig, error) {
	var cfg RestAPIConfig
	configPath, err := GetConfigPath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}
