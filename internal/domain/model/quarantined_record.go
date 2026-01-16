package model

import (
	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
)

type QuarantinedRecord struct {
	id             types.ID
	tenantID       types.Optional[types.ID]
	sourceID       types.ID
	sourceType     string
	rawDataID      string
	ingestRecordID types.ID

	rawData map[string]any

	reason       QuarantineReason
	reasonDetail string
	anomalies    []Anomaly
	confidence   ConfidenceScore

	resolution      QuarantineResolution
	resolvedBy      types.Optional[string]
	resolvedByDID   types.Optional[string]
	resolutionNotes types.Optional[string]
	modifiedData    map[string]any

	quarantinedAt types.Timestamp
	expiresAt     types.Optional[types.Timestamp]
	resolvedAt    types.Optional[types.Timestamp]

	ingestSigner   *provenance.SignatureInfo
	resolverSigner *provenance.SignatureInfo
}

func NewQuarantinedRecord(
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	sourceType string,
	rawDataID string,
	ingestRecordID types.ID,
	rawData map[string]any,
	reason QuarantineReason,
	reasonDetail string,
	anomalies []Anomaly,
	confidence ConfidenceScore,
	expiresAt types.Optional[types.Timestamp],
	ingestSigner *provenance.SignatureInfo,
) (*QuarantinedRecord, error) {
	if sourceID.IsEmpty() {
		return nil, domainerror.ErrSourceIDRequired
	}
	if sourceType == "" {
		return nil, domainerror.ErrSourceTypeRequired
	}
	if rawDataID == "" {
		return nil, domainerror.ErrRawDataIDRequired
	}
	if ingestRecordID.IsEmpty() {
		return nil, domainerror.ErrRecordIDRequired
	}
	if rawData == nil {
		return nil, domainerror.ErrPayloadRequired
	}
	if !reason.IsValid() {
		return nil, domainerror.ErrResolutionInvalid
	}

	return &QuarantinedRecord{
		id:              types.NewID(),
		tenantID:        tenantID,
		sourceID:        sourceID,
		sourceType:      sourceType,
		rawDataID:       rawDataID,
		ingestRecordID:  ingestRecordID,
		rawData:         rawData,
		reason:          reason,
		reasonDetail:    reasonDetail,
		anomalies:       anomalies,
		confidence:      confidence,
		resolution:      QuarantineResolutionPending,
		resolvedBy:      types.None[string](),
		resolvedByDID:   types.None[string](),
		resolutionNotes: types.None[string](),
		modifiedData:    nil,
		quarantinedAt:   types.Now(),
		expiresAt:       expiresAt,
		resolvedAt:      types.None[types.Timestamp](),
		ingestSigner:    ingestSigner,
		resolverSigner:  nil,
	}, nil
}

func ReconstructQuarantinedRecord(
	id types.ID,
	tenantID types.Optional[types.ID],
	sourceID types.ID,
	sourceType string,
	rawDataID string,
	ingestRecordID types.ID,
	rawData map[string]any,
	reason QuarantineReason,
	reasonDetail string,
	anomalies []Anomaly,
	confidence ConfidenceScore,
	resolution QuarantineResolution,
	resolvedBy types.Optional[string],
	resolvedByDID types.Optional[string],
	resolutionNotes types.Optional[string],
	modifiedData map[string]any,
	quarantinedAt types.Timestamp,
	expiresAt types.Optional[types.Timestamp],
	resolvedAt types.Optional[types.Timestamp],
	ingestSigner *provenance.SignatureInfo,
	resolverSigner *provenance.SignatureInfo,
) *QuarantinedRecord {
	return &QuarantinedRecord{
		id:              id,
		tenantID:        tenantID,
		sourceID:        sourceID,
		sourceType:      sourceType,
		rawDataID:       rawDataID,
		ingestRecordID:  ingestRecordID,
		rawData:         rawData,
		reason:          reason,
		reasonDetail:    reasonDetail,
		anomalies:       anomalies,
		confidence:      confidence,
		resolution:      resolution,
		resolvedBy:      resolvedBy,
		resolvedByDID:   resolvedByDID,
		resolutionNotes: resolutionNotes,
		modifiedData:    modifiedData,
		quarantinedAt:   quarantinedAt,
		expiresAt:       expiresAt,
		resolvedAt:      resolvedAt,
		ingestSigner:    ingestSigner,
		resolverSigner:  resolverSigner,
	}
}

