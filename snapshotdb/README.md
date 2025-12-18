# SnapshotDB

SnapshotDB periodically fetches Hacker News snapshots from a Parser service and stores them in SQLite. It provides a REST API for clients to query historical snapshot data and track how stories change over time.

## Building

```bash
go build -o snapshotdb .
```

Requires Go 1.21+ and CGO (for SQLite).

## Usage

```bash
./snapshotdb --api <port> --db <path> --parser <port> --freq <seconds>
```

### Required Arguments

| Argument | Description |
|----------|-------------|
| `--api`  | Port number for the REST API |
| `--db`   | Path to SQLite database file |
| `--parser` | Port number of the Parser service on localhost |
| `--freq` | Fetch interval in seconds |

### Example

```bash
./snapshotdb --api 8082 --db data.db --parser 8081 --freq 60
```

This starts SnapshotDB on port 8082, storing data in `data.db`, fetching from Parser at `localhost:8081` every 60 seconds.

## Startup Behavior

- **No existing snapshots**: Fetches immediately from Parser with exponential backoff retry (5 attempts, 100ms initial delay, 2x backoff, 5s max delay)
- **Existing snapshots**: Calculates time until next scheduled fetch based on the last snapshot's timestamp

## REST API

All time window parameters use **Unix timestamps** (seconds since epoch).

### GET /status

Returns operational statistics.

```bash
curl http://localhost:8082/status
```

Response:
```json
{
  "uptime_seconds": 3600,
  "started_at": 1702382400,
  "snapshots_total": 60,
  "snapshots_errors": 2,
  "last_snapshot_at": 1702386000,
  "next_snapshot_at": 1702386060
}
```

### GET /snapshots?from=&to=

Returns all snapshots within a time window, including full story data.

```bash
curl "http://localhost:8082/snapshots?from=1702382400&to=1702386000"
```

### GET /stories?from=&to=

Returns deduplicated story IDs within a time window.

```bash
curl "http://localhost:8082/stories?from=1702382400&to=1702386000"
```

Response:
```json
{
  "story_ids": ["46243904", "46174114", "46245923"],
  "count": 3
}
```

### GET /story/{id}?from=&to=

Returns all data for a specific story across snapshots in a time window. Useful for tracking how a story's points, comments, and rank change over time.

```bash
curl "http://localhost:8082/story/46243904?from=1702382400&to=1702386000"
```

### GET /doc

Returns full API documentation as JSON.

```bash
curl http://localhost:8082/doc
```

## Database Schema

### snapshots table
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| fetched_at | DATETIME | When snapshot was taken |
| num_pages | INTEGER | Number of HN pages fetched |
| total_stories | INTEGER | Total stories in snapshot |

### stories table
| Column | Type | Description |
|--------|------|-------------|
| id | INTEGER | Primary key |
| snapshot_id | INTEGER | Foreign key to snapshots |
| story_id | TEXT | Hacker News story ID |
| rank | INTEGER | Position on HN |
| headline | TEXT | Story title |
| url | TEXT | Article URL |
| username | TEXT | Submitter |
| points | INTEGER | Upvotes |
| comments | INTEGER | Comment count |
| discussion_url | TEXT | HN discussion link |
| age_value | INTEGER | Numeric age |
| age_unit | TEXT | Age unit (minutes/hours/days) |
| page | INTEGER | HN page number |

### Indexes
- `idx_snapshots_fetched_at` on snapshots(fetched_at)
- `idx_stories_snapshot_id` on stories(snapshot_id)
- `idx_stories_story_id` on stories(story_id)

## Error Logging

Errors are logged to `<db-path>.errors.jsonl` in JSONL format:

```json
{"timestamp":1702382400,"error":"connection refused","context":"fetch_snapshot"}
```

## Files

```
snapshotdb/
├── main.go                 # Entry point, orchestration
├── go.mod
├── config/
│   └── config.go           # CLI argument parsing
├── store/
│   └── store.go            # SQLite operations
├── parser/
│   └── client.go           # Parser API client with retry
├── api/
│   ├── handlers.go         # HTTP handlers
│   └── models.go           # Request/response types
├── scheduler/
│   └── scheduler.go        # Periodic fetching
└── logger/
    └── logger.go           # JSONL error logging
```

## Debugging

### Check service status
```bash
curl http://localhost:8082/status
```

### View error log
```bash
cat data.db.errors.jsonl
```

### Query database directly
```bash
sqlite3 data.db "SELECT COUNT(*) FROM snapshots;"
sqlite3 data.db "SELECT fetched_at, total_stories FROM snapshots ORDER BY fetched_at DESC LIMIT 5;"
```

### Check if Parser is reachable
```bash
curl -X POST http://localhost:8081/fetch | head -100
```

## Graceful Shutdown

SnapshotDB handles SIGINT and SIGTERM signals gracefully:
1. Stops the scheduler (no new fetches)
2. Shuts down HTTP server (allows in-flight requests to complete)
3. Closes database connection

Press Ctrl+C or send SIGTERM to shut down cleanly.
