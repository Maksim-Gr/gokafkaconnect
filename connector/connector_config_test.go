package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
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
			defer resp.Body.Close() //nolint:errcheck
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

	nw, err := network.New(ctx)

	require.NoError(t, err)
	t.Cleanup(func() { _ = nw.Remove(ctx) })

	zooReq := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-zookeeper:7.2.0",
		ExposedPorts: []string{"2181/tcp"},
		Networks:     []string{nw.Name},
		Env: map[string]string{
			"ZOOKEEPER_CLIENT_PORT": "2181",
			"ZOOKEEPER_TICK_TIME":   "2000",
		},
		WaitingFor: wait.ForListeningPort("2181/tcp").WithStartupTimeout(30 * time.Second),
		Name:       "zookeeper",
	}
	zooC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: zooReq,
		Started:          true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = zooC.Terminate(ctx) })

	kafkaReq := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-kafka:7.2.0",
		ExposedPorts: []string{"9092/tcp"},
		Networks:     []string{nw.Name},
		Env: map[string]string{
			"KAFKA_BROKER_ID":                        "1",
			"KAFKA_ZOOKEEPER_CONNECT":                "zookeeper:2181",
			"KAFKA_ADVERTISED_LISTENERS":             "PLAINTEXT://kafka:9092",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
			"KAFKA_LISTENERS":                        "PLAINTEXT://0.0.0.0:9092",
			"KAFKA_LOG_DIRS":                         "/tmp/kafka-logs",
		},
		WaitingFor: wait.ForListeningPort("9092/tcp").WithStartupTimeout(90 * time.Second),
		Name:       "kafka",
	}
	kafkaC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: kafkaReq,
		Started:          true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = kafkaC.Terminate(ctx) })

	connectReq := testcontainers.ContainerRequest{
		Image:        "confluentinc/cp-kafka-connect:7.5.0",
		ExposedPorts: []string{"8083/tcp"},
		Networks:     []string{nw.Name},
		Env: map[string]string{
			"CONNECT_BOOTSTRAP_SERVERS":                 "kafka:9092",
			"CONNECT_REST_PORT":                         "8083",
			"CONNECT_GROUP_ID":                          "quickstart",
			"CONNECT_CONFIG_STORAGE_TOPIC":              "docker-connect-configs",
			"CONNECT_OFFSET_STORAGE_TOPIC":              "docker-connect-offsets",
			"CONNECT_STATUS_STORAGE_TOPIC":              "docker-connect-status",
			"CONNECT_CONFIG_STORAGE_REPLICATION_FACTOR": "1",
			"CONNECT_OFFSET_STORAGE_REPLICATION_FACTOR": "1",
			"CONNECT_STATUS_STORAGE_REPLICATION_FACTOR": "1",
			"CONNECT_KEY_CONVERTER":                     "org.apache.kafka.connect.json.JsonConverter",
			"CONNECT_VALUE_CONVERTER":                   "org.apache.kafka.connect.json.JsonConverter",
			"CONNECT_INTERNAL_KEY_CONVERTER":            "org.apache.kafka.connect.json.JsonConverter",
			"CONNECT_INTERNAL_VALUE_CONVERTER":          "org.apache.kafka.connect.json.JsonConverter",
			"CONNECT_LOG4J_ROOT_LOGLEVEL":               "INFO",
			"CONNECT_PLUGIN_PATH":                       "/usr/share/java,/usr/share/confluent-hub-components",
			"CONNECT_REST_ADVERTISED_HOST_NAME":         "localhost",
		},
		WaitingFor: wait.ForHTTP("/").WithPort("8083/tcp").WithStartupTimeout(120 * time.Second),
		Name:       "kafkaconnect",
	}
	connectC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: connectReq,
		Started:          true,
	})

	require.NoError(t, err)
	t.Cleanup(func() { _ = connectC.Terminate(ctx) })

	connectHost, err := connectC.Host(ctx)
	require.NoError(t, err)
	connectPort, err := connectC.MappedPort(ctx, "8083")
	require.NoError(t, err)

	connectURL := fmt.Sprintf("http://%s:%s", connectHost, connectPort.Port())
	WaitForKafkaConnectStartUp(t, connectURL, 20*time.Second)

	return KafkaConnectContainer{
		Container: connectC,
		URL:       connectURL,
	}
}

func TestSubmitListAndListStatuses(t *testing.T) {
	kc := setupKafkaConnect(t)

	connectorConfig := `{
	"name": "test-connector",
	"config": {
		"connector.class": "org.apache.kafka.connect.tools.MockSinkConnector",
		"tasks.max": "1",
		"topics": "test-topic"
	}
	}`

	err := SubmitConnector(connectorConfig, kc.URL)
	require.NoError(t, err)

	connectors, err := ListConnectors(kc.URL)
	require.NoError(t, err)

	require.Contains(t, connectors, "test-connector")
}

func TestDumpConnectorConfig(t *testing.T) {
	kc := setupKafkaConnect(t)

}
