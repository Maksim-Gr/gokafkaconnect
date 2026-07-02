package util

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Maksim-Gr/kkon/internal/connector"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newListClient(t *testing.T, status int, body string) *connector.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return connector.NewClient(srv.URL)
}

func TestResolveConnectorName_FlagValuePassthrough(t *testing.T) {
	name, err := ResolveConnectorName(context.Background(), nil, "my-connector")
	require.NoError(t, err)
	assert.Equal(t, "my-connector", name)
}

func TestResolveConnectorName_NoConnectors(t *testing.T) {
	client := newListClient(t, http.StatusOK, "[]")
	_, err := ResolveConnectorName(context.Background(), client, "")
	assert.ErrorIs(t, err, ErrNothingToDo)
}

func TestResolveConnectorName_ListError(t *testing.T) {
	client := newListClient(t, http.StatusInternalServerError, "boom")
	_, err := ResolveConnectorName(context.Background(), client, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to list connectors")
}

func TestResolveConnectorNames_ArgsPassthrough(t *testing.T) {
	names, err := ResolveConnectorNames(context.Background(), nil, []string{"a", "b"}, false)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, names)
}

func TestResolveConnectorNames_ArgsAndAllConflict(t *testing.T) {
	_, err := ResolveConnectorNames(context.Background(), nil, []string{"a"}, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot combine --all with connector names")
}

func TestResolveConnectorNames_AllSorted(t *testing.T) {
	client := newListClient(t, http.StatusOK, `["zulu","alpha"]`)
	names, err := ResolveConnectorNames(context.Background(), client, nil, true)
	require.NoError(t, err)
	assert.Equal(t, []string{"alpha", "zulu"}, names)
}

func TestResolveConnectorNames_AllEmpty(t *testing.T) {
	client := newListClient(t, http.StatusOK, "[]")
	_, err := ResolveConnectorNames(context.Background(), client, nil, true)
	assert.ErrorIs(t, err, ErrNothingToDo)
}

func TestResolveTaskID_FlagValueFound(t *testing.T) {
	client := newListClient(t, http.StatusOK, `[{"id":{"connector":"alpha","task":0}},{"id":{"connector":"alpha","task":1}}]`)
	id, err := ResolveTaskID(context.Background(), client, "alpha", 1, false)
	require.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestResolveTaskID_FlagValueNotFound(t *testing.T) {
	client := newListClient(t, http.StatusOK, `[{"id":{"connector":"alpha","task":0}}]`)
	_, err := ResolveTaskID(context.Background(), client, "alpha", 7, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task 7 not found")
}

func TestResolveTaskID_NoTasks(t *testing.T) {
	client := newListClient(t, http.StatusOK, "[]")
	_, err := ResolveTaskID(context.Background(), client, "alpha", -1, false)
	assert.ErrorIs(t, err, ErrNothingToDo)
}

func TestResolveTaskID_DryRunSkipsLookup(t *testing.T) {
	id, err := ResolveTaskID(context.Background(), nil, "alpha", 3, true)
	require.NoError(t, err)
	assert.Equal(t, 3, id)
}
