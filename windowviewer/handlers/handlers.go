package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"windowviewer/client"
	"windowviewer/compute"
	"windowviewer/models"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	snapshotDB *client.SnapshotDBClient
}

// NewHandler creates a new Handler with the given SnapshotDB client
func NewHandler(snapshotDB *client.SnapshotDBClient) *Handler {
	return &Handler{
		snapshotDB: snapshotDB,
	}
}

// TopHandler handles GET /top requests
func (h *Handler) TopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Parse query parameters
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	criteria := r.URL.Query().Get("criteria")
	limitStr := r.URL.Query().Get("limit")

	// Validate required parameters
	if fromStr == "" {
		writeError(w, http.StatusBadRequest, "missing required parameter: from")
		return
	}
	if toStr == "" {
		writeError(w, http.StatusBadRequest, "missing required parameter: to")
		return
	}
	if criteria == "" {
		writeError(w, http.StatusBadRequest, "missing required parameter: criteria")
		return
	}

	// Parse from
	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid 'from' parameter: must be a Unix timestamp")
		return
	}

	// Parse to
	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid 'to' parameter: must be a Unix timestamp")
		return
	}

	// Validate time window
	if from > to {
		writeError(w, http.StatusBadRequest, "'from' must be less than or equal to 'to'")
		return
	}

	// Validate criteria
	if !compute.IsValidCriteria(criteria) {
		writeError(w, http.StatusBadRequest, "invalid criteria: must be one of "+strings.Join(compute.ValidCriteria(), ", "))
		return
	}

	// Parse limit (default 10)
	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 1 {
			writeError(w, http.StatusBadRequest, "invalid 'limit' parameter: must be a positive integer")
			return
		}
		limit = parsedLimit
	}

	// Fetch snapshots from SnapshotDB
	snapshotsResp, err := h.snapshotDB.GetSnapshots(from, to)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch snapshots: "+err.Error())
		return
	}

	// Compute top stories
	rankedStories, err := compute.ComputeTopStories(snapshotsResp.Snapshots, criteria, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to compute top stories: "+err.Error())
		return
	}

	// Build response
	response := models.TopStoriesResponse{
		From:     from,
		To:       to,
		Criteria: criteria,
		Stories:  rankedStories,
	}

	writeJSON(w, http.StatusOK, response)
}

// DocHandler handles GET /doc requests
func (h *Handler) DocHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	doc := map[string]interface{}{
		"name":        "Window Viewer API",
		"version":     "1.0.0",
		"description": "REST API for computing top Hacker News stories within a time window based on various criteria",
		"endpoints": []map[string]interface{}{
			{
				"method":      "GET",
				"path":        "/top",
				"description": "Returns the top stories within a time window based on the specified criteria",
				"parameters": []map[string]interface{}{
					{
						"name":        "from",
						"type":        "integer",
						"required":    true,
						"description": "Start of time window (Unix timestamp)",
					},
					{
						"name":        "to",
						"type":        "integer",
						"required":    true,
						"description": "End of time window (Unix timestamp)",
					},
					{
						"name":        "criteria",
						"type":        "string",
						"required":    true,
						"description": "Ranking criteria. One of: best_rank, max_points, max_comments, incremental_comments, incremental_points",
						"values": []map[string]string{
							{"value": "best_rank", "description": "Highest ranking (lowest rank number) achieved during window"},
							{"value": "max_points", "description": "Highest points value achieved during window"},
							{"value": "max_comments", "description": "Highest total comments achieved during window"},
							{"value": "incremental_comments", "description": "Number of comments added during the window"},
							{"value": "incremental_points", "description": "Number of points added during the window"},
						},
					},
					{
						"name":        "limit",
						"type":        "integer",
						"required":    false,
						"description": "Number of stories to return (default: 10)",
					},
				},
				"response": map[string]interface{}{
					"content_type": "application/json",
					"description":  "Top stories with their computed metrics",
				},
				"example": map[string]interface{}{
					"request": "GET /top?from=1702382400&to=1702386000&criteria=max_points&limit=3",
					"response": map[string]interface{}{
						"from":     1702382400,
						"to":       1702386000,
						"criteria": "max_points",
						"stories": []map[string]interface{}{
							{
								"story_id":            "46174114",
								"headline":            "4 billion if statements (2023)",
								"url":                 "https://example.com/article",
								"username":            "damethos",
								"discussion_url":      "https://news.ycombinator.com/item?id=46174114",
								"best_rank":           1,
								"max_points":          443,
								"max_comments":        156,
								"incremental_points":  87,
								"incremental_comments": 42,
							},
							{
								"story_id":            "46243904",
								"headline":            "SQLite JSON at Full Index Speed",
								"url":                 "https://example.com/sqlite",
								"username":            "upmostly",
								"discussion_url":      "https://news.ycombinator.com/item?id=46243904",
								"best_rank":           3,
								"max_points":          215,
								"max_comments":        89,
								"incremental_points":  45,
								"incremental_comments": 23,
							},
							{
								"story_id":            "46245398",
								"headline":            "Epic celebrates end of the Apple Tax",
								"url":                 "https://example.com/epic",
								"username":            "nobody9999",
								"discussion_url":      "https://news.ycombinator.com/item?id=46245398",
								"best_rank":           5,
								"max_points":          163,
								"max_comments":        72,
								"incremental_points":  31,
								"incremental_comments": 18,
							},
						},
					},
				},
			},
			{
				"method":      "GET",
				"path":        "/doc",
				"description": "Returns this API documentation",
				"response": map[string]interface{}{
					"content_type": "application/json",
					"description":  "Full API documentation with examples",
				},
			},
		},
	}

	writeJSON(w, http.StatusOK, doc)
}

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.ErrorResponse{Error: message})
}
