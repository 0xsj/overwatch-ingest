// services/agents/internal/domain/agent/errors.go
package agent

import (
	"fmt"

	"github.com/0xsj/scout/platform/pkg/errors"
)

// Domain error codes
const (
	ErrCodeInvalidAgentID       errors.Code = "AGENT:INVALID_ID"
	ErrCodeInvalidProvider      errors.Code = "AGENT:INVALID_PROVIDER"
	ErrCodeInvalidModel         errors.Code = "AGENT:INVALID_MODEL"
	ErrCodeInvalidStatus        errors.Code = "AGENT:INVALID_STATUS"
	ErrCodeAgentNotInitialized  errors.Code = "AGENT:NOT_INITIALIZED"
	ErrCodeAgentAlreadyActive   errors.Code = "AGENT:ALREADY_ACTIVE"
	ErrCodeAgentDeactivated     errors.Code = "AGENT:DEACTIVATED"
	ErrCodeAgentBusy            errors.Code = "AGENT:BUSY"
	ErrCodeNoTaskInProgress     errors.Code = "AGENT:NO_TASK_IN_PROGRESS"
	ErrCodeInvalidTaskTransition errors.Code = "AGENT:INVALID_TASK_TRANSITION"
)

// Domain errors
var (
	// ErrAgentNotInitialized is returned when trying to operate on a non-initialized agent
	ErrAgentNotInitialized = errors.Validation("agent has not been initialized")

	// ErrAgentAlreadyActive is returned when trying to activate an already active agent
	ErrAgentAlreadyActive = errors.Conflict("agent", "agent is already active")

	// ErrAgentDeactivated is returned when trying to operate on a deactivated agent
	ErrAgentDeactivated = errors.Validation("agent has been deactivated")

	// ErrAgentBusy is returned when trying to assign a task to a busy agent
	ErrAgentBusy = errors.Conflict("agent", "agent is currently busy processing a task")

	// ErrNoTaskInProgress is returned when trying to complete/fail a task with no task in progress
	ErrNoTaskInProgress = errors.Validation("no task currently in progress")
)

// NewInvalidProviderError creates an error for invalid provider.
func NewInvalidProviderError(provider string) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidProvider,
		fmt.Sprintf("invalid provider: %s", provider),
	).WithDetail("provider", provider)
}

// NewInvalidModelError creates an error for invalid model.
func NewInvalidModelError(model string) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidModel,
		fmt.Sprintf("invalid model: %s", model),
	).WithDetail("model", model)
}

// NewInvalidStatusError creates an error for invalid status.
func NewInvalidStatusError(status string) error {
	return errors.New(
		errors.ErrorTypeValidation,
		ErrCodeInvalidStatus,
		fmt.Sprintf("invalid status: %s", status),
	).WithDetail("status", status)
}