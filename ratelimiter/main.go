package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"ratelimiter/api"
	"ratelimiter/limiter"
)

func main() {
	rate := flag.Int("rate", 0, "Minimum number of seconds between HTTP requests (required, must be positive)")
	port := flag.Int("api", 0, "Port number for the REST API (required)")
	flag.Parse()

	// Validate required arguments
	if *rate <= 0 {
		fmt.Fprintln(os.Stderr, "Error: --rate must be a positive integer")
		flag.Usage()
		os.Exit(1)
	}

	if *port <= 0 || *port > 65535 {
		fmt.Fprintln(os.Stderr, "Error: --api must be a valid port number (1-65535)")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize the rate limiter
	rl := limiter.New(*rate)

	// Initialize the API handler
	handler := api.NewHandler(rl)

	// Set up routes
	http.HandleFunc("/fetch", handler.HandleFetch)
	http.HandleFunc("/doc", handler.HandleDoc)

	// Start the server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting Rate Limiter API server on port %d (rate limit: 1 request per %d seconds)", *port, *rate)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
