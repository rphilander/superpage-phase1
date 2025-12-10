package api

import (
	"encoding/json"
	"net/http"

	"querier/models"
	"querier/parser"
	"querier/query"
	"querier/store"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	store        *store.Store
	parserClient *parser.Client
}

// NewHandler creates a new Handler
func NewHandler(s *store.Store, p *parser.Client) *Handler {
	return &Handler{
		store:        s,
		parserClient: p,
	}
}

// HandleRefresh handles POST /refresh
func (h *Handler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	stories, fetchedAt, err := h.parserClient.Fetch()
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	h.store.Update(stories, fetchedAt)

	writeJSON(w, http.StatusOK, models.RefreshResponse{
		Message:    "Data refreshed successfully",
		StoryCount: len(stories),
		FetchedAt:  fetchedAt,
	})
}

// HandleQuery handles POST /query
func (h *Handler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Auto-refresh if store is empty
	if h.store.IsEmpty() {
		stories, fetchedAt, err := h.parserClient.Fetch()
		if err != nil {
			writeError(w, http.StatusBadGateway, "No data available and failed to fetch: "+err.Error())
			return
		}
		h.store.Update(stories, fetchedAt)
	}

	var req models.QueryRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
			return
		}
	}

	stories, fetchedAt := h.store.Get()

	// Apply filters
	stories = query.Filter(stories, req.Filters)

	// Apply sorting
	stories = query.Sort(stories, req.Sort)

	writeJSON(w, http.StatusOK, models.QueryResponse{
		Stories:   stories,
		Count:     len(stories),
		FetchedAt: fetchedAt,
	})
}

// HandleDoc handles GET /doc
func (h *Handler) HandleDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	doc := getAPIDocumentation()
	writeJSON(w, http.StatusOK, doc)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, models.ErrorResponse{Error: message})
}

