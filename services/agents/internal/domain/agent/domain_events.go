// services/agents/internal/domain/agent/domain_events.go
package agent

import (
	"time"
)

// Event type constants
const (
	EventTypeAgentInitialized = "AgentInitialized"
	EventTypeAgentActivated   = "AgentActivated"
	EventTypeAgentDeactivated = "AgentDeactivated"
	EventTypeAgentTaskReceived = "AgentTaskReceived"
	EventTypeAgentTaskCompleted = "AgentTaskCompleted"
	EventTypeAgentTaskFailed   = "AgentTaskFailed"
)

// AgentInitialized is emitted when a new agent is initialized.
type AgentInitialized struct {
	BaseEvent
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

// NewAgentInitialized creates a new AgentInitialized event.
func NewAgentInitialized(aggregateID, name, provider, model string, sequence int64) *AgentInitialized {
	return &AgentInitialized{
		BaseEvent: NewBaseEvent(aggregateID, EventTypeAgentInitialized, sequence),
		Name:      name,
		Provider:  provider,
		Model:     model,
	}
}

// AgentActivated is emitted when an agent is activated.
type AgentActivated struct {
	BaseEvent
	ActivatedAt time.Time `json:"activated_at"`
}

// NewAgentActivated creates a new AgentActivated event.
func NewAgentActivated(aggregateID string, sequence int64) *AgentActivated {
	return &AgentActivated{
		BaseEvent:   NewBaseEvent(aggregateID, EventTypeAgentActivated, sequence),
		ActivatedAt: time.Now(),
	}
}

// AgentDeactivated is emitted when an agent is deactivated.
type AgentDeactivated struct {
	BaseEvent
	Reason        string    `json:"reason,omitempty"`
	DeactivatedAt time.Time `json:"deactivated_at"`
}

// NewAgentDeactivated creates a new AgentDeactivated event.
func NewAgentDeactivated(aggregateID, reason string, sequence int64) *AgentDeactivated {
	return &AgentDeactivated{
		BaseEvent:     NewBaseEvent(aggregateID, EventTypeAgentDeactivated, sequence),
		Reason:        reason,
		DeactivatedAt: time.Now(),
	}
}

// AgentTaskReceived is emitted when an agent receives a task.
type AgentTaskReceived struct {
	BaseEvent
	TaskID   string                 `json:"task_id"`
	TaskType string                 `json:"task_type"`
	Input    map[string]interface{} `json:"input,omitempty"`
}

// NewAgentTaskReceived creates a new AgentTaskReceived event.
func NewAgentTaskReceived(aggregateID, taskID, taskType string, input map[string]interface{}, sequence int64) *AgentTaskReceived {
	return &AgentTaskReceived{
		BaseEvent: NewBaseEvent(aggregateID, EventTypeAgentTaskReceived, sequence),
		TaskID:    taskID,
		TaskType:  taskType,
		Input:     input,
	}
}

// AgentTaskCompleted is emitted when an agent completes a task.
type AgentTaskCompleted struct {
	BaseEvent
	TaskID      string                 `json:"task_id"`
	Output      map[string]interface{} `json:"output,omitempty"`
	CompletedAt time.Time              `json:"completed_at"`
}

// NewAgentTaskCompleted creates a new AgentTaskCompleted event.
func NewAgentTaskCompleted(aggregateID, taskID string, output map[string]interface{}, sequence int64) *AgentTaskCompleted {
	return &AgentTaskCompleted{
		BaseEvent:   NewBaseEvent(aggregateID, EventTypeAgentTaskCompleted, sequence),
		TaskID:      taskID,
		Output:      output,
		CompletedAt: time.Now(),
	}
}

// AgentTaskFailed is emitted when an agent task fails.
type AgentTaskFailed struct {
	BaseEvent
	TaskID   string    `json:"task_id"`
	Reason   string    `json:"reason"`
	FailedAt time.Time `json:"failed_at"`
}

// NewAgentTaskFailed creates a new AgentTaskFailed event.
func NewAgentTaskFailed(aggregateID, taskID, reason string, sequence int64) *AgentTaskFailed {
	return &AgentTaskFailed{
		BaseEvent: NewBaseEvent(aggregateID, EventTypeAgentTaskFailed, sequence),
		TaskID:    taskID,
		Reason:    reason,
		FailedAt:  time.Now(),
	}
}