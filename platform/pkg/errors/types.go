// platform/pkg/errors/types.go
package errors

import "fmt"

// ErrorType categorizes errors by their nature.
// Types are intentionally generic to support all platform services.
type ErrorType string

const (
	// Client Errors (4xx equivalent) - caused by invalid client input or state
	// These errors typically should not be retried without changing the request.
	ErrorTypeValidation    ErrorType = "VALIDATION"      // Invalid input data
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"       // Resource does not exist
	ErrorTypeAlreadyExists ErrorType = "ALREADY_EXISTS"  // Resource already exists (conflict)
	ErrorTypeUnauthorized  ErrorType = "UNAUTHORIZED"    // Missing or invalid authentication
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"       // Insufficient permissions
	ErrorTypeConflict      ErrorType = "CONFLICT"        // State conflict (e.g., optimistic locking)
	ErrorTypeRateLimit     ErrorType = "RATE_LIMIT"      // Rate limit exceeded

	// Server Errors (5xx equivalent) - caused by server-side issues
	// These errors may be retryable depending on the specific type.
	ErrorTypeInternal       ErrorType = "INTERNAL"        // Internal server error
	ErrorTypeUnavailable    ErrorType = "UNAVAILABLE"     // Service/dependency unavailable
	ErrorTypeTimeout        ErrorType = "TIMEOUT"         // Operation timeout
	ErrorTypeNotImplemented ErrorType = "NOT_IMPLEMENTED" // Feature not implemented

	// Infrastructure Errors - platform component failures
	ErrorTypeDatabase ErrorType = "DATABASE" // Database operation failed
	ErrorTypeCache    ErrorType = "CACHE"    // Cache operation failed
	ErrorTypeNetwork  ErrorType = "NETWORK"  // Network communication failed
	ErrorTypeEvent    ErrorType = "EVENT"    // Event bus operation failed
)

// String returns the string representation of the ErrorType.
func (t ErrorType) String() string {
	return string(t)
}

// IsValid returns true if this is a recognized ErrorType.
func (t ErrorType) IsValid() bool {
	switch t {
	case ErrorTypeValidation,
		ErrorTypeNotFound,
		ErrorTypeAlreadyExists,
		ErrorTypeUnauthorized,
		ErrorTypeForbidden,
		ErrorTypeConflict,
		ErrorTypeRateLimit,
		ErrorTypeInternal,
		ErrorTypeUnavailable,
		ErrorTypeTimeout,
		ErrorTypeNotImplemented,
		ErrorTypeDatabase,
		ErrorTypeCache,
		ErrorTypeNetwork,
		ErrorTypeEvent:
		return true
	default:
		return false
	}
}

// IsClientError returns true if the error is caused by client input (4xx equivalent).
// These errors typically should not be retried without changing the request.
func (t ErrorType) IsClientError() bool {
	switch t {
	case ErrorTypeValidation,
		ErrorTypeNotFound,
		ErrorTypeAlreadyExists,
		ErrorTypeUnauthorized,
		ErrorTypeForbidden,
		ErrorTypeConflict,
		ErrorTypeRateLimit:
		return true
	default:
		return false
	}
}

// IsServerError returns true if the error is caused by server-side issues (5xx equivalent).
// These errors may be retryable depending on the specific type.
func (t ErrorType) IsServerError() bool {
	switch t {
	case ErrorTypeInternal,
		ErrorTypeUnavailable,
		ErrorTypeTimeout,
		ErrorTypeNotImplemented,
		ErrorTypeDatabase,
		ErrorTypeCache,
		ErrorTypeNetwork,
		ErrorTypeEvent:
		return true
	default:
		return false
	}
}

// IsRetryable returns true if errors of this type can typically be retried.
// Note: Even retryable errors should use exponential backoff and respect retry limits.
func (t ErrorType) IsRetryable() bool {
	switch t {
	case ErrorTypeTimeout,
		ErrorTypeUnavailable,
		ErrorTypeRateLimit,
		ErrorTypeInternal,
		ErrorTypeNetwork,
		ErrorTypeCache,
		ErrorTypeDatabase,
		ErrorTypeEvent:
		return true
	default:
		return false
	}
}

// HTTPStatusCode returns the recommended HTTP status code for this error type.
func (t ErrorType) HTTPStatusCode() int {
	switch t {
	case ErrorTypeValidation:
		return 400 // Bad Request
	case ErrorTypeUnauthorized:
		return 401 // Unauthorized
	case ErrorTypeForbidden:
		return 403 // Forbidden
	case ErrorTypeNotFound:
		return 404 // Not Found
	case ErrorTypeConflict, ErrorTypeAlreadyExists:
		return 409 // Conflict
	case ErrorTypeRateLimit:
		return 429 // Too Many Requests
	case ErrorTypeInternal:
		return 500 // Internal Server Error
	case ErrorTypeNotImplemented:
		return 501 // Not Implemented
	case ErrorTypeUnavailable, ErrorTypeDatabase, ErrorTypeCache:
		return 503 // Service Unavailable
	case ErrorTypeTimeout, ErrorTypeNetwork:
		return 504 // Gateway Timeout
	default:
		return 500 // Internal Server Error (safe default)
	}
}

// ParseErrorType converts a string to an ErrorType.
// Returns zero value if the string is not a valid ErrorType.
func ParseErrorType(s string) ErrorType {
	t := ErrorType(s)
	if t.IsValid() {
		return t
	}
	return ""
}

// MustParseErrorType converts a string to an ErrorType.
// Panics if the string is not a valid ErrorType.
// Only use this for constants known at compile time.
func MustParseErrorType(s string) ErrorType {
	t := ErrorType(s)
	if !t.IsValid() {
		panic(fmt.Sprintf("invalid error type: %s", s))
	}
	return t
}

// AllErrorTypes returns all valid ErrorTypes.
// Useful for testing and validation.
func AllErrorTypes() []ErrorType {
	return []ErrorType{
		ErrorTypeValidation,
		ErrorTypeNotFound,
		ErrorTypeAlreadyExists,
		ErrorTypeUnauthorized,
		ErrorTypeForbidden,
		ErrorTypeConflict,
		ErrorTypeRateLimit,
		ErrorTypeInternal,
		ErrorTypeUnavailable,
		ErrorTypeTimeout,
		ErrorTypeNotImplemented,
		ErrorTypeDatabase,
		ErrorTypeCache,
		ErrorTypeNetwork,
		ErrorTypeEvent,
	}
}