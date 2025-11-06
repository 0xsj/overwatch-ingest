// platform/pkg/errors/stack.go
package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

const (
	// MaxStackDepth is the maximum number of stack frames to capture.
	MaxStackDepth = 32

	// SkipFrames is the number of frames to skip when capturing the stack.
	// Skips: runtime.Callers, captureStack, WithStack
	SkipFrames = 3
)

// captureStack captures the current stack trace.
// Returns program counter values for stack frames.
func captureStack(skip int) []uintptr {
	var pcs [MaxStackDepth]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[:n]
}

// WithStack captures the current stack trace and attaches it to the error.
// This is opt-in to avoid performance overhead when stack traces aren't needed.
// Returns the error for method chaining.
// Nil-safe: returns nil if e is nil.
//
// Example:
//
//	return errors.Internal("unexpected error").WithStack()
func (e *Error) WithStack() *Error {
	if e == nil {
		return nil
	}

	// Only capture if not already captured
	if e.stack == nil {
		e.stack = captureStack(SkipFrames)
	}

	return e
}

// HasStack returns true if a stack trace has been captured.
func (e *Error) HasStack() bool {
	return e != nil && len(e.stack) > 0
}

// StackTrace returns the stack trace as formatted strings.
// Each string is in the format: "function (file:line)"
// Returns nil if no stack trace has been captured.
//
// Example output:
//
//	main.processIncident (main.go:45)
//	main.handleRequest (main.go:30)
//	net/http.HandlerFunc.ServeHTTP (server.go:2109)
func (e *Error) StackTrace() []string {
	if !e.HasStack() {
		return nil
	}

	frames := runtime.CallersFrames(e.stack)
	var trace []string

	for {
		frame, more := frames.Next()

		// Format: function (file:line)
		trace = append(trace, fmt.Sprintf("%s (%s:%d)",
			frame.Function,
			frame.File,
			frame.Line,
		))

		if !more {
			break
		}
	}

	return trace
}

// StackTraceString returns the stack trace as a single formatted string.
// Each frame is on a new line with indentation.
// Returns empty string if no stack trace has been captured.
//
// Example output:
//
//	at main.processIncident (main.go:45)
//	at main.handleRequest (main.go:30)
//	at net/http.HandlerFunc.ServeHTTP (server.go:2109)
func (e *Error) StackTraceString() string {
	frames := e.StackTrace()
	if len(frames) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, frame := range frames {
		sb.WriteString("  at ")
		sb.WriteString(frame)
		sb.WriteString("\n")
	}

	return sb.String()
}

// ErrorWithStack returns the error message with stack trace appended.
// Useful for logging the full error context.
//
// Example output:
//
//	[INTERNAL:INTERNAL_ERROR] unexpected error
//	  at main.processIncident (main.go:45)
//	  at main.handleRequest (main.go:30)
func (e *Error) ErrorWithStack() string {
	if e == nil {
		return ""
	}

	msg := e.Error()

	if e.HasStack() {
		msg += "\n" + e.StackTraceString()
	}

	return msg
}

// stackJSON is the JSON representation of a stack frame.
type stackJSON struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// MarshalStackJSON marshals the stack trace to JSON.
// Returns nil if no stack trace has been captured.
//
// Example output:
//
//	[
//	  {"function": "main.processIncident", "file": "main.go", "line": 45},
//	  {"function": "main.handleRequest", "file": "main.go", "line": 30}
//	]
func (e *Error) MarshalStackJSON() ([]byte, error) {
	if !e.HasStack() {
		return []byte("null"), nil
	}

	frames := runtime.CallersFrames(e.stack)
	var stackFrames []stackJSON

	for {
		frame, more := frames.Next()

		stackFrames = append(stackFrames, stackJSON{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})

		if !more {
			break
		}
	}

	return json.Marshal(stackFrames)
}

// WrapWithStack is a convenience function that wraps an error and captures the stack.
// Equivalent to Wrap(...).WithStack()
func WrapWithStack(err error, errType ErrorType, code Code, message string) *Error {
	return Wrap(err, errType, code, message).WithStack()
}

// WrapfWithStack wraps an error with a formatted message and captures the stack.
// Equivalent to Wrapf(...).WithStack()
func WrapfWithStack(err error, errType ErrorType, code Code, format string, args ...interface{}) *Error {
	return Wrapf(err, errType, code, format, args...).WithStack()
}
