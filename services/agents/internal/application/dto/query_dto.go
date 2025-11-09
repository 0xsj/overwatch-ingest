// services/agents/internal/application/dto/query_dto.go
package dto

import (
	"time"
)

// GetAgentQuery represents a query to get a single agent by ID.
type GetAgentQuery struct {
	AgentID string
}

// ListAgentsQuery represents a query to list agents with filters.
type ListAgentsQuery struct {
	Status    string  // Optional: filter by status
	AgentType string  // Optional: filter by type
	Limit     int     // Pagination
	Offset    int     // Pagination
}

// AgentDTO is the data transfer object for an agent.
// This is what gets returned from queries.
type AgentDTO struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	AgentType     string     `json:"agent_type"`
	Status        string     `json:"status"`
	Location      *LocationDTO `json:"location,omitempty"`
	DeployedAt    *time.Time `json:"deployed_at,omitempty"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Version       int64      `json:"version"`
}

// LocationDTO is the data transfer object for location.
type LocationDTO struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// AgentListDTO is the response for listing agents.
type AgentListDTO struct {
	Agents     []AgentDTO `json:"agents"`
	TotalCount int        `json:"total_count"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}