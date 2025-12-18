# WindowUI

WindowUI is a web frontend for browsing top Hacker News stories. It provides an interactive interface for selecting time windows, ranking criteria, and filtering results.

## Architecture

```
Browser → WindowUI → Window Viewer → SnapshotDB
              ↓
         HTML/JS UI
```

WindowUI:
1. Serves an HTML page with interactive controls
2. Proxies API requests to Window Viewer's `/top` endpoint
3. Provides client-side fuzzy search filtering

## Building

```bash
go build -o windowui .
```

## Running

```bash
./windowui --browser <port> --windowviewer <port>
```

### Required Arguments

| Argument | Description |
|----------|-------------|
| `--browser` | Port for the HTTP server (browser access) |
| `--windowviewer` | Port where Window Viewer is listening |

### Example

```bash
./windowui --browser 3000 --windowviewer 8090
```

Then open `http://localhost:3000` in your browser.

## REST API

### GET /

Serves the main HTML page with the interactive UI.

### GET /api/stories

Proxies requests to Window Viewer and returns top stories.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `from` | integer | Yes | Start of time window (Unix timestamp) |
| `to` | integer | Yes | End of time window (Unix timestamp) |
| `criteria` | string | No | Ranking criteria (default: `max_points`) |
| `limit` | integer | No | Number of stories (default: 10) |

**Example:**
```bash
curl "http://localhost:3000/api/stories?from=1702382400&to=1702386000&criteria=max_points&limit=30"
```

## UI Features

### Time Window Selection

Two-click range selection for defining the time window:
1. Click a button to set the first boundary (pulses to indicate pending)
2. Click another button to set the second boundary
3. The range is highlighted between the two selections

Available presets: Now, 1h, 2h, 3h, 5h, 8h, 13h, 24h, 2d, 3d, 5d, 8d

### Ranking Criteria

| Button | Criteria | Description |
|--------|----------|-------------|
| Best Rank | `best_rank` | Highest position achieved on HN |
| Max Points | `max_points` | Highest points value |
| Max Comments | `max_comments` | Highest comment count |
| + Comments | `incremental_comments` | Comments added during window |
| + Points | `incremental_points` | Points added during window |

### Result Limit

Select how many stories to display: 30, 50, 100, or 200.

### Fuzzy Search

Client-side headline filtering using Fuse.js. Type in the search box to filter displayed stories by headline.

## Project Structure

```
windowui/
├── main.go              # Entry point, CLI parsing, HTTP server setup
├── handlers.go          # HTTP handlers for / and /api/stories
├── client.go            # Window Viewer API client
├── models.go            # Data structures (Story, TopStoriesResponse)
├── templates/
│   └── index.html       # Web UI with embedded CSS and JavaScript
├── go.mod
└── README.md
```

## Dependencies

- **Window Viewer** - Must be running and accessible at the specified port
- **Fuse.js** - Loaded from CDN for client-side fuzzy search

Go dependencies: standard library only (`net/http`, `html/template`, `encoding/json`, `flag`)

## Debugging

### Check if WindowUI is running
```bash
curl http://localhost:3000/api/stories?from=0&to=9999999999&criteria=max_points&limit=5
```

### Check Window Viewer connectivity
```bash
curl "http://localhost:8090/top?from=0&to=9999999999&criteria=max_points&limit=5"
```

### Verify data flow
```bash
# Check SnapshotDB has data
curl http://localhost:8082/status

# Check Window Viewer can query
curl "http://localhost:8090/top?from=0&to=9999999999&criteria=max_points&limit=5"

# Check WindowUI proxies correctly
curl "http://localhost:3000/api/stories?from=0&to=9999999999&criteria=max_points&limit=5"
```
