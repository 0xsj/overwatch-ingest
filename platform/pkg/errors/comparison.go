// platform/pkg/errors/comparison.go
package errors

// IsType checks if an error is of a specific ErrorType.
// Returns false if err is nil or not an *Error.
func IsType(err error, errType ErrorType) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	return e.Type == errType
}

// IsOneOfTypes checks if an error is one of the specified ErrorTypes.
// Returns false if err is nil or not an *Error.
func IsOneOfTypes(err error, types ...ErrorType) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	for _, t := range types {
		if e.Type == t {
			return true
		}
	}

	return false
}

// HasCode checks if an error has a specific Code.
// Returns false if err is nil or not an *Error.
func HasCode(err error, code Code) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	return e.Code == code
}

// HasOneOfCodes checks if an error has one of the specified Codes.
// Returns false if err is nil or not an *Error.
func HasOneOfCodes(err error, codes ...Code) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	for _, c := range codes {
		if e.Code == c {
			return true
		}
	}

	return false
}

// IsClientError checks if an error is a client error (4xx equivalent).
// Returns false if err is nil or not an *Error.
func IsClientError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	return e.IsClientError()
}

// IsServerError checks if an error is a server error (5xx equivalent).
// Returns false if err is nil or not an *Error.
func IsServerError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	return e.IsServerError()
}

// IsRetryable checks if an error can typically be retried.
// Returns false if err is nil or not an *Error.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	return e.IsRetryable()
}

// GetType extracts the ErrorType from an error.
// Returns empty ErrorType if err is nil or not an *Error.
func GetType(err error) ErrorType {
	if err == nil {
		return ""
	}

	e, ok := As(err)
	if !ok {
		return ""
	}

	return e.Type
}

// GetCode extracts the Code from an error.
// Returns empty Code if err is nil or not an *Error.
func GetCode(err error) Code {
	if err == nil {
		return ""
	}

	e, ok := As(err)
	if !ok {
		return ""
	}

	return e.Code
}

// GetMessage extracts the message from an error.
// Returns empty string if err is nil.
// For non-*Error types, returns err.Error().
func GetMessage(err error) string {
	if err == nil {
		return ""
	}

	e, ok := As(err)
	if !ok {
		return err.Error()
	}

	return e.Message
}

// GetDetails extracts the details map from an error.
// Returns nil if err is nil or not an *Error.
// Returns a copy to prevent external modification.
func GetDetails(err error) map[string]string {
	if err == nil {
		return nil
	}

	e, ok := As(err)
	if !ok {
		return nil
	}

	if e.Details == nil {
		return nil
	}

	// Return a copy to prevent external modification
	details := make(map[string]string, len(e.Details))
	for k, v := range e.Details {
		details[k] = v
	}

	return details
}

// GetDetail extracts a specific detail value by key.
// Returns empty string if err is nil, not an *Error, or key doesn't exist.
func GetDetail(err error, key string) string {
	if err == nil {
		return ""
	}

	e, ok := As(err)
	if !ok {
		return ""
	}

	return e.GetDetail(key)
}

// HasDetail checks if an error has a specific detail key.
// Returns false if err is nil or not an *Error.
func HasDetail(err error, key string) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	return e.HasDetail(key)
}

// GetCause extracts the underlying cause from an error.
// Returns nil if err is nil or has no cause.
// Works with both *Error and standard library error wrapping.
func GetCause(err error) error {
	if err == nil {
		return nil
	}

	// Try *Error first
	if e, ok := As(err); ok {
		return e.Cause
	}

	// Fallback to standard library Unwrap
	type unwrapper interface {
		Unwrap() error
	}

	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}

	return nil
}

// GetRootCause extracts the root cause by unwrapping all error layers.
// Returns the deepest error in the chain.
// Returns nil if err is nil.
func GetRootCause(err error) error {
	if err == nil {
		return nil
	}

	for {
		cause := GetCause(err)
		if cause == nil {
			return err
		}
		err = cause
	}
}

// Match checks if an error matches specific criteria.
// All non-zero criteria must match for the function to return true.
// Returns false if err is nil or not an *Error.
//
// Example:
//
//	if errors.Match(err, errors.ErrorTypeValidation, "REQUIRED_FIELD", "") {
//	    // Handle validation error with specific code
//	}
func Match(err error, errType ErrorType, code Code, messageContains string) bool {
	if err == nil {
		return false
	}

	e, ok := As(err)
	if !ok {
		return false
	}

	// Check type if specified
	if errType != "" && e.Type != errType {
		return false
	}

	// Check code if specified
	if code != "" && e.Code != code {
		return false
	}

	// Check message if specified
	if messageContains != "" {
		if !contains(e.Message, messageContains) {
			return false
		}
	}

	return true
}

// contains is a simple substring check helper.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexOfSubstring(s, substr) >= 0)
}

// indexOfSubstring finds the index of substr in s.
func indexOfSubstring(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	}
	if n > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-n; i++ {
		if s[i:i+n] == substr {
			return i
		}
	}
	return -1
}
