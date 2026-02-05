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

// ToJSON takes a map of string key-value pairs and returns a pretty-printed JSON string.
func ToJSON(config map[string]string) (string, error) {
	out, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func ToPrettyJSON(v interface{}) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func GetConfigPath() (string, error) {
	exe, err := getExecutablePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(exe, "config.yaml"), nil
}

func ValidateURL(input string) error {
	if input == "" {
		return errors.New("URL cannot be empty")
	}

	testURL := input
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		testURL = "http://" + input
	}

	parsed, err := url.ParseRequestURI(testURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	if parsed.Host == "" {
		return errors.New("URL must contain a host (e.g. localhost:8083 or example.com)")
	}

	host := parsed.Hostname()
	if !strings.Contains(host, ".") && !strings.Contains(host, ":") && host != "localhost" {
		return fmt.Errorf("invalid host: %s (must include '.' or ':' or be 'localhost')", host)
	}

	return nil
}

// SaveConfig saves config to file
func SaveConfig(cfg RestAPIConfig, configPath string) error {
	// Create a directory if not exists
	err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
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
