// platform/pkg/errors/constructors.go
package errors

import (
	"fmt"
	"time"
)

// =========================================
// Generic Error Constructors
// =========================================

// NotFound creates a not found error for a resource.
func NotFound(resourceType, resourceID string) *Error {
	return New(
		ErrorTypeNotFound,
		"RESOURCE_NOT_FOUND",
		fmt.Sprintf("%s not found: %s", resourceType, resourceID),
	).WithDetail("resource_type", resourceType).
		WithDetail("resource_id", resourceID)
}

// AlreadyExists creates a conflict error for a resource that already exists.
func AlreadyExists(resourceType, resourceID string) *Error {
	return New(
		ErrorTypeAlreadyExists,
		"RESOURCE_ALREADY_EXISTS",
		fmt.Sprintf("%s already exists: %s", resourceType, resourceID),
	).WithDetail("resource_type", resourceType).
		WithDetail("resource_id", resourceID)
}

// Validation creates a validation error with a message.
func Validation(message string) *Error {
	return New(
		ErrorTypeValidation,
		"VALIDATION_FAILED",
		message,
	)
}

// ValidationWithField creates a validation error for a specific field.
func ValidationWithField(field, message string) *Error {
	return New(
		ErrorTypeValidation,
		"VALIDATION_FAILED",
		fmt.Sprintf("validation failed for field '%s': %s", field, message),
	).WithDetail("field", field)
}

// RequiredField creates a validation error for a missing required field.
func RequiredField(field string) *Error {
	return New(
		ErrorTypeValidation,
		"REQUIRED_FIELD_MISSING",
		fmt.Sprintf("required field '%s' is missing", field),
	).WithDetail("field", field)
}

// InvalidField creates a validation error for an invalid field value.
func InvalidField(field, reason string) *Error {
	return New(
		ErrorTypeValidation,
		"INVALID_FIELD_VALUE",
		fmt.Sprintf("invalid value for field '%s': %s", field, reason),
	).WithDetail("field", field).
		WithDetail("reason", reason)
}

// =========================================
// Authorization Error Constructors
// =========================================

// Unauthorized creates an unauthorized error.
func Unauthorized(reason string) *Error {
	return New(
		ErrorTypeUnauthorized,
		"UNAUTHORIZED",
		fmt.Sprintf("unauthorized: %s", reason),
	).WithDetail("reason", reason)
}

// Forbidden creates a forbidden error.
func Forbidden(resource, action string) *Error {
	return New(
		ErrorTypeForbidden,
		"FORBIDDEN",
		fmt.Sprintf("forbidden: cannot %s %s", action, resource),
	).WithDetail("resource", resource).
		WithDetail("action", action)
}

// =========================================
// Internal Error Constructors
// =========================================

// Internal creates a generic internal error.
func Internal(message string) *Error {
	return New(
		ErrorTypeInternal,
		"INTERNAL_ERROR",
		message,
	)
}

// InternalWithCause creates an internal error with a wrapped cause.
func InternalWithCause(message string, cause error) *Error {
	return New(
		ErrorTypeInternal,
		"INTERNAL_ERROR",
		message,
	).WithCause(cause)
}

// Internalf creates an internal error with a formatted message.
func Internalf(format string, args ...interface{}) *Error {
	return New(
		ErrorTypeInternal,
		"INTERNAL_ERROR",
		fmt.Sprintf(format, args...),
	)
}

// NotImplemented creates a not implemented error.
func NotImplemented(feature string) *Error {
	return New(
		ErrorTypeNotImplemented,
		"NOT_IMPLEMENTED",
		fmt.Sprintf("feature not implemented: %s", feature),
	).WithDetail("feature", feature)
}

// =========================================
// Timeout Error Constructors
// =========================================

// Timeout creates a timeout error.
func Timeout(operation string, duration time.Duration) *Error {
	return New(
		ErrorTypeTimeout,
		"OPERATION_TIMEOUT",
		fmt.Sprintf("operation timed out: %s", operation),
	).WithDetail("operation", operation).
		WithDetail("timeout", duration.String())
}

// TimeoutWithDuration creates a timeout error with milliseconds.
func TimeoutWithDuration(operation string, durationMS uint64) *Error {
	return New(
		ErrorTypeTimeout,
		"OPERATION_TIMEOUT",
		fmt.Sprintf("operation timed out after %dms: %s", durationMS, operation),
	).WithDetail("operation", operation).
		WithDetail("timeout_ms", fmt.Sprintf("%d", durationMS))
}

// =========================================
// Service Availability Constructors
// =========================================

// Unavailable creates a service unavailable error.
func Unavailable(service string) *Error {
	return New(
		ErrorTypeUnavailable,
		"SERVICE_UNAVAILABLE",
		fmt.Sprintf("service unavailable: %s", service),
	).WithDetail("service", service)
}

// UnavailableWithCause creates a service unavailable error with a cause.
func UnavailableWithCause(service string, cause error) *Error {
	return New(
		ErrorTypeUnavailable,
		"SERVICE_UNAVAILABLE",
		fmt.Sprintf("service unavailable: %s", service),
	).WithDetail("service", service).
		WithCause(cause)
}

// =========================================
// Conflict Error Constructors
// =========================================

