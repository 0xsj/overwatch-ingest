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

// AgentStatus represents an agent's operational status.
type AgentStatus string

const (
	StatusInitialized AgentStatus = "initialized" // Created but not active
	StatusActive      AgentStatus = "active"      // Ready to receive tasks
	StatusBusy        AgentStatus = "busy"        // Processing a task
	StatusDeactivated AgentStatus = "deactivated" // Stopped
)

// String returns the string representation.
func (s AgentStatus) String() string {
	return string(s)
}

// IsValid checks if the status is valid.
func (s AgentStatus) IsValid() bool {
	switch s {
	case StatusInitialized, StatusActive, StatusBusy, StatusDeactivated:
		return true
	}
	return false
}

// ParseAgentStatus parses a string into an AgentStatus.
func ParseAgentStatus(s string) (AgentStatus, error) {
	status := AgentStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid agent status: %s", s)
	}
	return status, nil
}

// Provider represents the LLM provider.
type Provider string

const (
	ProviderAnthropic Provider = "anthropic"
	ProviderOpenAI    Provider = "openai"
	ProviderLocal     Provider = "local" // For testing
)

// String returns the string representation.
func (p Provider) String() string {
	return string(p)
}

// IsValid checks if the provider is valid.
func (p Provider) IsValid() bool {
	switch p {
	case ProviderAnthropic, ProviderOpenAI, ProviderLocal:
		return true
	}
	return false
}

// ParseProvider parses a string into a Provider.
func ParseProvider(s string) (Provider, error) {
	p := Provider(s)
	if !p.IsValid() {
		return "", fmt.Errorf("invalid provider: %s", s)
	}
	return p, nil
}

// Model represents the LLM model.
type Model string

const (
	ModelClaude35Sonnet Model = "claude-3-5-sonnet-20241022"
	ModelClaude3Opus    Model = "claude-3-opus-20240229"
	ModelGPT4           Model = "gpt-4"
	ModelGPT4Turbo      Model = "gpt-4-turbo"
	ModelGPT35Turbo     Model = "gpt-3.5-turbo"
)

// String returns the string representation.
func (m Model) String() string {
	return string(m)
}

// IsValid checks if the model is valid.
func (m Model) IsValid() bool {
	switch m {
	case ModelClaude35Sonnet, ModelClaude3Opus, ModelGPT4, ModelGPT4Turbo, ModelGPT35Turbo:
		return true
	}
	return false
}

// ParseModel parses a string into a Model.
func ParseModel(s string) (Model, error) {
	m := Model(s)
	if !m.IsValid() {
		return "", fmt.Errorf("invalid model: %s", s)
	}
	return m, nil
}