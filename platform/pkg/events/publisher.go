// platform/pkg/events/publisher.go
package events

import (
	"context"
)

// Publisher publishes events to an event bus.
// Implementations: NATS, Kafka, RabbitMQ
type Publisher interface {
	// Publish publishes raw bytes to a subject/topic.
	Publish(ctx context.Context, subject string, data []byte) error

	// PublishJSON publishes a JSON-serialized event to a subject/topic.
	PublishJSON(ctx context.Context, subject string, event interface{}) error

	// Close closes the publisher connection.
	Close() error
}