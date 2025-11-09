// services/agents/internal/domain/agent/errors.go
package agent

import (
	"fmt"

	"github.com/0xsj/scout/platform/pkg/errors"
)

// Domain error codes
const (
	ErrCodeInvalidAgentName        errors.Code = "AGENT:INVALID_NAME"
	ErrCodeInvalidAgentType        errors.Code = "AGENT:INVALID_TYPE"
	ErrCodeAgentAlreadyCreated     errors.Code = "AGENT:ALREADY_CREATED"
	ErrCodeAgentNotCreated         errors.Code = "AGENT:NOT_CREATED"
	ErrCodeAgentDeactivated        errors.Code = "AGENT:DEACTIVATED"
	ErrCodeInvalidLocation         errors.Code = "AGENT:INVALID_LOCATION"
	ErrCodeInvalidStatusTransition errors.Code = "AGENT:INVALID_STATUS_TRANSITION"
)

// Domain errors using platform constructors
var (
	// ErrInvalidAgentName is returned when agent name is invalid
	ErrInvalidAgentName = errors.Validation("agent name is required and cannot be empty")

	// ErrAgentAlreadyCreated is returned when trying to create an already created agent
	ErrAgentAlreadyCreated = errors.Conflict("agent", "agent has already been created")

	// ErrAgentNotCreated is returned when trying to operate on an agent that hasn't been created
	ErrAgentNotCreated = errors.Validation("agent has not been created yet")

	// ErrAgentDeactivated is returned when trying to operate on a deactivated agent
	ErrAgentDeactivated = errors.Validation("agent has been deactivated")
)

// NewInvalidAgentTypeError creates an error for invalid agent type with the provided value.
func NewInvalidAgentTypeError(agentType string) error {
	return errors.InvalidField("agent_type", agentType)
}

// NewInvalidStatusTransitionError creates an error for invalid status transitions.
func NewInvalidStatusTransitionError(from, to Status) error {
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