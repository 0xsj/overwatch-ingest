package mocks

import (
	"context"
	"sync"

	"github.com/0xsj/overwatch-ingest/internal/domain/event"
)

// EventPublisher is a mock implementation of messaging.EventPublisher.
type EventPublisher struct {
	mu     sync.RWMutex
	Events []event.Event

	Calls struct {
		Publish                         int
		PublishAll                      int
		PublishRecordReceived           int
		PublishRecordValidated          int
		PublishRecordAccepted           int
		PublishRecordRejected           int
		PublishRecordQuarantined        int
		PublishBatchReceived            int
		PublishBatchCompleted           int
		PublishQuarantineResolved       int
		PublishQuarantineExpired        int
		PublishSignatureVerified        int
		PublishSignatureFailed          int
		PublishSourceReliabilityUpdated int
		PublishIngestError              int
	}

	Errors struct {
		Publish                         error
		PublishAll                      error
		PublishRecordReceived           error
		PublishRecordValidated          error
		PublishRecordAccepted           error
		PublishRecordRejected           error
		PublishRecordQuarantined        error
		PublishBatchReceived            error
		PublishBatchCompleted           error
		PublishQuarantineResolved       error
		PublishQuarantineExpired        error
		PublishSignatureVerified        error
		PublishSignatureFailed          error
		PublishSourceReliabilityUpdated error
		PublishIngestError              error
	}
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{}
}

func (p *EventPublisher) Publish(_ context.Context, evt event.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.Publish++
	if p.Errors.Publish != nil {
		return p.Errors.Publish
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishAll(_ context.Context, events []event.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishAll++
	if p.Errors.PublishAll != nil {
		return p.Errors.PublishAll
	}
	p.Events = append(p.Events, events...)
	return nil
}

func (p *EventPublisher) PublishRecordReceived(_ context.Context, evt event.RecordReceived) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishRecordReceived++
	if p.Errors.PublishRecordReceived != nil {
		return p.Errors.PublishRecordReceived
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishRecordValidated(_ context.Context, evt event.RecordValidated) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishRecordValidated++
	if p.Errors.PublishRecordValidated != nil {
		return p.Errors.PublishRecordValidated
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishRecordAccepted(_ context.Context, evt event.RecordAccepted) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishRecordAccepted++
	if p.Errors.PublishRecordAccepted != nil {
		return p.Errors.PublishRecordAccepted
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishRecordRejected(_ context.Context, evt event.RecordRejected) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishRecordRejected++
	if p.Errors.PublishRecordRejected != nil {
		return p.Errors.PublishRecordRejected
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishRecordQuarantined(_ context.Context, evt event.RecordQuarantined) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishRecordQuarantined++
	if p.Errors.PublishRecordQuarantined != nil {
		return p.Errors.PublishRecordQuarantined
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishBatchReceived(_ context.Context, evt event.BatchReceived) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishBatchReceived++
	if p.Errors.PublishBatchReceived != nil {
		return p.Errors.PublishBatchReceived
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishBatchCompleted(_ context.Context, evt event.BatchCompleted) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishBatchCompleted++
	if p.Errors.PublishBatchCompleted != nil {
		return p.Errors.PublishBatchCompleted
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishQuarantineResolved(_ context.Context, evt event.QuarantineResolved) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishQuarantineResolved++
	if p.Errors.PublishQuarantineResolved != nil {
		return p.Errors.PublishQuarantineResolved
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishQuarantineExpired(_ context.Context, evt event.QuarantineExpired) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishQuarantineExpired++
	if p.Errors.PublishQuarantineExpired != nil {
		return p.Errors.PublishQuarantineExpired
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishSignatureVerified(_ context.Context, evt event.SignatureVerified) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishSignatureVerified++
	if p.Errors.PublishSignatureVerified != nil {
		return p.Errors.PublishSignatureVerified
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishSignatureFailed(_ context.Context, evt event.SignatureFailed) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishSignatureFailed++
	if p.Errors.PublishSignatureFailed != nil {
		return p.Errors.PublishSignatureFailed
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishSourceReliabilityUpdated(_ context.Context, evt event.SourceReliabilityUpdated) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishSourceReliabilityUpdated++
	if p.Errors.PublishSourceReliabilityUpdated != nil {
		return p.Errors.PublishSourceReliabilityUpdated
	}
	p.Events = append(p.Events, evt)
	return nil
}

func (p *EventPublisher) PublishIngestError(_ context.Context, evt event.IngestError) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Calls.PublishIngestError++
	if p.Errors.PublishIngestError != nil {
		return p.Errors.PublishIngestError
	}
	p.Events = append(p.Events, evt)
	return nil
}

// PublishedEvents returns a copy of all published events.
func (p *EventPublisher) PublishedEvents() []event.Event {
	p.mu.RLock()
	defer p.mu.RUnlock()
	cp := make([]event.Event, len(p.Events))
	copy(cp, p.Events)
	return cp
}

// LastEvent returns the most recently published event, or nil if none.
func (p *EventPublisher) LastEvent() event.Event {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.Events) == 0 {
		return nil
	}
	return p.Events[len(p.Events)-1]
}

// EventCount returns the total number of published events.
func (p *EventPublisher) EventCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.Events)
}

// Reset clears all published events and counters.
func (p *EventPublisher) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Events = nil
	p.Calls = struct {
		Publish                         int
		PublishAll                      int
		PublishRecordReceived           int
		PublishRecordValidated          int
		PublishRecordAccepted           int
		PublishRecordRejected           int
		PublishRecordQuarantined        int
		PublishBatchReceived            int
		PublishBatchCompleted           int
		PublishQuarantineResolved       int
		PublishQuarantineExpired        int
		PublishSignatureVerified        int
		PublishSignatureFailed          int
		PublishSourceReliabilityUpdated int
		PublishIngestError              int
	}{}
}
