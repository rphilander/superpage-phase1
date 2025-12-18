package main

// Story represents a single Hacker News story from Window Viewer
type Story struct {
	StoryID             string `json:"story_id"`
	Headline            string `json:"headline"`
	URL                 string `json:"url"`
	Username            string `json:"username"`
	DiscussionURL       string `json:"discussion_url"`
	BestRank            int    `json:"best_rank"`
	MaxPoints           int    `json:"max_points"`
	MaxComments         int    `json:"max_comments"`
	IncrementalComments int    `json:"incremental_comments"`
	IncrementalPoints   int    `json:"incremental_points"`
}

// TopStoriesResponse represents the response from Window Viewer's /top endpoint
type TopStoriesResponse struct {
	From     int64   `json:"from"`
	To       int64   `json:"to"`
	Criteria string  `json:"criteria"`
	Stories  []Story `json:"stories"`
}
