# Rate Limiter

A Go CLI application that fetches HTML documents from URLs with rate limiting to avoid burdening remote websites.

## Overview

Rate Limiter provides a REST API for fetching web pages. It enforces a configurable minimum interval between HTTP requests, ensuring that even when multiple clients are making requests concurrently, the remote servers are not overwhelmed.

## Building

```bash
go build -o ratelimiter .
```

## Usage

```bash
./ratelimiter --rate <num-sec> --api <port-no>
```

### Required Arguments

| Argument | Description |
|----------|-------------|
| `--rate` | Minimum number of seconds between HTTP requests (positive integer) |
| `--api` | Port number for the REST API (1-65535) |

### Example

Start the server with a rate limit of 5 seconds between requests on port 8080:

```bash
./ratelimiter --rate 5 --api 8080
```

## REST API

### POST /fetch

Fetches an HTML document from the specified URL.

**Request:**
```bash
curl -X POST http://localhost:8080/fetch \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

**Response:**
```json
{
  "url": "https://example.com",
  "html": "<!doctype html>...",
  "status_code": 200,
  "content_length": 1256,
  "fetched_at": "2025-12-05T10:30:00Z"
}
```

**Error Response:**
```json
{
  "error": "description of error"
}
```

### GET /doc

Returns detailed API documentation in JSON format.

```bash
curl http://localhost:8080/doc
```

## Architecture

```
ratelimiter/
├── main.go           # Entry point, CLI argument parsing, server startup
├── api/
│   └── handler.go    # REST API handlers for /fetch and /doc
├── limiter/
│   └── limiter.go    # Rate limiting logic
├── go.mod
└── README.md
```

### Package Descriptions

#### `main`
The entry point that:
- Parses and validates CLI arguments using the `flag` package
- Initializes the rate limiter with the configured interval
- Sets up HTTP routes and starts the server

#### `limiter`
Contains the `RateLimiter` type which:
- Tracks the timestamp of the last HTTP request
- Uses a mutex to ensure thread-safe access when multiple API requests arrive concurrently
- Blocks (waits) when necessary to respect the configured rate limit
- Fetches URLs using an HTTP client with a 30-second timeout

#### `api`
Contains the `Handler` type which:
- Implements the `/fetch` endpoint (accepts JSON, returns JSON with metadata)
- Implements the `/doc` endpoint (returns API documentation)
- Handles errors with appropriate HTTP status codes

### Rate Limiting Strategy

The rate limiter uses a **blocking queue strategy**:
1. When a `/fetch` request arrives, it acquires a mutex lock
2. If the time since the last request is less than the configured interval, it sleeps for the remaining time
3. It then performs the HTTP fetch and releases the lock

This means:
- Requests are processed in order (first-come, first-served)
- No requests are rejected; they wait in line
- The actual time between remote HTTP requests is always >= the configured interval

### Thread Safety

The rate limiter is thread-safe. Multiple concurrent API requests are safely serialized through a mutex, ensuring that:
- Only one HTTP request to external URLs happens at a time
- The timing between requests is properly enforced regardless of concurrent API calls

## Error Handling

| HTTP Status | Cause |
|-------------|-------|
| 200 | Successful fetch |
| 400 | Invalid JSON or missing URL in request body |
| 405 | Wrong HTTP method (e.g., GET on /fetch) |
| 502 | Failed to fetch the remote URL (connection error, timeout, etc.) |

## Dependencies

This project uses only the Go standard library:
- `net/http` - HTTP server and client
- `encoding/json` - JSON encoding/decoding
- `flag` - CLI argument parsing
- `sync` - Mutex for thread safety
- `time` - Time tracking and sleeping
- `io` - Reading response bodies
- `fmt`, `log`, `os` - Output and logging
