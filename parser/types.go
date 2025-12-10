package main

// Story represents a single Hacker News story with all its metadata
type Story struct {
	Rank          int    `json:"rank"`
	ID            string `json:"id"`
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

// FetchResponse is the response returned by POST /fetch
type FetchResponse struct {
	FetchedAt    string  `json:"fetched_at"`
	NumPages     int     `json:"num_pages"`
	TotalStories int     `json:"total_stories"`
	Stories      []Story `json:"stories"`
}

// RateLimiterRequest is the request body sent to the Rate Limiter
type RateLimiterRequest struct {
	URL string `json:"url"`
}

// RateLimiterResponse is the response from the Rate Limiter
type RateLimiterResponse struct {
	HTML          string `json:"html"`
	FetchedAt     string `json:"fetched_at"`
	StatusCode    int    `json:"status_code"`
	URL           string `json:"url"`
	ContentLength int    `json:"content_length"`
}

// ErrorResponse is returned when an error occurs
type ErrorResponse struct {
	Error string `json:"error"`
}
