// platform/pkg/errors/constructors_test.go
package errors

import (
	"errors"
	"testing"
	"time"
)

func TestNotFound(t *testing.T) {
	e := NotFound("user", "123")

	if e.Type != ErrorTypeNotFound {
		t.Errorf("NotFound() Type = %v, want %v", e.Type, ErrorTypeNotFound)
	}
	if e.Code != "RESOURCE_NOT_FOUND" {
		t.Errorf("NotFound() Code = %v, want %v", e.Code, "RESOURCE_NOT_FOUND")
	}
	if e.GetDetail("resource_type") != "user" {
		t.Errorf("NotFound() resource_type detail incorrect")
	}
	if e.GetDetail("resource_id") != "123" {
		t.Errorf("NotFound() resource_id detail incorrect")
	}
}

func TestAlreadyExists(t *testing.T) {
	e := AlreadyExists("user", "alice")

	if e.Type != ErrorTypeAlreadyExists {
		t.Errorf("AlreadyExists() Type = %v, want %v", e.Type, ErrorTypeAlreadyExists)
	}
	if e.Code != "RESOURCE_ALREADY_EXISTS" {
		t.Errorf("AlreadyExists() Code = %v, want %v", e.Code, "RESOURCE_ALREADY_EXISTS")
	}
}

func TestValidation(t *testing.T) {
	e := Validation("invalid input")

	if e.Type != ErrorTypeValidation {
		t.Errorf("Validation() Type = %v, want %v", e.Type, ErrorTypeValidation)
	}
	if e.Code != "VALIDATION_FAILED" {
		t.Errorf("Validation() Code = %v, want %v", e.Code, "VALIDATION_FAILED")
	}
}

func TestValidationWithField(t *testing.T) {
	e := ValidationWithField("email", "invalid format")

	if e.Type != ErrorTypeValidation {
		t.Errorf("ValidationWithField() Type incorrect")
	}
	if e.GetDetail("field") != "email" {
		t.Errorf("ValidationWithField() field detail incorrect")
	}
}

func TestRequiredField(t *testing.T) {
	e := RequiredField("username")

	if e.Code != "REQUIRED_FIELD_MISSING" {
		t.Errorf("RequiredField() Code = %v, want %v", e.Code, "REQUIRED_FIELD_MISSING")
	}
	if e.GetDetail("field") != "username" {
		t.Errorf("RequiredField() field detail incorrect")
	}
}

func TestInvalidField(t *testing.T) {
	e := InvalidField("age", "must be positive")

	if e.Code != "INVALID_FIELD_VALUE" {
		t.Errorf("InvalidField() Code incorrect")
	}
	if e.GetDetail("field") != "age" {
		t.Errorf("InvalidField() field detail incorrect")
	}
	if e.GetDetail("reason") != "must be positive" {
		t.Errorf("InvalidField() reason detail incorrect")
	}
}

func TestUnauthorized(t *testing.T) {
	e := Unauthorized("missing token")

	if e.Type != ErrorTypeUnauthorized {
		t.Errorf("Unauthorized() Type incorrect")
	}
	if e.GetDetail("reason") != "missing token" {
		t.Errorf("Unauthorized() reason detail incorrect")
	}
}

func TestForbidden(t *testing.T) {
	e := Forbidden("document", "delete")

	if e.Type != ErrorTypeForbidden {
		t.Errorf("Forbidden() Type incorrect")
	}
	if e.GetDetail("resource") != "document" {
		t.Errorf("Forbidden() resource detail incorrect")
	}
	if e.GetDetail("action") != "delete" {
		t.Errorf("Forbidden() action detail incorrect")
	}
}

func TestInternal(t *testing.T) {
	e := Internal("unexpected error")

	if e.Type != ErrorTypeInternal {
		t.Errorf("Internal() Type incorrect")
	}
	if e.Code != "INTERNAL_ERROR" {
		t.Errorf("Internal() Code incorrect")
	}
}

func TestInternalWithCause(t *testing.T) {
	cause := errors.New("root cause")
	e := InternalWithCause("wrapper", cause)

	if e.Cause != cause {
		t.Errorf("InternalWithCause() should wrap cause")
	}
}

