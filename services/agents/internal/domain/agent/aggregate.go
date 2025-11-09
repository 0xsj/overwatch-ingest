// services/agents/internal/domain/agent/aggregate.go
package agent

import (
	"time"
)

// Agent is the aggregate root for the Agent bounded context.
// It represents an emergency response agent and encapsulates all business rules.
// State is rebuilt from events (Event Sourcing).
type Agent struct {
	// Identity
	id AgentID

	// State (rebuilt from events)
	name          string
	agentType     AgentType
	status        Status
	location      *Location
	deployedAt    *time.Time
	deactivatedAt *time.Time

	// Event Sourcing metadata
	version           int64   // Current version (number of events applied)
	uncommittedEvents []Event // Events generated but not yet persisted
}

// NewAgent creates a new agent aggregate (for new agents).
func NewAgent(id AgentID) *Agent {
	return &Agent{
		id:                id,
		status:            "", // Will be set when created
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
	case *AgentCreated:
		a.name = e.Name
		a.agentType = AgentType(e.AgentType)
		a.status = StatusCreated

	case *AgentDeployed:
		loc := Location{Latitude: e.Latitude, Longitude: e.Longitude}
		a.location = &loc
		a.deployedAt = &e.DeployedAt
		a.status = StatusDeployed

	case *AgentStatusChanged:
		a.status = Status(e.NewStatus)

	case *AgentLocationUpdated:
		loc := Location{Latitude: e.Latitude, Longitude: e.Longitude}
		a.location = &loc

	case *AgentDeactivated:
		a.status = StatusDeactivated
		a.deactivatedAt = &e.DeactivatedAt
	}

	// Increment version
	a.version++

	// Record uncommitted event if needed
	if recordEvent {
		a.uncommittedEvents = append(a.uncommittedEvents, event)
	}
}

// === COMMANDS (Business Operations) ===

// Create creates a new agent.
func (a *Agent) Create(name string, agentType AgentType, createdBy string) error {
	// Business rules validation
	if a.status != "" {
		return ErrAgentAlreadyCreated
	}

	if name == "" {
		return ErrInvalidAgentName
	}

	if !agentType.IsValid() {
		return NewInvalidAgentTypeError(agentType.String())
	}

	// Generate and apply event
	event := NewAgentCreated(
		a.id.String(),
		name,
		agentType.String(),
		createdBy,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// Deploy deploys the agent to a specific location.
func (a *Agent) Deploy(lat, lon float64, deployedBy string) error {
	// Business rules validation
	if a.status == "" {
		return ErrAgentNotCreated
	}

	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	location, err := NewLocation(lat, lon)
	if err != nil {
		return err
	}

	// Generate and apply event
	event := NewAgentDeployed(
		a.id.String(),
		location.Latitude,
		location.Longitude,
		deployedBy,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// ChangeStatus changes the agent's operational status.
func (a *Agent) ChangeStatus(newStatus Status, reason string) error {
	// Business rules validation
	if a.status == "" {
		return ErrAgentNotCreated
	}

	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	if !newStatus.IsValid() {
		return NewInvalidAgentTypeError(newStatus.String())
	}

	// Check if transition is valid
	if !a.status.CanTransitionTo(newStatus) {
		return NewInvalidStatusTransitionError(a.status, newStatus)
	}

	// Generate and apply event
	event := NewAgentStatusChanged(
		a.id.String(),
		a.status.String(),
		newStatus.String(),
		reason,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// UpdateLocation updates the agent's current location.
func (a *Agent) UpdateLocation(lat, lon float64) error {
	// Business rules validation
	if a.status == "" {
		return ErrAgentNotCreated
	}

	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	location, err := NewLocation(lat, lon)
	if err != nil {
		return err
	}

	// Generate and apply event
	event := NewAgentLocationUpdated(
		a.id.String(),
		location.Latitude,
		location.Longitude,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// Deactivate deactivates the agent.
func (a *Agent) Deactivate(reason, deactivatedBy string) error {
	// Business rules validation
	if a.status == "" {
		return ErrAgentNotCreated
	}

	if a.status == StatusDeactivated {
		return ErrAgentDeactivated
	}

	// Generate and apply event
	event := NewAgentDeactivated(
		a.id.String(),
		reason,
		deactivatedBy,
		a.nextSequenceNumber(),
	)

	a.apply(event, true)

	return nil
}

// === GETTERS (Query State) ===

func (a *Agent) Name() string {
	return a.name
}

func (a *Agent) Type() AgentType {
	return a.agentType
}

func (a *Agent) Status() Status {
	return a.status
}

func (a *Agent) Location() *Location {
	return a.location
}

func (a *Agent) DeployedAt() *time.Time {
	return a.deployedAt
}

func (a *Agent) DeactivatedAt() *time.Time {
	return a.deactivatedAt
}

func (a *Agent) IsDeactivated() bool {
	return a.status == StatusDeactivated
}