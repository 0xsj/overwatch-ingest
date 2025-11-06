// platform/pkg/errors/stack_test.go
package errors

import (
	"strings"
	"testing"
)

func TestError_WithStack(t *testing.T) {
	e := New(ErrorTypeInternal, "TEST", "test error").WithStack()

	if !e.HasStack() {
		t.Errorf("WithStack() should capture stack trace")
	}

	if len(e.stack) == 0 {
		t.Errorf("WithStack() stack should not be empty")
	}

	// Test nil safety
	var nilErr *Error
	if nilErr.WithStack() != nil {
		t.Errorf("nil Error.WithStack() should return nil")
	}
}

func TestError_WithStack_OnlyOnce(t *testing.T) {
	e := New(ErrorTypeInternal, "TEST", "test error").WithStack()
	originalStack := e.stack

	// Call WithStack again
	e.WithStack()

	// Should not re-capture
	if len(e.stack) != len(originalStack) {
		t.Errorf("WithStack() should not re-capture if already captured")
	}
}

func TestError_HasStack(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected bool
	}{
		{"with_stack", New(ErrorTypeInternal, "TEST", "test").WithStack(), true},
		{"without_stack", New(ErrorTypeInternal, "TEST", "test"), false},
		{"nil_error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.HasStack(); got != tt.expected {
				t.Errorf("HasStack() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_StackTrace(t *testing.T) {
	e := New(ErrorTypeInternal, "TEST", "test error").WithStack()
	trace := e.StackTrace()

	if len(trace) == 0 {
		t.Errorf("StackTrace() should return frames")
	}

	// Check format
	for _, frame := range trace {
		if !strings.Contains(frame, "(") || !strings.Contains(frame, ":") {
			t.Errorf("StackTrace() frame format incorrect: %s", frame)
		}
	}

	// Test without stack
	e2 := New(ErrorTypeInternal, "TEST", "test")
	if e2.StackTrace() != nil {
		t.Errorf("StackTrace() without stack should return nil")
	}
}

func TestError_StackTraceString(t *testing.T) {
	e := New(ErrorTypeInternal, "TEST", "test error").WithStack()
	trace := e.StackTraceString()

	if trace == "" {
		t.Errorf("StackTraceString() should not be empty")
	}

	// Should have "at" prefix for each frame
	if !strings.Contains(trace, "  at ") {
		t.Errorf("StackTraceString() should contain '  at ' prefix")
	}

	// Test without stack
	e2 := New(ErrorTypeInternal, "TEST", "test")
	if e2.StackTraceString() != "" {
		t.Errorf("StackTraceString() without stack should return empty string")
	}
}

func TestError_ErrorWithStack(t *testing.T) {
	e := New(ErrorTypeInternal, "TEST_CODE", "test error").WithStack()
	msg := e.ErrorWithStack()

	// Should contain error message
	if !strings.Contains(msg, "[INTERNAL:TEST_CODE] test error") {
		t.Errorf("ErrorWithStack() should contain error message")
	}

	// Should contain stack trace
	if !strings.Contains(msg, "  at ") {
		t.Errorf("ErrorWithStack() should contain stack trace")
	}

	// Test without stack
	e2 := New(ErrorTypeInternal, "TEST", "test")
	msg2 := e2.ErrorWithStack()
	if strings.Contains(msg2, "  at ") {
		t.Errorf("ErrorWithStack() without stack should not contain trace")
	}
}

func TestError_MarshalStackJSON(t *testing.T) {
	e := New(ErrorTypeInternal, "TEST", "test error").WithStack()
	data, err := e.MarshalStackJSON()
	if err != nil {
		t.Fatalf("MarshalStackJSON() error = %v", err)
	}

	if len(data) == 0 {
		t.Errorf("MarshalStackJSON() should return data")
	}

	// Test without stack
	e2 := New(ErrorTypeInternal, "TEST", "test")
	data2, err := e2.MarshalStackJSON()
	if err != nil {
		t.Fatalf("MarshalStackJSON() without stack error = %v", err)
	}
	if string(data2) != "null" {
		t.Errorf("MarshalStackJSON() without stack = %s, want null", string(data2))
	}
}

func TestWrapWithStack(t *testing.T) {
	cause := New(ErrorTypeInternal, "CAUSE", "cause error")
	e := WrapWithStack(cause, ErrorTypeDatabase, "WRAPPER", "wrapped error")

	if !e.HasStack() {
		t.Errorf("WrapWithStack() should capture stack")
	}

	if e.Cause != cause {
		t.Errorf("WrapWithStack() should preserve cause")
	}
}

func TestWrapfWithStack(t *testing.T) {
	cause := New(ErrorTypeInternal, "CAUSE", "cause error")
	e := WrapfWithStack(cause, ErrorTypeDatabase, "WRAPPER", "wrapped for user %s", "alice")

	if !e.HasStack() {
		t.Errorf("WrapfWithStack() should capture stack")
	}

	expectedMsg := "wrapped for user alice"
	if e.Message != expectedMsg {
		t.Errorf("WrapfWithStack() Message = %v, want %v", e.Message, expectedMsg)
	}
}

// Helper function to test stack depth
func stackDepthHelper1() *Error {
	return stackDepthHelper2()
}

func stackDepthHelper2() *Error {
	return stackDepthHelper3()
}

func stackDepthHelper3() *Error {
	return New(ErrorTypeInternal, "TEST", "test").WithStack()
}

func TestStackDepth(t *testing.T) {
	e := stackDepthHelper1()
	trace := e.StackTrace()

	// Should have captured multiple frames
	if len(trace) < 3 {
		t.Errorf("Stack trace should contain at least 3 frames, got %d", len(trace))
	}

	// Check for our helper functions in the trace
	traceStr := strings.Join(trace, "\n")
	if !strings.Contains(traceStr, "stackDepthHelper") {
		t.Errorf("Stack trace should contain helper functions")
	}
}

// Benchmarks

func BenchmarkError_WithStack(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(ErrorTypeInternal, "TEST", "test").WithStack()
	}
}

func BenchmarkError_StackTrace(b *testing.B) {
	e := New(ErrorTypeInternal, "TEST", "test").WithStack()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.StackTrace()
	}
}

func BenchmarkError_StackTraceString(b *testing.B) {
	e := New(ErrorTypeInternal, "TEST", "test").WithStack()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.StackTraceString()
	}
}

func BenchmarkError_ErrorWithStack(b *testing.B) {
	e := New(ErrorTypeInternal, "TEST", "test").WithStack()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.ErrorWithStack()
	}
}

func BenchmarkWrapWithStack(b *testing.B) {
	cause := New(ErrorTypeInternal, "CAUSE", "cause")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WrapWithStack(cause, ErrorTypeDatabase, "WRAPPER", "wrapped")
	}
}
