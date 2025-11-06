// platform/pkg/errors/json.go
package errors

import (
	"encoding/json"
)

// errorJSON is the JSON representation of an Error.
// Uses snake_case for broader compatibility (REST, gRPC-Gateway, etc.)
type errorJSON struct {
	Type    string            `json:"type"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Cause   *errorJSON        `json:"cause,omitempty"`
}

// MarshalJSON implements json.Marshaler.
// Excludes the cause chain for security (prevents leaking internal errors).
// Use MarshalJSONVerbose() if you need the full error chain.
func (e *Error) MarshalJSON() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}

	ej := errorJSON{
		Type:    e.Type.String(),
		Code:    e.Code.String(),
		Message: e.Message,
		Details: e.Details,
		// Cause intentionally omitted
	}

	return json.Marshal(ej)
}

// MarshalJSONVerbose marshals the error including the full cause chain.
// Use this for internal debugging/logging, not for external APIs.
func (e *Error) MarshalJSONVerbose() ([]byte, error) {
	if e == nil {
		return []byte("null"), nil
	}

	ej := errorJSON{
		Type:    e.Type.String(),
		Code:    e.Code.String(),
		Message: e.Message,
		Details: e.Details,
	}

	// Include cause chain if present
	if e.Cause != nil {
		// If cause is *Error, marshal it recursively
		if causeErr, ok := e.Cause.(*Error); ok {
			causeBytes, err := causeErr.MarshalJSONVerbose()
			if err != nil {
				return nil, err
			}
			var causeJSON errorJSON
			if err := json.Unmarshal(causeBytes, &causeJSON); err != nil {
				return nil, err
			}
			ej.Cause = &causeJSON
		} else {
			// For non-*Error causes, create a simple representation
			ej.Cause = &errorJSON{
				Type:    ErrorTypeInternal.String(),
				Code:    "WRAPPED_ERROR",
				Message: e.Cause.Error(),
			}
		}
	}

	return json.Marshal(ej)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Error) UnmarshalJSON(data []byte) error {
	var ej errorJSON
	if err := json.Unmarshal(data, &ej); err != nil {
		return err
	}

	e.Type = ErrorType(ej.Type)
	e.Code = Code(ej.Code)
	e.Message = ej.Message
	e.Details = ej.Details

	// Reconstruct cause chain if present
	if ej.Cause != nil {
		causeErr := &Error{}
		causeBytes, err := json.Marshal(ej.Cause)
		if err != nil {
			return err
		}
		if err := causeErr.UnmarshalJSON(causeBytes); err != nil {
			return err
		}
		e.Cause = causeErr
	}

	return nil
}

// ToJSON converts an error to JSON bytes.
// Returns the default (safe) JSON representation without cause chain.
// Returns nil if err is nil.
func ToJSON(err error) ([]byte, error) {
	if err == nil {
		return []byte("null"), nil
	}

	if e, ok := err.(*Error); ok {
		return e.MarshalJSON()
	}

	// For non-*Error types, create a generic representation
	genericErr := Internal(err.Error())
	return genericErr.MarshalJSON()
}

// ToJSONVerbose converts an error to verbose JSON bytes including cause chain.
// Use for debugging/logging, not for external APIs.
// Returns nil if err is nil.
func ToJSONVerbose(err error) ([]byte, error) {
	if err == nil {
		return []byte("null"), nil
	}

	if e, ok := err.(*Error); ok {
		return e.MarshalJSONVerbose()
	}

	// For non-*Error types, create a generic representation
	genericErr := Internal(err.Error())
	return genericErr.MarshalJSONVerbose()
}

// FromJSON deserializes an Error from JSON bytes.
func FromJSON(data []byte) (*Error, error) {
	var e Error
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	return &e, nil
}
