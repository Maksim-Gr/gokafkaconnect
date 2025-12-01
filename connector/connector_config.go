package connector

import (
	"fmt"
	"io"
	"net/http"
)

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
