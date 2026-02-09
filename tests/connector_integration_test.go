package tests

import (
	"fmt"
	"os"
	"path/filepath"
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
func setupConnectors(t *testing.T, client *c.Client, names []string) {
	t.Helper()

	for i, name := range names {
		topic := fmt.Sprintf("test-topic-%c", 'A'+i)
		config := fmt.Sprintf(mockSinkConfigTemplate, name, topic)

		err := client.SubmitConnector(config)
		require.NoError(t, err, "creating connector: %s", name)
	}
}

func cleanupConnectors(t *testing.T, client *c.Client, names []string) {
	t.Helper()
	for _, name := range names {
		_ = client.DeleteConnector(name)
	}
}

func TestConnectorLifecycle(t *testing.T) {
	kc := KafkaConnectFixture(t)
	client := c.NewClient(kc.URL)

	connectorNames := []string{
		"test-op-connector-1",
		"test-op-connector-2",
		"test-op-connector-3",
	}

	cleanupConnectors(t, client, connectorNames)
	defer cleanupConnectors(t, client, connectorNames)

	t.Run("CreateAndList", func(t *testing.T) {
		setupConnectors(t, client, connectorNames)

		got, err := client.ListConnectors()
		require.NoError(t, err)

		for _, name := range connectorNames {
			require.Contains(t, got, name)
		}
	})

	t.Run("ListStatuses", func(t *testing.T) {
		// Wait briefly for status propagation.
		time.Sleep(2 * time.Second)

		statuses, err := client.ListConnectorStatuses()
		require.NoError(t, err)

		for _, name := range connectorNames {
			_, ok := statuses[name]
			require.True(t, ok, "status missing for connector: %s", name)
		}
	})

	t.Run("BackupConnectorConfig", func(t *testing.T) {
		outputDir := os.TempDir()

		backupFile, err := c.BackupConnectorConfig(client, connectorNames[:2], outputDir)
		require.NoError(t, err, "BackupConnectorConfig should not return an error")

		require.FileExists(t, backupFile, "Backup file should exist")
		require.Contains(t, filepath.Base(backupFile), "config_", "Backup file name should contain 'config_' prefix")

		defer func() {
			if err := os.Remove(backupFile); err != nil && !os.IsNotExist(err) {
				t.Logf("failed to remove backup file: %v", err)
			}
		}()
	})

	t.Run("DeleteOne", func(t *testing.T) {
		target := connectorNames[2]

		err := client.DeleteConnector(target)
		require.NoError(t, err)

		got, err := client.ListConnectors()
		require.NoError(t, err)
		require.NotContains(t, got, target)
	})

	t.Run("DeleteNonExistentReturnsError", func(t *testing.T) {
		err := client.DeleteConnector("non-existent-connector")
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to delete connector")
	})
}

func TestGetConnectorConfig(t *testing.T) {
	kc := KafkaConnectFixture(t)
	client := c.NewClient(kc.URL)

	connectorName := "test-getconfig"
	topic := "topic-gc"
	defer cleanupConnectors(t, client, []string{connectorName})

	config := fmt.Sprintf(mockSinkConfigTemplate, connectorName, topic)

	require.NoError(t, client.SubmitConnector(config))

	jsonConfig, err := client.GetConnectorConfig(connectorName)
	require.NoError(t, err)

	require.Contains(t, jsonConfig, "connector.class")
	require.Contains(t, jsonConfig, "MockSinkConnector")
}

func TestSubmitAndList(t *testing.T) {
	kc := KafkaConnectFixture(t)
	client := c.NewClient(kc.URL)

	connectorName := "test-submit-list"
	topic := "test-topic-sl"
	defer cleanupConnectors(t, client, []string{connectorName})

	config := fmt.Sprintf(mockSinkConfigTemplate, connectorName, topic)

	require.NoError(t, client.SubmitConnector(config))

	got, err := client.ListConnectors()
	require.NoError(t, err)
	require.Contains(t, got, connectorName)
}
