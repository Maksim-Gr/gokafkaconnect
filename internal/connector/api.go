package connector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (c *Client) ListConnectors() ([]string, error) {
	body, status, err := c.doRequest(http.MethodGet, "/connectors", nil)
	if err != nil {
		return nil, err
	}
	if !isSuccess(status) {
		return nil, fmt.Errorf("failed to list connectors: %s", string(body))
	}
	var connectors []string
	if err := json.Unmarshal(body, &connectors); err != nil {
		return nil, err
	}
	return connectors, nil
}

func (c *Client) ListConnectorStatuses() (ConnectorsStatusResponse, error) {
	body, status, err := c.doRequest(http.MethodGet, "/connectors?expand=status", nil)
	if err != nil {
		return nil, err
	}
	if !isSuccess(status) {
		return nil, fmt.Errorf("failed to list connector statuses: %s", string(body))
	}
	var connectorsStatus ConnectorsStatusResponse
	if err := json.Unmarshal(body, &connectorsStatus); err != nil {
		return nil, err
	}
	return connectorsStatus, nil
}

func (c *Client) DeleteConnector(name string) error {
	body, status, err := c.doRequest(http.MethodDelete, fmt.Sprintf("/connectors/%s", name), nil)
	if err != nil {
		return err
	}
	if status == http.StatusConflict {
		return fmt.Errorf("failed to delete connector %s: a rebalance is in process", name)
	}
	if !isSuccess(status) {
		return fmt.Errorf("failed to delete connector %s: %s", name, string(body))
	}
	return nil
}

func (c *Client) SubmitConnector(configJson string) error {
	body, status, err := c.doRequest(http.MethodPost, "/connectors", []byte(configJson))
	if err != nil {
		return err
	}
	if !isSuccess(status) {
		return fmt.Errorf("failed to submit connector configuration: %s", string(body))
	}
	return nil
}

func (c *Client) GetConnectorConfig(name string) (string, error) {
	body, status, err := c.doRequest(
		http.MethodGet,
		"/connectors/"+name+"/config",
		nil,
	)
	if err != nil {
		return "", err
	}
	if !isSuccess(status) {
		return "", fmt.Errorf("failed to get connector config: %s", body)
	}
	return string(body), nil
}

func (c *Client) GetConnectorConfigJSON(name string) (map[string]interface{}, error) {
	body, status, err := c.doRequest(
		http.MethodGet,
		"/connectors/"+name+"/config",
		nil,
	)
	if err != nil {
		return nil, err
	}
	if !isSuccess(status) {
		return nil, fmt.Errorf("failed to get config for %s: %s", name, body)
	}

	var cfg map[string]interface{}
	if err := json.Unmarshal(body, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func BackupConnectorConfig(
	client *Client,
	connectors []string,
	outputDir string,
) (string, error) {

	dumpConfig := make(map[string]map[string]interface{})

	for _, name := range connectors {
		cfg, err := client.GetConnectorConfigJSON(name)
		if err != nil {
			return "", err
		}
		dumpConfig[name] = cfg
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
	defer file.Close() //nolint:errcheck

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dumpConfig); err != nil {
		return "", fmt.Errorf("failed to encode config: %w", err)
	}

	return outputFile, nil
}
