package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RateLimiterClient handles communication with the Rate Limiter service
type RateLimiterClient struct {
	baseURL string
	client  *http.Client
}

// NewRateLimiterClient creates a new Rate Limiter client
func NewRateLimiterClient(port int) *RateLimiterClient {
	return &RateLimiterClient{
		baseURL: fmt.Sprintf("http://localhost:%d", port),
		client:  &http.Client{},
	}
}

// FetchURL requests the Rate Limiter to fetch the given URL
func (r *RateLimiterClient) FetchURL(url string) (*RateLimiterResponse, error) {
	reqBody := RateLimiterRequest{URL: url}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", r.baseURL+"/fetch", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach rate limiter: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
			return nil, fmt.Errorf("rate limiter error: %s", errResp.Error)
		}
		return nil, fmt.Errorf("rate limiter returned status %d", resp.StatusCode)
	}

	var result RateLimiterResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