func TestInternalf(t *testing.T) {
	e := Internalf("error code: %d", 500)

	expected := "error code: 500"
	if e.Message != expected {
		t.Errorf("Internalf() Message = %v, want %v", e.Message, expected)
	}
}

func TestNotImplemented(t *testing.T) {
	e := NotImplemented("GraphQL API")

	if e.Type != ErrorTypeNotImplemented {
		t.Errorf("NotImplemented() Type incorrect")
	}
	if e.GetDetail("feature") != "GraphQL API" {
		t.Errorf("NotImplemented() feature detail incorrect")
	}
}

func TestTimeout(t *testing.T) {
	duration := 5 * time.Second
	e := Timeout("database query", duration)

	if e.Type != ErrorTypeTimeout {
		t.Errorf("Timeout() Type incorrect")
	}
	if e.GetDetail("operation") != "database query" {
		t.Errorf("Timeout() operation detail incorrect")
	}
	if e.GetDetail("timeout") != "5s" {
		t.Errorf("Timeout() timeout detail incorrect")
	}
}

func TestTimeoutWithDuration(t *testing.T) {
	e := TimeoutWithDuration("api call", 3000)

	if e.GetDetail("timeout_ms") != "3000" {
		t.Errorf("TimeoutWithDuration() timeout_ms detail incorrect")
	}
}

func TestUnavailable(t *testing.T) {
	e := Unavailable("payment-service")

	if e.Type != ErrorTypeUnavailable {
		t.Errorf("Unavailable() Type incorrect")
	}
	if e.GetDetail("service") != "payment-service" {
		t.Errorf("Unavailable() service detail incorrect")
	}
}

func TestUnavailableWithCause(t *testing.T) {
	cause := errors.New("connection refused")
	e := UnavailableWithCause("auth-service", cause)

	if e.Cause != cause {
		t.Errorf("UnavailableWithCause() should wrap cause")
	}
	if e.GetDetail("service") != "auth-service" {
		t.Errorf("UnavailableWithCause() service detail incorrect")
	}
}

func TestConflict(t *testing.T) {
	e := Conflict("reservation", "time slot already booked")

	if e.Type != ErrorTypeConflict {
		t.Errorf("Conflict() Type incorrect")
	}
	if e.GetDetail("resource") != "reservation" {
		t.Errorf("Conflict() resource detail incorrect")
	}
	if e.GetDetail("reason") != "time slot already booked" {
		t.Errorf("Conflict() reason detail incorrect")
	}
}

func TestRateLimit(t *testing.T) {
	e := RateLimit(100, "minute")

	if e.Type != ErrorTypeRateLimit {
		t.Errorf("RateLimit() Type incorrect")
	}
	if e.GetDetail("limit") != "100" {
		t.Errorf("RateLimit() limit detail incorrect")
	}
	if e.GetDetail("window") != "minute" {
		t.Errorf("RateLimit() window detail incorrect")
	}
}

func TestRateLimitWithRetry(t *testing.T) {
	e := RateLimitWithRetry(100, "minute", 60)

	if e.GetDetail("retry_after_seconds") != "60" {
		t.Errorf("RateLimitWithRetry() retry_after_seconds detail incorrect")
	}
}

func TestDatabaseError(t *testing.T) {
	cause := errors.New("connection timeout")
	e := DatabaseError("SELECT", cause)

	if e.Type != ErrorTypeDatabase {
		t.Errorf("DatabaseError() Type incorrect")
	}
	if e.Cause != cause {
		t.Errorf("DatabaseError() should wrap cause")
	}
	if e.GetDetail("operation") != "SELECT" {
		t.Errorf("DatabaseError() operation detail incorrect")
	}
}

func TestDatabaseErrorWithTable(t *testing.T) {
	cause := errors.New("deadlock")
	e := DatabaseErrorWithTable("UPDATE", "users", cause)

	if e.GetDetail("table") != "users" {
		t.Errorf("DatabaseErrorWithTable() table detail incorrect")
	}
}

func TestCacheError(t *testing.T) {
	cause := errors.New("redis unavailable")
	e := CacheError("GET", cause)

	if e.Type != ErrorTypeCache {
		t.Errorf("CacheError() Type incorrect")
	}
	if e.Cause != cause {
		t.Errorf("CacheError() should wrap cause")
	}
}

