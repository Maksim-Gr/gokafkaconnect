package connector

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to connect to %s: %s", url, string(body))
		}
	}
}
