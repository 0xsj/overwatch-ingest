services/agents/
├── api/v1/                           # gRPC contract
│
├── cmd/
│   ├── server/                       # Command/Query gRPC server
│   │   └── main.go
│   ├── projector/                    # Read model projector (separate process)
│   │   └── main.go
│   └── worker/                       # Event handlers/sagas
│       └── main.go
│
├── internal/
│   ├── domain/                       # Write Model (Event Sourced)
│   │   └── agent/
│   │       ├── aggregate.go          # Aggregate Root (holds uncommitted events)
│   │       ├── commands.go           # Command definitions
│   │       ├── events.go             # Domain events (THE source of truth)
│   │       ├── errors.go             # Domain errors
│   │       ├── value_objects.go      # AgentID, Status, AgentType
│   │       └── event_store.go        # Event store interface
│   │
│   ├── application/
│   │   ├── command/                  # Command Handlers (Write Side)
│   │   │   ├── handler.go            # Base command handler
│   │   │   ├── create_agent.go       # CreateAgentHandler
│   │   │   ├── deploy_agent.go       # DeployAgentHandler
│   │   │   ├── deactivate_agent.go   # DeactivateAgentHandler
│   │   │   └── update_location.go    # UpdateAgentLocationHandler
│   │   │
│   │   ├── query/                    # Query Handlers (Read Side)
│   │   │   ├── handler.go            # Base query handler
│   │   │   ├── get_agent.go          # GetAgentHandler
│   │   │   ├── list_agents.go        # ListAgentsHandler
│   │   │   └── get_agent_status.go   # GetAgentStatusHandler
│   │   │
│   │   ├── projection/               # Read Model Projections
│   │   │   ├── projector.go          # Base projector
│   │   │   ├── agent_list.go         # Agent list projection
│   │   │   ├── agent_status.go       # Agent status projection
│   │   │   └── agent_location.go     # Agent location projection
│   │   │
│   │   └── dto/                      # Data Transfer Objects
│   │       ├── command_dto.go
│   │       └── query_dto.go
│   │
│   ├── infrastructure/
│   │   ├── eventstore/               # Event Store Implementation
│   │   │   ├── postgres/
│   │   │   │   ├── event_store.go    # Store events in Postgres
│   │   │   │   ├── snapshots.go      # Snapshot storage
│   │   │   │   └── migrations/
│   │   │   │       ├── 001_events.sql
│   │   │   │       └── 002_snapshots.sql
│   │   │   └── stream.go             # Event stream interface
│   │   │
│   │   ├── readmodel/                # Read Model Storage
│   │   │   ├── postgres/
│   │   │   │   ├── agent_repo.go     # Query-optimized tables
│   │   │   │   └── migrations/
│   │   │   └── redis/
│   │   │       └── agent_cache.go    # Hot data cache
│   │   │
│   │   ├── messaging/                # Event Bus
│   │   │   ├── nats/
│   │   │   │   ├── publisher.go      # Publish events
│   │   │   │   └── subscriber.go     # Subscribe to events
│   │   │   └── rabbitmq/
│   │   │       └── publisher.go
│   │   │
│   │   └── serialization/            # Event serialization
│   │       ├── json.go
│   │       └── protobuf.go
│   │
│   └── interfaces/
│       ├── grpc/
│       │   ├── server.go
│       │   ├── command_handler.go    # Routes commands to command handlers
│       │   ├── query_handler.go      # Routes queries to query handlers
│       │   └── mappers.go            # Proto ↔ DTO
│       │
│       └── events/                   # Event Handlers (Process Managers/Sagas)
│           ├── handler.go
│           └── agent_events.go       # Reacts to agent events
│
├── config.go
├── go.mod
└── README.md