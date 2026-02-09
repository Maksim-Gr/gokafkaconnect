package connector

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

func NewClient(kafkaConnectURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(kafkaConnectURL, "/"),
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
}

func (c *Client) doRequest(method, path string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequest(
		method,
		c.baseURL+path,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close() //nolint:errcheck

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respBody, resp.StatusCode, nil
}

func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