func (q *QuarantinedRecord) ID() types.ID                                { return q.id }
func (q *QuarantinedRecord) TenantID() types.Optional[types.ID]          { return q.tenantID }
func (q *QuarantinedRecord) SourceID() types.ID                          { return q.sourceID }
func (q *QuarantinedRecord) SourceType() string                          { return q.sourceType }
func (q *QuarantinedRecord) RawDataID() string                           { return q.rawDataID }
func (q *QuarantinedRecord) IngestRecordID() types.ID                    { return q.ingestRecordID }
func (q *QuarantinedRecord) RawData() map[string]any                     { return q.rawData }
func (q *QuarantinedRecord) Reason() QuarantineReason                    { return q.reason }
func (q *QuarantinedRecord) ReasonDetail() string                        { return q.reasonDetail }
func (q *QuarantinedRecord) Anomalies() []Anomaly                        { return q.anomalies }
func (q *QuarantinedRecord) Confidence() ConfidenceScore                 { return q.confidence }
func (q *QuarantinedRecord) Resolution() QuarantineResolution            { return q.resolution }
func (q *QuarantinedRecord) ResolvedBy() types.Optional[string]          { return q.resolvedBy }
func (q *QuarantinedRecord) ResolvedByDID() types.Optional[string]       { return q.resolvedByDID }
func (q *QuarantinedRecord) ResolutionNotes() types.Optional[string]     { return q.resolutionNotes }
func (q *QuarantinedRecord) ModifiedData() map[string]any                { return q.modifiedData }
func (q *QuarantinedRecord) QuarantinedAt() types.Timestamp              { return q.quarantinedAt }
func (q *QuarantinedRecord) ExpiresAt() types.Optional[types.Timestamp]  { return q.expiresAt }
func (q *QuarantinedRecord) ResolvedAt() types.Optional[types.Timestamp] { return q.resolvedAt }
func (q *QuarantinedRecord) IngestSigner() *provenance.SignatureInfo     { return q.ingestSigner }
func (q *QuarantinedRecord) ResolverSigner() *provenance.SignatureInfo   { return q.resolverSigner }

func (q *QuarantinedRecord) Approve(
	resolvedBy string,
	resolvedByDID string,
	notes string,
	resolverSigner *provenance.SignatureInfo,
) error {
	if !q.resolution.IsPending() {
		return domainerror.ErrQuarantineAlreadyResolved
	}

	q.resolution = QuarantineResolutionApproved
	q.resolvedBy = types.Some(resolvedBy)
	q.resolvedByDID = types.Some(resolvedByDID)
	if notes != "" {
		q.resolutionNotes = types.Some(notes)
	}
	q.resolvedAt = types.Some(types.Now())
	q.resolverSigner = resolverSigner

	return nil
}

func (q *QuarantinedRecord) ApproveWithModifications(
	resolvedBy string,
	resolvedByDID string,
	notes string,
	modifiedData map[string]any,
	resolverSigner *provenance.SignatureInfo,
) error {
	if !q.resolution.IsPending() {
		return domainerror.ErrQuarantineAlreadyResolved
	}
	if modifiedData == nil {
		return domainerror.ErrPayloadRequired
	}

	q.resolution = QuarantineResolutionModified
	q.resolvedBy = types.Some(resolvedBy)
	q.resolvedByDID = types.Some(resolvedByDID)
	if notes != "" {
		q.resolutionNotes = types.Some(notes)
	}
	q.modifiedData = modifiedData
	q.resolvedAt = types.Some(types.Now())
	q.resolverSigner = resolverSigner

	return nil
}

func (q *QuarantinedRecord) Reject(
	resolvedBy string,
	resolvedByDID string,
	notes string,
	resolverSigner *provenance.SignatureInfo,
) error {
	if !q.resolution.IsPending() {
		return domainerror.ErrQuarantineAlreadyResolved
	}

	q.resolution = QuarantineResolutionRejected
	q.resolvedBy = types.Some(resolvedBy)
	q.resolvedByDID = types.Some(resolvedByDID)
	if notes != "" {
		q.resolutionNotes = types.Some(notes)
	}
	q.resolvedAt = types.Some(types.Now())
	q.resolverSigner = resolverSigner

	return nil
}

func (q *QuarantinedRecord) Expire() error {
	if !q.resolution.IsPending() {
		return domainerror.ErrQuarantineAlreadyResolved
	}

	q.resolution = QuarantineResolutionExpired
	q.resolvedAt = types.Some(types.Now())

	return nil
}

func (q *QuarantinedRecord) IsPending() bool {
	return q.resolution.IsPending()
}

func (q *QuarantinedRecord) IsResolved() bool {
	return q.resolution.IsResolved()
}

func (q *QuarantinedRecord) IsExpired() bool {
	if q.resolution == QuarantineResolutionExpired {
		return true
	}

	if q.resolution.IsPending() && q.expiresAt.IsPresent() {
		return types.Now().After(q.expiresAt.MustGet())
	}

	return false
}

func (q *QuarantinedRecord) IsAccepted() bool {
	return q.resolution.IsAccepted()
}

func (q *QuarantinedRecord) IsRejected() bool {
	return q.resolution.IsRejected()
}

func (q *QuarantinedRecord) HasModifiedData() bool {
	return len(q.modifiedData) > 0
}

func (q *QuarantinedRecord) GetEffectiveData() map[string]any {
	if q.HasModifiedData() {
		return q.modifiedData
	}
	return q.rawData
}

func (q *QuarantinedRecord) AnomalyCount() int {
	return len(q.anomalies)
}

func (q *QuarantinedRecord) MaxAnomalySeverity() AnomalySeverity {
	max := AnomalySeverityUnspecified
	for _, a := range q.anomalies {
		if a.Severity() > max {
			max = a.Severity()
		}
	}
	return max
}

func (q *QuarantinedRecord) BelongsToTenant(tenantID types.ID) bool {
	if q.tenantID.IsEmpty() {
		return false
	}
	return q.tenantID.MustGet() == tenantID
}

func (q *QuarantinedRecord) RequiresSecurityReview() bool {
	return q.reason.RequiresSecurityReview()
}

func (q *QuarantinedRecord) Priority() int {
	return q.reason.Priority()
}
