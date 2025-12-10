package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"webui/client"
)

//go:embed templates/index.html
var templatesFS embed.FS

// Server handles HTTP requests for the WebUI
type Server struct {
	browserPort   int
	querierClient *client.QuerierClient
}

// NewServer creates a new WebUI server
func NewServer(browserPort, querierPort int) *Server {
	return &Server{
		browserPort:   browserPort,
		querierClient: client.NewQuerierClient(querierPort),
	}
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.handleIndex)
	mux.HandleFunc("POST /api/query", s.handleQuery)
	mux.HandleFunc("POST /api/refresh", s.handleRefresh)

	addr := fmt.Sprintf(":%d", s.browserPort)
	log.Printf("WebUI listening on http://localhost%s", addr)

	return http.ListenAndServe(addr, mux)
}

// handleIndex serves the main HTML page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	html, err := templatesFS.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html)
}

// handleQuery proxies query requests to the Querier
func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendJSONError(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	result, err := s.querierClient.Query(body)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadGateway)
		return
	}

	sendJSON(w, result)
}

// handleRefresh proxies refresh requests to the Querier
func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	result, err := s.querierClient.Refresh()
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadGateway)
		return
	}

	sendJSON(w, result)
}

// sendJSON writes a JSON response
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// sendJSONError writes a JSON error response
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
