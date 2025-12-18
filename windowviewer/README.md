# Window Viewer

Window Viewer is a REST API service that computes "top" Hacker News stories within a time window based on various ranking criteria. It fetches snapshot data from SnapshotDB and performs aggregation and ranking computations.

## Building

```bash
go build -o windowviewer .
```

## Running

```bash
./windowviewer --api <port> --snapshotdb <port>
```

### Required Arguments

- `--api <port>` - Port number for Window Viewer's HTTP API
- `--snapshotdb <port>` - Port number where SnapshotDB is listening on localhost

### Example

```bash
./windowviewer --api 8090 --snapshotdb 8082
```

## REST API

### GET /top

Returns the top stories within a time window based on the specified criteria.

#### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `from` | integer | Yes | Start of time window (Unix timestamp) |
| `to` | integer | Yes | End of time window (Unix timestamp) |
| `criteria` | string | Yes | Ranking criteria (see below) |
| `limit` | integer | No | Number of stories to return (default: 10) |

#### Criteria Values

| Value | Description |
|-------|-------------|
| `best_rank` | Highest ranking (lowest rank number) achieved during window |
| `max_points` | Highest points value achieved during window |
| `max_comments` | Highest total comments achieved during window |
| `incremental_comments` | Number of comments added during the window |
| `incremental_points` | Number of points added during the window |

#### Example Request

```bash
curl "http://localhost:8090/top?from=1702382400&to=1702386000&criteria=max_points&limit=5"
```

#### Example Response

```json
{
  "from": 1702382400,
  "to": 1702386000,
  "criteria": "max_points",
  "stories": [
    {
      "story_id": "46174114",
      "headline": "4 billion if statements (2023)",
      "url": "https://example.com/article",
      "username": "damethos",
      "discussion_url": "https://news.ycombinator.com/item?id=46174114",
      "best_rank": 1,
      "max_points": 443,
      "max_comments": 156,
      "incremental_points": 87,
      "incremental_comments": 42
    }
  ]
}
```

### GET /doc

Returns API documentation in JSON format.

#### Example Request

```bash
curl "http://localhost:8090/doc"
```

## Architecture

```
windowviewer/
├── main.go              # Entry point, CLI parsing, HTTP server setup
├── client/
│   └── snapshotdb.go    # HTTP client for SnapshotDB API
├── handlers/
│   └── handlers.go      # HTTP request handlers for /top and /doc
├── models/
│   └── models.go        # Data structures for API requests/responses
├── compute/
│   └── ranking.go       # Ranking computation logic
├── go.mod
└── README.md
```

### Package Descriptions

#### main.go
Entry point that parses command-line flags, creates the SnapshotDB client, sets up HTTP routes, and starts the server.

#### client/snapshotdb.go
HTTP client for communicating with SnapshotDB. Provides methods:
- `GetSnapshots(from, to int64)` - Fetches all snapshots within a time window

#### handlers/handlers.go
HTTP handlers for the REST API:
- `TopHandler` - Handles GET /top requests, validates parameters, calls compute logic
- `DocHandler` - Returns API documentation JSON

#### models/models.go
Data structures:
- `Story` - Story data from SnapshotDB
- `Snapshot` - Snapshot containing stories
- `RankedStory` - Story with computed ranking metric
- `TopStoriesResponse` - Response for GET /top

#### compute/ranking.go
Core ranking computation logic:
- Aggregates story data across all snapshots in the window
- Computes the appropriate metric for each story based on criteria
- Sorts and returns top N stories

## Ranking Algorithm

### best_rank
For each story, finds the minimum (best) rank achieved across all snapshots. Stories are sorted by rank ascending (rank 1 is best).

### max_points
For each story, finds the maximum points value across all snapshots. Stories are sorted by points descending.

### max_comments
For each story, finds the maximum comment count across all snapshots. Stories are sorted by comments descending.

### incremental_comments
Calculates `last_comments - first_comments` for each story within the window. If a story first appears mid-window, its initial value serves as the baseline. Stories are sorted by increment descending.

### incremental_points
Calculates `last_points - first_points` for each story within the window. If a story first appears mid-window, its initial value serves as the baseline. Stories are sorted by increment descending.

## Dependencies

- **SnapshotDB** - Must be running and accessible on localhost at the port specified by `--snapshotdb`

## Error Handling

The API returns JSON error responses with appropriate HTTP status codes:

```json
{
  "error": "error message here"
}
```

Common errors:
- `400 Bad Request` - Missing or invalid parameters
- `500 Internal Server Error` - Failed to fetch from SnapshotDB or compute results

## Debugging

### Check if Window Viewer is running
```bash
curl http://localhost:8090/doc
```

### Check SnapshotDB connectivity
```bash
curl http://localhost:8082/status
```

### Test with a broad time range
```bash
curl "http://localhost:8090/top?from=0&to=9999999999&criteria=max_points&limit=5"
```

### Verify SnapshotDB has data
```bash
curl "http://localhost:8082/snapshots?from=0&to=9999999999"
```

## Logs

Window Viewer logs to stdout:
- Startup message with API port
- SnapshotDB connection port

For detailed request logging, consider adding middleware or using a reverse proxy.
