package connector

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectorOperations(t *testing.T) {
	kc := setupKafkaConnect(t)
	connectorNames := []string{
		"test-connector-1",
		"test-connector-2",
		"test-connector-to-delete",
	}

	for _, name := range connectorNames {
		err := DeleteConnector(kc.URL, name)
		if err != nil {
			return
		}
	}
	for i, name := range connectorNames {
		topicLetter := string(rune('A' + i))
		connectorConfig := fmt.Sprintf(`{
			"name": "%s",
			"config": {
				"connector.class": "org.apache.kafka.connect.tools.MockSinkConnector",
				"tasks.max": "1",
				"topics": "test-topic-%s"
			}
		}`, name, topicLetter)

		err := SubmitConnector(connectorConfig, kc.URL)
		require.NoError(t, err, fmt.Sprintf("failed to create connector  %s", name))
	}

	connectors, err := ListConnectors(kc.URL)
	require.NoError(t, err)
	for _, name := range connectorNames {
		require.Contains(t, connectors, name, fmt.Sprintf("Connector %s should be in the list", name))
	}

	statuses, err := ListConnectorStatuses(kc.URL)
	require.NoError(t, err)
	for _, name := range connectorNames {
		status, ok := statuses[name]
		require.True(t, ok, fmt.Sprintf("Connector %s should be in the list", name))
		require.Equal(t, name, status.Name)
	}

	tempFile := os.TempDir() + "/test-connector-config.json"
	err = DumpConnectorConfig(kc.URL, connectorNames, tempFile)
	require.NoError(t, err)
	require.FileExists(t, tempFile)

	defer os.Remove(tempFile)

	connectorToDelete := connectorNames[2]
	err = DeleteConnector(kc.URL, connectorToDelete)
	require.NoError(t, err, fmt.Sprintf("Failed to delete connector %s", connectorToDelete))

	connectors, err = ListConnectors(kc.URL)
	require.NoError(t, err)
	require.NotContains(t, connectors, connectorToDelete, fmt.Sprintf("Deleted connector %s should not be in the list", connectorToDelete))

	err = DeleteConnector(kc.URL, "non-existent-connector")
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to delete connector")

	for _, name := range connectorNames[:2] {
		err := DeleteConnector(kc.URL, name)
		require.NoError(t, err, fmt.Sprintf("Failed to clean up connector %s", name))
	}

	connectors, err = ListConnectors(kc.URL)
	require.NoError(t, err)
	for _, name := range connectorNames {
		require.NotContains(t, connectors, name, fmt.Sprintf("Connector %s should be deleted", name))
	}

}
