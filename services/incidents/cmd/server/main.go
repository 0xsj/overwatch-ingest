package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Incidents server starting...")
	
	// TODO: Initialize dependencies (config, logger, database, etc.)
	// TODO: Set up gRPC/HTTP server
	// TODO: Register handlers
	// TODO: Start server
	
	fmt.Println("Incidents server running on :8082")
	
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	<-sigChan
	log.Println("Shutdown signal received, stopping server...")
}