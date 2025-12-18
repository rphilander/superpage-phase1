package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Handlers holds HTTP handlers and their dependencies
type Handlers struct {
	client    *WindowViewerClient
	templates *template.Template
}

// NewHandlers creates a new Handlers instance
func NewHandlers(client *WindowViewerClient) (*Handlers, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Handlers{
		client:    client,
		templates: tmpl,
	}, nil
}

// HandleIndex serves the main HTML page
func (h *Handlers) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// HandleAPIStories proxies story requests to Window Viewer
func (h *Handlers) HandleAPIStories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	from, err := strconv.ParseInt(query.Get("from"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
		return
	}

	to, err := strconv.ParseInt(query.Get("to"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'to' parameter", http.StatusBadRequest)
		return
	}

	criteria := query.Get("criteria")
	if criteria == "" {
		criteria = "max_points"
	}

	limit := 10
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	stories, err := h.client.GetTopStories(from, to, criteria, limit)
	if err != nil {
		log.Printf("Error fetching stories: %v", err)
		http.Error(w, "Failed to fetch stories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stories); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
