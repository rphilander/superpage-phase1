package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"windowviewer/models"
)

// SnapshotDBClient is an HTTP client for SnapshotDB
type SnapshotDBClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewSnapshotDBClient creates a new SnapshotDB client
func NewSnapshotDBClient(port int) *SnapshotDBClient {
	return &SnapshotDBClient{
		baseURL: fmt.Sprintf("http://localhost:%d", port),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetSnapshots fetches all snapshots within the given time window
func (c *SnapshotDBClient) GetSnapshots(from, to int64) (*models.SnapshotsResponse, error) {
	url := fmt.Sprintf("%s/snapshots?from=%d&to=%d", c.baseURL, from, to)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch snapshots: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.SnapshotsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetStory fetches all occurrences of a story within the given time window
func (c *SnapshotDBClient) GetStory(storyID string, from, to int64) (*models.StoryResponse, error) {
	url := fmt.Sprintf("%s/story/%s?from=%d&to=%d", c.baseURL, storyID, from, to)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch story: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result models.StoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