func getAPIDocumentation() map[string]interface{} {
	return map[string]interface{}{
		"name":        "Querier API",
		"version":     "1.0.0",
		"description": "Queries and filters Hacker News data obtained from the Parser service",
		"endpoints": []map[string]interface{}{
			{
				"method":      "POST",
				"path":        "/refresh",
				"description": "Fetches fresh data from the Parser service and stores it in memory",
				"request": map[string]interface{}{
					"body": nil,
				},
				"response": map[string]interface{}{
					"success": map[string]interface{}{
						"status_code":  200,
						"content_type": "application/json",
						"body": map[string]interface{}{
							"message": map[string]interface{}{
								"type":        "string",
								"description": "Success message",
							},
							"story_count": map[string]interface{}{
								"type":        "integer",
								"description": "Number of stories fetched",
							},
							"fetched_at": map[string]interface{}{
								"type":        "string",
								"format":      "RFC3339",
								"description": "Timestamp when data was fetched",
							},
						},
						"example": map[string]interface{}{
							"message":     "Data refreshed successfully",
							"story_count": 60,
							"fetched_at":  "2025-12-06T10:30:00Z",
						},
					},
					"error": map[string]interface{}{
						"status_codes": []int{405, 502},
						"body": map[string]interface{}{
							"error": map[string]interface{}{
								"type":        "string",
								"description": "Error description",
							},
						},
					},
				},
			},
			{
				"method":      "POST",
				"path":        "/query",
				"description": "Queries the stored data with optional filters and sorting. If no data is stored, automatically fetches from Parser first.",
				"request": map[string]interface{}{
					"content_type": "application/json",
					"body": map[string]interface{}{
						"filters": map[string]interface{}{
							"type":        "object",
							"optional":    true,
							"description": "Filter criteria (all optional, combined with AND logic)",
							"properties": map[string]interface{}{
								"headline": map[string]interface{}{
									"type":        "FuzzyFilter",
									"description": "Fuzzy match against story headline",
								},
								"username": map[string]interface{}{
									"type":        "FuzzyFilter",
									"description": "Fuzzy match against submitter username",
								},
								"url": map[string]interface{}{
									"type":        "FuzzyFilter",
									"description": "Fuzzy match against story URL",
								},
								"points": map[string]interface{}{
									"type":        "RangeInt",
									"description": "Filter by points (upvotes) range",
								},
								"comments": map[string]interface{}{
									"type":        "RangeInt",
									"description": "Filter by comment count range",
								},
								"rank": map[string]interface{}{
									"type":        "RangeInt",
									"description": "Filter by story rank (1-indexed position)",
								},
								"page": map[string]interface{}{
									"type":        "RangeInt",
									"description": "Filter by HN page number",
								},
								"age": map[string]interface{}{
									"type":        "AgeRange",
									"description": "Filter by story age",
								},
							},
						},
						"sort": map[string]interface{}{
							"type":        "array",
							"optional":    true,
							"description": "Sort specifications (applied in order)",
							"items": map[string]interface{}{
								"type": "SortSpec",
							},
						},
					},
				},
				"response": map[string]interface{}{
					"success": map[string]interface{}{
						"status_code":  200,
						"content_type": "application/json",
						"body": map[string]interface{}{
							"stories": map[string]interface{}{
								"type":        "array",
								"description": "Array of matching stories",
								"items":       "Story",
							},
							"count": map[string]interface{}{
								"type":        "integer",
								"description": "Number of matching stories",
							},
							"fetched_at": map[string]interface{}{
								"type":        "string",
								"format":      "RFC3339",
								"description": "When the underlying data was fetched",
							},
						},
					},
					"error": map[string]interface{}{
						"status_codes": []int{400, 405, 502},
					},
				},
				"examples": []map[string]interface{}{
					{
						"name":        "Filter by headline and points",
						"description": "Find stories with 'linux' in headline and at least 50 points, sorted by points descending",
						"request": map[string]interface{}{
							"filters": map[string]interface{}{
								"headline": map[string]interface{}{
									"match": "linux",
								},
								"points": map[string]interface{}{
									"min": 50,
								},
							},
							"sort": []map[string]interface{}{
								{"field": "points", "direction": "desc"},
							},
						},
					},
					{
						"name":        "Filter with custom fuzzy threshold",
						"description": "Loose fuzzy match (threshold 30) for username",
						"request": map[string]interface{}{
							"filters": map[string]interface{}{
								"username": map[string]interface{}{
									"match":     "john",
									"threshold": 30,
								},
							},
						},
					},
					{
						"name":        "Filter by age range",
						"description": "Stories posted between 2-6 hours ago",
						"request": map[string]interface{}{
							"filters": map[string]interface{}{
								"age": map[string]interface{}{
									"unit": "hours",
									"min":  2,
									"max":  6,
								},
							},
						},
					},
					{
						"name":        "No filters - return all",
						"description": "Return all stories without filtering",
						"request":     map[string]interface{}{},
					},
				},
			},
			{
				"method":      "GET",
				"path":        "/doc",
				"description": "Returns this API documentation",
				"response": map[string]interface{}{
					"success": map[string]interface{}{
						"status_code":  200,
						"content_type": "application/json",
						"description":  "This documentation object",
					},
				},
			},
		},
		"types": map[string]interface{}{
			"Story": map[string]interface{}{
				"description": "A Hacker News story",
				"properties": map[string]interface{}{
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
					"discussion_url": map[string]interface{}{
						"type":        "string",
						"description": "URL to the Hacker News discussion page",
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
						"description": "Number of comments",
					},
					"rank": map[string]interface{}{
						"type":        "integer",
						"description": "Story's position on Hacker News (1-indexed)",
					},
					"page": map[string]interface{}{
						"type":        "integer",
						"description": "Which Hacker News page this story appeared on",
					},
					"age_value": map[string]interface{}{
						"type":        "integer",
						"description": "Numeric value of the story's age",
					},
					"age_unit": map[string]interface{}{
						"type":        "string",
						"description": "Unit of the age (minutes, hours, days)",
					},
				},
			},
			"FuzzyFilter": map[string]interface{}{
				"description": "Filter for fuzzy string matching",
				"properties": map[string]interface{}{
					"match": map[string]interface{}{
						"type":        "string",
						"required":    true,
						"description": "The string to match against",
					},
					"threshold": map[string]interface{}{
						"type":        "integer",
						"optional":    true,
						"default":     50,
						"description": "Match threshold (0-100). Higher values require closer matches. Default is 50.",
					},
				},
			},
			"RangeInt": map[string]interface{}{
				"description": "Filter for integer range (one-sided or two-sided)",
				"properties": map[string]interface{}{
					"min": map[string]interface{}{
						"type":        "integer",
						"optional":    true,
						"description": "Minimum value (inclusive). Omit for no lower bound.",
					},
					"max": map[string]interface{}{
						"type":        "integer",
						"optional":    true,
						"description": "Maximum value (inclusive). Omit for no upper bound.",
					},
				},
			},
			"AgeRange": map[string]interface{}{
				"description": "Filter for story age",
				"properties": map[string]interface{}{
					"unit": map[string]interface{}{
						"type":        "string",
						"required":    true,
						"description": "Time unit: 'minutes', 'hours', or 'days'",
					},
					"min": map[string]interface{}{
						"type":        "integer",
						"optional":    true,
						"description": "Minimum age in specified units (inclusive)",
					},
					"max": map[string]interface{}{
						"type":        "integer",
						"optional":    true,
						"description": "Maximum age in specified units (inclusive)",
					},
				},
			},
			"SortSpec": map[string]interface{}{
				"description": "Specification for sorting",
				"properties": map[string]interface{}{
					"field": map[string]interface{}{
						"type":        "string",
						"required":    true,
						"description": "Field to sort by: id, headline, url, discussion_url, username, points, comments, rank, page, age_value, or age (normalized)",
					},
					"direction": map[string]interface{}{
						"type":        "string",
						"required":    true,
						"description": "'asc' for ascending or 'desc' for descending",
					},
				},
			},
		},
	}
}
