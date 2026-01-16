package event

import (
	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

const (
	EventTypeRecordReceived    = "ingest.record.received"
	EventTypeRecordValidated   = "ingest.record.validated"
	EventTypeRecordAccepted    = "ingest.record.accepted"
	EventTypeRecordRejected    = "ingest.record.rejected"
	EventTypeRecordQuarantined = "ingest.record.quarantined"

	EventTypeBatchReceived  = "ingest.batch.received"
	EventTypeBatchCompleted = "ingest.batch.completed"

	EventTypeQuarantineResolved = "ingest.quarantine.resolved"
	EventTypeQuarantineExpired  = "ingest.quarantine.expired"

	EventTypeSignatureVerified = "ingest.signature.verified"
	EventTypeSignatureFailed   = "ingest.signature.failed"

	EventTypeSourceReliabilityUpdated = "ingest.source_reliability.updated"

	EventTypeIngestError = "ingest.error"
)

type Event interface {
	EventType() string
	OccurredAt() types.Timestamp
	TenantID() types.Optional[types.ID]
	SourceID() types.ID
}

type baseEvent struct {
	eventType  string
	occurredAt types.Timestamp
	tenantID   types.Optional[types.ID]
	sourceID   types.ID
}

func (e baseEvent) EventType() string                  { return e.eventType }
func (e baseEvent) OccurredAt() types.Timestamp        { return e.occurredAt }
func (e baseEvent) TenantID() types.Optional[types.ID] { return e.tenantID }
func (e baseEvent) SourceID() types.ID                 { return e.sourceID }

type RecordReceived struct {
	baseEvent
	IngestRecordID types.ID
	RawDataID      string
	SourceType     string
	SourceDID      types.Optional[string]
	CollectorDID   types.Optional[string]
}

func NewRecordReceived(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	ingestRecordID types.ID,
	rawDataID string,
	sourceType string,
	sourceDID types.Optional[string],
	collectorDID types.Optional[string],
) RecordReceived {
	return RecordReceived{
		baseEvent: baseEvent{
			eventType:  EventTypeRecordReceived,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		IngestRecordID: ingestRecordID,
		RawDataID:      rawDataID,
		SourceType:     sourceType,
		SourceDID:      sourceDID,
		CollectorDID:   collectorDID,
	}
}

type RecordValidated struct {
	baseEvent
	IngestRecordID  types.ID
	RawDataID       string
	SourceType      string
	Valid           bool
	AnomalyCount    int
	ConfidenceScore float64
}

func NewRecordValidated(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	ingestRecordID types.ID,
	rawDataID string,
	sourceType string,
	valid bool,
	anomalyCount int,
	confidenceScore float64,
) RecordValidated {
	return RecordValidated{
		baseEvent: baseEvent{
			eventType:  EventTypeRecordValidated,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		IngestRecordID:  ingestRecordID,
		RawDataID:       rawDataID,
		SourceType:      sourceType,
		Valid:           valid,
		AnomalyCount:    anomalyCount,
		ConfidenceScore: confidenceScore,
	}
}

type RecordAccepted struct {
	baseEvent
	IngestRecordID  types.ID
	RawDataID       string
	SourceType      string
	EntityType      string
	EntityID        string
	EventIDs        []string
	ConfidenceScore float64
	SourceSigner    *provenance.SignatureInfo
	CollectorSigner *provenance.SignatureInfo
	IngestSigner    *provenance.SignatureInfo
}

func NewRecordAccepted(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	ingestRecordID types.ID,
	rawDataID string,
	sourceType string,
	entityType string,
	entityID string,
	eventIDs []string,
	confidenceScore float64,
	sourceSigner *provenance.SignatureInfo,
	collectorSigner *provenance.SignatureInfo,
	ingestSigner *provenance.SignatureInfo,
) RecordAccepted {
	return RecordAccepted{
		baseEvent: baseEvent{
			eventType:  EventTypeRecordAccepted,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		IngestRecordID:  ingestRecordID,
		RawDataID:       rawDataID,
		SourceType:      sourceType,
		EntityType:      entityType,
		EntityID:        entityID,
		EventIDs:        eventIDs,
		ConfidenceScore: confidenceScore,
		SourceSigner:    sourceSigner,
		CollectorSigner: collectorSigner,
		IngestSigner:    ingestSigner,
	}
}

type RecordRejected struct {
	baseEvent
	IngestRecordID  types.ID
	RawDataID       string
	SourceType      string
	Reason          string
	Anomalies       []model.Anomaly
	ConfidenceScore float64
}

func NewRecordRejected(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	ingestRecordID types.ID,
	rawDataID string,
	sourceType string,
	reason string,
	anomalies []model.Anomaly,
	confidenceScore float64,
) RecordRejected {
	return RecordRejected{
		baseEvent: baseEvent{
			eventType:  EventTypeRecordRejected,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		IngestRecordID:  ingestRecordID,
		RawDataID:       rawDataID,
		SourceType:      sourceType,
		Reason:          reason,
		Anomalies:       anomalies,
		ConfidenceScore: confidenceScore,
	}
}

type RecordQuarantined struct {
	baseEvent
	IngestRecordID  types.ID
	RawDataID       string
	QuarantineID    types.ID
	SourceType      string
	Reason          model.QuarantineReason
	ReasonDetail    string
	Anomalies       []model.Anomaly
	ConfidenceScore float64
	ExpiresAt       types.Optional[types.Timestamp]
}

func NewRecordQuarantined(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	ingestRecordID types.ID,
	rawDataID string,
	quarantineID types.ID,
	sourceType string,
	reason model.QuarantineReason,
	reasonDetail string,
	anomalies []model.Anomaly,
	confidenceScore float64,
	expiresAt types.Optional[types.Timestamp],
) RecordQuarantined {
	return RecordQuarantined{
		baseEvent: baseEvent{
			eventType:  EventTypeRecordQuarantined,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		IngestRecordID:  ingestRecordID,
		RawDataID:       rawDataID,
		QuarantineID:    quarantineID,
		SourceType:      sourceType,
		Reason:          reason,
		ReasonDetail:    reasonDetail,
		Anomalies:       anomalies,
		ConfidenceScore: confidenceScore,
		ExpiresAt:       expiresAt,
	}
}

type BatchReceived struct {
	baseEvent
	BatchID      string
	SourceType   string
	RecordCount  int
	CollectorDID types.Optional[string]
}

func NewBatchReceived(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	batchID string,
	sourceType string,
	recordCount int,
	collectorDID types.Optional[string],
) BatchReceived {
	return BatchReceived{
		baseEvent: baseEvent{
			eventType:  EventTypeBatchReceived,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		BatchID:      batchID,
		SourceType:   sourceType,
		RecordCount:  recordCount,
		CollectorDID: collectorDID,
	}
}

type BatchCompleted struct {
	baseEvent
	BatchID            string
	SourceType         string
	TotalRecords       int
	AcceptedRecords    int
	RejectedRecords    int
	QuarantinedRecords int
	AverageConfidence  float64
	ProcessingTimeMs   int
}

func NewBatchCompleted(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	batchID string,
	sourceType string,
	totalRecords int,
	acceptedRecords int,
	rejectedRecords int,
	quarantinedRecords int,
	averageConfidence float64,
	processingTimeMs int,
) BatchCompleted {
	return BatchCompleted{
		baseEvent: baseEvent{
			eventType:  EventTypeBatchCompleted,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		BatchID:            batchID,
		SourceType:         sourceType,
		TotalRecords:       totalRecords,
		AcceptedRecords:    acceptedRecords,
		RejectedRecords:    rejectedRecords,
		QuarantinedRecords: quarantinedRecords,
		AverageConfidence:  averageConfidence,
		ProcessingTimeMs:   processingTimeMs,
	}
}

type QuarantineResolved struct {
	baseEvent
	QuarantineID   types.ID
	IngestRecordID types.ID
	SourceType     string
	Resolution     model.QuarantineResolution
	ResolvedBy     string
	ResolvedByDID  types.Optional[string]
	Notes          types.Optional[string]
	EntityType     types.Optional[string]
	EntityID       types.Optional[string]
	EventIDs       []string
}

func NewQuarantineResolved(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	quarantineID types.ID,
	ingestRecordID types.ID,
	sourceType string,
	resolution model.QuarantineResolution,
	resolvedBy string,
	resolvedByDID types.Optional[string],
	notes types.Optional[string],
	entityType types.Optional[string],
	entityID types.Optional[string],
	eventIDs []string,
) QuarantineResolved {
	return QuarantineResolved{
		baseEvent: baseEvent{
			eventType:  EventTypeQuarantineResolved,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		QuarantineID:   quarantineID,
		IngestRecordID: ingestRecordID,
		SourceType:     sourceType,
		Resolution:     resolution,
		ResolvedBy:     resolvedBy,
		ResolvedByDID:  resolvedByDID,
		Notes:          notes,
		EntityType:     entityType,
		EntityID:       entityID,
		EventIDs:       eventIDs,
	}
}

type QuarantineExpired struct {
	baseEvent
	QuarantineID   types.ID
	IngestRecordID types.ID
	SourceType     string
	QuarantinedAt  types.Timestamp
	ExpiredAt      types.Timestamp
}

func NewQuarantineExpired(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	quarantineID types.ID,
	ingestRecordID types.ID,
	sourceType string,
	quarantinedAt types.Timestamp,
) QuarantineExpired {
	return QuarantineExpired{
		baseEvent: baseEvent{
			eventType:  EventTypeQuarantineExpired,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		QuarantineID:   quarantineID,
		IngestRecordID: ingestRecordID,
		SourceType:     sourceType,
		QuarantinedAt:  quarantinedAt,
		ExpiredAt:      types.Now(),
	}
}

type SignatureVerified struct {
	baseEvent
	IngestRecordID types.ID
	RawDataID      string
	SignerDID      string
	SignerType     string
	SignedAt       types.Timestamp
}

func NewSignatureVerified(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	ingestRecordID types.ID,
	rawDataID string,
	signerDID string,
	signerType string,
	signedAt types.Timestamp,
) SignatureVerified {
	return SignatureVerified{
		baseEvent: baseEvent{
			eventType:  EventTypeSignatureVerified,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		IngestRecordID: ingestRecordID,
		RawDataID:      rawDataID,
		SignerDID:      signerDID,
		SignerType:     signerType,
		SignedAt:       signedAt,
	}
}

type SignatureFailed struct {
	baseEvent
	IngestRecordID types.ID
	RawDataID      string
	SignerDID      string
	SignerType     string
	Error          string
}

func NewSignatureFailed(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	ingestRecordID types.ID,
	rawDataID string,
	signerDID string,
	signerType string,
	err string,
) SignatureFailed {
	return SignatureFailed{
		baseEvent: baseEvent{
			eventType:  EventTypeSignatureFailed,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		IngestRecordID: ingestRecordID,
		RawDataID:      rawDataID,
		SignerDID:      signerDID,
		SignerType:     signerType,
		Error:          err,
	}
}

type SourceReliabilityUpdated struct {
	baseEvent
	PreviousScore   float64
	NewScore        float64
	TotalRecords    int64
	AcceptedRecords int64
	RejectedRecords int64
	Reason          string
}

func NewSourceReliabilityUpdated(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	previousScore float64,
	newScore float64,
	totalRecords int64,
	acceptedRecords int64,
	rejectedRecords int64,
	reason string,
) SourceReliabilityUpdated {
	return SourceReliabilityUpdated{
		baseEvent: baseEvent{
			eventType:  EventTypeSourceReliabilityUpdated,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		PreviousScore:   previousScore,
		NewScore:        newScore,
		TotalRecords:    totalRecords,
		AcceptedRecords: acceptedRecords,
		RejectedRecords: rejectedRecords,
		Reason:          reason,
	}
}

type IngestError struct {
	baseEvent
	ErrorID        string
	ErrorType      string
	Error          string
	IngestRecordID types.Optional[types.ID]
	RawDataID      types.Optional[string]
	Context        map[string]any
	Recoverable    bool
}

func NewIngestError(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	errorID string,
	errorType string,
	err string,
	ingestRecordID types.Optional[types.ID],
	rawDataID types.Optional[string],
	context map[string]any,
	recoverable bool,
) IngestError {
	return IngestError{
		baseEvent: baseEvent{
			eventType:  EventTypeIngestError,
			occurredAt: types.Now(),
			tenantID:   tenantID,
			sourceID:   sourceID,
		},
		ErrorID:        errorID,
		ErrorType:      errorType,
		Error:          err,
		IngestRecordID: ingestRecordID,
		RawDataID:      rawDataID,
		Context:        context,
		Recoverable:    recoverable,
	}
}
