// services/agents/internal/domain/agent/events.go
package agent

import (
	"time"
)

// Event type constants
const (
	EventTypeAgentCreated         = "AgentCreated"
	EventTypeAgentDeployed        = "AgentDeployed"
	EventTypeAgentStatusChanged   = "AgentStatusChanged"
	EventTypeAgentLocationUpdated = "AgentLocationUpdated"
	EventTypeAgentDeactivated     = "AgentDeactivated"
)

// AgentCreated is emitted when a new agent is created.
type AgentCreated struct {
	BaseEvent
	Name      string `json:"name"`
	AgentType string `json:"agent_type"`
	CreatedBy string `json:"created_by,omitempty"`
}

// NewAgentCreated creates a new AgentCreated event.
func NewAgentCreated(aggregateID, name, agentType, createdBy string, sequence int64) *AgentCreated {
	return &AgentCreated{
		BaseEvent: NewBaseEvent(aggregateID, EventTypeAgentCreated, sequence),
		Name:      name,
		AgentType: agentType,
		CreatedBy: createdBy,
	}
}

// AgentDeployed is emitted when an agent is deployed to the field.
type AgentDeployed struct {
	BaseEvent
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	DeployedAt time.Time `json:"deployed_at"`
	DeployedBy string    `json:"deployed_by,omitempty"`
}

// NewAgentDeployed creates a new AgentDeployed event.
func NewAgentDeployed(aggregateID string, lat, lon float64, deployedBy string, sequence int64) *AgentDeployed {
	return &AgentDeployed{
		BaseEvent:  NewBaseEvent(aggregateID, EventTypeAgentDeployed, sequence),
		Latitude:   lat,
		Longitude:  lon,
		DeployedAt: time.Now(),
		DeployedBy: deployedBy,
	}
}

// AgentStatusChanged is emitted when an agent's status changes.
type AgentStatusChanged struct {
	BaseEvent
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Reason    string    `json:"reason,omitempty"`
	ChangedAt time.Time `json:"changed_at"`
}

// NewAgentStatusChanged creates a new AgentStatusChanged event.
func NewAgentStatusChanged(aggregateID, oldStatus, newStatus, reason string, sequence int64) *AgentStatusChanged {
	return &AgentStatusChanged{
		BaseEvent: NewBaseEvent(aggregateID, EventTypeAgentStatusChanged, sequence),
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Reason:    reason,
		ChangedAt: time.Now(),
	}
}

// AgentLocationUpdated is emitted when an agent's location changes.
type AgentLocationUpdated struct {
	BaseEvent
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewAgentLocationUpdated creates a new AgentLocationUpdated event.
func NewAgentLocationUpdated(aggregateID string, lat, lon float64, sequence int64) *AgentLocationUpdated {
	return &AgentLocationUpdated{
		BaseEvent: NewBaseEvent(aggregateID, EventTypeAgentLocationUpdated, sequence),
		Latitude:  lat,
		Longitude: lon,
		UpdatedAt: time.Now(),
	}
}

// AgentDeactivated is emitted when an agent is deactivated.
type AgentDeactivated struct {
	BaseEvent
	Reason        string    `json:"reason,omitempty"`
	DeactivatedAt time.Time `json:"deactivated_at"`
	DeactivatedBy string    `json:"deactivated_by,omitempty"`
}

// NewAgentDeactivated creates a new AgentDeactivated event.
func NewAgentDeactivated(aggregateID, reason, deactivatedBy string, sequence int64) *AgentDeactivated {
	return &AgentDeactivated{
		BaseEvent:     NewBaseEvent(aggregateID, EventTypeAgentDeactivated, sequence),
		Reason:        reason,
		DeactivatedAt: time.Now(),
		DeactivatedBy: deactivatedBy,
	}
}