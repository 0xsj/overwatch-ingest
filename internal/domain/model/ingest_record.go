package model

import (
	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
)

type IngestRecord struct {
	id         types.ID
	tenantID   types.Optional[types.ID]
	sourceID   types.ID
	sourceType string
	rawDataID  string

	status IngestStatus

	validation ValidationResult
	confidence ConfidenceScore

	entityType types.Optional[string]
	entityID   types.Optional[string]
	eventIDs   []string

	rejectionReason types.Optional[string]
	quarantineID    types.Optional[types.ID]

	sourceSigner    *provenance.SignatureInfo
	collectorSigner *provenance.SignatureInfo
	ingestSigner    *provenance.SignatureInfo

	sourceSignatureVerified    types.Optional[bool]
	collectorSignatureVerified types.Optional[bool]

	receivedAt  types.Timestamp
	processedAt types.Timestamp
}

func NewIngestRecord(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	sourceType string,
	rawDataID string,
	receivedAt types.Timestamp,
) (*IngestRecord, error) {
	if sourceID.IsEmpty() {
		return nil, domainerror.ErrSourceIDRequired
	}
	if sourceType == "" {
		return nil, domainerror.ErrSourceTypeRequired
	}
	if rawDataID == "" {
		return nil, domainerror.ErrRawDataIDRequired
	}

	return &IngestRecord{
		id:                         types.NewID(),
		tenantID:                   tenantID,
		sourceID:                   sourceID,
		sourceType:                 sourceType,
		rawDataID:                  rawDataID,
		status:                     IngestStatusPending,
		entityType:                 types.None[string](),
		entityID:                   types.None[string](),
		eventIDs:                   nil,
		rejectionReason:            types.None[string](),
		quarantineID:               types.None[types.ID](),
		sourceSigner:               nil,
		collectorSigner:            nil,
		ingestSigner:               nil,
		sourceSignatureVerified:    types.None[bool](),
		collectorSignatureVerified: types.None[bool](),
		receivedAt:                 receivedAt,
		processedAt:                types.Now(),
	}, nil
}

func ReconstructIngestRecord(
	id types.ID,
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	sourceType string,
	rawDataID string,
	status IngestStatus,
	validation ValidationResult,
	confidence ConfidenceScore,
	entityType types.Optional[string],
	entityID types.Optional[string],
	eventIDs []string,
	rejectionReason types.Optional[string],
	quarantineID types.Optional[types.ID],
	sourceSigner *provenance.SignatureInfo,
	collectorSigner *provenance.SignatureInfo,
	ingestSigner *provenance.SignatureInfo,
	sourceSignatureVerified types.Optional[bool],
	collectorSignatureVerified types.Optional[bool],
	receivedAt types.Timestamp,
	processedAt types.Timestamp,
) *IngestRecord {
	return &IngestRecord{
		id:                         id,
		tenantID:                   tenantID,
		sourceID:                   sourceID,
		sourceType:                 sourceType,
		rawDataID:                  rawDataID,
		status:                     status,
		validation:                 validation,
		confidence:                 confidence,
		entityType:                 entityType,
		entityID:                   entityID,
		eventIDs:                   eventIDs,
		rejectionReason:            rejectionReason,
		quarantineID:               quarantineID,
		sourceSigner:               sourceSigner,
		collectorSigner:            collectorSigner,
		ingestSigner:               ingestSigner,
		sourceSignatureVerified:    sourceSignatureVerified,
		collectorSignatureVerified: collectorSignatureVerified,
		receivedAt:                 receivedAt,
		processedAt:                processedAt,
	}
}

func (r *IngestRecord) ID() types.ID                               { return r.id }
func (r *IngestRecord) TenantID() types.Optional[types.ID]         { return r.tenantID }
func (r *IngestRecord) SourceID() types.ID                         { return r.sourceID }
func (r *IngestRecord) SourceType() string                         { return r.sourceType }
func (r *IngestRecord) RawDataID() string                          { return r.rawDataID }
func (r *IngestRecord) Status() IngestStatus                       { return r.status }
func (r *IngestRecord) Validation() ValidationResult               { return r.validation }
func (r *IngestRecord) Confidence() ConfidenceScore                { return r.confidence }
func (r *IngestRecord) EntityType() types.Optional[string]         { return r.entityType }
func (r *IngestRecord) EntityID() types.Optional[string]           { return r.entityID }
func (r *IngestRecord) EventIDs() []string                         { return r.eventIDs }
func (r *IngestRecord) RejectionReason() types.Optional[string]    { return r.rejectionReason }
func (r *IngestRecord) QuarantineID() types.Optional[types.ID]     { return r.quarantineID }
func (r *IngestRecord) SourceSigner() *provenance.SignatureInfo    { return r.sourceSigner }
func (r *IngestRecord) CollectorSigner() *provenance.SignatureInfo { return r.collectorSigner }
func (r *IngestRecord) IngestSigner() *provenance.SignatureInfo    { return r.ingestSigner }
func (r *IngestRecord) SourceSignatureVerified() types.Optional[bool] {
	return r.sourceSignatureVerified
}

