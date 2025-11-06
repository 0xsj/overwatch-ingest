// platform/pkg/errors/error_test.go
package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test message")

	if e.Type != ErrorTypeValidation {
		t.Errorf("New() Type = %v, want %v", e.Type, ErrorTypeValidation)
	}
	if e.Code != "TEST_CODE" {
		t.Errorf("New() Code = %v, want %v", e.Code, "TEST_CODE")
	}
	if e.Message != "test message" {
		t.Errorf("New() Message = %v, want %v", e.Message, "test message")
	}
	if e.Details != nil {
		t.Errorf("New() Details should be nil, got %v", e.Details)
	}
	if e.Cause != nil {
		t.Errorf("New() Cause should be nil, got %v", e.Cause)
	}
}

func TestNewf(t *testing.T) {
	e := Newf(ErrorTypeInternal, "TEST_CODE", "error for user %s with id %d", "alice", 123)

	expected := "error for user alice with id 123"
	if e.Message != expected {
		t.Errorf("Newf() Message = %v, want %v", e.Message, expected)
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name:     "simple_error",
			err:      New(ErrorTypeValidation, "TEST_CODE", "test message"),
			expected: "[VALIDATION:TEST_CODE] test message",
		},
		{
			name:     "with_cause",
			err:      New(ErrorTypeInternal, "TEST_CODE", "outer error").WithCause(errors.New("inner error")),
			expected: "[INTERNAL:TEST_CODE] outer error | caused by: inner error",
		},
		{
			name:     "nil_error",
			err:      nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	e := New(ErrorTypeInternal, "TEST_CODE", "wrapped error").WithCause(cause)

	unwrapped := e.Unwrap()
	if unwrapped != cause {
		t.Errorf("Error.Unwrap() = %v, want %v", unwrapped, cause)
	}

	// Test nil case
	var nilErr *Error
	if nilErr.Unwrap() != nil {
		t.Errorf("nil Error.Unwrap() should return nil")
	}

	// Test no cause
	noCause := New(ErrorTypeInternal, "TEST_CODE", "no cause")
	if noCause.Unwrap() != nil {
		t.Errorf("Error with no cause should Unwrap() to nil")
	}
}

func TestError_Is(t *testing.T) {
	e1 := New(ErrorTypeValidation, "TEST_CODE", "message 1")
	e2 := New(ErrorTypeValidation, "TEST_CODE", "message 2")
	e3 := New(ErrorTypeInternal, "TEST_CODE", "message 3")
	e4 := New(ErrorTypeValidation, "OTHER_CODE", "message 4")

	tests := []struct {
		name     string
		err      *Error
		target   error
		expected bool
	}{
		{"same_type_and_code", e1, e2, true},
		{"different_type", e1, e3, false},
		{"different_code", e1, e4, false},
		{"nil_error", nil, e1, false},
		{"nil_target", e1, nil, false},
		{"both_nil", nil, nil, true},
		{"non_error_type", e1, errors.New("std error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Is(tt.target); got != tt.expected {
				t.Errorf("Error.Is() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_WithCause(t *testing.T) {
	cause := errors.New("root cause")
	e := New(ErrorTypeInternal, "TEST_CODE", "wrapper").WithCause(cause)

	if e.Cause != cause {
		t.Errorf("WithCause() Cause = %v, want %v", e.Cause, cause)
	}

	// Test chaining
	e2 := e.WithCause(errors.New("new cause"))
	if e2.Cause.Error() != "new cause" {
		t.Errorf("WithCause() chaining failed")
	}

	// Test nil safety
	var nilErr *Error
	if nilErr.WithCause(cause) != nil {
		t.Errorf("nil Error.WithCause() should return nil")
	}
}

func TestError_WithDetail(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test").WithDetail("key1", "value1")

	if e.Details["key1"] != "value1" {
		t.Errorf("WithDetail() did not set key1")
	}

	// Test chaining
	e.WithDetail("key2", "value2")
	if e.Details["key2"] != "value2" {
		t.Errorf("WithDetail() chaining failed")
	}

	// Test nil safety
	var nilErr *Error
	if nilErr.WithDetail("key", "value") != nil {
		t.Errorf("nil Error.WithDetail() should return nil")
	}
}

func TestError_WithDetails(t *testing.T) {
	details := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	e := New(ErrorTypeValidation, "TEST_CODE", "test").WithDetails(details)

	if len(e.Details) != 2 {
		t.Errorf("WithDetails() should set 2 details, got %d", len(e.Details))
	}

	for k, v := range details {
		if e.Details[k] != v {
			t.Errorf("WithDetails() did not set %s correctly", k)
		}
	}

	// Test nil safety
	var nilErr *Error
	if nilErr.WithDetails(details) != nil {
		t.Errorf("nil Error.WithDetails() should return nil")
	}
}

func TestError_GetDetail(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test").
		WithDetail("existing", "value")

	tests := []struct {
		name     string
		err      *Error
		key      string
		expected string
	}{
		{"existing_key", e, "existing", "value"},
		{"missing_key", e, "missing", ""},
		{"nil_error", nil, "key", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.GetDetail(tt.key); got != tt.expected {
				t.Errorf("GetDetail(%s) = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func TestError_HasDetail(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test").
		WithDetail("existing", "value")

	tests := []struct {
		name     string
		err      *Error
		key      string
		expected bool
	}{
		{"existing_key", e, "existing", true},
		{"missing_key", e, "missing", false},
		{"nil_error", nil, "key", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.HasDetail(tt.key); got != tt.expected {
				t.Errorf("HasDetail(%s) = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func TestError_HasDetails(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected bool
	}{
		{"with_details", New(ErrorTypeValidation, "TEST_CODE", "test").WithDetail("key", "value"), true},
		{"no_details", New(ErrorTypeValidation, "TEST_CODE", "test"), false},
		{"nil_error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.HasDetails(); got != tt.expected {
				t.Errorf("HasDetails() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_IsClientError(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected bool
	}{
		{"validation_is_client", New(ErrorTypeValidation, "CODE", "msg"), true},
		{"not_found_is_client", New(ErrorTypeNotFound, "CODE", "msg"), true},
		{"internal_not_client", New(ErrorTypeInternal, "CODE", "msg"), false},
		{"nil_error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsClientError(); got != tt.expected {
				t.Errorf("IsClientError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_IsServerError(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected bool
	}{
		{"internal_is_server", New(ErrorTypeInternal, "CODE", "msg"), true},
		{"database_is_server", New(ErrorTypeDatabase, "CODE", "msg"), true},
		{"validation_not_server", New(ErrorTypeValidation, "CODE", "msg"), false},
		{"nil_error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsServerError(); got != tt.expected {
				t.Errorf("IsServerError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_IsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected bool
	}{
		{"timeout_retryable", New(ErrorTypeTimeout, "CODE", "msg"), true},
		{"unavailable_retryable", New(ErrorTypeUnavailable, "CODE", "msg"), true},
		{"validation_not_retryable", New(ErrorTypeValidation, "CODE", "msg"), false},
		{"nil_error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsRetryable(); got != tt.expected {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_HTTPStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected int
	}{
		{"validation_400", New(ErrorTypeValidation, "CODE", "msg"), 400},
		{"not_found_404", New(ErrorTypeNotFound, "CODE", "msg"), 404},
		{"internal_500", New(ErrorTypeInternal, "CODE", "msg"), 500},
		{"nil_error_500", nil, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.HTTPStatusCode(); got != tt.expected {
				t.Errorf("HTTPStatusCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAs(t *testing.T) {
	customErr := New(ErrorTypeValidation, "TEST_CODE", "test")
	stdErr := errors.New("standard error")

	tests := []struct {
		name      string
		err       error
		expectOk  bool
		expectNil bool
	}{
		{"custom_error", customErr, true, false},
		{"wrapped_custom", errors.Join(stdErr, customErr), true, false},
		{"standard_error", stdErr, false, true},
		{"nil_error", nil, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, ok := As(tt.err)
			if ok != tt.expectOk {
				t.Errorf("As() ok = %v, want %v", ok, tt.expectOk)
			}
			if (e == nil) != tt.expectNil {
				t.Errorf("As() nil = %v, want %v", e == nil, tt.expectNil)
			}
		})
	}
}

func TestIs(t *testing.T) {
	e1 := New(ErrorTypeValidation, "TEST_CODE", "msg1")
	e2 := New(ErrorTypeValidation, "TEST_CODE", "msg2")
	stdErr := errors.New("standard error")

	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{"matching_errors", e1, e2, true},
		{"wrapped_matching", errors.Join(stdErr, e1), e2, true},
		{"non_matching", e1, stdErr, false},
		{"nil_err", nil, e1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err, tt.target); got != tt.expected {
				t.Errorf("Is() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Benchmarks

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(ErrorTypeValidation, "TEST_CODE", "test message")
	}
}

func BenchmarkNewf(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Newf(ErrorTypeValidation, "TEST_CODE", "error for user %s", "alice")
	}
}

func BenchmarkError_Error(b *testing.B) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test message")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkError_WithCause(b *testing.B) {
	cause := errors.New("cause")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(ErrorTypeInternal, "TEST_CODE", "test").WithCause(cause)
	}
}

func BenchmarkError_WithDetail(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(ErrorTypeValidation, "TEST_CODE", "test").WithDetail("key", "value")
	}
}

func BenchmarkError_WithDetails(b *testing.B) {
	details := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(ErrorTypeValidation, "TEST_CODE", "test").WithDetails(details)
	}
}

func BenchmarkError_ChainedBuilding(b *testing.B) {
	cause := errors.New("cause")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(ErrorTypeDatabase, "DB_ERROR", "database failed").
			WithCause(cause).
			WithDetail("table", "users").
			WithDetail("operation", "SELECT")
	}
}
