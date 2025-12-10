# Querier

Querier is a Go service that fetches Hacker News data from the Parser service, stores it in memory, and provides filtering and sorting capabilities via a REST API.

## Prerequisites

- Go 1.21 or later
- Parser service running (provides the underlying HN data)

## Building

```bash
go build -o querier
```

## Running

```bash
./querier --api <port> --parser <port>
```

**Required arguments:**
- `--api <port>`: Port number for the Querier REST API
- `--parser <port>`: Port number where the Parser service is listening

**Example:**
```bash
./querier --api 8082 --parser 8081
```

## REST API

### GET /doc

Returns comprehensive API documentation in JSON format.

```bash
curl http://localhost:8082/doc
```

### POST /refresh

Fetches fresh data from the Parser service and stores it in memory.

```bash
curl -X POST http://localhost:8082/refresh
```

**Response:**
```json
{
  "message": "Data refreshed successfully",
  "story_count": 60,
  "fetched_at": "2025-12-06T10:30:00Z"
}
```

### POST /query

Queries the stored data with optional filters and sorting. If no data is stored, automatically fetches from Parser first.

```bash
curl -X POST http://localhost:8082/query \
  -H "Content-Type: application/json" \
  -d '{"filters": {...}, "sort": [...]}'
```

**Response:**
```json
{
  "stories": [...],
  "count": 10,
  "fetched_at": "2025-12-06T10:30:00Z"
}
```

## Query Filters

All filters are optional and combined with AND logic.

### Fuzzy String Filters

Available for: `headline`, `username`, `url`

```json
{
  "filters": {
    "headline": {
      "match": "linux",
      "threshold": 50
    }
  }
}
```

- `match`: String to match against (required)
- `threshold`: Match threshold 0-100 (optional, default: 50). Higher values require closer matches.

### Integer Range Filters

Available for: `points`, `comments`, `rank`, `page`

```json
{
  "filters": {
    "points": {"min": 100},
    "comments": {"min": 10, "max": 500},
    "rank": {"max": 30}
  }
}
```

- `min`: Minimum value inclusive (optional)
- `max`: Maximum value inclusive (optional)

### Age Range Filter

```json
{
  "filters": {
    "age": {
      "unit": "hours",
      "min": 2,
      "max": 6
    }
  }
}
```

- `unit`: Time unit - "minutes", "hours", or "days" (required)
- `min`: Minimum age in specified units (optional)
- `max`: Maximum age in specified units (optional)

## Sorting

Specify an array of sort specifications. Multiple fields are applied in order.

```json
{
  "sort": [
    {"field": "points", "direction": "desc"},
    {"field": "comments", "direction": "asc"}
  ]
}
```

**Sortable fields:** `id`, `headline`, `url`, `discussion_url`, `username`, `points`, `comments`, `rank`, `page`, `age_value`, `age` (normalized to minutes)

**Directions:** `asc` (ascending) or `desc` (descending)

## Examples

### Get all stories sorted by points
```bash
curl -X POST http://localhost:8082/query \
  -H "Content-Type: application/json" \
  -d '{"sort": [{"field": "points", "direction": "desc"}]}'
```

### Find stories with "AI" in headline, 50+ points
```bash
curl -X POST http://localhost:8082/query \
  -H "Content-Type: application/json" \
  -d '{
    "filters": {
      "headline": {"match": "AI"},
      "points": {"min": 50}
    },
    "sort": [{"field": "points", "direction": "desc"}]
  }'
```

### Get recent stories (last 2 hours)
```bash
curl -X POST http://localhost:8082/query \
  -H "Content-Type: application/json" \
  -d '{
    "filters": {
      "age": {"unit": "hours", "max": 2}
    }
  }'
```

## Architecture

```
querier/
├── main.go              # Entry point, CLI parsing, HTTP server setup
├── models/
│   └── models.go        # Data types (Story, QueryRequest, etc.)
├── api/
│   └── handlers.go      # HTTP handlers for /query, /refresh, /doc
├── store/
│   └── store.go         # Thread-safe in-memory data store
├── parser/
│   └── client.go        # HTTP client for Parser API
└── query/
    ├── filter.go        # Filter logic (fuzzy match, range filters)
    └── sort.go          # Multi-field sorting
```

### Key Components

**Store** (`store/store.go`): Thread-safe in-memory storage using `sync.RWMutex`. Stores stories and the timestamp when they were fetched.

**Parser Client** (`parser/client.go`): HTTP client that calls `POST /fetch` on the Parser service to retrieve HN stories.

**Filter** (`query/filter.go`): Implements filtering logic:
- Fuzzy string matching using `github.com/sahilm/fuzzy`
- Integer range comparisons
- Age range filtering with unit normalization to minutes

**Sort** (`query/sort.go`): Multi-field stable sorting using `sort.SliceStable`.

**Handlers** (`api/handlers.go`): HTTP request handlers that orchestrate fetching, filtering, sorting, and response formatting.

### Data Flow

1. Client calls `POST /refresh` or `POST /query` (auto-refresh if empty)
2. Querier fetches data from Parser via `POST /fetch`
3. Stories are stored in memory
4. For queries, filters are applied sequentially (AND logic)
5. Sorting is applied to filtered results
6. Response includes stories, count, and fetched_at timestamp

## Dependencies

- `github.com/sahilm/fuzzy` - Fuzzy string matching library
- Go standard library (net/http, encoding/json, sync, sort, flag)
