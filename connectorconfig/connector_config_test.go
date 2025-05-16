package connectorconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"testing"
	"time"
)

func WaitForKafkaConnectStartUp(t *testing.T, baseURL string, timeout time.Duration) {
	t.Helper()

	type Plugin struct {
		Class string `json:"class"`
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(fmt.Sprintf("%s/connector-plugins", baseURL))
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var plugins []Plugin
			if json.NewDecoder(resp.Body).Decode(&plugins) == nil {
				if len(plugins) > 0 {
					return
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
	t.Fatalf("Kafka Connect container is not ready within %s", deadline)

}

type KafkaConnectContainer struct {
	Container testcontainers.Container
	URL       string
}

func setupKafkaConnect(t *testing.T) KafkaConnectContainer {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-kafka-connect:7.5.0",
		ExposedPorts: []string{"8083/tcp"},
		Env: map[string]string{
			"CONNECT_BOOTSTRAP_SERVERS":         "dummy:9092",
			"CONNECT_REST_PORT":                 "8083",
			"CONNECT_GROUP_ID":                  "quickstart",
			"CONNECT_CONFIG_STORAGE_TOPIC":      "docker-connect-configs",
			"CONNECT_OFFSET_STORAGE_TOPIC":      "docker-connect-offsets",
			"CONNECT_STATUS_STORAGE_TOPIC":      "docker-connect-status",
			"CONNECT_KEY_CONVERTER":             "org.apache.kafka.connect.json.JsonConverter",
			"CONNECT_VALUE_CONVERTER":           "org.apache.kafka.connect.json.JsonConverter",
			"CONNECT_INTERNAL_KEY_CONVERTER":    "org.apache.kafka.connect.json.JsonConverter",
			"CONNECT_INTERNAL_VALUE_CONVERTER":  "org.apache.kafka.connect.json.JsonConverter",
			"CONNECT_LOG4J_ROOT_LOGLEVEL":       "INFO",
			"CONNECT_PLUGIN_PATH":               "/usr/share/java",
			"CONNECT_REST_ADVERTISED_HOST_NAME": "localhost",
		},
		WaitingFor: wait.ForHTTP("/").WithPort("8083/tcp").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "8083")
	require.NoError(t, err)
	url := fmt.Sprintf("http://%s:%s", host, port)
	WaitForKafkaConnectStartUp(t, url, 30*time.Second)

	return KafkaConnectContainer{
		Container: container,
		URL:       url,
	}
}

func TestSubmitListAndListStatuses(t *testing.T) {
	kc := setupKafkaConnect(t)

	connectorConfig := `{
	"name": "test-connector",
	"config": {
		"connector.class": "FileStreamSink",
		"tasks.max": "1",
		"file": "/tmp/test.sink.txt",
		"topics": "test-topic"
		}
	}`

	err := SubmitConnector(connectorConfig, kc.URL)
	require.NoError(t, err)

	connectors, err := ListConnectors(kc.URL)
	require.NoError(t, err)

	require.Contains(t, connectors, "test-connector")
}
