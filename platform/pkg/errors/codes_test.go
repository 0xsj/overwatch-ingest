// platform/pkg/errors/codes_test.go
package errors

import (
	"testing"
)

func TestCode_String(t *testing.T) {
	tests := []struct {
		name     string
		code     Code
		expected string
	}{
		{"simple_code", Code("USER_NOT_FOUND"), "USER_NOT_FOUND"},
		{"empty_code", Code(""), ""},
		{"complex_code", Code("DATABASE_CONNECTION_FAILED"), "DATABASE_CONNECTION_FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.String(); got != tt.expected {
				t.Errorf("Code.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCode_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		code     Code
		expected bool
	}{
		{"empty_code", Code(""), true},
		{"non_empty_code", Code("USER_NOT_FOUND"), false},
		{"whitespace_not_empty", Code(" "), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.IsEmpty(); got != tt.expected {
				t.Errorf("Code.IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Benchmarks

func BenchmarkCode_String(b *testing.B) {
	code := Code("USER_NOT_FOUND")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = code.String()
	}
}

func BenchmarkCode_IsEmpty(b *testing.B) {
	code := Code("USER_NOT_FOUND")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = code.IsEmpty()
	}
}
