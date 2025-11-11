// services/agents/internal/domain/agent/agent.go
package agent

import (
	"time"

	"github.com/0xsj/scout/platform/pkg/errors"
)

// Agent is the aggregate root for the Agent bounded context.
// It represents an AI/ML-powered assistant that processes tasks.
type Agent struct {
	id            AgentID
	name          string
	provider      Provider
	model         Model
	status        AgentStatus
	currentTaskID *string
	activatedAt   *time.Time
	deactivatedAt *time.Time
	createdAt     time.Time
	updatedAt     time.Time
}

// NewAgent creates a new agent aggregate.
func NewAgent(id AgentID, name string, provider Provider, model Model) (*Agent, error) {
	// Validation
	if name == "" {
		return nil, errors.RequiredField("name")
	}
	if !provider.IsValid() {
		return nil, NewInvalidProviderError(provider.String())
	}
	if !model.IsValid() {
		return nil, NewInvalidModelError(model.String())
	}

	now := time.Now()
	return &Agent{
		id:        id,
		name:      name,
		provider:  provider,
		model:     model,
		status:    StatusInitialized,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// LoadAgent reconstructs an agent from the database.
// Used by the repository to load existing agents.
func LoadAgent(
	id AgentID,
	name string,
	provider Provider,
	model Model,
	status AgentStatus,
	currentTaskID *string,
	activatedAt, deactivatedAt *time.Time,
	createdAt, updatedAt time.Time,
) *Agent {
	return &Agent{
		id:            id,
		name:          name,
		provider:      provider,
		model:         model,
		status:        status,
		currentTaskID: currentTaskID,
		activatedAt:   activatedAt,
		deactivatedAt: deactivatedAt,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

// === COMMANDS (Business Operations) ===

// Activate activates the agent, making it ready to receive tasks.
func (a *Agent) Activate() error {
	if a.status == StatusActive {
		return ErrAgentAlreadyActive
	}
	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	now := time.Now()
	a.status = StatusActive
	a.activatedAt = &now
	a.updatedAt = now
	return nil
}

// ReceiveTask assigns a task to the agent.
func (a *Agent) ReceiveTask(taskID string) error {
	if a.status != StatusActive {
		return errors.Validation("agent must be active to receive tasks")
	}
	if a.status == StatusBusy {
		return ErrAgentBusy
	}
	if taskID == "" {
		return errors.RequiredField("task_id")
	}

	a.status = StatusBusy
	a.currentTaskID = &taskID
	a.updatedAt = time.Now()
	return nil
}

// CompleteTask marks the current task as completed.
func (a *Agent) CompleteTask(taskID string) error {
	if a.currentTaskID == nil {
		return ErrNoTaskInProgress
	}
	if *a.currentTaskID != taskID {
		return errors.New(
			errors.ErrorTypeValidation,
			"TASK_MISMATCH",
			"task ID does not match current task",
		).WithDetail("expected", *a.currentTaskID).
			WithDetail("provided", taskID)
	}

	a.status = StatusActive
	a.currentTaskID = nil
	a.updatedAt = time.Now()
	return nil
}

// FailTask marks the current task as failed.
func (a *Agent) FailTask(taskID string) error {
	if a.currentTaskID == nil {
		return ErrNoTaskInProgress
	}
	if *a.currentTaskID != taskID {
		return errors.New(
			errors.ErrorTypeValidation,
			"TASK_MISMATCH",
			"task ID does not match current task",
		).WithDetail("expected", *a.currentTaskID).
			WithDetail("provided", taskID)
	}

	a.status = StatusActive
	a.currentTaskID = nil
	a.updatedAt = time.Now()
	return nil
}

// Deactivate deactivates the agent.
func (a *Agent) Deactivate() error {
	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	now := time.Now()
	a.status = StatusDeactivated
	a.deactivatedAt = &now
	a.currentTaskID = nil
	a.updatedAt = now
	return nil
}

// === GETTERS (Query State) ===

func (a *Agent) ID() AgentID               { return a.id }
func (a *Agent) Name() string              { return a.name }
func (a *Agent) Provider() Provider        { return a.provider }
func (a *Agent) Model() Model              { return a.model }
func (a *Agent) Status() AgentStatus       { return a.status }
func (a *Agent) CurrentTaskID() *string    { return a.currentTaskID }
func (a *Agent) ActivatedAt() *time.Time   { return a.activatedAt }
func (a *Agent) DeactivatedAt() *time.Time { return a.deactivatedAt }
func (a *Agent) CreatedAt() time.Time      { return a.createdAt }
func (a *Agent) UpdatedAt() time.Time      { return a.updatedAt }
func (a *Agent) IsActive() bool            { return a.status == StatusActive }
func (a *Agent) IsBusy() bool              { return a.status == StatusBusy }
func (a *Agent) IsDeactivated() bool       { return a.status == StatusDeactivated }