// Package client provides HTTP client functionality for Vikunja API
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Config holds Vikunja API configuration
type Config struct {
	BaseURL string
	Token   string
}

// Client wraps HTTP client for Vikunja API
type Client struct {
	config     Config
	httpClient *http.Client
}

// New creates a new Vikunja API client from environment
func New() (*Client, error) {
	baseURL := os.Getenv("VIKUNJA_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("VIKUNJA_URL environment variable is required")
	}

	token := os.Getenv("VIKUNJA_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("VIKUNJA_TOKEN environment variable is required")
	}

	return &Client{
		config: Config{
			BaseURL: baseURL,
			Token:   token,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// NewWithConfig creates a client with explicit config
func NewWithConfig(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Request performs an HTTP request to the Vikunja API
func (c *Client) Request(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	url := c.config.BaseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get performs a GET request
func (c *Client) Get(endpoint string) ([]byte, error) {
	return c.Request(http.MethodGet, endpoint, nil)
}

// Post performs a POST request
func (c *Client) Post(endpoint string, body interface{}) ([]byte, error) {
	return c.Request(http.MethodPost, endpoint, body)
}

// Put performs a PUT request
func (c *Client) Put(endpoint string, body interface{}) ([]byte, error) {
	return c.Request(http.MethodPut, endpoint, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(endpoint string) ([]byte, error) {
	return c.Request(http.MethodDelete, endpoint, nil)
}
