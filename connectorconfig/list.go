package connectorconfig

import (
	"fmt"
	"io"
	"net/http"
)

func ListConnectors(kafkaConnectURL string) ([]byte, error) {
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
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, _ := io.ReadAll(resp.Body)
		return body, nil
	}
	body, _ := io.ReadAll(resp.Body)
	return nil, fmt.Errorf("failed to list connectors: %s", string(body))
}
