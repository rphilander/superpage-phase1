package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler holds the dependencies for HTTP handlers
type Handler struct {
	rateLimiter *RateLimiterClient
	numPages    int
}

// NewHandler creates a new Handler
func NewHandler(rateLimiter *RateLimiterClient, numPages int) *Handler {
	return &Handler{
		rateLimiter: rateLimiter,
		numPages:    numPages,
	}
}

// HandleFetch handles POST /fetch requests
func (h *Handler) HandleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var allStories []Story
	var firstFetchedAt string

	// Fetch each page sequentially
	for page := 1; page <= h.numPages; page++ {
		url := buildHNURL(page)

		resp, err := h.rateLimiter.FetchURL(url)
		if err != nil {
			writeError(w, http.StatusBadGateway, fmt.Sprintf("Failed to fetch page %d: %v", page, err))
			return
		}

		// Track the first fetch time for metadata
		if firstFetchedAt == "" {
			firstFetchedAt = resp.FetchedAt
		}

		// Parse the HTML
		stories, err := ParseHNPage(resp.HTML, page)
		if err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse page %d: %v", page, err))
			return
		}

		allStories = append(allStories, stories...)
	}

	// Build response
	response := FetchResponse{
		FetchedAt:    firstFetchedAt,
		NumPages:     h.numPages,
		TotalStories: len(allStories),
		Stories:      allStories,
	}

	writeJSON(w, http.StatusOK, response)
}

// HandleDoc handles GET /doc requests
func (h *Handler) HandleDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	doc := map[string]interface{}{
		"name":        "Parser API",
		"version":     "1.0.0",
		"description": "Fetches and parses Hacker News top stories via a rate-limited fetcher",
		"endpoints": []map[string]interface{}{
			{
				"method":      "POST",
				"path":        "/fetch",
				"description": "Fetches N pages of Hacker News top stories and returns parsed data",
				"request": map[string]interface{}{
					"body":         nil,
					"content_type": nil,
				},
				"response": map[string]interface{}{
					"success": map[string]interface{}{
						"status_code":  200,
						"content_type": "application/json",
						"body": map[string]interface{}{
							"fetched_at": map[string]interface{}{
								"type":        "string",
								"format":      "RFC3339",
								"description": "Timestamp when the first page was fetched",
							},
							"num_pages": map[string]interface{}{
								"type":        "integer",
								"description": "Number of Hacker News pages fetched",
							},
							"total_stories": map[string]interface{}{
								"type":        "integer",
								"description": "Total number of stories parsed across all pages",
							},
							"stories": map[string]interface{}{
								"type":        "array",
								"description": "Array of story objects",
								"items": map[string]interface{}{
									"rank": map[string]interface{}{
										"type":        "integer",
										"description": "Story's position on Hacker News (1-indexed)",
									},
									"id": map[string]interface{}{
										"type":        "string",
										"description": "Hacker News story ID",
									},
									"headline": map[string]interface{}{
										"type":        "string",
										"description": "Story title/headline",
									},
									"url": map[string]interface{}{
										"type":        "string",
										"description": "URL of the linked article",
									},
									"username": map[string]interface{}{
										"type":        "string",
										"description": "Username of the story submitter",
									},
									"points": map[string]interface{}{
										"type":        "integer",
										"description": "Number of upvotes/points",
									},
									"comments": map[string]interface{}{
										"type":        "integer",
										"description": "Number of comments on the story",
									},
									"discussion_url": map[string]interface{}{
										"type":        "string",
										"description": "URL to the Hacker News discussion page",
									},
									"age_value": map[string]interface{}{
										"type":        "integer",
										"description": "Numeric value of the story's age",
									},
									"age_unit": map[string]interface{}{
										"type":        "string",
										"description": "Unit of the age (minutes, hours, days)",
									},
									"page": map[string]interface{}{
										"type":        "integer",
										"description": "Which Hacker News page this story appeared on",
									},
								},
							},
						},
						"example": map[string]interface{}{
							"fetched_at":    "2025-12-06T10:30:00Z",
							"num_pages":     2,
							"total_stories": 60,
							"stories": []map[string]interface{}{
								{
									"rank":           1,
									"id":             "46173547",
									"headline":       "Tiny Core Linux: a 23 MB Linux distro with graphical desktop",
									"url":            "http://www.tinycorelinux.net/",
									"username":       "LorenDB",
									"points":         221,
									"comments":       114,
									"discussion_url": "https://news.ycombinator.com/item?id=46173547",
									"age_value":      4,
									"age_unit":       "hours",
									"page":           1,
								},
							},
						},
					},
					"error": map[string]interface{}{
						"status_codes": []int{405, 500, 502},
						"content_type": "application/json",
						"body": map[string]interface{}{
							"error": map[string]interface{}{
								"type":        "string",
								"description": "A description of the error",
							},
						},
						"examples": []map[string]interface{}{
							{
								"status_code": 502,
								"body": map[string]interface{}{
									"error": "Failed to fetch page 1: failed to reach rate limiter: connection refused",
								},
							},
							{
								"status_code": 500,
								"body": map[string]interface{}{
									"error": "Failed to parse page 1: failed to parse HTML: unexpected EOF",
								},
							},
						},
					},
				},
			},
			{
				"method":      "GET",
				"path":        "/doc",
				"description": "Returns this API documentation",
				"request": map[string]interface{}{
					"body": nil,
				},
				"response": map[string]interface{}{
					"success": map[string]interface{}{
						"status_code":  200,
						"content_type": "application/json",
						"description":  "Returns this documentation object",
					},
				},
			},
		},
	}

	writeJSON(w, http.StatusOK, doc)
}

// buildHNURL builds the Hacker News URL for a given page number
func buildHNURL(page int) string {
	if page == 1 {
		return "https://news.ycombinator.com/"
	}
	return fmt.Sprintf("https://news.ycombinator.com/?p=%d", page)
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
