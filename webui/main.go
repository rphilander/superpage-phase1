package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"webui/server"
)

func main() {
	browserPort := flag.Int("browser", 0, "Port number for browser HTTP requests (required)")
	querierPort := flag.Int("querier", 0, "Port number where Querier is listening (required)")

	flag.Parse()

	if *browserPort == 0 {
		fmt.Fprintln(os.Stderr, "Error: --browser port is required")
		flag.Usage()
		os.Exit(1)
	}

	if *querierPort == 0 {
		fmt.Fprintln(os.Stderr, "Error: --querier port is required")
		flag.Usage()
		os.Exit(1)
	}

	srv := server.NewServer(*browserPort, *querierPort)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
