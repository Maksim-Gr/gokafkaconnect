package connector

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

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
