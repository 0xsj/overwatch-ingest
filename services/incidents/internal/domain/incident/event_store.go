// services/incidents/internal/domain/incident/event_store.go
package incident

import (
	"context"
)

// EventStore defines the interface for persisting and retrieving events.
// This is part of the domain layer but will be implemented in the infrastructure layer.
type EventStore interface {
	// Save persists uncommitted events for an aggregate.
	// Events are saved atomically - either all succeed or all fail.
	// expectedVersion is used for optimistic concurrency control.
	// After successful save, events are typically published to the event bus.
	Save(ctx context.Context, aggregateID string, events []Event, expectedVersion int64) error

	// Load retrieves all events for an aggregate by its ID.
	// Events are returned in chronological order (oldest first).
	// Returns empty slice if aggregate doesn't exist.
	Load(ctx context.Context, aggregateID string) ([]Event, error)

	// LoadFromVersion retrieves events for an aggregate starting from a specific version.
	// Useful for incremental updates or catching up projections.
	// Returns empty slice if no events exist from that version.
	LoadFromVersion(ctx context.Context, aggregateID string, fromVersion int64) ([]Event, error)

	// Exists checks if an aggregate with the given ID has any events.
	// Useful for checking if an aggregate exists before attempting operations.
	Exists(ctx context.Context, aggregateID string) (bool, error)
}

// SnapshotStore defines the interface for storing and retrieving aggregate snapshots.
// Snapshots are optional optimizations to avoid replaying all events.
type SnapshotStore interface {
	// SaveSnapshot saves a snapshot of the aggregate state.
	// The snapshot represents the state at a specific version.
	SaveSnapshot(ctx context.Context, aggregateID string, version int64, snapshot interface{}) error

	// LoadSnapshot retrieves the most recent snapshot for an aggregate.
	// Returns nil if no snapshot exists.
	LoadSnapshot(ctx context.Context, aggregateID string) (*Snapshot, error)
}

// Snapshot represents a saved aggregate state at a specific version.
type Snapshot struct {
	AggregateID string
	Version     int64
	State       interface{} // The serialized aggregate state
}

// EventStream defines the interface for subscribing to event streams.
// Used by projectors and process managers to react to events.
type EventStream interface {
	// Subscribe subscribes to events of a specific aggregate type.
	// The handler is called for each event in order.
	Subscribe(ctx context.Context, aggregateType string, handler EventHandler) error

	// SubscribeFromVersion subscribes to events starting from a specific version.
	// Useful for catching up a projection that's behind.
	SubscribeFromVersion(ctx context.Context, aggregateType string, fromVersion int64, handler EventHandler) error
}

// EventHandler is a function that handles an event.
// Returns error if the event cannot be processed (will retry based on implementation).
type EventHandler func(ctx context.Context, event Event) error