package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Define command line flags
	apiPort := flag.Int("api", 0, "Port number for the Parser REST API (required)")
	rateLimiterPort := flag.Int("ratelimiter", 0, "Port number where the Rate Limiter is listening (required)")
	numPages := flag.Int("num-pages", 0, "Number of Hacker News pages to fetch (required, must be positive)")

	flag.Parse()

	// Validate required arguments
	var errors []string

	if *apiPort == 0 {
		errors = append(errors, "--api port is required")
	} else if *apiPort < 1 || *apiPort > 65535 {
		errors = append(errors, "--api port must be between 1 and 65535")
	}

	if *rateLimiterPort == 0 {
		errors = append(errors, "--ratelimiter port is required")
	} else if *rateLimiterPort < 1 || *rateLimiterPort > 65535 {
		errors = append(errors, "--ratelimiter port must be between 1 and 65535")
	}

	if *numPages == 0 {
		errors = append(errors, "--num-pages is required")
	} else if *numPages < 1 {
		errors = append(errors, "--num-pages must be a positive integer")
	}

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr, "Error: Invalid arguments")
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create Rate Limiter client
	rateLimiter := NewRateLimiterClient(*rateLimiterPort)

	// Create handler
	handler := NewHandler(rateLimiter, *numPages)

	// Set up routes
	http.HandleFunc("/fetch", handler.HandleFetch)
	http.HandleFunc("/doc", handler.HandleDoc)

	// Start server
	addr := fmt.Sprintf(":%d", *apiPort)
	log.Printf("Parser starting on port %d", *apiPort)
	log.Printf("Rate Limiter configured at localhost:%d", *rateLimiterPort)
	log.Printf("Configured to fetch %d page(s) from Hacker News", *numPages)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
