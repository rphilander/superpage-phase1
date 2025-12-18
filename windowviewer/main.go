package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"windowviewer/client"
	"windowviewer/handlers"
)

func main() {
	// Define command line flags
	apiPort := flag.Int("api", 0, "Port number for Window Viewer's HTTP API (required)")
	snapshotDBPort := flag.Int("snapshotdb", 0, "Port number where SnapshotDB is listening on localhost (required)")

	flag.Parse()

	// Validate required flags
	if *apiPort == 0 {
		fmt.Fprintln(os.Stderr, "Error: --api flag is required")
		flag.Usage()
		os.Exit(1)
	}
	if *snapshotDBPort == 0 {
		fmt.Fprintln(os.Stderr, "Error: --snapshotdb flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create SnapshotDB client
	snapshotDB := client.NewSnapshotDBClient(*snapshotDBPort)

	// Create handlers
	h := handlers.NewHandler(snapshotDB)

	// Setup routes
	http.HandleFunc("/top", h.TopHandler)
	http.HandleFunc("/doc", h.DocHandler)

	// Start server
	addr := fmt.Sprintf(":%d", *apiPort)
	log.Printf("Window Viewer starting on port %d", *apiPort)
	log.Printf("Connected to SnapshotDB on port %d", *snapshotDBPort)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
