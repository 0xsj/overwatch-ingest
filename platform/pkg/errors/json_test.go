// platform/pkg/errors/json_test.go
package errors

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestError_MarshalJSON(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test message").
		WithDetail("field", "email")

	data, err := e.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["type"] != "VALIDATION" {
		t.Errorf("MarshalJSON() type = %v, want VALIDATION", result["type"])
	}
	if result["code"] != "TEST_CODE" {
		t.Errorf("MarshalJSON() code = %v, want TEST_CODE", result["code"])
	}
	if result["message"] != "test message" {
		t.Errorf("MarshalJSON() message = %v, want 'test message'", result["message"])
	}

	// Should NOT include cause by default
	if _, hasCause := result["cause"]; hasCause {
		t.Errorf("MarshalJSON() should not include cause by default")
	}

	// Check details
	details, ok := result["details"].(map[string]interface{})
	if !ok {
		t.Fatalf("MarshalJSON() details not a map")
	}
	if details["field"] != "email" {
		t.Errorf("MarshalJSON() details field incorrect")
	}
}

func TestError_MarshalJSON_Nil(t *testing.T) {
	var e *Error
	data, err := e.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() on nil should not error: %v", err)
	}

	if string(data) != "null" {
		t.Errorf("MarshalJSON() on nil = %s, want null", string(data))
	}
}

