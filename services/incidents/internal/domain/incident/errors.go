// services/incidents/internal/domain/incident/errors.go
package incident

import (
	"fmt"

	"github.com/0xsj/scout/platform/pkg/errors"
)

// Domain error codes
const (
	ErrCodeInvalidIncidentID       errors.Code = "INCIDENT:INVALID_ID"
	ErrCodeInvalidSeverity         errors.Code = "INCIDENT:INVALID_SEVERITY"
	ErrCodeInvalidIncidentType     errors.Code = "INCIDENT:INVALID_TYPE"
	ErrCodeInvalidStatus           errors.Code = "INCIDENT:INVALID_STATUS"
	ErrCodeInvalidLocation         errors.Code = "INCIDENT:INVALID_LOCATION"
	ErrCodeInvalidReportSource     errors.Code = "INCIDENT:INVALID_REPORT_SOURCE"
	ErrCodeIncidentNotReported     errors.Code = "INCIDENT:NOT_REPORTED"
	ErrCodeIncidentAlreadyReported errors.Code = "INCIDENT:ALREADY_REPORTED"
	ErrCodeIncidentAlreadyVerified errors.Code = "INCIDENT:ALREADY_VERIFIED"
	ErrCodeIncidentNotVerified     errors.Code = "INCIDENT:NOT_VERIFIED"
	ErrCodeIncidentAlreadyClosed   errors.Code = "INCIDENT:ALREADY_CLOSED"
	ErrCodeInvalidStatusTransition errors.Code = "INCIDENT:INVALID_STATUS_TRANSITION"
	ErrCodeNoRespondersProvided    errors.Code = "INCIDENT:NO_RESPONDERS"
)

// Domain errors
var (
	// ErrIncidentNotReported is returned when trying to operate on a non-reported incident
	ErrIncidentNotReported = errors.Validation("incident has not been reported")

	// ErrIncidentAlreadyReported is returned when trying to report an already reported incident
	ErrIncidentAlreadyReported = errors.Conflict("incident", "incident has already been reported")

	// ErrIncidentAlreadyVerified is returned when trying to verify an already verified incident
	ErrIncidentAlreadyVerified = errors.Conflict("incident", "incident has already been verified")

	// ErrIncidentNotVerified is returned when trying to dispatch an unverified incident
	ErrIncidentNotVerified = errors.Validation("incident must be verified before dispatch")

	// ErrIncidentAlreadyClosed is returned when trying to operate on a closed incident
	ErrIncidentAlreadyClosed = errors.Validation("incident has been closed")

	// ErrNoRespondersProvided is returned when dispatching with no responders
	ErrNoRespondersProvided = errors.Validation("at least one responder must be provided for dispatch")
)

// NewInvalidSeverityError creates an error for invalid severity.
func NewInvalidSeverityError(severity string) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidSeverity,
		fmt.Sprintf("invalid severity: %s", severity),
	).WithDetail("severity", severity)
}

// NewInvalidIncidentTypeError creates an error for invalid incident type.
func NewInvalidIncidentTypeError(incidentType string) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidIncidentType,
		fmt.Sprintf("invalid incident type: %s", incidentType),
	).WithDetail("incident_type", incidentType)
}

// NewInvalidReportSourceError creates an error for invalid report source.
func NewInvalidReportSourceError(source string) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidReportSource,
		fmt.Sprintf("invalid report source: %s", source),
	).WithDetail("report_source", source)
}

// NewInvalidStatusError creates an error for invalid status.
func NewInvalidStatusError(status string) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidStatus,
		fmt.Sprintf("invalid status: %s", status),
	).WithDetail("status", status)
}

// NewInvalidStatusTransitionError creates an error for invalid status transitions.
func NewInvalidStatusTransitionError(from, to IncidentStatus) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidStatusTransition,
		"invalid status transition",
	).WithDetail("from_status", from.String()).
		WithDetail("to_status", to.String())
}

// NewInvalidLocationError creates an error for invalid location with coordinates.
func NewInvalidLocationError(lat, lon float64, reason string) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidLocation,
		reason,
	).WithDetail("latitude", fmt.Sprintf("%f", lat)).
		WithDetail("longitude", fmt.Sprintf("%f", lon))
}