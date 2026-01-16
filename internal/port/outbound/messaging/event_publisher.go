package messaging

import (
	"context"

	"github.com/0xsj/overwatch-ingest/internal/domain/event"
)

type EventPublisher interface {
	Publish(ctx context.Context, evt event.Event) error
	PublishAll(ctx context.Context, events []event.Event) error

	PublishRecordReceived(ctx context.Context, evt event.RecordReceived) error
	PublishRecordValidated(ctx context.Context, evt event.RecordValidated) error
	PublishRecordAccepted(ctx context.Context, evt event.RecordAccepted) error
	PublishRecordRejected(ctx context.Context, evt event.RecordRejected) error
	PublishRecordQuarantined(ctx context.Context, evt event.RecordQuarantined) error
	PublishBatchReceived(ctx context.Context, evt event.BatchReceived) error
	PublishBatchCompleted(ctx context.Context, evt event.BatchCompleted) error
	PublishQuarantineResolved(ctx context.Context, evt event.QuarantineResolved) error
	PublishQuarantineExpired(ctx context.Context, evt event.QuarantineExpired) error
	PublishSignatureVerified(ctx context.Context, evt event.SignatureVerified) error
	PublishSignatureFailed(ctx context.Context, evt event.SignatureFailed) error
	PublishSourceReliabilityUpdated(ctx context.Context, evt event.SourceReliabilityUpdated) error
	PublishIngestError(ctx context.Context, evt event.IngestError) error
}

const (
	TopicRecordEvents      = "ingest.record"
	TopicBatchEvents       = "ingest.batch"
	TopicQuarantineEvents  = "ingest.quarantine"
	TopicSignatureEvents   = "ingest.signature"
	TopicReliabilityEvents = "ingest.reliability"
	TopicErrorEvents       = "ingest.error"
)

func TopicForEvent(evt event.Event) string {
	switch evt.EventType() {
	case event.EventTypeRecordReceived,
		event.EventTypeRecordValidated,
		event.EventTypeRecordAccepted,
		event.EventTypeRecordRejected,
		event.EventTypeRecordQuarantined:
		return TopicRecordEvents

	case event.EventTypeBatchReceived,
		event.EventTypeBatchCompleted:
		return TopicBatchEvents

	case event.EventTypeQuarantineResolved,
		event.EventTypeQuarantineExpired:
		return TopicQuarantineEvents

	case event.EventTypeSignatureVerified,
		event.EventTypeSignatureFailed:
		return TopicSignatureEvents

	case event.EventTypeSourceReliabilityUpdated:
		return TopicReliabilityEvents

	case event.EventTypeIngestError:
		return TopicErrorEvents

	default:
		return TopicRecordEvents
	}
}
