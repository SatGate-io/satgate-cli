package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/SatGate-io/satgate-cli/internal/config"
)

// Client wraps HTTP calls to the SatGate Admin API
type Client struct {
	cfg  *config.Config
	http *http.Client
}

// Surface returns the configured surface
func (c *Client) Surface() string {
	return c.cfg.Surface
}

// New creates a new API client from current config
func New() (*Client, error) {
	cfg := config.Get()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{
		cfg: cfg,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Get performs a GET request to the given path
func (c *Client) Get(path string) ([]byte, int, error) {
	return c.do("GET", path, "")
}

// Post performs a POST request with a JSON body
func (c *Client) Post(path string, body interface{}) ([]byte, int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, 0, fmt.Errorf("marshaling request: %w", err)
	}
	return c.do("POST", path, string(data))
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) ([]byte, int, error) {
	return c.do("DELETE", path, "")
}

func (c *Client) do(method, path, body string) ([]byte, int, error) {
	url := strings.TrimRight(c.cfg.Gateway, "/") + path

	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	// Set auth header based on surface
	headerKey, headerVal := c.cfg.AuthHeader()
	req.Header.Set(headerKey, headerVal)

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set tenant header for cloud surface
	if c.cfg.Surface == "cloud" && c.cfg.Tenant != "" {
		req.Header.Set("X-SatGate-Tenant", c.cfg.Tenant)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}

	return data, resp.StatusCode, nil
}
