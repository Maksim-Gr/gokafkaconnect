package connector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ConnectorStatus struct {
	Name      string `json:"name"`
	Connector struct {
		State    string `json:"state"`
		WorkerID string `json:"worker_id"`
	} `json:"connector"`
	Tasks []struct {
		ID       int    `json:"id"`
		State    string `json:"state"`
		WorkerID string `json:"worker_id"`
	} `json:"tasks"`
	Type string `json:"type"`
}

type ConnectorsStatusResponse map[string]ConnectorStatus

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
