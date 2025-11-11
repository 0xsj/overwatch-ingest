// services/agents/internal/domain/agent/events.go
package agent

import "time"

// Domain events for integration between services.
// These are published to NATS/RabbitMQ for inter-service communication.

// AgentCreatedEvent is published when a new agent is created.
type AgentCreatedEvent struct {
	AgentID   string    `json:"agent_id"`
	Name      string    `json:"name"`
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
}

// AgentActivatedEvent is published when an agent is activated.
type AgentActivatedEvent struct {
	AgentID     string    `json:"agent_id"`
	ActivatedAt time.Time `json:"activated_at"`
}

// AgentTaskReceivedEvent is published when an agent receives a task.
type AgentTaskReceivedEvent struct {
	AgentID    string    `json:"agent_id"`
	TaskID     string    `json:"task_id"`
	ReceivedAt time.Time `json:"received_at"`
}

// AgentTaskCompletedEvent is published when an agent completes a task.
type AgentTaskCompletedEvent struct {
	AgentID     string    `json:"agent_id"`
	TaskID      string    `json:"task_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// AgentTaskFailedEvent is published when an agent task fails.
type AgentTaskFailedEvent struct {
	AgentID  string    `json:"agent_id"`
	TaskID   string    `json:"task_id"`
	FailedAt time.Time `json:"failed_at"`
}

// AgentDeactivatedEvent is published when an agent is deactivated.
type AgentDeactivatedEvent struct {
	AgentID       string    `json:"agent_id"`
	DeactivatedAt time.Time `json:"deactivated_at"`
}