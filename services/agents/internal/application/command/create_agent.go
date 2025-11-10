package command

// // services/agents/internal/application/command/create_agent.go
// package command

// import (
// 	"context"

// 	"github.com/0xsj/scout/platform/pkg/errors"
// 	"github.com/0xsj/scout/services/agents/internal/application/dto"
// 	"github.com/0xsj/scout/services/agents/internal/domain/agent"
// )

// // CreateAgentHandler handles the CreateAgent command.
// type CreateAgentHandler struct {
// 	eventStore agent.EventStore
// }

// // NewCreateAgentHandler creates a new CreateAgentHandler.
// func NewCreateAgentHandler(eventStore agent.EventStore) *CreateAgentHandler {
// 	return &CreateAgentHandler{
// 		eventStore: eventStore,
// 	}
// }

// // Handle executes the CreateAgent command.
// func (h *CreateAgentHandler) Handle(ctx context.Context, cmd dto.CreateAgentCommand) (*dto.AgentDTO, error) {
// 	// 1. Generate or parse agent ID
// 	var agentID agent.AgentID
// 	var err error

// 	if cmd.AgentID == "" {
// 		agentID = agent.NewAgentID()
// 	} else {
// 		agentID, err = agent.ParseAgentID(cmd.AgentID)
// 		if err != nil {
// 			return nil, errors.Wrap(err, errors.ErrorTypeValidation, "INVALID_AGENT_ID", "invalid agent ID")
// 		}
// 	}

// 	// 2. Check if agent already exists
// 	exists, err := h.eventStore.Exists(ctx, agentID.String())
// 	if err != nil {
// 		return nil, errors.Wrap(err, errors.ErrorTypeDatabase, "DB_ERROR", "failed to check agent existence")
// 	}
// 	if exists {
// 		return nil, errors.AlreadyExists("agent", agentID.String())
// 	}

// 	// 3. Parse agent type
// 	agentType, err := agent.ParseAgentType(cmd.AgentType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 4. Create aggregate
// 	agg := agent.NewAgent(agentID)

// 	// 5. Execute domain command
// 	if err := agg.Create(cmd.Name, agentType, cmd.CreatedBy); err != nil {
// 		return nil, err
// 	}

// 	// 6. Get uncommitted events
// 	events := agg.UncommittedEvents()
// 	if len(events) == 0 {
// 		return nil, errors.Internal("no events generated")
// 	}

// 	// 7. Save events to event store
// 	if err := h.eventStore.Save(ctx, agentID.String(), events, 0); err != nil {
// 		return nil, errors.Wrap(err, errors.ErrorTypeDatabase, "EVENT_SAVE_FAILED", "failed to save events")
// 	}

// 	// 8. Clear uncommitted events
// 	agg.ClearUncommittedEvents()

// 	// 9. Convert aggregate to DTO for response
// 	return aggregateToDTO(agg), nil
// }

// // aggregateToDTO converts an Agent aggregate to AgentDTO.
// func aggregateToDTO(agg *agent.Agent) *dto.AgentDTO {
// 	var location *dto.LocationDTO
// 	if loc := agg.Location(); loc != nil {
// 		location = &dto.LocationDTO{
// 			Latitude:  loc.Latitude,
// 			Longitude: loc.Longitude,
// 		}
// 	}

// 	return &dto.AgentDTO{
// 		ID:            agg.ID().String(),
// 		Name:          agg.Name(),
// 		AgentType:     agg.Type().String(),
// 		Status:        agg.Status().String(),
// 		Location:      location,
// 		DeployedAt:    agg.DeployedAt(),
// 		DeactivatedAt: agg.DeactivatedAt(),
// 		Version:       agg.Version(),
// 	}
// }