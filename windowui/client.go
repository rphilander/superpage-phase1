package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WindowViewerClient is an HTTP client for the Window Viewer API
type WindowViewerClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewWindowViewerClient creates a new Window Viewer client
func NewWindowViewerClient(port int) *WindowViewerClient {
	return &WindowViewerClient{
		baseURL: fmt.Sprintf("http://localhost:%d", port),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetTopStories fetches top stories from Window Viewer
func (c *WindowViewerClient) GetTopStories(from, to int64, criteria string, limit int) (*TopStoriesResponse, error) {
	url := fmt.Sprintf("%s/top?from=%d&to=%d&criteria=%s&limit=%d",
		c.baseURL, from, to, criteria, limit)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stories: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result TopStoriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