func (r *IngestRecord) CollectorSignatureVerified() types.Optional[bool] {
	return r.collectorSignatureVerified
}
func (r *IngestRecord) ReceivedAt() types.Timestamp  { return r.receivedAt }
func (r *IngestRecord) ProcessedAt() types.Timestamp { return r.processedAt }

func (r *IngestRecord) SetValidation(validation ValidationResult) {
	r.validation = validation
	r.processedAt = types.Now()
}

func (r *IngestRecord) SetConfidence(confidence ConfidenceScore) {
	r.confidence = confidence
	r.processedAt = types.Now()
}

func (r *IngestRecord) SetSourceSigner(signer *provenance.SignatureInfo, verified bool) {
	r.sourceSigner = signer
	r.sourceSignatureVerified = types.Some(verified)
}

func (r *IngestRecord) SetCollectorSigner(signer *provenance.SignatureInfo, verified bool) {
	r.collectorSigner = signer
	r.collectorSignatureVerified = types.Some(verified)
}

func (r *IngestRecord) SetIngestSigner(signer *provenance.SignatureInfo) {
	r.ingestSigner = signer
}

func (r *IngestRecord) Accept(entityType, entityID string, eventIDs []string) error {
	if !r.status.CanTransitionTo(IngestStatusAccepted) {
		return domainerror.ErrInvalidStatusTransition
	}

	r.status = IngestStatusAccepted
	r.entityType = types.Some(entityType)
	r.entityID = types.Some(entityID)
	r.eventIDs = eventIDs
	r.processedAt = types.Now()

	return nil
}

func (r *IngestRecord) Reject(reason string) error {
	if !r.status.CanTransitionTo(IngestStatusRejected) {
		return domainerror.ErrInvalidStatusTransition
	}

	r.status = IngestStatusRejected
	r.rejectionReason = types.Some(reason)
	r.processedAt = types.Now()

	return nil
}

func (r *IngestRecord) Quarantine(quarantineID types.ID) error {
	if !r.status.CanTransitionTo(IngestStatusQuarantined) {
		return domainerror.ErrInvalidStatusTransition
	}

	r.status = IngestStatusQuarantined
	r.quarantineID = types.Some(quarantineID)
	r.processedAt = types.Now()

	return nil
}

func (r *IngestRecord) ResolveFromQuarantine(resolution QuarantineResolution, entityType, entityID string, eventIDs []string) error {
	if !r.status.IsQuarantined() {
		return domainerror.ErrRecordNotQuarantined
	}

	targetStatus := resolution.ToIngestStatus()
	if !r.status.CanTransitionTo(targetStatus) {
		return domainerror.ErrInvalidStatusTransition
	}

	r.status = targetStatus

	if resolution.IsAccepted() {
		r.entityType = types.Some(entityType)
		r.entityID = types.Some(entityID)
		r.eventIDs = eventIDs
	}

	r.processedAt = types.Now()

	return nil
}

func (r *IngestRecord) IsPending() bool     { return r.status.IsPending() }
func (r *IngestRecord) IsAccepted() bool    { return r.status.IsAccepted() }
func (r *IngestRecord) IsRejected() bool    { return r.status.IsRejected() }
func (r *IngestRecord) IsQuarantined() bool { return r.status.IsQuarantined() }
func (r *IngestRecord) IsTerminal() bool    { return r.status.IsTerminal() }

func (r *IngestRecord) HasEntity() bool {
	return r.entityType.IsPresent() && r.entityID.IsPresent()
}

func (r *IngestRecord) HasEvents() bool {
	return len(r.eventIDs) > 0
}

func (r *IngestRecord) HasSourceSigner() bool {
	return r.sourceSigner != nil && r.sourceSigner.IsValid()
}

func (r *IngestRecord) HasCollectorSigner() bool {
	return r.collectorSigner != nil && r.collectorSigner.IsValid()
}

func (r *IngestRecord) HasIngestSigner() bool {
	return r.ingestSigner != nil && r.ingestSigner.IsValid()
}

func (r *IngestRecord) HasValidProvenance() bool {
	collectorVerified := r.collectorSignatureVerified.OrElse(false)

	if r.HasSourceSigner() {
		sourceVerified := r.sourceSignatureVerified.OrElse(false)
		return collectorVerified && sourceVerified
	}

	return collectorVerified
}

func (r *IngestRecord) BelongsToTenant(tenantID types.ID) bool {
	if r.tenantID.IsEmpty() {
		return false
	}
	return r.tenantID.MustGet() == tenantID
}
