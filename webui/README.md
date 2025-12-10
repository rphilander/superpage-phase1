# WebUI

A Go web application that provides a browser-based interface for viewing and filtering Hacker News stories from the Querier service.

## Overview

WebUI is the frontend component of the Hacker News aggregation system. It serves an HTML page with interactive filtering and sorting controls, proxying requests to the Querier service to fetch and query story data.

## Architecture

```
Browser → WebUI → Querier → Parser → Rate Limiter → Hacker News
```

## Building

```bash
go build -o webui
```

## Running

```bash
./webui --browser <port> --querier <port>
```

### Required Arguments

| Argument | Description |
|----------|-------------|
| `--browser` | Port number for browser HTTP requests |
| `--querier` | Port number where the Querier service is listening |

### Example

```bash
# Start WebUI on port 8083, connecting to Querier on port 8082
./webui --browser 8083 --querier 8082
```

Then open `http://localhost:8083` in your browser.

## HTTP Endpoints

### GET /

Serves the main HTML page with the story viewer interface.

### POST /api/query

Proxies query requests to the Querier's `/query` endpoint. Accepts the same JSON body format as the Querier.

**Request:**
```bash
curl -X POST http://localhost:8083/api/query \
  -H "Content-Type: application/json" \
  -d '{"filters": {"points": {"min": 100}}, "sort": [{"field": "points", "direction": "desc"}]}'
```

**Response:**
```json
{
  "count": 15,
  "fetched_at": "2025-12-06T10:30:00Z",
  "stories": [...]
}
```

### POST /api/refresh

Proxies refresh requests to the Querier's `/refresh` endpoint, triggering a fresh fetch from Hacker News.

**Request:**
```bash
curl -X POST http://localhost:8083/api/refresh
```

**Response:**
```json
{
  "message": "Data refreshed successfully",
  "story_count": 60,
  "fetched_at": "2025-12-06T10:30:00Z"
}
```

## Project Structure

```
webui/
├── main.go              # Entry point, CLI argument parsing
├── client/
│   └── querier.go       # HTTP client for Querier API
├── server/
│   ├── server.go        # HTTP server and route handlers
│   └── templates/
│       └── index.html   # Main HTML page (embedded at build time)
├── go.mod               # Go module definition
└── README.md            # This file
```

## Code Overview

### main.go
- Parses `--browser` and `--querier` CLI arguments
- Validates that both ports are provided
- Creates and starts the server

### client/querier.go
`QuerierClient` handles communication with the Querier API:
- `NewQuerierClient(port)` - Creates a new client with 30-second timeout
- `Query(requestBody)` - Sends a query request, returns `*QueryResponse`
- `Refresh()` - Triggers data refresh, returns `*RefreshResponse`

Data types:
- `Story` - Represents a single HN story
- `QueryResponse` - Response from `/query` endpoint
- `RefreshResponse` - Response from `/refresh` endpoint
- `ErrorResponse` - Error response format

### server/server.go
`Server` handles HTTP requests:
- `NewServer(browserPort, querierPort)` - Creates a new server instance
- `Start()` - Sets up routes and starts listening
- `handleIndex()` - Serves the embedded HTML template
- `handleQuery()` - Proxies requests to Querier's `/query`
- `handleRefresh()` - Proxies requests to Querier's `/refresh`

The HTML template is embedded using Go's `//go:embed` directive, so the binary is self-contained.

### server/templates/index.html
Single-page application with:
- **Header**: Title, status display, refresh button
- **Filter controls**: Fuzzy search (headline, username, URL), range filters (points, comments, rank, page, age), sort options
- **Story table**: Displays rank, headline with domain, points, comments (linked to discussion), username, age

JavaScript functions:
- `fetchStories()` - Calls `/api/query` with current filters
- `refreshData()` - Calls `/api/refresh` then fetches stories
- `buildQueryBody()` - Constructs query JSON from form inputs
- `applyFilters()` / `clearFilters()` - Filter management
- `renderStories(data)` - Renders the story table
- `escapeHtml(text)` - Prevents XSS in rendered content

## UI Features

### Filters
All filters are combined with AND logic:

| Filter | Type | Description |
|--------|------|-------------|
| Headline | Fuzzy | Search story headlines |
| Username | Fuzzy | Filter by author |
| URL | Fuzzy | Filter by story URL |
| Points | Range | Min/max point values |
| Comments | Range | Min/max comment counts |
| Rank | Range | Min/max rank positions |
| Page | Range | Min/max HN page numbers |
| Age | Range | Min/max age with unit (minutes/hours/days) |

### Sorting
Single-field sorting with ascending/descending direction. Sortable fields: rank, points, comments, age, headline, username, page.

## Error Handling

| HTTP Status | Cause |
|-------------|-------|
| 200 | Success |
| 400 | Invalid request body |
| 502 | Querier unreachable or returned error |

Errors are displayed in the UI status area with the error message.

## Dependencies

This project uses only the Go standard library:
- `net/http` - HTTP server
- `encoding/json` - JSON encoding/decoding
- `embed` - Embedding HTML template
- `flag` - CLI argument parsing
- `io` - Reading request bodies
- `fmt`, `log` - Output and logging

## Debugging

### Verify Querier connectivity
```bash
curl -X POST http://localhost:8082/query \
  -H "Content-Type: application/json" \
  -d '{}'
```

### Test WebUI API endpoints
```bash
# Test query proxy
curl -X POST http://localhost:8083/api/query \
  -H "Content-Type: application/json" \
  -d '{}' | jq .

# Test refresh proxy
curl -X POST http://localhost:8083/api/refresh | jq .
```

### Browser DevTools
- Open Network tab to inspect `/api/query` and `/api/refresh` requests
- Check Console for JavaScript errors
- Verify request/response payloads match expected format

### Common Issues

| Symptom | Likely Cause |
|---------|--------------|
| "Failed to connect to Querier" | Querier service not running on configured port |
| "Loading..." never completes | JavaScript error or network issue |
| Empty results | Filters too restrictive, or data not yet fetched |
| Stories not updating | Need to click "Refresh Data" to fetch fresh data |
