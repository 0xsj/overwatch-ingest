package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Gateway worker starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("shutdown signal received, stopping worker...")
		cancel()
	}()

	// init dependencies

	// worker tick simulation
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Gateway worker running, processing every 5 seconds....")

	for {
		select {
		case <-ctx.Done():
			log.Println("Gateway worker stopped")
			return
		case <-ticker.C:
			log.Println("Gateway worker: processing tick...")
		}
	}
}
