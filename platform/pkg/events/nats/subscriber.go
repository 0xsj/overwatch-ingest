// platform/pkg/events/nats/subscriber.go
package nats

import (
	"context"
	"fmt"

	"github.com/0xsj/scout/platform/pkg/events"
	"github.com/nats-io/nats.go"
)

// Subscriber implements events.Subscriber for NATS.
type Subscriber struct {
	conn *nats.Conn
}

// NewSubscriber creates a new NATS subscriber.
func NewSubscriber(cfg events.Config) (*Subscriber, error) {
	opts := []nats.Option{
		nats.MaxReconnects(cfg.MaxReconnects()),
		nats.ReconnectWait(cfg.ReconnectWait()),
	}

	conn, err := nats.Connect(cfg.URL(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &Subscriber{conn: conn}, nil
}

// Subscribe subscribes to a subject.
func (s *Subscriber) Subscribe(ctx context.Context, subject string, handler events.MessageHandler) (events.Subscription, error) {
	sub, err := s.conn.Subscribe(subject, func(msg *nats.Msg) {
		eventMsg := events.NewMessage(subject, msg.Data)
		if err := handler(ctx, eventMsg); err != nil {
			// Log error but don't stop processing
			// TODO: Add proper logging when logger is available
			fmt.Printf("error handling message on subject %s: %v\n", subject, err)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}

	return &subscription{sub: sub}, nil
}

// SubscribeQueue subscribes to a subject with queue group (load balancing).
func (s *Subscriber) SubscribeQueue(ctx context.Context, subject, queue string, handler events.MessageHandler) (events.Subscription, error) {
	sub, err := s.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		eventMsg := events.NewMessage(subject, msg.Data)
		if err := handler(ctx, eventMsg); err != nil {
			fmt.Printf("error handling message on subject %s: %v\n", subject, err)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to queue subscribe to subject %s: %w", subject, err)
	}

	return &subscription{sub: sub}, nil
}

// Close closes the subscriber connection.
func (s *Subscriber) Close() error {
	s.conn.Close()
	return nil
}

// Ensure Subscriber implements events.Subscriber interface
var _ events.Subscriber = (*Subscriber)(nil)

// subscription wraps a NATS subscription.
type subscription struct {
	sub *nats.Subscription
}

// Unsubscribe stops receiving messages.
func (s *subscription) Unsubscribe() error {
	return s.sub.Unsubscribe()
}

// Subject returns the subject being subscribed to.
func (s *subscription) Subject() string {
	return s.sub.Subject
}

// Ensure subscription implements events.Subscription interface
var _ events.Subscription = (*subscription)(nil)