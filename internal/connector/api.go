package connector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func ListConnectors(kafkaConnectURL string) ([]string, error) {
	url := fmt.Sprintf("%s/connectors", kafkaConnectURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, _ := io.ReadAll(resp.Body)
		var connectors []string
		if err := json.Unmarshal(body, &connectors); err != nil {
			return nil, err
		}
		return connectors, nil
	}
	body, _ := io.ReadAll(resp.Body)
	return nil, fmt.Errorf("failed to list connector: %s", string(body))
}

func ListConnectorStatuses(kafkaConnectURL string) (ConnectorsStatusResponse, error) {
	url := fmt.Sprintf("%s/connectors?expand=status", kafkaConnectURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var result ConnectorsStatusResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse status response: %w", err)
		}
		return result, nil
	}

	return nil, fmt.Errorf("failed to list connector statuses: %s", string(body))
}

// DeleteConnector delete connector from Kafka Connect API
func DeleteConnector(kafkaConnectURL string, connector string) error {
	url := fmt.Sprintf("%s/connectors/%s", kafkaConnectURL, connector)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	if resp.StatusCode == 409 {
		return fmt.Errorf("failed to delete connector %s: a rebalance is in process", connector)
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to delete connector: %s", string(body))
}

func SubmitConnector(configJson string, kafkaConnectURL string) error {

	url := fmt.Sprintf("%s/connectors", kafkaConnectURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(configJson)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to submit connector configuration: %s", string(body))
}

func BackupConnectorConfig(kafkaConnectURL string, connectors []string, outputDir string) (string, error) {
	dumpConfig := make(map[string]map[string]interface{})

	for _, name := range connectors {
		url := fmt.Sprintf("%s/connectors/%s/config", kafkaConnectURL, name)

		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("failed to connect to %s: %s", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("failed to connect to %s: %s", url, string(body))
		}

		var config map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
			return "", fmt.Errorf("failed to decode config for %s: %w", name, err)
		}

		dumpConfig[name] = config
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	outputFile := filepath.Join(outputDir, fmt.Sprintf("config_%s.json", timestamp))

	file, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dumpConfig); err != nil {
		return "", fmt.Errorf("failed to encode config: %w", err)
	}

	return outputFile, nil
}

func GetConnectorConfig(kafkaConnectURL, name string) (string, error) {
	url := fmt.Sprintf("%s/connectors/%s/config", kafkaConnectURL, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return string(body), nil
	}
	return "", fmt.Errorf("failed to get connector config: %s", string(body))

}
