package connector

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validateFixture is a trimmed but realistic validate-endpoint payload: a
// required field without a default, a field with typed (boolean) default and
// recommended values, and an errored field.
const validateFixture = `{
  "name": "com.example.MySink",
  "error_count": 1,
  "configs": [
    {
      "definition": {
        "name": "connection.url",
        "type": "STRING",
        "required": true,
        "default_value": null,
        "importance": "HIGH",
        "documentation": "JDBC connection URL.",
        "group": "Connection"
      },
      "value": {"name": "connection.url", "value": null, "recommended_values": [], "errors": []}
    },
    {
      "definition": {
        "name": "auto.create",
        "type": "BOOLEAN",
        "required": false,
        "default_value": false,
        "importance": "MEDIUM",
        "documentation": "Auto-create destination tables.",
        "group": "Writes"
      },
      "value": {"name": "auto.create", "value": null, "recommended_values": [true, false], "errors": []}
    },
    {
      "definition": {
        "name": "topics",
        "type": "LIST",
        "required": true,
        "default_value": "",
        "importance": "HIGH",
        "documentation": "Topics to consume.",
        "group": "Common"
      },
      "value": {"name": "topics", "value": "", "recommended_values": [], "errors": ["Missing required configuration \"topics\""]}
    }
  ]
}`

func TestValidateConnectorConfig_ParsesDefinitions(t *testing.T) {
	client, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(validateFixture))
	})

	res, err := client.ValidateConnectorConfig(context.Background(), "com.example.MySink", map[string]string{
		"connector.class": "com.example.MySink",
	})
	require.NoError(t, err)
	assert.Equal(t, "/connector-plugins/com.example.MySink/config/validate", rec.path)
	assert.Equal(t, http.MethodPut, rec.method)

	assert.Equal(t, 1, res.ErrorCount)
	require.Len(t, res.Configs, 3)

	url := res.Configs[0]
	assert.Equal(t, "connection.url", url.Definition.Name)
	assert.True(t, url.Definition.Required)
	assert.Nil(t, url.Definition.DefaultValue)
	assert.Equal(t, "JDBC connection URL.", url.Definition.Documentation)

	auto := res.Configs[1]
	assert.Equal(t, false, auto.Definition.DefaultValue, "typed defaults must survive parsing")
	assert.Equal(t, []any{true, false}, auto.Value.RecommendedValues)

	topics := res.Configs[2]
	require.Len(t, topics.Value.Errors, 1)
	assert.Contains(t, topics.Value.Errors[0], "Missing required configuration")
}
