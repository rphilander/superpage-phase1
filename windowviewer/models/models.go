package models

// Story represents a Hacker News story as returned by SnapshotDB
type Story struct {
	StoryID       string `json:"story_id"`
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

// Snapshot represents a snapshot from SnapshotDB
type Snapshot struct {
	ID           int     `json:"id"`
	FetchedAt    int64   `json:"fetched_at"`
	NumPages     int     `json:"num_pages"`
	TotalStories int     `json:"total_stories"`
	Stories      []Story `json:"stories"`
}

// SnapshotsResponse represents the response from GET /snapshots
type SnapshotsResponse struct {
	Snapshots []Snapshot `json:"snapshots"`
}

// StoryOccurrence represents a story's data at a specific snapshot time
type StoryOccurrence struct {
	SnapshotID    int    `json:"snapshot_id"`
	FetchedAt     int64  `json:"fetched_at"`
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

// StoryResponse represents the response from GET /story/{id}
type StoryResponse struct {
	StoryID     string            `json:"story_id"`
	Occurrences []StoryOccurrence `json:"occurrences"`
}

// RankedStory represents a story with all computed ranking metrics
type RankedStory struct {
	StoryID             string `json:"story_id"`
	Headline            string `json:"headline"`
	URL                 string `json:"url"`
	Username            string `json:"username"`
	DiscussionURL       string `json:"discussion_url"`
	BestRank            int    `json:"best_rank"`
	MaxPoints           int    `json:"max_points"`
	MaxComments         int    `json:"max_comments"`
	IncrementalPoints   int    `json:"incremental_points"`
	IncrementalComments int    `json:"incremental_comments"`
}

// TopStoriesResponse is the API response for GET /top
type TopStoriesResponse struct {
	From     int64         `json:"from"`
	To       int64         `json:"to"`
	Criteria string        `json:"criteria"`
	Stories  []RankedStory `json:"stories"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error string `json:"error"`
}
