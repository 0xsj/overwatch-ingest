// platform/pkg/observability/logger/noop.go
package logger

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorWhite  = "\033[97m"
)

// noopLogger is a simple console logger with colored output.
// Useful for local development and debugging.
type noopLogger struct {
	fields []any
}

// NewNoop creates a new console logger with colored output.
// Great for local development - logs are human-readable with colors.
func NewNoop() Logger {
	return &noopLogger{}
}

func (l *noopLogger) Debug(msg string, keysAndValues ...any) {
	l.log(DebugLevel, msg, keysAndValues...)
}

func (l *noopLogger) Info(msg string, keysAndValues ...any) {
	l.log(InfoLevel, msg, keysAndValues...)
}

func (l *noopLogger) Warn(msg string, keysAndValues ...any) {
	l.log(WarnLevel, msg, keysAndValues...)
}

func (l *noopLogger) Error(msg string, keysAndValues ...any) {
	l.log(ErrorLevel, msg, keysAndValues...)
}

func (l *noopLogger) With(keysAndValues ...any) Logger {
	newFields := make([]any, len(l.fields)+len(keysAndValues))
	copy(newFields, l.fields)
	copy(newFields[len(l.fields):], keysAndValues)
	return &noopLogger{fields: newFields}
}

func (l *noopLogger) WithContext(ctx context.Context) Logger {
	// Could extract values from context if needed
	return l
}

func (l *noopLogger) WithError(err error) Logger {
	return l.With("error", err.Error())
}

// log outputs a colored log message to stdout
func (l *noopLogger) log(level Level, msg string, keysAndValues ...any) {
	// Timestamp
	timestamp := time.Now().Format("15:04:05.000")
	
	// Level with color
	levelStr := l.colorizeLevel(level)
	
	// Message
	msgStr := fmt.Sprintf("%s%s%s", colorWhite, msg, colorReset)
	
	// Combine persistent fields with new fields
	allFields := append(l.fields, keysAndValues...)
	
	// Format fields
	fieldsStr := l.formatFields(allFields)
	
	// Output
	if fieldsStr != "" {
		fmt.Printf("%s%s%s %s %s %s\n",
			colorGray, timestamp, colorReset,
			levelStr,
			msgStr,
			fieldsStr,
		)
	} else {
		fmt.Printf("%s%s%s %s %s\n",
			colorGray, timestamp, colorReset,
			levelStr,
			msgStr,
		)
	}
}

// colorizeLevel returns a colored level string
func (l *noopLogger) colorizeLevel(level Level) string {
	switch level {
	case DebugLevel:
		return fmt.Sprintf("%sDEBUG%s", colorCyan, colorReset)
	case InfoLevel:
		return fmt.Sprintf("%sINFO %s", colorGreen, colorReset)
	case WarnLevel:
		return fmt.Sprintf("%sWARN %s", colorYellow, colorReset)
	case ErrorLevel:
		return fmt.Sprintf("%sERROR%s", colorRed, colorReset)
	default:
		return "UNKNOWN"
	}
}

// formatFields formats key-value pairs with colors
func (l *noopLogger) formatFields(keysAndValues []any) string {
	if len(keysAndValues) == 0 {
		return ""
	}
	
	var parts []string
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 >= len(keysAndValues) {
			break
		}
		
		key := keysAndValues[i]
		value := keysAndValues[i+1]
		
		// Color the key
		keyStr := fmt.Sprintf("%s%v%s", colorPurple, key, colorReset)
		
		// Format the value
		valueStr := fmt.Sprintf("%v", value)
		
		parts = append(parts, fmt.Sprintf("%s=%s", keyStr, valueStr))
	}
	
	return strings.Join(parts, " ")
}