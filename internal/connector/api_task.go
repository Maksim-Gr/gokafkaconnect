package connector

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TaskRef represents an entry returned by GET /connectors/{name}/tasks.
type TaskRef struct {
	Connector string `json:"connector"`
	Task      int    `json:"task"`
}

// TaskStatus represents the response from GET /connectors/{name}/tasks/{id}/status.
type TaskStatus struct {
	ID       int    `json:"id"`
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
	Trace    string `json:"trace,omitempty"`
}

// ListConnectorTasks lists tasks for a connector.
func (c *Client) ListConnectorTasks(connectorName string) ([]TaskRef, error) {
	path := fmt.Sprintf("/connectors/%s/tasks", connectorName)

	body, status, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if !isSuccess(status) {
		return nil, fmt.Errorf("failed to list tasks for %s: %s", connectorName, string(body))
	}

	var tasks []TaskRef
	if err := json.Unmarshal(body, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetConnectorTaskStatus fetches the status of a single task.
func (c *Client) GetConnectorTaskStatus(connectorName string, taskID int) (TaskStatus, error) {
	path := fmt.Sprintf("/connectors/%s/tasks/%d/status", connectorName, taskID)

	body, status, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return TaskStatus{}, err
	}
	if !isSuccess(status) {
		return TaskStatus{}, fmt.Errorf("failed to get task status for %s task %d: %s", connectorName, taskID, string(body))
	}

	var ts TaskStatus
	if err := json.Unmarshal(body, &ts); err != nil {
		return TaskStatus{}, err
	}
	return ts, nil
}

// RestartConnectorTask restarts a single task.
func (c *Client) RestartConnectorTask(connectorName string, taskID int) error {
	path := fmt.Sprintf("/connectors/%s/tasks/%d/restart", connectorName, taskID)

	body, status, err := c.doRequest(http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	if !isSuccess(status) {
		return fmt.Errorf("failed to restart %s task %d: %s", connectorName, taskID, string(body))
	}
	return nil
}
