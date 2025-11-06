// platform/pkg/errors/doc.go

/*
Package errors provides a rich, flexible error handling system for the Scout platform.

# Overview

This package offers structured error types with metadata, cause chaining, stack traces,
and utilities for error comparison and serialization. It is designed to be generic and
unopinionated, allowing services to build domain-specific errors on top of the primitives.

# Core Concepts

The error system is built on three primitives:

1. ErrorType - Categorizes errors (validation, not_found, internal, etc.)
2. Code - Specific error identifier (service-defined)
3. Error - Rich error struct with metadata, cause chain, and optional stack trace

# Basic Usage

Creating errors:

	// Simple error
	err := errors.New(errors.ErrorTypeValidation, "INVALID_EMAIL", "email format is invalid")

	// With formatted message
	err := errors.Newf(errors.ErrorTypeNotFound, "USER_NOT_FOUND", "user %s not found", userID)

	// Using constructors
	err := errors.NotFound("user", userID)
	err := errors.Validation("invalid input")
	err := errors.Internal("unexpected error")

# Builder Pattern

Errors support method chaining for adding context:

	err := errors.New(errors.ErrorTypeDatabase, "QUERY_FAILED", "failed to query users").
		WithDetail("table", "users").
		WithDetail("operation", "SELECT").
		WithCause(sqlErr).
		WithStack()

# Error Wrapping

Wrap standard library errors or other errors:

	if err := db.Query(...); err != nil {
		return errors.Wrap(err, errors.ErrorTypeDatabase, "DB_ERROR", "query failed")
	}

	// With formatted message
	return errors.Wrapf(err, errors.ErrorTypeInternal, "PROCESS_FAILED",
		"failed to process user %s", userID)

	// With stack trace
	return errors.WrapWithStack(err, errors.ErrorTypeInternal, "ERROR", "wrapped")

# Error Types

Client Errors (4xx - caused by invalid input):
  - ErrorTypeValidation    - Invalid input data
  - ErrorTypeNotFound      - Resource doesn't exist
  - ErrorTypeAlreadyExists - Resource already exists
  - ErrorTypeUnauthorized  - Missing/invalid authentication
  - ErrorTypeForbidden     - Insufficient permissions
  - ErrorTypeConflict      - State conflict
  - ErrorTypeRateLimit     - Rate limit exceeded

Server Errors (5xx - caused by server issues):
  - ErrorTypeInternal       - Internal server error
  - ErrorTypeUnavailable    - Service unavailable
  - ErrorTypeTimeout        - Operation timeout
  - ErrorTypeNotImplemented - Feature not implemented

Infrastructure Errors:
  - ErrorTypeDatabase - Database operation failed
  - ErrorTypeCache    - Cache operation failed
  - ErrorTypeNetwork  - Network call failed
  - ErrorTypeEvent    - Event bus operation failed

# Common Constructors

The package provides convenient constructors for common error scenarios:

	// Resource errors
	errors.NotFound(resourceType, resourceID)
	errors.AlreadyExists(resourceType, resourceID)

	// Validation errors
	errors.Validation(message)
	errors.ValidationWithField(field, message)
	errors.RequiredField(field)
	errors.InvalidField(field, reason)

	// Authorization errors
	errors.Unauthorized(reason)
	errors.Forbidden(resource, action)

	// Internal errors
	errors.Internal(message)
	errors.InternalWithCause(message, cause)
	errors.Internalf(format, args...)

	// Infrastructure errors
	errors.DatabaseError(operation, cause)
	errors.DatabaseErrorWithTable(operation, table, cause)
	errors.CacheError(operation, cause)
	errors.NetworkError(operation, cause)
	errors.EventError(operation, cause)

	// Timeout/availability
	errors.Timeout(operation, duration)
	errors.Unavailable(service)

# Error Comparison

Check error properties:

	if errors.IsType(err, errors.ErrorTypeValidation) {
		// Handle validation error
	}

	if errors.IsRetryable(err) {
		// Retry logic
	}

	if errors.HasCode(err, "USER_NOT_FOUND") {
		// Specific error handling
	}

Extract error information:

	errType := errors.GetType(err)
	code := errors.GetCode(err)
	details := errors.GetDetails(err)
	cause := errors.GetCause(err)
	root := errors.GetRootCause(err)

# Stack Traces

Opt-in stack trace capture:

	err := errors.Internal("unexpected error").WithStack()

	// Get stack frames
	frames := err.StackTrace()  // []string

	// Format as string
	trace := err.StackTraceString()

	// Error message with stack
	log.Error(err.ErrorWithStack())

# JSON Serialization

Standard serialization (excludes cause for security):

	data, err := error.MarshalJSON()
	// or
	data, err := errors.ToJSON(err)

Verbose serialization (includes full cause chain):

	data, err := error.MarshalJSONVerbose()
	// or
	data, err := errors.ToJSONVerbose(err)

Deserialization:

	var e errors.Error
	json.Unmarshal(data, &e)
	// or
	e, err := errors.FromJSON(data)

JSON format:

	{
		"type": "VALIDATION",
		"code": "REQUIRED_FIELD_MISSING",
		"message": "required field 'email' is missing",
		"details": {
			"field": "email"
		}
	}

# Service-Specific Errors

Services should define their own error codes and optionally extend the error type:

	// services/incidents/internal/errors/codes.go
	package errors

	import "github.com/0xsj/scout/platform/pkg/errors"

	const (
		CodeIncidentNotFound     errors.Code = "INCIDENT_NOT_FOUND"
		CodeInvalidSeverity      errors.Code = "INVALID_SEVERITY"
		CodeGeocodingFailed      errors.Code = "GEOCODING_FAILED"
	)

	// Constructor
	func IncidentNotFound(id string) *errors.Error {
		return errors.NotFound("incident", id).
			WithDetail("incident_id", id)
	}

# HTTP Status Codes

Get recommended HTTP status code:

	statusCode := err.HTTPStatusCode()  // Returns int (400, 404, 500, etc.)

	// Or check categories
	if err.IsClientError() {
		// 4xx error
	}
	if err.IsServerError() {
		// 5xx error
	}

# Best Practices

1. Use specific error types and codes for better observability
2. Add contextual details for debugging (table names, IDs, etc.)
3. Wrap errors at boundaries to add context
4. Capture stack traces for unexpected errors
5. Use constructors for common patterns
6. Define service-specific codes as constants
7. Never expose internal errors to external APIs (use standard JSON marshaling)

# Examples

Complete error handling example:

	func (s *Service) GetUser(ctx context.Context, userID string) (*User, error) {
		// Validate input
		if userID == "" {
			return nil, errors.RequiredField("user_id")
		}

		// Database query with error wrapping
		user, err := s.db.Query(ctx, "SELECT * FROM users WHERE id = $1", userID)
		if err != nil {
			if isNotFound(err) {
				return nil, errors.NotFound("user", userID)
			}
			return nil, errors.DatabaseErrorWithTable("SELECT", "users", err).
				WithDetail("user_id", userID).
				WithStack()
		}

		return user, nil
	}

Error handling in HTTP handler:

	func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("id")

		user, err := h.service.GetUser(r.Context(), userID)
		if err != nil {
			statusCode := errors.HTTPStatusCode(err)

			response := ErrorResponse{
				Error: ErrorDetail{
					Code:    string(errors.GetCode(err)),
					Message: errors.GetMessage(err),
					Details: errors.GetDetails(err),
				},
			}

			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(response)
			return
		}

		json.NewEncoder(w).Encode(user)
	}

Retry logic based on error type:

	func (c *Client) CallWithRetry(ctx context.Context) error {
		for attempt := 0; attempt < maxRetries; attempt++ {
			err := c.call(ctx)
			if err == nil {
				return nil
			}

			if !errors.IsRetryable(err) {
				return err
			}

			// Exponential backoff
			time.Sleep(backoff(attempt))
		}
		return errors.Timeout("call", timeout)
	}
*/
package errors
