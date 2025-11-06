// platform/pkg/errors/comparison_test.go
package errors

import (
	"errors"
	"reflect"
	"testing"
)

func TestIsType(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST", "test")

	tests := []struct {
		name     string
		err      error
		errType  ErrorType
		expected bool
	}{
		{"matching_type", e, ErrorTypeValidation, true},
		{"different_type", e, ErrorTypeInternal, false},
		{"nil_error", nil, ErrorTypeValidation, false},
		{"standard_error", errors.New("std"), ErrorTypeValidation, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsType(tt.err, tt.errType); got != tt.expected {
				t.Errorf("IsType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsOneOfTypes(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST", "test")

	tests := []struct {
		name     string
		err      error
		types    []ErrorType
		expected bool
	}{
		{"matches_first", e, []ErrorType{ErrorTypeValidation, ErrorTypeInternal}, true},
		{"matches_second", e, []ErrorType{ErrorTypeInternal, ErrorTypeValidation}, true},
		{"no_match", e, []ErrorType{ErrorTypeInternal, ErrorTypeTimeout}, false},
		{"empty_types", e, []ErrorType{}, false},
		{"nil_error", nil, []ErrorType{ErrorTypeValidation}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOneOfTypes(tt.err, tt.types...); got != tt.expected {
				t.Errorf("IsOneOfTypes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasCode(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test")

	tests := []struct {
		name     string
		err      error
		code     Code
		expected bool
	}{
		{"matching_code", e, "TEST_CODE", true},
		{"different_code", e, "OTHER_CODE", false},
		{"nil_error", nil, "TEST_CODE", false},
		{"standard_error", errors.New("std"), "TEST_CODE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasCode(tt.err, tt.code); got != tt.expected {
				t.Errorf("HasCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasOneOfCodes(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test")

	tests := []struct {
		name     string
		err      error
		codes    []Code
		expected bool
	}{
		{"matches_first", e, []Code{"TEST_CODE", "OTHER_CODE"}, true},
		{"matches_second", e, []Code{"OTHER_CODE", "TEST_CODE"}, true},
		{"no_match", e, []Code{"CODE1", "CODE2"}, false},
		{"empty_codes", e, []Code{}, false},
		{"nil_error", nil, []Code{"TEST_CODE"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasOneOfCodes(tt.err, tt.codes...); got != tt.expected {
				t.Errorf("HasOneOfCodes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsClientError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"validation_is_client", New(ErrorTypeValidation, "CODE", "msg"), true},
		{"not_found_is_client", New(ErrorTypeNotFound, "CODE", "msg"), true},
		{"internal_not_client", New(ErrorTypeInternal, "CODE", "msg"), false},
		{"nil_error", nil, false},
		{"standard_error", errors.New("std"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsClientError(tt.err); got != tt.expected {
				t.Errorf("IsClientError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsServerError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"internal_is_server", New(ErrorTypeInternal, "CODE", "msg"), true},
		{"database_is_server", New(ErrorTypeDatabase, "CODE", "msg"), true},
		{"validation_not_server", New(ErrorTypeValidation, "CODE", "msg"), false},
		{"nil_error", nil, false},
		{"standard_error", errors.New("std"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsServerError(tt.err); got != tt.expected {
				t.Errorf("IsServerError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"timeout_retryable", New(ErrorTypeTimeout, "CODE", "msg"), true},
		{"unavailable_retryable", New(ErrorTypeUnavailable, "CODE", "msg"), true},
		{"validation_not_retryable", New(ErrorTypeValidation, "CODE", "msg"), false},
		{"nil_error", nil, false},
		{"standard_error", errors.New("std"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.expected {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorType
	}{
		{"custom_error", New(ErrorTypeValidation, "CODE", "msg"), ErrorTypeValidation},
		{"nil_error", nil, ""},
		{"standard_error", errors.New("std"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetType(tt.err); got != tt.expected {
				t.Errorf("GetType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"custom_error", New(ErrorTypeValidation, "CODE", "test message"), "test message"},
		{"standard_error", errors.New("std error"), "std error"},
		{"nil_error", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMessage(tt.err); got != tt.expected {
				t.Errorf("GetMessage() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetDetails(t *testing.T) {
	e := New(ErrorTypeValidation, "CODE", "msg").
		WithDetail("key1", "value1").
		WithDetail("key2", "value2")

	tests := []struct {
		name     string
		err      error
		expected map[string]string
	}{
		{"with_details", e, map[string]string{"key1": "value1", "key2": "value2"}},
		{"no_details", New(ErrorTypeValidation, "CODE", "msg"), nil},
		{"nil_error", nil, nil},
		{"standard_error", errors.New("std"), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDetails(tt.err)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("GetDetails() = %v, want %v", got, tt.expected)
			}

			// Verify it returns a copy (mutation doesn't affect original)
			if got != nil {
				got["new_key"] = "new_value"
				if tt.err != nil {
					if customErr, ok := tt.err.(*Error); ok {
						if customErr.HasDetail("new_key") {
							t.Errorf("GetDetails() should return a copy")
						}
					}
				}
			}
		})
	}
}

func TestGetDetail(t *testing.T) {
	e := New(ErrorTypeValidation, "CODE", "msg").
		WithDetail("existing", "value")

	tests := []struct {
		name     string
		err      error
		key      string
		expected string
	}{
		{"existing_key", e, "existing", "value"},
		{"missing_key", e, "missing", ""},
		{"nil_error", nil, "key", ""},
		{"standard_error", errors.New("std"), "key", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDetail(tt.err, tt.key); got != tt.expected {
				t.Errorf("GetDetail() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasDetail(t *testing.T) {
	e := New(ErrorTypeValidation, "CODE", "msg").
		WithDetail("existing", "value")

	tests := []struct {
		name     string
		err      error
		key      string
		expected bool
	}{
		{"existing_key", e, "existing", true},
		{"missing_key", e, "missing", false},
		{"nil_error", nil, "key", false},
		{"standard_error", errors.New("std"), "key", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasDetail(tt.err, tt.key); got != tt.expected {
				t.Errorf("HasDetail() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetCause(t *testing.T) {
	stdCause := errors.New("standard cause")
	customCause := New(ErrorTypeInternal, "CAUSE", "cause")

	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{"with_custom_cause", New(ErrorTypeDatabase, "CODE", "msg").WithCause(customCause), customCause},
		{"with_std_cause", New(ErrorTypeDatabase, "CODE", "msg").WithCause(stdCause), stdCause},
		{"no_cause", New(ErrorTypeValidation, "CODE", "msg"), nil},
		{"nil_error", nil, nil},
		{"standard_error", stdCause, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCause(tt.err)
			if got != tt.expected {
				t.Errorf("GetCause() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetRootCause(t *testing.T) {
	root := errors.New("root cause")
	middle := New(ErrorTypeInternal, "MIDDLE", "middle").WithCause(root)
	top := New(ErrorTypeDatabase, "TOP", "top").WithCause(middle)

	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{"deep_chain", top, root},
		{"single_layer", middle, root},
		{"no_cause", New(ErrorTypeValidation, "CODE", "msg"), New(ErrorTypeValidation, "CODE", "msg")},
		{"nil_error", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRootCause(tt.err)

			// For no_cause case, compare by value not reference
			if tt.name == "no_cause" {
				if got == nil || got.Error() != tt.expected.Error() {
					t.Errorf("GetRootCause() = %v, want %v", got, tt.expected)
				}
			} else {
				if got != tt.expected {
					t.Errorf("GetRootCause() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestMatch(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test message contains keyword")

	tests := []struct {
		name            string
		err             error
		errType         ErrorType
		code            Code
		messageContains string
		expected        bool
	}{
		{"all_match", e, ErrorTypeValidation, "TEST_CODE", "keyword", true},
		{"type_only", e, ErrorTypeValidation, "", "", true},
		{"code_only", e, "", "TEST_CODE", "", true},
		{"message_only", e, "", "", "keyword", true},
		{"type_mismatch", e, ErrorTypeInternal, "", "", false},
		{"code_mismatch", e, "", "OTHER_CODE", "", false},
		{"message_mismatch", e, "", "", "missing", false},
		{"nil_error", nil, ErrorTypeValidation, "", "", false},
		{"standard_error", errors.New("std"), ErrorTypeValidation, "", "", false},
		{"all_empty", e, "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Match(tt.err, tt.errType, tt.code, tt.messageContains); got != tt.expected {
				t.Errorf("Match() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMatch_EdgeCases(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST", "exact")

	// Test exact message match
	if !Match(e, "", "", "exact") {
		t.Errorf("Match() should match exact message")
	}

	// Test substring in message
	e2 := New(ErrorTypeValidation, "TEST", "this is a test message")
	if !Match(e2, "", "", "test") {
		t.Errorf("Match() should match message substring")
	}

	// Test empty message search (should match any error)
	if !Match(e, "", "", "") {
		t.Errorf("Match() with empty message should match")
	}
}

// Benchmarks

func BenchmarkIsType(b *testing.B) {
	e := New(ErrorTypeValidation, "CODE", "msg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsType(e, ErrorTypeValidation)
	}
}

func BenchmarkIsOneOfTypes(b *testing.B) {
	e := New(ErrorTypeValidation, "CODE", "msg")
	types := []ErrorType{ErrorTypeInternal, ErrorTypeValidation, ErrorTypeTimeout}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsOneOfTypes(e, types...)
	}
}

func BenchmarkHasCode(b *testing.B) {
	e := New(ErrorTypeValidation, "TEST_CODE", "msg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HasCode(e, "TEST_CODE")
	}
}

func BenchmarkIsClientError(b *testing.B) {
	e := New(ErrorTypeValidation, "CODE", "msg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsClientError(e)
	}
}

func BenchmarkIsRetryable(b *testing.B) {
	e := New(ErrorTypeTimeout, "CODE", "msg")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsRetryable(e)
	}
}

func BenchmarkGetDetails(b *testing.B) {
	e := New(ErrorTypeValidation, "CODE", "msg").
		WithDetail("key1", "value1").
		WithDetail("key2", "value2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetDetails(e)
	}
}

func BenchmarkGetRootCause(b *testing.B) {
	root := errors.New("root")
	middle := New(ErrorTypeInternal, "MIDDLE", "middle").WithCause(root)
	top := New(ErrorTypeDatabase, "TOP", "top").WithCause(middle)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetRootCause(top)
	}
}

func BenchmarkMatch(b *testing.B) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test message with keyword")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Match(e, ErrorTypeValidation, "TEST_CODE", "keyword")
	}
}
