package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xsj/scout/platform/pkg/observability/logger"
	"github.com/0xsj/scout/platform/pkg/observability/logger/zap"
	"github.com/0xsj/scout/services/agents/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load(false)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	var appLogger logger.Logger
	if cfg.Environment() == "development" {
		appLogger = logger.NewNoop()
	} else {
		level := logger.ParseLevel(cfg.LogLevel())
		appLogger, err = zap.New(level)
		if err != nil {
			log.Fatalf("Failed to create logger: %v", err)
		}
	}

	// Add service-level fields
	appLogger = appLogger.With(
		"service", "agents",
		"environment", cfg.Environment(),
		"version", "0.1.0",
	)

	// Log startup
	appLogger.Info("agents server starting",
		"port", cfg.Port(),
		"log_level", cfg.LogLevel(),
		"postgres_host", cfg.Postgres().Host(),
		"redis_host", cfg.Redis().Host(),
		"nats_url", cfg.NATS().URL(),
	)

	appLogger.Info("agents server running",
		"port", cfg.Port(),
		"address", fmt.Sprintf("http://localhost:%d", cfg.Port()),
	)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	appLogger.Info("shutdown signal received", "signal", sig.String())

	appLogger.Info("agents server stopped")
}