func TestError_MarshalJSONVerbose(t *testing.T) {
	cause := New(ErrorTypeInternal, "INNER", "inner error")
	e := New(ErrorTypeDatabase, "OUTER", "outer error").WithCause(cause)

	data, err := e.MarshalJSONVerbose()
	if err != nil {
		t.Fatalf("MarshalJSONVerbose() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Should include cause in verbose mode
	causeData, hasCause := result["cause"].(map[string]interface{})
	if !hasCause {
		t.Fatalf("MarshalJSONVerbose() should include cause")
	}

	if causeData["type"] != "INTERNAL" {
		t.Errorf("MarshalJSONVerbose() cause type incorrect")
	}
	if causeData["code"] != "INNER" {
		t.Errorf("MarshalJSONVerbose() cause code incorrect")
	}
}

func TestError_MarshalJSONVerbose_StandardError(t *testing.T) {
	stdErr := errors.New("standard error")
	e := New(ErrorTypeInternal, "WRAPPER", "wrapped").WithCause(stdErr)

	data, err := e.MarshalJSONVerbose()
	if err != nil {
		t.Fatalf("MarshalJSONVerbose() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	causeData, hasCause := result["cause"].(map[string]interface{})
	if !hasCause {
		t.Fatalf("MarshalJSONVerbose() should include standard error cause")
	}

	if causeData["message"] != "standard error" {
		t.Errorf("MarshalJSONVerbose() should preserve standard error message")
	}
}

func TestError_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"type": "VALIDATION",
		"code": "TEST_CODE",
		"message": "test message",
		"details": {
			"field": "email"
		}
	}`

	var e Error
	if err := json.Unmarshal([]byte(jsonData), &e); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if e.Type != ErrorTypeValidation {
		t.Errorf("UnmarshalJSON() Type = %v, want VALIDATION", e.Type)
	}
	if e.Code != "TEST_CODE" {
		t.Errorf("UnmarshalJSON() Code = %v, want TEST_CODE", e.Code)
	}
	if e.Message != "test message" {
		t.Errorf("UnmarshalJSON() Message = %v, want 'test message'", e.Message)
	}
	if e.GetDetail("field") != "email" {
		t.Errorf("UnmarshalJSON() details field incorrect")
	}
}

func TestError_UnmarshalJSON_WithCause(t *testing.T) {
	jsonData := `{
		"type": "DATABASE",
		"code": "OUTER",
		"message": "outer error",
		"cause": {
			"type": "INTERNAL",
			"code": "INNER",
			"message": "inner error"
		}
	}`

	var e Error
	if err := json.Unmarshal([]byte(jsonData), &e); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if e.Cause == nil {
		t.Fatalf("UnmarshalJSON() should restore cause")
	}

	causeErr, ok := e.Cause.(*Error)
	if !ok {
		t.Fatalf("UnmarshalJSON() cause should be *Error")
	}

	if causeErr.Type != ErrorTypeInternal {
		t.Errorf("UnmarshalJSON() cause type incorrect")
	}
	if causeErr.Code != "INNER" {
		t.Errorf("UnmarshalJSON() cause code incorrect")
	}
}

func TestToJSON(t *testing.T) {
	e := New(ErrorTypeValidation, "TEST", "test")
	data, err := ToJSON(e)
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if result["type"] != "VALIDATION" {
		t.Errorf("ToJSON() type incorrect")
	}

	// Test nil case
	data, err = ToJSON(nil)
	if err != nil {
		t.Fatalf("ToJSON(nil) should not error")
	}
	if string(data) != "null" {
		t.Errorf("ToJSON(nil) = %s, want null", string(data))
	}

	// Test standard error
	stdErr := errors.New("standard")
	data, err = ToJSON(stdErr)
	if err != nil {
		t.Fatalf("ToJSON() with standard error should not fail: %v", err)
	}
}

func TestToJSONVerbose(t *testing.T) {
	cause := New(ErrorTypeInternal, "INNER", "inner")
	e := New(ErrorTypeDatabase, "OUTER", "outer").WithCause(cause)

	data, err := ToJSONVerbose(e)
	if err != nil {
		t.Fatalf("ToJSONVerbose() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if _, hasCause := result["cause"]; !hasCause {
		t.Errorf("ToJSONVerbose() should include cause")
	}
}

func TestFromJSON(t *testing.T) {
	jsonData := `{
		"type": "VALIDATION",
		"code": "TEST_CODE",
		"message": "test message"
	}`

	e, err := FromJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	if e.Type != ErrorTypeValidation {
		t.Errorf("FromJSON() Type incorrect")
	}
	if e.Code != "TEST_CODE" {
		t.Errorf("FromJSON() Code incorrect")
	}
}

func TestJSON_RoundTrip(t *testing.T) {
	original := New(ErrorTypeValidation, "TEST_CODE", "test message").
		WithDetail("field", "email").
		WithDetail("reason", "invalid format")

	// Marshal
	data, err := original.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Unmarshal
	restored, err := FromJSON(data)
	if err != nil {
		t.Fatalf("FromJSON() error = %v", err)
	}

	// Compare
	if restored.Type != original.Type {
		t.Errorf("Round trip Type mismatch")
	}
	if restored.Code != original.Code {
		t.Errorf("Round trip Code mismatch")
	}
	if restored.Message != original.Message {
		t.Errorf("Round trip Message mismatch")
	}
	if restored.GetDetail("field") != original.GetDetail("field") {
		t.Errorf("Round trip details mismatch")
	}
}

// Benchmarks

func BenchmarkError_MarshalJSON(b *testing.B) {
	e := New(ErrorTypeValidation, "TEST_CODE", "test message").
		WithDetail("field", "email")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.MarshalJSON()
	}
}

func BenchmarkError_MarshalJSONVerbose(b *testing.B) {
	cause := New(ErrorTypeInternal, "INNER", "inner error")
	e := New(ErrorTypeDatabase, "OUTER", "outer error").WithCause(cause)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.MarshalJSONVerbose()
	}
}

func BenchmarkError_UnmarshalJSON(b *testing.B) {
	jsonData := []byte(`{"type":"VALIDATION","code":"TEST_CODE","message":"test message"}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var e Error
		_ = json.Unmarshal(jsonData, &e)
	}
}

func BenchmarkToJSON(b *testing.B) {
	e := New(ErrorTypeValidation, "TEST", "test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ToJSON(e)
	}
}
