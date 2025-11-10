// services/agents/internal/domain/agent/aggregate.go
package agent

import (
	"time"

	"github.com/0xsj/scout/platform/pkg/errors"
)

// Agent is the aggregate root for the Agent bounded context.
// It represents an AI/ML-powered assistant that processes tasks using LLMs and tools.
// State is rebuilt from events (Event Sourcing).
type Agent struct {
	// Identity
	id AgentID

	// Configuration
	name     string
	provider Provider
	model    Model

	// State
	status        AgentStatus
	currentTaskID *string // Task currently being processed (if busy)

	// Metadata
	activatedAt   *time.Time
	deactivatedAt *time.Time

	// Event Sourcing metadata
	version           int64   // Current version (number of events applied)
	uncommittedEvents []Event // Events generated but not yet persisted
}

// NewAgent creates a new agent aggregate (for new agents).
func NewAgent(id AgentID) *Agent {
	return &Agent{
		id:                id,
		status:            "", // Will be set when initialized
		version:           0,
		uncommittedEvents: make([]Event, 0),
	}
}

// LoadFromHistory rebuilds an agent aggregate from historical events.
func LoadFromHistory(id AgentID, events []Event) *Agent {
	agent := NewAgent(id)
	for _, event := range events {
		agent.apply(event, false) // false = don't add to uncommitted events
	}
	return agent
}

// ID returns the agent's ID.
func (a *Agent) ID() AgentID {
	return a.id
}

// Version returns the current version.
func (a *Agent) Version() int64 {
	return a.version
}

// UncommittedEvents returns events that have been generated but not persisted.
func (a *Agent) UncommittedEvents() []Event {
	return a.uncommittedEvents
}

// ClearUncommittedEvents clears the uncommitted events (called after persistence).
func (a *Agent) ClearUncommittedEvents() {
	a.uncommittedEvents = make([]Event, 0)
}

// nextSequenceNumber returns the next sequence number for a new event.
func (a *Agent) nextSequenceNumber() int64 {
	return a.version + 1
}

// apply applies an event to the aggregate's state.
// If recordEvent is true, adds to uncommittedEvents.
func (a *Agent) apply(event Event, recordEvent bool) {
	// Apply state changes based on event type
	switch e := event.(type) {
	case *AgentInitialized:
		a.name = e.Name
		a.provider = Provider(e.Provider)
		a.model = Model(e.Model)
		a.status = StatusInitialized

	case *AgentActivated:
		a.status = StatusActive
		a.activatedAt = &e.ActivatedAt

	case *AgentDeactivated:
		a.status = StatusDeactivated
		a.deactivatedAt = &e.DeactivatedAt
		a.currentTaskID = nil // Clear any task in progress

	case *AgentTaskReceived:
		a.status = StatusBusy
		taskID := e.TaskID
		a.currentTaskID = &taskID

	case *AgentTaskCompleted:
		a.status = StatusActive
		a.currentTaskID = nil

	case *AgentTaskFailed:
		a.status = StatusActive
		a.currentTaskID = nil
	}

	// Increment version
	a.version++

	// Record uncommitted event if needed
	if recordEvent {
		a.uncommittedEvents = append(a.uncommittedEvents, event)
	}
}

// === COMMANDS (Business Operations) ===

// Initialize initializes a new agent with configuration.
func (a *Agent) Initialize(name string, provider Provider, model Model) error {
	// Business rules validation
	if a.status != "" {
		return errors.Conflict("agent", "agent has already been initialized")
	}

	if name == "" {
		return errors.RequiredField("name")
	}

	if !provider.IsValid() {
		return NewInvalidProviderError(provider.String())
	}

	if !model.IsValid() {
		return NewInvalidModelError(model.String())
	}

	// Generate and apply event
	event := NewAgentInitialized(
		a.id.String(),
		name,
		provider.String(),
		model.String(),
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// Activate activates the agent, making it ready to receive tasks.
func (a *Agent) Activate() error {
	// Business rules validation
	if a.status == "" {
		return ErrAgentNotInitialized
	}

	if a.status == StatusActive {
		return ErrAgentAlreadyActive
	}

	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	// Generate and apply event
	event := NewAgentActivated(
		a.id.String(),
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// Deactivate deactivates the agent.
func (a *Agent) Deactivate(reason string) error {
	// Business rules validation
	if a.status == "" {
		return ErrAgentNotInitialized
	}

	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	// Generate and apply event
	event := NewAgentDeactivated(
		a.id.String(),
		reason,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// ReceiveTask assigns a task to the agent.
func (a *Agent) ReceiveTask(taskID, taskType string, input map[string]interface{}) error {
	// Business rules validation
	if a.status == "" {
		return ErrAgentNotInitialized
	}

	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	if a.status == StatusBusy {
		return ErrAgentBusy
	}

	if taskID == "" {
		return errors.RequiredField("task_id")
	}

	if taskType == "" {
		return errors.RequiredField("task_type")
	}

	// Generate and apply event
	event := NewAgentTaskReceived(
		a.id.String(),
		taskID,
		taskType,
		input,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// CompleteTask marks the current task as completed.
func (a *Agent) CompleteTask(taskID string, output map[string]interface{}) error {
	// Business rules validation
	if a.currentTaskID == nil {
		return ErrNoTaskInProgress
	}

	if *a.currentTaskID != taskID {
		return errors.New(
			errors.ErrorTypeValidation,
			ErrCodeInvalidTaskTransition,
			"task ID mismatch",
		).WithDetail("expected_task_id", *a.currentTaskID).
			WithDetail("provided_task_id", taskID)
	}

	// Generate and apply event
	event := NewAgentTaskCompleted(
		a.id.String(),
		taskID,
		output,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// FailTask marks the current task as failed.
func (a *Agent) FailTask(taskID, reason string) error {
	// Business rules validation
	if a.currentTaskID == nil {
		return ErrNoTaskInProgress
	}

	if *a.currentTaskID != taskID {
		return errors.New(
			errors.ErrorTypeValidation,
			ErrCodeInvalidTaskTransition,
			"task ID mismatch",
		).WithDetail("expected_task_id", *a.currentTaskID).
			WithDetail("provided_task_id", taskID)
	}

	if reason == "" {
		return errors.RequiredField("reason")
	}

	// Generate and apply event
	event := NewAgentTaskFailed(
		a.id.String(),
		taskID,
		reason,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// === GETTERS (Query State) ===

func (a *Agent) Name() string {
	return a.name
}

func (a *Agent) Provider() Provider {
	return a.provider
}

func (a *Agent) Model() Model {
	return a.model
}

func (a *Agent) Status() AgentStatus {
	return a.status
}

func (a *Agent) CurrentTaskID() *string {
	return a.currentTaskID
}

func (a *Agent) ActivatedAt() *time.Time {
	return a.activatedAt
}

func (a *Agent) DeactivatedAt() *time.Time {
	return a.deactivatedAt
}

func (a *Agent) IsActive() bool {
	return a.status == StatusActive
}

func (a *Agent) IsBusy() bool {
	return a.status == StatusBusy
}

func (a *Agent) IsDeactivated() bool {
	return a.status == StatusDeactivated
}