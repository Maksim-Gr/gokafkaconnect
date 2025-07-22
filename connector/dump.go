package connector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DumpConnectorConfig(kafkaConnectURL string, connectors []string, outPutFile string) error {
	dumpConfig := make(map[string]map[string]interface{})

	for _, name := range connectors {
		url := fmt.Sprintf("%s/connector/%s", kafkaConnectURL, name)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %s", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %s", url, err)
		}
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to connect to %s: %s", url, string(body))
		}

		var config map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
			return fmt.Errorf("failed to decode config: %s: %w", name, err)
		}
		dumpConfig[name] = config
	}
	file, err := os.Create(outPutFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	defer file.Close() //nolint:errcheck

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(dumpConfig); err != nil {
		return fmt.Errorf("failed to encode config: %s", err)
	}
	return nil
}
