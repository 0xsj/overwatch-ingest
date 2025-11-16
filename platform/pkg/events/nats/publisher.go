// platform/pkg/events/nats/publisher.go
package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0xsj/scout/platform/pkg/events"
	"github.com/nats-io/nats.go"
)

// Publisher implements events.Publisher for NATS.
type Publisher struct {
	conn *nats.Conn
}

// NewPublisher creates a new NATS publisher.
func NewPublisher(cfg events.Config) (*Publisher, error) {
	opts := []nats.Option{
		nats.MaxReconnects(cfg.MaxReconnects()),
		nats.ReconnectWait(cfg.ReconnectWait()),
	}

	conn, err := nats.Connect(cfg.URL(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &Publisher{conn: conn}, nil
}

// Publish publishes raw bytes to a subject.
func (p *Publisher) Publish(ctx context.Context, subject string, data []byte) error {
	if err := p.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish to NATS subject %s: %w", subject, err)
	}
	return nil
}

// PublishJSON publishes a JSON-serialized event to a subject.
func (p *Publisher) PublishJSON(ctx context.Context, subject string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.Publish(ctx, subject, data)
}

// Close closes the NATS connection.
func (p *Publisher) Close() error {
	p.conn.Close()
	return nil
}

// Ensure Publisher implements events.Publisher interface
var _ events.Publisher = (*Publisher)(nil)