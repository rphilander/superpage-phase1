package api

import (
	"encoding/json"
	"net/http"

	"ratelimiter/limiter"
)

// Handler holds the dependencies for the API handlers.
type Handler struct {
	rateLimiter *limiter.RateLimiter
}

// NewHandler creates a new Handler with the given rate limiter.
func NewHandler(rl *limiter.RateLimiter) *Handler {
	return &Handler{
		rateLimiter: rl,
	}
}

// FetchRequest represents the request body for POST /fetch.
type FetchRequest struct {
	URL string `json:"url"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleFetch handles POST /fetch requests.
func (h *Handler) HandleFetch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req FetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		h.sendError(w, "URL is required", http.StatusBadRequest)
		return
	}

	result, err := h.rateLimiter.Fetch(req.URL)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// HandleDoc handles GET /doc requests.
func (h *Handler) HandleDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	doc := map[string]interface{}{
		"name":        "Rate Limiter API",
		"version":     "1.0.0",
		"description": "A rate-limited HTTP fetcher API that retrieves HTML documents from URLs while respecting rate limits to avoid burdening remote websites.",
		"endpoints": []map[string]interface{}{
			{
				"method":      "POST",
				"path":        "/fetch",
				"description": "Fetches an HTML document from the specified URL. Requests are queued and processed according to the configured rate limit.",
				"request": map[string]interface{}{
					"content_type": "application/json",
					"body": map[string]interface{}{
						"url": map[string]string{
							"type":        "string",
							"required":    "true",
							"description": "The URL to fetch",
						},
					},
					"example": map[string]string{
						"url": "https://example.com",
					},
				},
				"response": map[string]interface{}{
					"success": map[string]interface{}{
						"status_code":  200,
						"content_type": "application/json",
						"body": map[string]interface{}{
							"url": map[string]string{
								"type":        "string",
								"description": "The URL that was fetched",
							},
							"html": map[string]string{
								"type":        "string",
								"description": "The HTML content retrieved from the URL",
							},
							"status_code": map[string]string{
								"type":        "integer",
								"description": "The HTTP status code returned by the remote server",
							},
							"content_length": map[string]string{
								"type":        "integer",
								"description": "The size of the response body in bytes",
							},
							"fetched_at": map[string]string{
								"type":        "string",
								"format":      "RFC3339",
								"description": "The timestamp when the fetch was performed",
							},
						},
						"example": map[string]interface{}{
							"url":            "https://example.com",
							"html":           "<!doctype html><html>...</html>",
							"status_code":    200,
							"content_length": 1256,
							"fetched_at":     "2025-12-05T10:30:00Z",
						},
					},
					"error": map[string]interface{}{
						"status_codes": []int{400, 405, 502},
						"content_type": "application/json",
						"body": map[string]interface{}{
							"error": map[string]string{
								"type":        "string",
								"description": "A description of the error",
							},
						},
						"examples": []map[string]interface{}{
							{
								"status_code": 400,
								"body": map[string]string{
									"error": "URL is required",
								},
							},
							{
								"status_code": 400,
								"body": map[string]string{
									"error": "Invalid JSON in request body",
								},
							},
							{
								"status_code": 502,
								"body": map[string]string{
									"error": "failed to fetch URL: connection refused",
								},
							},
						},
					},
				},
			},
			{
				"method":      "GET",
				"path":        "/doc",
				"description": "Returns this API documentation.",
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(doc)
}

// sendError sends a JSON error response.
func (h *Handler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
