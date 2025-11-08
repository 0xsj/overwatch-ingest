// platform/pkg/observability/logger/doc.go

/*
Package logger provides structured logging for the Scout platform.

# Overview

This package defines a logging interface and provides implementations for
different environments:

- noop.Logger - Colorized console output for local development
- zap.Logger - Production-grade structured logging with JSON output

# Basic Usage

Create a logger based on your environment:

	import (
		"github.com/0xsj/scout/platform/pkg/observability/logger"
		"github.com/0xsj/scout/platform/pkg/observability/logger/zap"
	)

	// Development: colorized console
	log := logger.NewNoop()

	// Production: JSON structured logging
	log, err := zap.New(logger.InfoLevel)
	if err != nil {
		panic(err)
	}

	// Log messages with structured fields
	log.Info("server started",
		"port", 8080,
		"environment", "production",
	)

# Structured Logging

All loggers use key-value pairs for structured logging:

	log.Info("user created",
		"user_id", "123",
		"email", "user@example.com",
		"tenant", "acme",
	)

	// Output (JSON in production):
	// {"level":"info","timestamp":"2025-01-01T12:00:00Z","message":"user created","user_id":"123","email":"user@example.com","tenant":"acme"}

# Log Levels

The package supports four log levels:

	logger.DebugLevel - Detailed debugging information
	logger.InfoLevel  - General informational messages
	logger.WarnLevel  - Warning messages for potentially harmful situations
	logger.ErrorLevel - Error messages for serious problems

Set the level when creating the logger:

	log, _ := zap.New(logger.DebugLevel)  // Show all messages
	log, _ := zap.New(logger.ErrorLevel)  // Only errors

# Attaching Fields

Create child loggers with persistent fields:

	// Request-scoped logger
	reqLogger := log.With(
		"request_id", "req-123",
		"user_id", "user-456",
	)

	reqLogger.Info("processing request")  // Includes request_id and user_id
	reqLogger.Info("request completed")   // Still includes those fields

# Error Logging

Attach errors to log messages:

	if err != nil {
		log.WithError(err).Error("failed to process request",
			"user_id", userID,
		)
	}

# Context Integration

Loggers can extract values from context (useful for trace IDs, request IDs):

	ctxLogger := log.WithContext(ctx)
	ctxLogger.Info("handling request")  // Includes values from context

# Production vs Development

For production environments, use zap with JSON encoding:

	log, err := zap.New(logger.InfoLevel)

For development/local testing, use the colorized noop logger:

	log := logger.NewNoop()

Or use zap's development mode:

	log, err := zap.NewDevelopment(logger.DebugLevel)

# Example: Service Logger

	package main

	import (
		"github.com/0xsj/scout/platform/pkg/observability/logger"
		"github.com/0xsj/scout/platform/pkg/observability/logger/zap"
	)

	func main() {
		// Create production logger
		log, err := zap.New(logger.InfoLevel)
		if err != nil {
			panic(err)
		}

		// Add service-level fields
		log = log.With(
			"service", "gateway",
			"version", "1.0.0",
		)

		log.Info("service starting", "port", 8080)

		// Create request-scoped logger
		handleRequest(log, "req-123")

		log.Info("service stopped")
	}

	func handleRequest(log logger.Logger, requestID string) {
		reqLog := log.With("request_id", requestID)

		reqLog.Info("handling request")
		reqLog.Debug("validating input")
		reqLog.Info("request completed", "duration_ms", 42)
	}
*/
package logger