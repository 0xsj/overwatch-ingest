// platform/pkg/errors/error.go
package errors

import (
	"errors"
	"fmt"
)

// Error is the core error type with rich metadata.
// It implements the error interface and supports error wrapping.
type Error struct {
	// Type categorizes the error (required)
	Type ErrorType

	// Code provides specific identification (required)
	Code Code

	// Message is a human-readable description (required)
	Message string

	// Details contains additional context as key-value pairs (optional)
	Details map[string]string

	// Cause is the underlying error that caused this error (optional)
	Cause error

	// stack contains captured stack trace program counters (optional)
	stack []uintptr
}

// New creates a new Error with the given type, code, and message.
// This is the primary constructor for creating errors.
func New(errType ErrorType, code Code, message string) *Error {
	return &Error{
		Type:    errType,
		Code:    code,
		Message: message,
	}
}

// Newf creates a new Error with a formatted message.
func Newf(errType ErrorType, code Code, format string, args ...interface{}) *Error {
	return &Error{
		Type:    errType,
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Error implements the error interface.
// Returns a formatted string representation of the error.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	// Build error string
	msg := fmt.Sprintf("[%s:%s] %s", e.Type, e.Code, e.Message)

	// Add cause if present
	if e.Cause != nil {
		msg += fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	return msg
}

// Unwrap returns the underlying cause.
// This enables errors.Is and errors.As to work with wrapped errors.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// Is checks if this error matches the target error.
// Matches if Type and Code are the same.
func (e *Error) Is(target error) bool {
	// Both nil should match
	if e == nil && target == nil {
		return true
	}

	// One nil, one not nil - no match
	if e == nil || target == nil {
		return false
	}

	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.Type == t.Type && e.Code == t.Code
}

// WithCause wraps an underlying error as the cause.
// Returns the error for method chaining.
// Nil-safe: returns nil if e is nil.
func (e *Error) WithCause(cause error) *Error {
	if e == nil {
		return nil
	}
	e.Cause = cause
	return e
}

// WithDetail adds a single key-value detail.
// Returns the error for method chaining.
// Nil-safe: returns nil if e is nil.
func (e *Error) WithDetail(key, value string) *Error {
	if e == nil {
		return nil
	}

	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	e.Details[key] = value
	return e
}

// WithDetails adds multiple details at once.
// Returns the error for method chaining.
// Nil-safe: returns nil if e is nil.
func (e *Error) WithDetails(details map[string]string) *Error {
	if e == nil {
		return nil
	}

	if e.Details == nil {
		e.Details = make(map[string]string, len(details))
	}
	for k, v := range details {
		e.Details[k] = v
	}
	return e
}

// GetDetail retrieves a detail value by key.
// Returns empty string if key doesn't exist.
func (e *Error) GetDetail(key string) string {
	if e == nil || e.Details == nil {
		return ""
	}
	return e.Details[key]
}

// HasDetail returns true if the error has a detail with the given key.
func (e *Error) HasDetail(key string) bool {
	if e == nil || e.Details == nil {
		return false
	}
	_, exists := e.Details[key]
	return exists
}

// HasDetails returns true if the error has any details attached.
func (e *Error) HasDetails() bool {
	return e != nil && len(e.Details) > 0
}

// IsClientError returns true if this is a client error (4xx).
func (e *Error) IsClientError() bool {
	if e == nil {
		return false
	}
	return e.Type.IsClientError()
}

// IsServerError returns true if this is a server error (5xx).
func (e *Error) IsServerError() bool {
	if e == nil {
		return false
	}
	return e.Type.IsServerError()
}

// IsRetryable returns true if this error can typically be retried.
func (e *Error) IsRetryable() bool {
	if e == nil {
		return false
	}
	return e.Type.IsRetryable()
}

// HTTPStatusCode returns the recommended HTTP status code.
func (e *Error) HTTPStatusCode() int {
	if e == nil {
		return 500
	}
	return e.Type.HTTPStatusCode()
}

// As attempts to convert a standard error to *Error.
// Returns the Error and true if successful, nil and false otherwise.
func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// Is checks if err matches target using errors.Is semantics.
// This is a convenience wrapper around the standard library errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}
