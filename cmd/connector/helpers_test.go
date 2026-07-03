package connector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Maksim-Gr/kkon/internal/connector"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectorHealthy(t *testing.T) {
	tests := []struct {
		name   string
		status connector.Status
		want   bool
	}{
		{
			name: "connector and tasks running",
			status: connector.Status{
				Connector: connector.ConnectorState{State: "RUNNING"},
				Tasks:     []connector.TaskState{{ID: 0, State: "RUNNING"}, {ID: 1, State: "RUNNING"}},
			},
			want: true,
		},
		{
			name: "running with no tasks",
			status: connector.Status{
				Connector: connector.ConnectorState{State: "RUNNING"},
			},
			want: true,
		},
		{
			name: "connector failed",
			status: connector.Status{
				Connector: connector.ConnectorState{State: "FAILED"},
				Tasks:     []connector.TaskState{{ID: 0, State: "RUNNING"}},
			},
			want: false,
		},
		{
			name: "task failed",
			status: connector.Status{
				Connector: connector.ConnectorState{State: "RUNNING"},
				Tasks:     []connector.TaskState{{ID: 0, State: "RUNNING"}, {ID: 1, State: "FAILED"}},
			},
			want: false,
		},
		{
			name: "connector paused",
			status: connector.Status{
				Connector: connector.ConnectorState{State: "PAUSED"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, connectorHealthy(tt.status))
		})
	}
}

// newStatusClient returns a client whose status endpoint replies with the JSON
// produced by body(callCount) on each call, so tests can simulate a connector
// that changes state over successive polls.
func newStatusClient(t *testing.T, body func(call int) (int, string)) *connector.Client {
	t.Helper()
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		code, payload := body(calls)
		calls++
		w.WriteHeader(code)
		_, _ = w.Write([]byte(payload))
	}))
	t.Cleanup(srv.Close)
	return connector.NewClient(srv.URL)
}

func TestWaitForConnectorRunning_BecomesHealthy(t *testing.T) {
	client := newStatusClient(t, func(call int) (int, string) {
		if call < 2 {
			return http.StatusOK, `{"name":"alpha","connector":{"state":"RESTARTING"},"tasks":[{"id":0,"state":"RESTARTING"}]}`
		}
		return http.StatusOK, `{"name":"alpha","connector":{"state":"RUNNING"},"tasks":[{"id":0,"state":"RUNNING"}]}`
	})

	status, healthy := waitForConnectorRunning(context.Background(), client, "alpha", 5, 0)
	require.True(t, healthy)
	assert.Equal(t, "RUNNING", status.Connector.State)
}

func TestWaitForConnectorRunning_NeverHealthy(t *testing.T) {
	client := newStatusClient(t, func(_ int) (int, string) {
		return http.StatusOK, `{"name":"alpha","connector":{"state":"FAILED"},"tasks":[]}`
	})

	status, healthy := waitForConnectorRunning(context.Background(), client, "alpha", 3, 0)
	assert.False(t, healthy)
	assert.Equal(t, "FAILED", status.Connector.State)
	assert.Equal(t, "alpha", status.Name)
}

func TestFieldsToPrompt(t *testing.T) {
	res := connector.ConfigValidationResponse{
		Configs: []connector.ConfigEntry{
			{Definition: connector.ConfigDefinition{Name: "connector.class", Required: true}},
			{Definition: connector.ConfigDefinition{Name: "name", Required: true}},
			{Definition: connector.ConfigDefinition{Name: "connection.url", Required: true, Documentation: "URL."}},
			{Definition: connector.ConfigDefinition{Name: "auto.create", Required: false, DefaultValue: false}},
			{Definition: connector.ConfigDefinition{Name: "tasks.max", Required: true, DefaultValue: float64(1)}},
			{
				Definition: connector.ConfigDefinition{Name: "topics", Required: true},
				Value:      connector.ConfigValue{Name: "topics", Value: "bad topic", Errors: []string{"invalid"}},
			},
			{
				Definition: connector.ConfigDefinition{Name: "connection.password", Required: true},
				Value:      connector.ConfigValue{Name: "connection.password"},
			},
			{
				Definition: connector.ConfigDefinition{Name: "insert.mode", Required: true},
				Value:      connector.ConfigValue{Name: "insert.mode", RecommendedValues: []any{"insert", "upsert"}},
			},
		},
	}

	fields := fieldsToPrompt(res)
	names := make([]string, 0, len(fields))
	for _, f := range fields {
		names = append(names, f.Name)
	}
	// connector.class and name are wizard-managed and always skipped;
	// non-required fields (auto.create) are skipped; every other required
	// field is included, even ones with a default (tasks.max) — the default
	// becomes the suggested answer instead of being silently accepted.
	assert.Equal(t, []string{"connection.url", "tasks.max", "topics", "connection.password", "insert.mode"}, names)

	assert.Equal(t, "URL.", fields[0].Doc)
	assert.Equal(t, "1", fields[1].Default, "required field with a default is suggested, not silently accepted")
	assert.Equal(t, "bad topic", fields[2].Default, "errored field re-prompts with its current value")
	assert.Equal(t, []string{"invalid"}, fields[2].Errors)
	assert.True(t, fields[3].Secret, "password fields must be masked")
	assert.Equal(t, []string{"insert", "upsert"}, fields[4].Recommended)
}

func TestIsSecretField(t *testing.T) {
	assert.True(t, isSecretField("connection.password"))
	assert.True(t, isSecretField("aws.secret.access.key"))
	assert.False(t, isSecretField("connection.url"))
}

func TestRunBulk_AllSucceed(t *testing.T) {
	var attempted []string
	err := runBulk(context.Background(), false, "pause", []string{"a", "b"}, func(_ context.Context, name string) error {
		attempted = append(attempted, name)
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, attempted)
}

func TestRunBulk_ContinuesPastFailure(t *testing.T) {
	var attempted []string
	err := runBulk(context.Background(), false, "pause", []string{"a", "bad", "c"}, func(_ context.Context, name string) error {
		attempted = append(attempted, name)
		if name == "bad" {
			return assert.AnError
		}
		return nil
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to pause 1 of 3 connector(s)")
	assert.Equal(t, []string{"a", "bad", "c"}, attempted, "all names must still be attempted")
}

func TestRunBulk_CanceledContextSkipsRemaining(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var attempted []string
	err := runBulk(ctx, false, "pause", []string{"a", "b"}, func(_ context.Context, name string) error {
		attempted = append(attempted, name)
		return nil
	})
	require.Error(t, err)
	assert.Empty(t, attempted, "canceled context must not invoke the operation")
}

func TestWaitForConnectorRunning_AllErrors(t *testing.T) {
	client := newStatusClient(t, func(_ int) (int, string) {
		return http.StatusInternalServerError, "boom"
	})

	status, healthy := waitForConnectorRunning(context.Background(), client, "alpha", 3, 0)
	assert.False(t, healthy)
	assert.Empty(t, status.Name)
}
