package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xsj/scout/gateway/config"
	"github.com/0xsj/scout/platform/pkg/observability/logger"
	"github.com/0xsj/scout/platform/pkg/observability/logger/zap"
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
		"service", "gateway",
		"environment", cfg.Environment(),
		"version", "0.1.0",
	)

	// Create HTTP server
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", handleHealth(appLogger))
	
	// Root endpoint
	mux.HandleFunc("/", handleRoot(appLogger))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port()),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		appLogger.Info("gateway server starting",
			"port", cfg.Port(),
			"address", fmt.Sprintf("http://localhost:%d", cfg.Port()),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("server failed to start", "error", err.Error())
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	appLogger.Info("shutdown signal received", "signal", sig.String())

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("server shutdown failed", "error", err.Error())
	}

	appLogger.Info("gateway server stopped")
}

// handleHealth returns a health check handler
func handleHealth(log logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("health check requested",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		response := map[string]interface{}{
			"status":  "healthy",
			"service": "gateway",
			"version": "0.1.0",
		}
		
		json.NewEncoder(w).Encode(response)
		
		log.Info("health check completed", "status", "healthy")
	}
}

// handleRoot returns a root handler
func handleRoot(log logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("request received",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		response := map[string]string{
			"message": "Scout Gateway API",
			"version": "0.1.0",
		}
		
		json.NewEncoder(w).Encode(response)
	}
}