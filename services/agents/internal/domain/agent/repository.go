// services/agents/internal/domain/agent/repository.go
package agent

import (
	"context"
)

// Repository defines the interface for agent persistence.
// Implemented by infrastructure layer (e.g., PostgresRepository).
type Repository interface {
	// Save creates or updates an agent.
	// Returns error if the operation fails.
	Save(ctx context.Context, agent *Agent) error

	// FindByID retrieves an agent by its ID.
	// Returns nil, nil if not found.
	FindByID(ctx context.Context, id AgentID) (*Agent, error)

	// FindAll retrieves all agents with pagination.
	FindAll(ctx context.Context, limit, offset int) ([]*Agent, error)

	// FindByStatus retrieves agents by status.
	FindByStatus(ctx context.Context, status AgentStatus, limit, offset int) ([]*Agent, error)

	// FindByProvider retrieves agents by provider.
	FindByProvider(ctx context.Context, provider Provider, limit, offset int) ([]*Agent, error)

	// Exists checks if an agent exists by ID.
	Exists(ctx context.Context, id AgentID) (bool, error)

	// Delete removes an agent.
	// Returns error if the agent doesn't exist or deletion fails.
	Delete(ctx context.Context, id AgentID) error

	// Count returns the total number of agents.
	Count(ctx context.Context) (int, error)

	// CountByStatus returns the number of agents with a specific status.
	CountByStatus(ctx context.Context, status AgentStatus) (int, error)
}