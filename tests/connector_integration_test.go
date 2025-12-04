package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	c "gokafkaconnect/internal/connector"

	"github.com/stretchr/testify/require"
)

const mockSinkConfigTemplate = `{
	"name": "%s",
	"config": {
		"connector.class": "org.apache.kafka.connect.tools.MockSinkConnector",
		"tasks.max": "1",
		"topics": "%s"
	}
}`

// setupConnectors creates a set of connectors for testing.
func setupConnectors(t *testing.T, kc *KafkaConnectTestFixture, names []string) {
	t.Helper()

	for i, name := range names {
		topic := fmt.Sprintf("test-topic-%c", 'A'+i)
		config := fmt.Sprintf(mockSinkConfigTemplate, name, topic)

		err := c.SubmitConnector(config, kc.URL)
		require.NoError(t, err, "creating connector: %s", name)
	}
}

func cleanupConnectors(t *testing.T, kc *KafkaConnectTestFixture, names []string) {
	t.Helper()
	for _, name := range names {
		_ = c.DeleteConnector(kc.URL, name)
	}
}

func TestConnectorLifecycle(t *testing.T) {
	kc := KafkaConnectFixture(t)

	connectorNames := []string{
		"test-op-connector-1",
		"test-op-connector-2",
		"test-op-connector-3",
	}

	cleanupConnectors(t, kc, connectorNames)
	defer cleanupConnectors(t, kc, connectorNames)

	t.Run("CreateAndList", func(t *testing.T) {
		setupConnectors(t, kc, connectorNames)

		got, err := c.ListConnectors(kc.URL)
		require.NoError(t, err)

		for _, name := range connectorNames {
			require.Contains(t, got, name)
		}
	})

	t.Run("ListStatuses", func(t *testing.T) {
		// Wait briefly for status propagation.
		time.Sleep(2 * time.Second)

		statuses, err := c.ListConnectorStatuses(kc.URL)
		require.NoError(t, err)

		for _, name := range connectorNames {
			_, ok := statuses[name]
			require.True(t, ok, "status missing for connector: %s", name)
		}
	})

	t.Run("DumpConfig", func(t *testing.T) {
		tempFile := os.TempDir() + "/connector-config.json"
		defer os.Remove(tempFile)

		err := c.DumpConnectorConfig(kc.URL, connectorNames[:2], tempFile)
		require.NoError(t, err)
		require.FileExists(t, tempFile)
	})

	t.Run("DeleteOne", func(t *testing.T) {
		target := connectorNames[2]

		err := c.DeleteConnector(kc.URL, target)
		require.NoError(t, err)

		got, err := c.ListConnectors(kc.URL)
		require.NoError(t, err)
		require.NotContains(t, got, target)
	})

	t.Run("DeleteNonExistentReturnsError", func(t *testing.T) {
		err := c.DeleteConnector(kc.URL, "non-existent-connector")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to delete connector")
	})
}

func TestDumpConnectorConfig_Single(t *testing.T) {
	kc := KafkaConnectFixture(t)

	connectorName := "test-dump-single"
	topic := "test-topic-ds"
	defer cleanupConnectors(t, kc, []string{connectorName})

	config := fmt.Sprintf(mockSinkConfigTemplate, connectorName, topic)

	require.NoError(t, c.SubmitConnector(config, kc.URL))

	tempFile := os.TempDir() + "/connector-dump.json"
	defer os.Remove(tempFile)

	require.NoError(t, c.DumpConnectorConfig(kc.URL, []string{connectorName}, tempFile))
	require.FileExists(t, tempFile)
}

func TestGetConnectorConfig(t *testing.T) {
	kc := KafkaConnectFixture(t)

	connectorName := "test-getconfig"
	topic := "topic-gc"
	defer cleanupConnectors(t, kc, []string{connectorName})

	config := fmt.Sprintf(mockSinkConfigTemplate, connectorName, topic)

	require.NoError(t, c.SubmitConnector(config, kc.URL))

	jsonConfig, err := c.GetConnectorConfig(kc.URL, connectorName)
	require.NoError(t, err)

	require.Contains(t, jsonConfig, "connector.class")
	require.Contains(t, jsonConfig, "MockSinkConnector")
}

func TestSubmitAndList(t *testing.T) {
	kc := KafkaConnectFixture(t)

	connectorName := "test-submit-list"
	topic := "test-topic-sl"
	defer cleanupConnectors(t, kc, []string{connectorName})

	config := fmt.Sprintf(mockSinkConfigTemplate, connectorName, topic)

	require.NoError(t, c.SubmitConnector(config, kc.URL))

	got, err := c.ListConnectors(kc.URL)
	require.NoError(t, err)
	require.Contains(t, got, connectorName)
}
