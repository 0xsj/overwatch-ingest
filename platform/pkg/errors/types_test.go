// platform/pkg/errors/types_test.go
package errors

import (
	"testing"
)

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected string
	}{
		{"validation", ErrorTypeValidation, "VALIDATION"},
		{"not_found", ErrorTypeNotFound, "NOT_FOUND"},
		{"internal", ErrorTypeInternal, "INTERNAL"},
		{"timeout", ErrorTypeTimeout, "TIMEOUT"},
		{"database", ErrorTypeDatabase, "DATABASE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errType.String(); got != tt.expected {
				t.Errorf("ErrorType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected bool
	}{
		{"validation_valid", ErrorTypeValidation, true},
		{"not_found_valid", ErrorTypeNotFound, true},
		{"internal_valid", ErrorTypeInternal, true},
		{"unknown_invalid", ErrorType("UNKNOWN"), false},
		{"empty_invalid", ErrorType(""), false},
		{"custom_invalid", ErrorType("CUSTOM_TYPE"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errType.IsValid(); got != tt.expected {
				t.Errorf("ErrorType.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorType_IsClientError(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected bool
	}{
		{"validation_is_client", ErrorTypeValidation, true},
		{"not_found_is_client", ErrorTypeNotFound, true},
		{"unauthorized_is_client", ErrorTypeUnauthorized, true},
		{"forbidden_is_client", ErrorTypeForbidden, true},
		{"rate_limit_is_client", ErrorTypeRateLimit, true},
		{"internal_not_client", ErrorTypeInternal, false},
		{"timeout_not_client", ErrorTypeTimeout, false},
		{"database_not_client", ErrorTypeDatabase, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errType.IsClientError(); got != tt.expected {
				t.Errorf("ErrorType.IsClientError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorType_IsServerError(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected bool
	}{
		{"internal_is_server", ErrorTypeInternal, true},
		{"timeout_is_server", ErrorTypeTimeout, true},
		{"unavailable_is_server", ErrorTypeUnavailable, true},
		{"database_is_server", ErrorTypeDatabase, true},
		{"network_is_server", ErrorTypeNetwork, true},
		{"validation_not_server", ErrorTypeValidation, false},
		{"not_found_not_server", ErrorTypeNotFound, false},
		{"forbidden_not_server", ErrorTypeForbidden, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errType.IsServerError(); got != tt.expected {
				t.Errorf("ErrorType.IsServerError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorType_IsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected bool
	}{
		{"timeout_retryable", ErrorTypeTimeout, true},
		{"unavailable_retryable", ErrorTypeUnavailable, true},
		{"rate_limit_retryable", ErrorTypeRateLimit, true},
		{"network_retryable", ErrorTypeNetwork, true},
		{"database_retryable", ErrorTypeDatabase, true},
		{"validation_not_retryable", ErrorTypeValidation, false},
		{"not_found_not_retryable", ErrorTypeNotFound, false},
		{"forbidden_not_retryable", ErrorTypeForbidden, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errType.IsRetryable(); got != tt.expected {
				t.Errorf("ErrorType.IsRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestErrorType_HTTPStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected int
	}{
		{"validation_400", ErrorTypeValidation, 400},
		{"unauthorized_401", ErrorTypeUnauthorized, 401},
		{"forbidden_403", ErrorTypeForbidden, 403},
		{"not_found_404", ErrorTypeNotFound, 404},
		{"conflict_409", ErrorTypeConflict, 409},
		{"already_exists_409", ErrorTypeAlreadyExists, 409},
		{"rate_limit_429", ErrorTypeRateLimit, 429},
		{"internal_500", ErrorTypeInternal, 500},
		{"not_implemented_501", ErrorTypeNotImplemented, 501},
		{"unavailable_503", ErrorTypeUnavailable, 503},
		{"database_503", ErrorTypeDatabase, 503},
		{"timeout_504", ErrorTypeTimeout, 504},
		{"network_504", ErrorTypeNetwork, 504},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errType.HTTPStatusCode(); got != tt.expected {
				t.Errorf("ErrorType.HTTPStatusCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseErrorType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ErrorType
	}{
		{"valid_validation", "VALIDATION", ErrorTypeValidation},
		{"valid_not_found", "NOT_FOUND", ErrorTypeNotFound},
		{"valid_internal", "INTERNAL", ErrorTypeInternal},
		{"invalid_returns_empty", "INVALID", ""},
		{"empty_returns_empty", "", ""},
		{"lowercase_returns_empty", "validation", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseErrorType(tt.input); got != tt.expected {
				t.Errorf("ParseErrorType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMustParseErrorType(t *testing.T) {
	t.Run("valid_type_succeeds", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustParseErrorType() panicked on valid input: %v", r)
			}
		}()

		got := MustParseErrorType("VALIDATION")
		if got != ErrorTypeValidation {
			t.Errorf("MustParseErrorType() = %v, want %v", got, ErrorTypeValidation)
		}
	})

	t.Run("invalid_type_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustParseErrorType() should panic on invalid input")
			}
		}()

		MustParseErrorType("INVALID")
	})
}

func TestAllErrorTypes(t *testing.T) {
	types := AllErrorTypes()

	// Check we have all types
	expectedCount := 15 // Update if you add/remove types
	if len(types) != expectedCount {
		t.Errorf("AllErrorTypes() returned %d types, want %d", len(types), expectedCount)
	}

	// Check all are valid
	for _, et := range types {
		if !et.IsValid() {
			t.Errorf("AllErrorTypes() contains invalid type: %v", et)
		}
	}

	// Check for duplicates
	seen := make(map[ErrorType]bool)
	for _, et := range types {
		if seen[et] {
			t.Errorf("AllErrorTypes() contains duplicate: %v", et)
		}
		seen[et] = true
	}
}

func TestErrorType_ClientServerMutuallyExclusive(t *testing.T) {
	// Ensure no type is both client and server error
	for _, et := range AllErrorTypes() {
		isClient := et.IsClientError()
		isServer := et.IsServerError()

		if isClient && isServer {
			t.Errorf("ErrorType %v is both client and server error", et)
		}
	}
}
