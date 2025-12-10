package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// QuerierClient handles communication with the Querier API
type QuerierClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewQuerierClient creates a new client for the Querier API
func NewQuerierClient(port int) *QuerierClient {
	return &QuerierClient{
		baseURL: fmt.Sprintf("http://localhost:%d", port),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Story represents a Hacker News story from the Querier
type Story struct {
	ID            string `json:"id"`
	Headline      string `json:"headline"`
	URL           string `json:"url"`
	DiscussionURL string `json:"discussion_url"`
	Username      string `json:"username"`
	Points        int    `json:"points"`
	Comments      int    `json:"comments"`
	Rank          int    `json:"rank"`
	Page          int    `json:"page"`
	AgeValue      int    `json:"age_value"`
	AgeUnit       string `json:"age_unit"`
}

// QueryResponse is the response from the /query endpoint
type QueryResponse struct {
	Count     int     `json:"count"`
	FetchedAt string  `json:"fetched_at"`
	Stories   []Story `json:"stories"`
}

// RefreshResponse is the response from the /refresh endpoint
type RefreshResponse struct {
	Message    string `json:"message"`
	FetchedAt  string `json:"fetched_at"`
	StoryCount int    `json:"story_count"`
}

// ErrorResponse represents an error from the Querier API
type ErrorResponse struct {
	Error string `json:"error"`
}

// Query sends a query request to the Querier and returns the results
// The requestBody should be the raw JSON body to forward to the Querier
func (c *QuerierClient) Query(requestBody []byte) (*QueryResponse, error) {
	resp, err := c.httpClient.Post(
		c.baseURL+"/query",
		"application/json",
		bytes.NewReader(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Querier: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
			return nil, fmt.Errorf("Querier error: %s", errResp.Error)
		}
		return nil, fmt.Errorf("Querier returned status %d: %s", resp.StatusCode, string(body))
	}

	var result QueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// Refresh tells the Querier to fetch fresh data from the Parser
func (c *QuerierClient) Refresh() (*RefreshResponse, error) {
	resp, err := c.httpClient.Post(
		c.baseURL+"/refresh",
		"application/json",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Querier: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
			return nil, fmt.Errorf("Querier error: %s", errResp.Error)
		}
		return nil, fmt.Errorf("Querier returned status %d: %s", resp.StatusCode, string(body))
	}

	var result RefreshResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