func TestCacheErrorWithKey(t *testing.T) {
	cause := errors.New("key not found")
	e := CacheErrorWithKey("GET", "user:123", cause)

	if e.GetDetail("key") != "user:123" {
		t.Errorf("CacheErrorWithKey() key detail incorrect")
	}
}

func TestNetworkError(t *testing.T) {
	cause := errors.New("connection refused")
	e := NetworkError("HTTP GET", cause)

	if e.Type != ErrorTypeNetwork {
		t.Errorf("NetworkError() Type incorrect")
	}
	if e.Cause != cause {
		t.Errorf("NetworkError() should wrap cause")
	}
}

func TestNetworkErrorWithURL(t *testing.T) {
	cause := errors.New("timeout")
	e := NetworkErrorWithURL("GET", "https://api.example.com", cause)

	if e.GetDetail("url") != "https://api.example.com" {
		t.Errorf("NetworkErrorWithURL() url detail incorrect")
	}
}

func TestEventError(t *testing.T) {
	cause := errors.New("nats unavailable")
	e := EventError("publish", cause)

	if e.Type != ErrorTypeEvent {
		t.Errorf("EventError() Type incorrect")
	}
	if e.Cause != cause {
		t.Errorf("EventError() should wrap cause")
	}
}

func TestEventErrorWithSubject(t *testing.T) {
	cause := errors.New("no subscribers")
	e := EventErrorWithSubject("publish", "incidents.created", cause)

	if e.GetDetail("subject") != "incidents.created" {
		t.Errorf("EventErrorWithSubject() subject detail incorrect")
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("original error")
	e := Wrap(cause, ErrorTypeInternal, "WRAP_TEST", "wrapped message")

	if e.Type != ErrorTypeInternal {
		t.Errorf("Wrap() Type incorrect")
	}
	if e.Code != "WRAP_TEST" {
		t.Errorf("Wrap() Code incorrect")
	}
	if e.Cause != cause {
		t.Errorf("Wrap() should preserve cause")
	}

	// Test nil case
	if Wrap(nil, ErrorTypeInternal, "CODE", "msg") != nil {
		t.Errorf("Wrap(nil) should return nil")
	}

	// Test wrapping *Error
	customErr := New(ErrorTypeValidation, "ORIGINAL", "original")
	wrapped := Wrap(customErr, ErrorTypeInternal, "WRAPPED", "wrapped")
	if wrapped.Cause != customErr {
		t.Errorf("Wrap() should preserve *Error as cause")
	}
}

func TestWrapf(t *testing.T) {
	cause := errors.New("original")
	e := Wrapf(cause, ErrorTypeInternal, "TEST", "error for user %s", "alice")

	expected := "error for user alice"
	if e.Message != expected {
		t.Errorf("Wrapf() Message = %v, want %v", e.Message, expected)
	}

	// Test nil case
	if Wrapf(nil, ErrorTypeInternal, "CODE", "msg") != nil {
		t.Errorf("Wrapf(nil) should return nil")
	}
}

func TestWrapWithDetails(t *testing.T) {
	cause := errors.New("original")
	details := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	e := WrapWithDetails(cause, ErrorTypeInternal, "TEST", "wrapped", details)

	if e.GetDetail("key1") != "value1" {
		t.Errorf("WrapWithDetails() should preserve details")
	}
	if e.GetDetail("key2") != "value2" {
		t.Errorf("WrapWithDetails() should preserve details")
	}

	// Test nil case
	if WrapWithDetails(nil, ErrorTypeInternal, "CODE", "msg", details) != nil {
		t.Errorf("WrapWithDetails(nil) should return nil")
	}
}

// Benchmarks

func BenchmarkNotFound(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NotFound("user", "123")
	}
}

func BenchmarkValidationWithField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidationWithField("email", "invalid format")
	}
}

func BenchmarkDatabaseErrorWithTable(b *testing.B) {
	cause := errors.New("deadlock")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DatabaseErrorWithTable("UPDATE", "users", cause)
	}
}

func BenchmarkWrap(b *testing.B) {
	cause := errors.New("original")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Wrap(cause, ErrorTypeInternal, "TEST", "wrapped")
	}
}

func BenchmarkWrapf(b *testing.B) {
	cause := errors.New("original")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Wrapf(cause, ErrorTypeInternal, "TEST", "error for user %s", "alice")
	}
}
