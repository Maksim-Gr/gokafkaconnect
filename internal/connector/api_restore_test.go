package connector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestoreConnectorConfigs(t *testing.T) {
	client, rec := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	configs := map[string]map[string]string{
		"alpha": {"name": "alpha", "connector.class": "com.example.Alpha"},
	}

	results := RestoreConnectorConfigs(context.Background(), client, configs)
	require.Len(t, results, 1)
	assert.Equal(t, RestoreResult{Name: "alpha"}, results[0])
	assert.Equal(t, http.MethodPut, rec.method)
	assert.Equal(t, "/connectors/alpha/config", rec.path)

	var sent map[string]string
	require.NoError(t, json.Unmarshal(rec.body, &sent))
	assert.Equal(t, "com.example.Alpha", sent["connector.class"])
}

func TestRestoreConnectorConfigs_ContinuesPastFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/connectors/bad/") {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("boom"))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	client := NewClient(srv.URL)

	configs := map[string]map[string]string{
		"zulu":  {"name": "zulu"},
		"bad":   {"name": "bad"},
		"alpha": {"name": "alpha"},
	}

	results := RestoreConnectorConfigs(context.Background(), client, configs)
	require.Len(t, results, 3)

	// Sorted name order, and the failure does not stop the rest.
	assert.Equal(t, "alpha", results[0].Name)
	assert.Empty(t, results[0].Error)
	assert.Equal(t, "bad", results[1].Name)
	assert.Contains(t, results[1].Error, "boom")
	assert.Equal(t, "zulu", results[2].Name)
	assert.Empty(t, results[2].Error)
}

func TestRestoreConnectorConfigs_CanceledContext(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	client := NewClient(srv.URL)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	configs := map[string]map[string]string{
		"alpha": {"name": "alpha"},
		"beta":  {"name": "beta"},
	}

	results := RestoreConnectorConfigs(ctx, client, configs)
	require.Len(t, results, 2)
	for _, r := range results {
		assert.Contains(t, r.Error, "context canceled")
	}
	assert.Zero(t, calls, "canceled context must not contact the API")
}
