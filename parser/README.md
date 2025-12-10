# Parser

A Go REST API service that fetches Hacker News top stories via a Rate Limiter service, parses the HTML, and returns structured JSON data.

## Overview

Parser is part of a larger system designed to obtain and structure Hacker News content. It does not interact with Hacker News directly; instead, it uses a Rate Limiter service to fetch HTML pages, which it then parses into structured data.

## Architecture

```
Client → Parser → Rate Limiter → Hacker News
                      ↓
               HTML Response
                      ↓
            Parser (parse HTML)
                      ↓
              JSON Response
```

## Building

```bash
go build -o parser
```

## Running

```bash
./parser --api <port> --ratelimiter <port> --num-pages <N>
```

### Required Arguments

| Argument | Description |
|----------|-------------|
| `--api` | Port number for the Parser REST API |
| `--ratelimiter` | Port number where the Rate Limiter is listening |
| `--num-pages` | Number of Hacker News pages to fetch (must be positive) |

### Example

```bash
# Start the Parser on port 8081, connecting to Rate Limiter on port 8080, fetching 2 pages
./parser --api 8081 --ratelimiter 8080 --num-pages 2
```

## REST API

### POST /fetch

Fetches N pages of Hacker News top stories and returns parsed data.

**Request:**
```bash
curl -X POST http://localhost:8081/fetch
```

**Response (200 OK):**
```json
{
  "fetched_at": "2025-12-06T10:30:00Z",
  "num_pages": 2,
  "total_stories": 60,
  "stories": [
    {
      "rank": 1,
      "id": "46173547",
      "headline": "Tiny Core Linux: a 23 MB Linux distro with graphical desktop",
      "url": "http://www.tinycorelinux.net/",
      "username": "LorenDB",
      "points": 224,
      "comments": 115,
      "discussion_url": "https://news.ycombinator.com/item?id=46173547",
      "age_value": 5,
      "age_unit": "hours",
      "page": 1
    }
  ]
}
```

**Error Response (502 Bad Gateway):**
```json
{
  "error": "Failed to fetch page 1: failed to reach rate limiter: connection refused"
}
```

### GET /doc

Returns API documentation in JSON format.

**Request:**
```bash
curl http://localhost:8081/doc
```

## Project Structure

```
parser/
├── main.go           # Entry point, CLI argument parsing, server setup
├── types.go          # Data structures (Story, FetchResponse, etc.)
├── handler.go        # HTTP handlers for /fetch and /doc endpoints
├── ratelimiter.go    # Rate Limiter client for fetching URLs
├── parser.go         # HTML parsing logic for Hacker News pages
├── go.mod            # Go module definition
├── go.sum            # Dependency checksums
└── README.md         # This file
```

## Code Overview

### main.go
- Parses and validates command line arguments
- Creates the Rate Limiter client
- Sets up HTTP routes and starts the server

### types.go
Defines data structures:
- `Story` - Represents a single HN story with all metadata
- `FetchResponse` - Response format for POST /fetch
- `RateLimiterRequest/Response` - Communication with Rate Limiter
- `ErrorResponse` - Error response format

### handler.go
HTTP handlers:
- `HandleFetch` - Orchestrates fetching pages, parsing, and responding
- `HandleDoc` - Returns API documentation

### ratelimiter.go
Rate Limiter client:
- `NewRateLimiterClient(port)` - Creates a new client
- `FetchURL(url)` - Requests the Rate Limiter to fetch a URL

### parser.go
HTML parsing using `golang.org/x/net/html`:
- `ParseHNPage(html, pageNum)` - Parses a full HN page into stories
- Helper functions for DOM traversal and text extraction

## HTML Parsing Details

The parser extracts data from Hacker News HTML structure:

1. **Story rows** are `<tr class="athing submission">` elements
2. **Story ID** comes from the row's `id` attribute
3. **Rank** is in `<span class="rank">`
4. **Headline and URL** are in `<span class="titleline"><a>`
5. **Subtext row** (next sibling) contains:
   - Points: `<span class="score">`
   - Username: `<a class="hnuser">`
   - Age: `<span class="age"><a>`
   - Comments: Last `<a>` containing "comment" or "discuss"

## Dependencies

- `golang.org/x/net/html` - HTML parsing

## Debugging

### Verify Rate Limiter connectivity
```bash
curl -X POST http://localhost:8080/fetch \
  -H "Content-Type: application/json" \
  -d '{"url": "https://news.ycombinator.com/"}'
```

### Test Parser endpoints
```bash
# Test /doc endpoint
curl http://localhost:8081/doc | jq .

# Test /fetch endpoint
curl -X POST http://localhost:8081/fetch | jq .

# Get specific story data
curl -X POST http://localhost:8081/fetch | jq '.stories[0]'
```

### Verify parsed data against actual HN
```bash
# Fetch actual HN HTML
curl -s "https://news.ycombinator.com/" > hn.html

# Compare specific story IDs, headlines, points, etc.
grep 'id="score_' hn.html | head -5
```

## Error Handling

| Status | Cause |
|--------|-------|
| 405 | Method not allowed (e.g., GET on /fetch) |
| 500 | HTML parsing error |
| 502 | Rate Limiter unreachable or returned error |

## Testing Checklist

1. Verify Parser starts without errors
2. Verify GET /doc returns valid JSON
3. Verify POST /fetch returns stories
4. Verify story count matches `num-pages * 30` (approximately)
5. Verify parsed data matches actual HN page content
6. Verify error handling when Rate Limiter is down
