// platform/pkg/errors/codes.go
package errors

// Code represents a specific error code.
// Codes are string identifiers that uniquely identify error conditions.
//
// Naming convention (recommended but not enforced):
//   - Use UPPER_SNAKE_CASE
//   - Format: ENTITY_CONDITION or OPERATION_FAILURE
//   - Examples: "USER_NOT_FOUND", "DATABASE_CONNECTION_FAILED"
//
// Services should define their own code constants:
//
//	const (
//	    CodeIncidentNotFound errors.Code = "INCIDENT_NOT_FOUND"
//	    CodeInvalidSeverity  errors.Code = "INVALID_SEVERITY"
//	)
type Code string

// String returns the string representation of the Code.
func (c Code) String() string {
	return string(c)
}

// IsEmpty returns true if the code is empty.
func (c Code) IsEmpty() bool {
	return c == ""
}
