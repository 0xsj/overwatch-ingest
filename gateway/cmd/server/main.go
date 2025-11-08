package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/scout/gateway/config"
	"github.com/0xsj/scout/platform/pkg/observability/logger"
	"github.com/0xsj/scout/platform/pkg/observability/logger/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load(false) // false = no prefix (matches current docker-compose)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger based on environment
	var appLogger logger.Logger
	if cfg.Environment() == "development" {
		// Development: colorized console output
		appLogger = logger.NewNoop()
	} else {
		// Production: structured JSON logging
		level := logger.ParseLevel(cfg.LogLevel())
		appLogger, err = zap.New(level)
		if err != nil {
			log.Fatalf("Failed to create logger: %v", err)
		}
	}

	// Add service-level fields
	appLogger = appLogger.With(
		"service", "gateway",
		"environment", cfg.Environment(),
		"version", "0.1.0",
	)

	// Log startup
	appLogger.Info("gateway server starting",
		"port", cfg.Port(),
		"log_level", cfg.LogLevel(),
		"postgres_host", cfg.Postgres().Host(),
		"redis_host", cfg.Redis().Host(),
		"nats_url", cfg.NATS().URL(),
	)

	// TODO: Initialize dependencies (database, cache, message bus)
	// TODO: Set up HTTP server
	// TODO: Register routes
	// TODO: Start server

	appLogger.Info("gateway server running",
		"port", cfg.Port(),
		"address", fmt.Sprintf("http://localhost:%d", cfg.Port()),
	)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	appLogger.Info("shutdown signal received",
		"signal", sig.String(),
	)

	// TODO: Graceful shutdown
	// - Stop accepting new requests
	// - Wait for in-flight requests to complete
	// - Close database connections
	// - Close message bus connections

	appLogger.Info("gateway server stopped")
}