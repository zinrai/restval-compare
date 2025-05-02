package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	verbose    bool
}

func NewClient(timeout time.Duration, verbose bool) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		verbose: verbose,
	}
}

// Fetches JSON data from a URL
func (c *Client) FetchJSON(url string, headers map[string]string) (interface{}, error) {
	if c.verbose {
		fmt.Printf("GET %s\n", url)
	}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if c.verbose {
		fmt.Printf("Response status: %s\n", resp.Status)
		fmt.Printf("Response size: %d bytes\n", len(body))
	}

	// Parse JSON data
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	return result, nil
}