// Conflict creates a conflict error.
func Conflict(resource, reason string) *Error {
	return New(
		ErrorTypeConflict,
		"RESOURCE_CONFLICT",
		fmt.Sprintf("conflict on %s: %s", resource, reason),
	).WithDetail("resource", resource).
		WithDetail("reason", reason)
}

// =========================================
// Rate Limiting Constructors
// =========================================

// RateLimit creates a rate limit exceeded error.
func RateLimit(limit uint64, window string) *Error {
	return New(
		ErrorTypeRateLimit,
		"RATE_LIMIT_EXCEEDED",
		fmt.Sprintf("rate limit exceeded: %d requests per %s", limit, window),
	).WithDetail("limit", fmt.Sprintf("%d", limit)).
		WithDetail("window", window)
}

// RateLimitWithRetry creates a rate limit error with retry-after information.
func RateLimitWithRetry(limit uint64, window string, retryAfterSeconds uint64) *Error {
	return New(
		ErrorTypeRateLimit,
		"RATE_LIMIT_EXCEEDED",
		fmt.Sprintf("rate limit exceeded: %d requests per %s", limit, window),
	).WithDetail("limit", fmt.Sprintf("%d", limit)).
		WithDetail("window", window).
		WithDetail("retry_after_seconds", fmt.Sprintf("%d", retryAfterSeconds))
}

// =========================================
// Infrastructure Error Constructors
// =========================================

// DatabaseError creates a database error.
func DatabaseError(operation string, cause error) *Error {
	return New(
		ErrorTypeDatabase,
		"DATABASE_ERROR",
		fmt.Sprintf("database operation failed: %s", operation),
	).WithDetail("operation", operation).
		WithCause(cause)
}

// DatabaseErrorWithTable creates a database error with table context.
func DatabaseErrorWithTable(operation, table string, cause error) *Error {
	return New(
		ErrorTypeDatabase,
		"DATABASE_ERROR",
		fmt.Sprintf("database operation failed: %s on table %s", operation, table),
	).WithDetail("operation", operation).
		WithDetail("table", table).
		WithCause(cause)
}

// CacheError creates a cache error.
func CacheError(operation string, cause error) *Error {
	return New(
		ErrorTypeCache,
		"CACHE_ERROR",
		fmt.Sprintf("cache operation failed: %s", operation),
	).WithDetail("operation", operation).
		WithCause(cause)
}

// CacheErrorWithKey creates a cache error with key context.
func CacheErrorWithKey(operation, key string, cause error) *Error {
	return New(
		ErrorTypeCache,
		"CACHE_ERROR",
		fmt.Sprintf("cache operation failed: %s for key %s", operation, key),
	).WithDetail("operation", operation).
		WithDetail("key", key).
		WithCause(cause)
}

// NetworkError creates a network error.
func NetworkError(operation string, cause error) *Error {
	return New(
		ErrorTypeNetwork,
		"NETWORK_ERROR",
		fmt.Sprintf("network operation failed: %s", operation),
	).WithDetail("operation", operation).
		WithCause(cause)
}

// NetworkErrorWithURL creates a network error with URL context.
func NetworkErrorWithURL(operation, url string, cause error) *Error {
	return New(
		ErrorTypeNetwork,
		"NETWORK_ERROR",
		fmt.Sprintf("network operation failed: %s to %s", operation, url),
	).WithDetail("operation", operation).
		WithDetail("url", url).
		WithCause(cause)
}

// EventError creates an event bus error.
func EventError(operation string, cause error) *Error {
	return New(
		ErrorTypeEvent,
		"EVENT_ERROR",
		fmt.Sprintf("event operation failed: %s", operation),
	).WithDetail("operation", operation).
		WithCause(cause)
}

// EventErrorWithSubject creates an event error with subject context.
func EventErrorWithSubject(operation, subject string, cause error) *Error {
	return New(
		ErrorTypeEvent,
		"EVENT_ERROR",
		fmt.Sprintf("event operation failed: %s for subject %s", operation, subject),
	).WithDetail("operation", operation).
		WithDetail("subject", subject).
		WithCause(cause)
}

// =========================================
// Error Wrapping Functions
// =========================================

// Wrap wraps a standard error with error metadata.
// If err is nil, returns nil.
// If err is already *Error, it is returned as-is with the cause set.
func Wrap(err error, errType ErrorType, code Code, message string) *Error {
	if err == nil {
		return nil
	}

	// If already *Error, preserve it but update fields
	if e, ok := err.(*Error); ok {
		// Create new error with updated fields but preserve existing cause
		wrapped := New(errType, code, message)
		wrapped.Cause = e
		return wrapped
	}

	return New(errType, code, message).WithCause(err)
}

// Wrapf wraps a standard error with a formatted message.
// If err is nil, returns nil.
func Wrapf(err error, errType ErrorType, code Code, format string, args ...interface{}) *Error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, args...)

	// If already *Error, preserve it but update fields
	if e, ok := err.(*Error); ok {
		wrapped := New(errType, code, message)
		wrapped.Cause = e
		return wrapped
	}

	return New(errType, code, message).WithCause(err)
}

// WrapWithDetails wraps an error and adds details in one call.
func WrapWithDetails(err error, errType ErrorType, code Code, message string, details map[string]string) *Error {
	if err == nil {
		return nil
	}

	wrapped := Wrap(err, errType, code, message)
	if wrapped != nil {
		wrapped.WithDetails(details)
	}
	return wrapped
}
