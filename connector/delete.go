package connector

import (
	"fmt"
	"io"
	"net/http"
)

// DeleteConnector delete connector from Kafka Connect API
func DeleteConnector(kafkaConnectURL string, connector string) error {
	url := fmt.Sprintf("%s/connector/%s", kafkaConnectURL, connector)
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
