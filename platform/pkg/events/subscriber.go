// platform/pkg/events/subscriber.go
package events

import (
	"context"
)

// MessageHandler is a function that handles incoming messages.
type MessageHandler func(ctx context.Context, msg *Message) error

// Subscriber subscribes to events from an event bus.
// Implementations: NATS, Kafka, RabbitMQ
type Subscriber interface {
	// Subscribe subscribes to a subject/topic with a handler.
	// Returns a subscription that can be used to unsubscribe.
	Subscribe(ctx context.Context, subject string, handler MessageHandler) (Subscription, error)

	// SubscribeQueue subscribes to a subject/topic as part of a queue group.
	// Multiple subscribers with the same queue group will load-balance messages.
	SubscribeQueue(ctx context.Context, subject, queue string, handler MessageHandler) (Subscription, error)

	// Close closes all subscriptions and the subscriber connection.
	Close() error
}

// Subscription represents an active subscription.
type Subscription interface {
	// Unsubscribe stops receiving messages.
	Unsubscribe() error

	// Subject returns the subject/topic being subscribed to.
	Subject() string
}