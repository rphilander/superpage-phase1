package models

// Story represents a Hacker News story from the Parser
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

// QueryRequest is the body for POST /query
type QueryRequest struct {
	Filters *Filters   `json:"filters,omitempty"`
	Sort    []SortSpec `json:"sort,omitempty"`
}

// Filters contains all possible filter criteria
type Filters struct {
	Headline *FuzzyFilter `json:"headline,omitempty"`
	Username *FuzzyFilter `json:"username,omitempty"`
	URL      *FuzzyFilter `json:"url,omitempty"`
	Points   *RangeInt    `json:"points,omitempty"`
	Comments *RangeInt    `json:"comments,omitempty"`
	Rank     *RangeInt    `json:"rank,omitempty"`
	Page     *RangeInt    `json:"page,omitempty"`
	Age      *AgeRange    `json:"age,omitempty"`
}

// FuzzyFilter for string matching with optional threshold
type FuzzyFilter struct {
	Match     string `json:"match"`
	Threshold *int   `json:"threshold,omitempty"` // 0-100, default 50
}

// RangeInt for integer range filters (one-sided or two-sided)
type RangeInt struct {
	Min *int `json:"min,omitempty"`
	Max *int `json:"max,omitempty"`
}

// AgeRange for age-based filtering
type AgeRange struct {
	Unit string `json:"unit"` // minutes, hours, days
	Min  *int   `json:"min,omitempty"`
	Max  *int   `json:"max,omitempty"`
}

// SortSpec specifies a field and direction for sorting
type SortSpec struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // "asc" or "desc"
}

// QueryResponse is the response for POST /query
type QueryResponse struct {
	Stories   []Story `json:"stories"`
	Count     int     `json:"count"`
	FetchedAt string  `json:"fetched_at"`
}

// RefreshResponse is the response for POST /refresh
type RefreshResponse struct {
	Message    string `json:"message"`
	StoryCount int    `json:"story_count"`
	FetchedAt  string `json:"fetched_at"`
}

// ErrorResponse for error responses
type ErrorResponse struct {
	Error string `json:"error"`
}

// ParserResponse represents the response from the Parser API
type ParserResponse struct {
	FetchedAt    string  `json:"fetched_at"`
	NumPages     int     `json:"num_pages"`
	Stories      []Story `json:"stories"`
	TotalStories int     `json:"total_stories"`
}
