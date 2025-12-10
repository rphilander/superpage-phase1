package parser

import (
	"encoding/json"
	"fmt"
	"net/http"

	"querier/models"
)

// Client is an HTTP client for the Parser API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Parser client
func NewClient(parserPort int) *Client {
	return &Client{
		baseURL:    fmt.Sprintf("http://localhost:%d", parserPort),
		httpClient: &http.Client{},
	}
}

// Fetch calls POST /fetch on the Parser and returns stories with fetched_at timestamp
func (c *Client) Fetch() ([]models.Story, string, error) {
	resp, err := c.httpClient.Post(c.baseURL+"/fetch", "application/json", nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to reach parser: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp models.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, "", fmt.Errorf("parser returned status %d", resp.StatusCode)
		}
		return nil, "", fmt.Errorf("parser error: %s", errResp.Error)
	}

	var parserResp models.ParserResponse
	if err := json.NewDecoder(resp.Body).Decode(&parserResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode parser response: %w", err)
	}

	return parserResp.Stories, parserResp.FetchedAt, nil
}
