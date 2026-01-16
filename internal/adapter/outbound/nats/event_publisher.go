package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/provenance/middleware"

	"github.com/0xsj/overwatch-ingest/internal/domain/event"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/messaging"
)

type eventPublisher struct {
	conn          *nats.Conn
	subjectPrefix string
	publisher     *middleware.Publisher
}

func NewEventPublisher(conn *nats.Conn, subjectPrefix string) messaging.EventPublisher {
	if subjectPrefix == "" {
		subjectPrefix = "overwatch"
	}
	return &eventPublisher{
		conn:          conn,
		subjectPrefix: subjectPrefix,
		publisher:     nil,
	}
}

func NewSignedEventPublisher(conn *nats.Conn, subjectPrefix string, identity provenance.Identity) (messaging.EventPublisher, error) {
	if subjectPrefix == "" {
		subjectPrefix = "overwatch"
	}

	publisher, err := middleware.NewPublisher(conn, identity)
	if err != nil {
		return nil, fmt.Errorf("failed to create signed publisher: %w", err)
	}

	return &eventPublisher{
		conn:          conn,
		subjectPrefix: subjectPrefix,
		publisher:     publisher,
	}, nil
}

func (p *eventPublisher) Publish(ctx context.Context, evt event.Event) error {
	subject := p.subjectForEvent(evt)

	if p.publisher != nil {
		return p.publishSigned(ctx, subject, evt)
	}
	return p.publishLegacy(ctx, subject, evt)
}

func (p *eventPublisher) PublishAll(ctx context.Context, events []event.Event) error {
	for _, evt := range events {
		if err := p.Publish(ctx, evt); err != nil {
			return err
		}
	}
	return nil
}

func (p *eventPublisher) publishSigned(ctx context.Context, subject string, evt event.Event) error {
	if err := p.publisher.Publish(subject, evt.EventType(), evt); err != nil {
		return fmt.Errorf("failed to publish signed event: %w", err)
	}
	return nil
}

func (p *eventPublisher) publishLegacy(ctx context.Context, subject string, evt event.Event) error {
	envelope := legacyEnvelope{
		EventType:  evt.EventType(),
		SourceID:   string(evt.SourceID()),
		OccurredAt: evt.OccurredAt().Time().Unix(),
		Payload:    evt,
	}

	if evt.TenantID().IsPresent() {
		tid := string(evt.TenantID().MustGet())
		envelope.TenantID = &tid
	}

	data, err := provenance.BuildPayload(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := p.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (p *eventPublisher) subjectForEvent(evt event.Event) string {
	topic := messaging.TopicForEvent(evt)
	return fmt.Sprintf("%s.%s", p.subjectPrefix, topic)
}

// Typed publish methods

func (p *eventPublisher) PublishRecordReceived(ctx context.Context, evt event.RecordReceived) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishRecordValidated(ctx context.Context, evt event.RecordValidated) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishRecordAccepted(ctx context.Context, evt event.RecordAccepted) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishRecordRejected(ctx context.Context, evt event.RecordRejected) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishRecordQuarantined(ctx context.Context, evt event.RecordQuarantined) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishBatchReceived(ctx context.Context, evt event.BatchReceived) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishBatchCompleted(ctx context.Context, evt event.BatchCompleted) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishQuarantineResolved(ctx context.Context, evt event.QuarantineResolved) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishQuarantineExpired(ctx context.Context, evt event.QuarantineExpired) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishSignatureVerified(ctx context.Context, evt event.SignatureVerified) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishSignatureFailed(ctx context.Context, evt event.SignatureFailed) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishSourceReliabilityUpdated(ctx context.Context, evt event.SourceReliabilityUpdated) error {
	return p.Publish(ctx, evt)
}

func (p *eventPublisher) PublishIngestError(ctx context.Context, evt event.IngestError) error {
	return p.Publish(ctx, evt)
}

type legacyEnvelope struct {
	EventType  string  `json:"event_type"`
	TenantID   *string `json:"tenant_id,omitempty"`
	SourceID   string  `json:"source_id"`
	OccurredAt int64   `json:"occurred_at"`
	Payload    any     `json:"payload"`
}
