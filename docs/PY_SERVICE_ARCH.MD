services/tools/                           # or services/analytics/
├── api/
│   └── v1/
│       ├── tools.proto                   # gRPC contract (source)
│       ├── tools_pb2.py                  # Generated messages
│       └── tools_pb2_grpc.py             # Generated service
│
├── app/
│   ├── __init__.py
│   ├── main.py                           # Command/Query gRPC server
│   ├── projector.py                      # Read model projector
│   └── worker.py                         # Background workers
│
├── internal/
│   ├── __init__.py
│   │
│   ├── domain/                           # Domain Layer (Event Sourced)
│   │   ├── __init__.py
│   │   └── tool/
│   │       ├── __init__.py
│   │       ├── event.py                  # Event base class
│   │       ├── events.py                 # Domain events (ToolCreated, etc)
│   │       ├── aggregate.py              # Tool aggregate root
│   │       ├── value_objects.py          # ToolID, Status, ToolType
│   │       ├── errors.py                 # Domain errors
│   │       └── event_store.py            # Event store protocol
│   │
│   ├── application/                      # Application Layer
│   │   ├── __init__.py
│   │   │
│   │   ├── command/                      # Write Side (CQRS Commands)
│   │   │   ├── __init__.py
│   │   │   ├── handler.py                # Base command handler
│   │   │   ├── create_tool.py            # CreateToolHandler
│   │   │   ├── process_data.py           # ProcessDataHandler
│   │   │   └── deactivate_tool.py        # DeactivateToolHandler
│   │   │
│   │   ├── query/                        # Read Side (CQRS Queries)
│   │   │   ├── __init__.py
│   │   │   ├── handler.py                # Base query handler
│   │   │   ├── get_tool.py               # GetToolHandler
│   │   │   └── list_tools.py             # ListToolsHandler
│   │   │
│   │   ├── projection/                   # Read Model Projections
│   │   │   ├── __init__.py
│   │   │   ├── projector.py              # Base projector protocol
│   │   │   ├── tool_list.py              # Tool list projection
│   │   │   └── tool_status.py            # Tool status projection
│   │   │
│   │   └── dto/                          # Data Transfer Objects
│   │       ├── __init__.py
│   │       ├── command_dto.py
│   │       └── query_dto.py
│   │
│   ├── infrastructure/                   # Infrastructure Layer
│   │   ├── __init__.py
│   │   │
│   │   ├── eventstore/                   # Event Store Implementation
│   │   │   ├── __init__.py
│   │   │   ├── event_store.py            # Event store protocol
│   │   │   ├── postgres/
│   │   │   │   ├── __init__.py
│   │   │   │   ├── event_store.py        # Postgres implementation
│   │   │   │   ├── snapshots.py          # Snapshot repository
│   │   │   │   ├── serializer.py         # Event serialization
│   │   │   │   └── migrations/
│   │   │   │       ├── 001_events.sql
│   │   │   │       ├── 002_snapshots.sql
│   │   │   │       └── 003_read_models.sql
│   │   │   └── stream.py                 # Event stream reader
│   │   │
│   │   ├── readmodel/                    # Read Model Storage
│   │   │   ├── __init__.py
│   │   │   ├── repository.py             # Repository protocol
│   │   │   ├── postgres/
│   │   │   │   ├── __init__.py
│   │   │   │   ├── tool_repository.py
│   │   │   │   └── queries.py
│   │   │   └── redis/
│   │   │       ├── __init__.py
│   │   │       └── cache.py
│   │   │
│   │   ├── messaging/                    # Event Bus
│   │   │   ├── __init__.py
│   │   │   ├── publisher.py              # Publisher protocol
│   │   │   ├── subscriber.py             # Subscriber protocol
│   │   │   ├── nats/
│   │   │   │   ├── __init__.py
│   │   │   │   ├── publisher.py
│   │   │   │   └── subscriber.py
│   │   │   └── rabbitmq/
│   │   │       ├── __init__.py
│   │   │       ├── publisher.py
│   │   │       └── subscriber.py
│   │   │
│   │   └── serialization/                # Event Serialization
│   │       ├── __init__.py
│   │       ├── serializer.py             # Serializer protocol
│   │       ├── json_serializer.py
│   │       └── protobuf_serializer.py
│   │
│   └── interfaces/                       # Interface Adapters Layer
│       ├── __init__.py
│       │
│       ├── grpc/
│       │   ├── __init__.py
│       │   ├── server.py                 # gRPC server setup
│       │   ├── command_handler.py        # Command routing
│       │   ├── query_handler.py          # Query routing
│       │   └── mappers.py                # Conversions
│       │
│       └── events/                       # Event Handlers
│           ├── __init__.py
│           ├── handler.py
│           └── tool_saga.py
│
├── config.py                             # Service configuration
├── pyproject.toml
├── poetry.lock
└── README.md