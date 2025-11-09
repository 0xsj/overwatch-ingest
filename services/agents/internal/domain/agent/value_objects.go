// services/agents/internal/domain/agent/value_objects.go
package agent

import (
	"fmt"

	"github.com/google/uuid"
)

// AgentID is a unique identifier for an agent.
type AgentID struct {
	value string
}

// NewAgentID creates a new AgentID.
func NewAgentID() AgentID {
	return AgentID{value: uuid.New().String()}
}

// ParseAgentID parses a string into an AgentID.
func ParseAgentID(id string) (AgentID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return AgentID{}, fmt.Errorf("invalid agent id: %w", err)
	}
	return AgentID{value: id}, nil
}

func (id AgentID) String() string {
	return id.value
}

func (id AgentID) Equals(other AgentID) bool {
	return id.value == other.value
}

// Status represents an agent's operational status.
type Status string

const (
	StatusCreated     Status = "created"
	StatusDeployed    Status = "deployed"
	StatusActive      Status = "active"
	StatusInactive    Status = "inactive"
	StatusOffline     Status = "offline"
	StatusDeactivated Status = "deactivated"
)

// String returns the string representation.
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is valid.
func (s Status) IsValid() bool {
	switch s {
	case StatusCreated, StatusDeployed, StatusActive, StatusInactive, StatusOffline, StatusDeactivated:
		return true
	}
	return false
}

// CanTransitionTo checks if a status transition is valid.
func (s Status) CanTransitionTo(newStatus Status) bool {
	transitions := map[Status][]Status{
		StatusCreated:     {StatusDeployed, StatusDeactivated},
		StatusDeployed:    {StatusActive, StatusInactive, StatusDeactivated},
		StatusActive:      {StatusInactive, StatusOffline, StatusDeactivated},
		StatusInactive:    {StatusActive, StatusOffline, StatusDeactivated},
		StatusOffline:     {StatusActive, StatusInactive, StatusDeactivated},
		StatusDeactivated: {}, // Terminal state
	}

	validTransitions, ok := transitions[s]
	if !ok {
		return false
	}

	for _, validStatus := range validTransitions {
		if validStatus == newStatus {
			return true
		}
	}

	return false
}

// AgentType represents the type/role of an agent.
type AgentType string

const (
	AgentTypeFieldResponder  AgentType = "field_responder"
	AgentTypeIncidentManager AgentType = "incident_manager"
	AgentTypeCoordinator     AgentType = "coordinator"
	AgentTypeMedic           AgentType = "medic"
	AgentTypeSpecialist      AgentType = "specialist"
)

// String returns the string representation.
func (t AgentType) String() string {
	return string(t)
}

// IsValid checks if the agent type is valid.
func (t AgentType) IsValid() bool {
	switch t {
	case AgentTypeFieldResponder, AgentTypeIncidentManager, AgentTypeCoordinator, AgentTypeMedic, AgentTypeSpecialist:
		return true
	}
	return false
}

// ParseAgentType parses a string into an AgentType.
func ParseAgentType(s string) (AgentType, error) {
	t := AgentType(s)
	if !t.IsValid() {
		return "", fmt.Errorf("invalid agent type: %s", s)
	}
	return t, nil
}

// Location represents a geographic location.
type Location struct {
	Latitude  float64
	Longitude float64
}

// NewLocation creates a new Location.
func NewLocation(lat, lon float64) (Location, error) {
	if lat < -90 || lat > 90 {
		return Location{}, fmt.Errorf("invalid latitude: %f", lat)
	}
	if lon < -180 || lon > 180 {
		return Location{}, fmt.Errorf("invalid longitude: %f", lon)
	}
	return Location{Latitude: lat, Longitude: lon}, nil
}

// Equals checks if two locations are equal.
func (l Location) Equals(other Location) bool {
	return l.Latitude == other.Latitude && l.Longitude == other.Longitude
}