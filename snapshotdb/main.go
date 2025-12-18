package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"snapshotdb/api"
	"snapshotdb/config"
	"snapshotdb/logger"
	"snapshotdb/parser"
	"snapshotdb/scheduler"
	"snapshotdb/store"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: snapshotdb --api <port> --db <path> --parser <port> --freq <seconds>\n")
		os.Exit(1)
	}

	log.Printf("Starting SnapshotDB...")
	log.Printf("  API port: %d", cfg.APIPort)
	log.Printf("  Database: %s", cfg.DBPath)
	log.Printf("  Parser: localhost:%d", cfg.ParserPort)
	log.Printf("  Frequency: %d seconds", cfg.FreqSecs)

	st, err := store.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}
	defer st.Close()

	errLogger := logger.New(cfg.ErrorLogPath())
	parserClient := parser.NewClient(cfg.ParserURL())
	sched := scheduler.New(st, parserClient, errLogger, cfg.FreqSecs)

	handler := api.NewHandler(st, sched)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.APIPort),
		Handler: mux,
	}

	sched.Start()
	log.Printf("Scheduler started")

	go func() {
		log.Printf("HTTP server listening on port %d", cfg.APIPort)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("Received signal %v, shutting down...", sig)

	sched.Stop()
	log.Printf("Scheduler stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	log.Printf("HTTP server stopped")

	log.Printf("SnapshotDB shutdown complete")
}
