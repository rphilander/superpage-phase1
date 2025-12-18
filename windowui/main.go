package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	browserPort := flag.Int("browser", 0, "Port for browser HTTP requests (required)")
	windowViewerPort := flag.Int("windowviewer", 0, "Port for Window Viewer API (required)")
	flag.Parse()

	if *browserPort == 0 || *windowViewerPort == 0 {
		fmt.Fprintln(os.Stderr, "Error: --browser and --windowviewer flags are required")
		flag.Usage()
		os.Exit(1)
	}

	client := NewWindowViewerClient(*windowViewerPort)

	handlers, err := NewHandlers(client)
	if err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/api/stories", handlers.HandleAPIStories)

	addr := fmt.Sprintf(":%d", *browserPort)
	log.Printf("Starting WindowUI on http://localhost%s", addr)
	log.Printf("Using Window Viewer at localhost:%d", *windowViewerPort)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
