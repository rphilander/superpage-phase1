package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"querier/api"
	"querier/parser"
	"querier/store"
)

func main() {
	apiPort := flag.Int("api", 0, "Port number for the Querier API (required)")
	parserPort := flag.Int("parser", 0, "Port number where the Parser is listening (required)")
	flag.Parse()

	if *apiPort == 0 || *parserPort == 0 {
		fmt.Fprintln(os.Stderr, "Error: both --api and --parser arguments are required")
		fmt.Fprintln(os.Stderr, "Usage: querier --api <port> --parser <port>")
		os.Exit(1)
	}

	// Initialize components
	dataStore := store.New()
	parserClient := parser.NewClient(*parserPort)
	handler := api.NewHandler(dataStore, parserClient)

	// Register routes
	http.HandleFunc("/refresh", handler.HandleRefresh)
	http.HandleFunc("/query", handler.HandleQuery)
	http.HandleFunc("/doc", handler.HandleDoc)

	addr := fmt.Sprintf(":%d", *apiPort)
	log.Printf("Querier starting on port %d (Parser at port %d)", *apiPort, *parserPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
