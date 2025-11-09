// services/agents/internal/application/dto/command_dto.go
package dto

// CreateAgentCommand represents a command to create a new agent.
type CreateAgentCommand struct {
	AgentID   string  // Optional: if empty, will be generated
	Name      string
	AgentType string
	CreatedBy string  // User/system that created the agent
}

// DeployAgentCommand represents a command to deploy an agent.
type DeployAgentCommand struct {
	AgentID    string
	Latitude   float64
	Longitude  float64
	DeployedBy string
}

// UpdateAgentStatusCommand represents a command to change agent status.
type UpdateAgentStatusCommand struct {
	AgentID   string
	NewStatus string
	Reason    string  // Why the status is changing
}

// UpdateAgentLocationCommand represents a command to update agent location.
type UpdateAgentLocationCommand struct {
	AgentID   string
	Latitude  float64
	Longitude float64
}

// DeactivateAgentCommand represents a command to deactivate an agent.
type DeactivateAgentCommand struct {
	AgentID       string
	Reason        string
	DeactivatedBy string
}