package logger

import "context"

// Logger defines the logging interface for the Scout platform.
// It provides structured logging with levels and context awareness.
//
// All logging methods accept key-value pairs for structured logging.
// Keys should be strings, values can be any type.
//
// Example:
//
//	logger.Info("user created",
//	    "user_id", "123",
//	    "email", "user@example.com",
//	)
type Logger interface {
	// Debug logs a debug-level message with optional key-value pairs.
	Debug(msg string, keysAndValues ...any)

	// Info logs an info-level message with optional key-value pairs.
	Info(msg string, keysAndValues ...any)

	// Warn logs a warning-level message with optional key-value pairs.
	Warn(msg string, keysAndValues ...any)

	// Error logs an error-level message with optional key-value pairs.
	// Typically used for errors that should be investigated.
	Error(msg string, keysAndValues ...any)

	// With returns a new logger with the given key-value pairs attached.
	// These fields will be included in all subsequent log entries.
	//
	// Example:
	//
	//	userLogger := logger.With("user_id", "123", "tenant_id", "abc")
	//	userLogger.Info("action performed") // Will include user_id and tenant_id
	With(keysAndValues ...any) Logger

	// WithContext returns a new logger that extracts values from context.
	// This is useful for propagating request IDs, trace IDs, etc.
	WithContext(ctx context.Context) Logger

	// WithError returns a new logger with an error field attached.
	// The error will be logged with a standard "error" key.
	WithError(err error) Logger
}

// Field represents a key-value pair for structured logging.
// This is used internally by logger implementations.
type Field struct {
	Key   string
	Value any
}

// Level represents the log level.
type Level int

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "unknown"
	}
}

// ParseLevel parses a string into a log level.
func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}
