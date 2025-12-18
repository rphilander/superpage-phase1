package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"snapshotdb/store"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type FetchResponse struct {
	FetchedAt    string  `json:"fetched_at"`
	NumPages     int     `json:"num_pages"`
	TotalStories int     `json:"total_stories"`
	Stories      []Story `json:"stories"`
}

type Story struct {
	ID            string `json:"id"`
	Rank          int    `json:"rank"`
	Headline      string `json:"headline"`
	URL           string `json:"url"`
	Username      string `json:"username"`
	Points        int    `json:"points"`
	Comments      int    `json:"comments"`
	DiscussionURL string `json:"discussion_url"`
	AgeValue      int    `json:"age_value"`
	AgeUnit       string `json:"age_unit"`
	Page          int    `json:"page"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Fetch() (*store.Snapshot, error) {
	return c.FetchWithRetry(5)
}

func (c *Client) FetchWithRetry(maxRetries int) (*store.Snapshot, error) {
	var lastErr error
	delay := 100 * time.Millisecond
	maxDelay := 5 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		snapshot, err := c.doFetch()
		if err == nil {
			return snapshot, nil
		}

		lastErr = err

		if attempt < maxRetries-1 {
			time.Sleep(delay)
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

func (c *Client) doFetch() (*store.Snapshot, error) {
	req, err := http.NewRequest("POST", c.baseURL+"/fetch", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var fetchResp FetchResponse
	if err := json.NewDecoder(resp.Body).Decode(&fetchResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	fetchedAt, err := time.Parse(time.RFC3339, fetchResp.FetchedAt)
	if err != nil {
		return nil, fmt.Errorf("parsing fetched_at: %w", err)
	}

	stories := make([]store.Story, len(fetchResp.Stories))
	for i, s := range fetchResp.Stories {
		stories[i] = store.Story{
			StoryID:       s.ID,
			Rank:          s.Rank,
			Headline:      s.Headline,
			URL:           s.URL,
			Username:      s.Username,
			Points:        s.Points,
			Comments:      s.Comments,
			DiscussionURL: s.DiscussionURL,
			AgeValue:      s.AgeValue,
			AgeUnit:       s.AgeUnit,
			Page:          s.Page,
		}
	}

	return &store.Snapshot{
		FetchedAt:    fetchedAt,
		NumPages:     fetchResp.NumPages,
		TotalStories: fetchResp.TotalStories,
		Stories:      stories,
	}, nil
}
