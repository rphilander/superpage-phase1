package api

type StatusResponse struct {
	UptimeSeconds   int64  `json:"uptime_seconds"`
	StartedAt       int64  `json:"started_at"`
	SnapshotsTotal  int    `json:"snapshots_total"`
	SnapshotsErrors int    `json:"snapshots_errors"`
	LastSnapshotAt  *int64 `json:"last_snapshot_at"`
	NextSnapshotAt  *int64 `json:"next_snapshot_at"`
}

type SnapshotsResponse struct {
	Snapshots []SnapshotDTO `json:"snapshots"`
}

type SnapshotDTO struct {
	ID           int64      `json:"id"`
	FetchedAt    int64      `json:"fetched_at"`
	NumPages     int        `json:"num_pages"`
	TotalStories int        `json:"total_stories"`
	Stories      []StoryDTO `json:"stories"`
}

type StoryDTO struct {
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

type StoriesResponse struct {
	StoryIDs []string `json:"story_ids"`
	Count    int      `json:"count"`
}

type StoryResponse struct {
	StoryID     string               `json:"story_id"`
	Occurrences []StoryOccurrenceDTO `json:"occurrences"`
}

type StoryOccurrenceDTO struct {
	SnapshotID    int64  `json:"snapshot_id"`
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

type ErrorResponse struct {
	Error string `json:"error"`
}

type DocResponse struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Endpoints   []EndpointDoc `json:"endpoints"`
}

type EndpointDoc struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Description string            `json:"description"`
	Parameters  []ParameterDoc    `json:"parameters,omitempty"`
	Response    ResponseDoc       `json:"response"`
	Example     *EndpointExample  `json:"example,omitempty"`
}

type ParameterDoc struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

type ResponseDoc struct {
	ContentType string `json:"content_type"`
	Description string `json:"description"`
}

type EndpointExample struct {
	Request  string      `json:"request"`
	Response interface{} `json:"response"`
}
