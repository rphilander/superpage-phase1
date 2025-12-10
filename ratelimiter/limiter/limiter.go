package limiter

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// FetchResult contains the result of a URL fetch operation.
type FetchResult struct {
	URL           string    `json:"url"`
	HTML          string    `json:"html"`
	StatusCode    int       `json:"status_code"`
	ContentLength int64     `json:"content_length"`
	FetchedAt     time.Time `json:"fetched_at"`
}

// RateLimiter enforces a minimum interval between HTTP requests.
type RateLimiter struct {
	interval   time.Duration
	lastFetch  time.Time
	mu         sync.Mutex
	httpClient *http.Client
}

// New creates a new RateLimiter with the specified interval in seconds.
func New(intervalSeconds int) *RateLimiter {
	return &RateLimiter{
		interval: time.Duration(intervalSeconds) * time.Second,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Fetch retrieves the content from the specified URL, respecting the rate limit.
// If a request was made recently, this method blocks until the rate limit allows.
func (r *RateLimiter) Fetch(url string) (*FetchResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Wait if necessary to respect the rate limit
	if !r.lastFetch.IsZero() {
		elapsed := time.Since(r.lastFetch)
		if elapsed < r.interval {
			waitTime := r.interval - elapsed
			time.Sleep(waitTime)
		}
	}

	// Record the fetch time
	r.lastFetch = time.Now()

	// Make the HTTP request
	resp, err := r.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &FetchResult{
		URL:           url,
		HTML:          string(body),
		StatusCode:    resp.StatusCode,
		ContentLength: int64(len(body)),
		FetchedAt:     r.lastFetch,
	}, nil
}
