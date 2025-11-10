// services/incidents/internal/domain/incident/event.go
package incident

import (
	"time"
)

// Event represents a domain event that has occurred.
// Events are immutable facts about what happened in the past.
type Event interface {
	// AggregateID returns the ID of the aggregate this event belongs to
	AggregateID() string

	// EventType returns the type of event (e.g., "IncidentReported")
	EventType() string

	// OccurredAt returns when the event occurred
	OccurredAt() time.Time

	// Version returns the event schema version
	Version() int

	// SequenceNumber returns the position in the aggregate's event stream
	SequenceNumber() int64
}

// BaseEvent provides common event fields.
// Embed this in specific events to avoid repetition.
type BaseEvent struct {
	aggregateID    string
	eventType      string
	occurredAt     time.Time
	version        int
	sequenceNumber int64
}

func (e BaseEvent) AggregateID() string    { return e.aggregateID }
func (e BaseEvent) EventType() string      { return e.eventType }
func (e BaseEvent) OccurredAt() time.Time  { return e.occurredAt }
func (e BaseEvent) Version() int           { return e.version }
func (e BaseEvent) SequenceNumber() int64  { return e.sequenceNumber }

// NewBaseEvent creates a new base event.
func NewBaseEvent(aggregateID, eventType string, sequenceNumber int64) BaseEvent {
	return BaseEvent{
		aggregateID:    aggregateID,
		eventType:      eventType,
		occurredAt:     time.Now(),
		version:        1, // Start with version 1
		sequenceNumber: sequenceNumber,
	}
}