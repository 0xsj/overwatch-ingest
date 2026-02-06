package testutil

import (
	"time"

	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/security"
	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

// Fixtures provides builders for domain models in tests.
var Fixtures = &fixtures{}

type fixtures struct{}

// ─────────────────────────────────────────────────────────────────
// Crypto / Provenance
// ─────────────────────────────────────────────────────────────────

// KeyPair generates a new Ed25519 keypair.
func (f *fixtures) KeyPair() security.KeyPair {
	kp, err := security.GenerateEd25519()
	if err != nil {
		panic("fixtures: failed to generate keypair: " + err.Error())
	}
	return kp
}

// DID generates a valid DID from a new keypair.
func (f *fixtures) DID() *security.DID {
	kp := f.KeyPair()
	did, err := security.DIDFromKeyPair(kp)
	if err != nil {
		panic("fixtures: failed to create DID: " + err.Error())
	}
	return did
}

// Signer generates a valid SignatureInfo for ingest provenance.
func (f *fixtures) Signer() *provenance.SignatureInfo {
	did := f.DID()
	return &provenance.SignatureInfo{
		DID:        did.String(),
		SignerType: provenance.SignerTypeService,
		Signature:  Fake.Hex(64),
		SignedAt:   time.Now(),
	}
}

// SourceSigner generates a valid SignatureInfo for a source.
func (f *fixtures) SourceSigner() *provenance.SignatureInfo {
	did := f.DID()
	return &provenance.SignatureInfo{
		DID:        did.String(),
		SignerType: provenance.SignerTypeSource,
		Signature:  Fake.Hex(64),
		SignedAt:   time.Now(),
	}
}

// ─────────────────────────────────────────────────────────────────
// IngestRecord
// ─────────────────────────────────────────────────────────────────

// IngestRecord creates a new IngestRecord with default values.
func (f *fixtures) IngestRecord() *model.IngestRecord {
	rec, err := model.NewIngestRecord(
		types.None[types.ID](),
		types.NewID(),
		"ais",
		Fake.RawDataID(),
		types.Now(),
	)
	if err != nil {
		panic("fixtures: failed to create ingest record: " + err.Error())
	}
	return rec
}

// IngestRecordBuilder returns a builder for customizing IngestRecord creation.
func (f *fixtures) IngestRecordBuilder() *IngestRecordBuilder {
	return &IngestRecordBuilder{
		tenantID:   types.None[types.ID](),
		sourceID:   types.NewID(),
		sourceType: "ais",
		rawDataID:  Fake.RawDataID(),
		receivedAt: types.Now(),
	}
}

// IngestRecordBuilder builds customized IngestRecord instances.
type IngestRecordBuilder struct {
	tenantID   types.Optional[types.ID]
	sourceID   types.ID
	sourceType string
	rawDataID  string
	receivedAt types.Timestamp
}

func (b *IngestRecordBuilder) WithTenantID(id types.ID) *IngestRecordBuilder {
	b.tenantID = types.Some(id)
	return b
}

func (b *IngestRecordBuilder) WithSourceID(id types.ID) *IngestRecordBuilder {
	b.sourceID = id
	return b
}

func (b *IngestRecordBuilder) WithSourceType(st string) *IngestRecordBuilder {
	b.sourceType = st
	return b
}

func (b *IngestRecordBuilder) WithRawDataID(id string) *IngestRecordBuilder {
	b.rawDataID = id
	return b
}

func (b *IngestRecordBuilder) WithReceivedAt(t types.Timestamp) *IngestRecordBuilder {
	b.receivedAt = t
	return b
}

func (b *IngestRecordBuilder) Build() *model.IngestRecord {
	rec, err := model.NewIngestRecord(b.tenantID, b.sourceID, b.sourceType, b.rawDataID, b.receivedAt)
	if err != nil {
		panic("fixtures: failed to create ingest record: " + err.Error())
	}
	return rec
}

// AcceptedIngestRecord creates an IngestRecord in accepted status.
func (f *fixtures) AcceptedIngestRecord() *model.IngestRecord {
	rec := f.IngestRecord()
	rec.SetValidation(f.ValidValidationResult())
	rec.SetConfidence(f.HighConfidenceScore())
	if err := rec.Accept("vessel", "mmsi:123456789", []string{}); err != nil {
		panic("fixtures: failed to accept record: " + err.Error())
	}
	return rec
}

// RejectedIngestRecord creates an IngestRecord in rejected status.
func (f *fixtures) RejectedIngestRecord() *model.IngestRecord {
	rec := f.IngestRecord()
	rec.SetValidation(f.InvalidValidationResult())
	rec.SetConfidence(f.LowConfidenceScore())
	if err := rec.Reject("validation failed"); err != nil {
		panic("fixtures: failed to reject record: " + err.Error())
	}
	return rec
}

// QuarantinedIngestRecord creates an IngestRecord in quarantined status.
func (f *fixtures) QuarantinedIngestRecord() *model.IngestRecord {
	rec := f.IngestRecord()
	rec.SetValidation(f.WarningValidationResult())
	rec.SetConfidence(f.MediumConfidenceScore())
	if err := rec.Quarantine(types.NewID()); err != nil {
		panic("fixtures: failed to quarantine record: " + err.Error())
	}
	return rec
}

// ─────────────────────────────────────────────────────────────────
// QuarantinedRecord
// ─────────────────────────────────────────────────────────────────

// QuarantinedRecord creates a new QuarantinedRecord with default values.
func (f *fixtures) QuarantinedRecord() *model.QuarantinedRecord {
	qr, err := model.NewQuarantinedRecord(
		types.None[types.ID](),
		types.NewID(),
		"ais",
		Fake.RawDataID(),
		types.NewID(),
		Fake.VesselPayload(),
		model.QuarantineReasonAnomalyDetected,
		"anomalies detected in data",
		[]model.Anomaly{f.WarningAnomaly()},
		f.MediumConfidenceScore(),
		types.Some(types.FromTime(time.Now().Add(72*time.Hour))),
		f.Signer(),
	)
	if err != nil {
		panic("fixtures: failed to create quarantined record: " + err.Error())
	}
	return qr
}

// QuarantinedRecordBuilder returns a builder for customizing QuarantinedRecord creation.
func (f *fixtures) QuarantinedRecordBuilder() *QuarantinedRecordBuilder {
	return &QuarantinedRecordBuilder{
		tenantID:       types.None[types.ID](),
		sourceID:       types.NewID(),
		sourceType:     "ais",
		rawDataID:      Fake.RawDataID(),
		ingestRecordID: types.NewID(),
		rawData:        Fake.VesselPayload(),
		reason:         model.QuarantineReasonAnomalyDetected,
		reasonDetail:   "anomalies detected in data",
		anomalies:      []model.Anomaly{Fixtures.WarningAnomaly()},
		confidence:     Fixtures.MediumConfidenceScore(),
		expiresAt:      types.Some(types.FromTime(time.Now().Add(72 * time.Hour))),
		ingestSigner:   Fixtures.Signer(),
	}
}

// QuarantinedRecordBuilder builds customized QuarantinedRecord instances.
type QuarantinedRecordBuilder struct {
	tenantID       types.Optional[types.ID]
	sourceID       types.ID
	sourceType     string
	rawDataID      string
	ingestRecordID types.ID
	rawData        map[string]any
	reason         model.QuarantineReason
	reasonDetail   string
	anomalies      []model.Anomaly
	confidence     model.ConfidenceScore
	expiresAt      types.Optional[types.Timestamp]
	ingestSigner   *provenance.SignatureInfo
}

func (b *QuarantinedRecordBuilder) WithTenantID(id types.ID) *QuarantinedRecordBuilder {
	b.tenantID = types.Some(id)
	return b
}

func (b *QuarantinedRecordBuilder) WithSourceID(id types.ID) *QuarantinedRecordBuilder {
	b.sourceID = id
	return b
}

func (b *QuarantinedRecordBuilder) WithSourceType(st string) *QuarantinedRecordBuilder {
	b.sourceType = st
	return b
}

func (b *QuarantinedRecordBuilder) WithIngestRecordID(id types.ID) *QuarantinedRecordBuilder {
	b.ingestRecordID = id
	return b
}

func (b *QuarantinedRecordBuilder) WithRawData(data map[string]any) *QuarantinedRecordBuilder {
	b.rawData = data
	return b
}

func (b *QuarantinedRecordBuilder) WithReason(reason model.QuarantineReason) *QuarantinedRecordBuilder {
	b.reason = reason
	return b
}

func (b *QuarantinedRecordBuilder) WithReasonDetail(detail string) *QuarantinedRecordBuilder {
	b.reasonDetail = detail
	return b
}

func (b *QuarantinedRecordBuilder) WithAnomalies(anomalies []model.Anomaly) *QuarantinedRecordBuilder {
	b.anomalies = anomalies
	return b
}

func (b *QuarantinedRecordBuilder) WithConfidence(cs model.ConfidenceScore) *QuarantinedRecordBuilder {
	b.confidence = cs
	return b
}

func (b *QuarantinedRecordBuilder) WithExpiresAt(t types.Timestamp) *QuarantinedRecordBuilder {
	b.expiresAt = types.Some(t)
	return b
}

func (b *QuarantinedRecordBuilder) WithNoExpiry() *QuarantinedRecordBuilder {
	b.expiresAt = types.None[types.Timestamp]()
	return b
}

func (b *QuarantinedRecordBuilder) Build() *model.QuarantinedRecord {
	qr, err := model.NewQuarantinedRecord(
		b.tenantID,
		b.sourceID,
		b.sourceType,
		b.rawDataID,
		b.ingestRecordID,
		b.rawData,
		b.reason,
		b.reasonDetail,
		b.anomalies,
		b.confidence,
		b.expiresAt,
		b.ingestSigner,
	)
	if err != nil {
		panic("fixtures: failed to create quarantined record: " + err.Error())
	}
	return qr
}

// ─────────────────────────────────────────────────────────────────
// ValidationResult
// ─────────────────────────────────────────────────────────────────

// ValidValidationResult creates a passing validation result with no anomalies.
func (f *fixtures) ValidValidationResult() model.ValidationResult {
	return model.ValidResult(
		[]string{"mmsi", "lat", "lon", "speed", "course", "timestamp"},
		"v1.0.0",
	)
}

// InvalidValidationResult creates a failing validation result.
func (f *fixtures) InvalidValidationResult() model.ValidationResult {
	return model.InvalidSchemaResult(
		[]string{"mmsi", "lat"},
		[]model.Anomaly{f.CriticalAnomaly()},
		"v1.0.0",
	)
}

// WarningValidationResult creates a validation result with warning-level anomalies.
func (f *fixtures) WarningValidationResult() model.ValidationResult {
	return model.NewValidationResult(
		true,
		[]string{"mmsi", "lat", "lon", "speed", "timestamp"},
		[]string{"course"},
		[]model.Anomaly{f.WarningAnomaly()},
		"v1.0.0",
	)
}

// ErrorValidationResult creates a validation result with error-level anomalies.
func (f *fixtures) ErrorValidationResult() model.ValidationResult {
	return model.NewValidationResult(
		true,
		[]string{"mmsi", "lat", "lon"},
		[]string{"speed", "course"},
		[]model.Anomaly{f.ErrorAnomaly()},
		"v1.0.0",
	)
}

// ─────────────────────────────────────────────────────────────────
// Anomaly
// ─────────────────────────────────────────────────────────────────

// InfoAnomaly creates an info-level anomaly.
func (f *fixtures) InfoAnomaly() model.Anomaly {
	return model.NewAnomaly("speed", model.AnomalyTypeOutOfRange, model.AnomalySeverityInfo, "speed slightly above typical range")
}

// WarningAnomaly creates a warning-level anomaly.
func (f *fixtures) WarningAnomaly() model.Anomaly {
	return model.NewAnomaly("lat", model.AnomalyTypeOutOfRange, model.AnomalySeverityWarning, "latitude near boundary")
}

// ErrorAnomaly creates an error-level anomaly.
func (f *fixtures) ErrorAnomaly() model.Anomaly {
	return model.NewAnomaly("timestamp", model.AnomalyTypeTemporal, model.AnomalySeverityError, "timestamp is in the future")
}

// CriticalAnomaly creates a critical-level anomaly.
func (f *fixtures) CriticalAnomaly() model.Anomaly {
	return model.NewAnomaly("mmsi", model.AnomalyTypeSuspicious, model.AnomalySeverityCritical, "MMSI matches known spoofing pattern")
}

// ─────────────────────────────────────────────────────────────────
// SourceReliability
// ─────────────────────────────────────────────────────────────────

// SourceReliability creates a new SourceReliability with default values.
func (f *fixtures) SourceReliability(sourceID types.ID) *model.SourceReliability {
	return model.NewSourceReliability(sourceID, types.None[types.ID]())
}

// ReliableSource creates a SourceReliability with high acceptance rate.
func (f *fixtures) ReliableSource(sourceID types.ID) *model.SourceReliability {
	rel := model.NewSourceReliability(sourceID, types.None[types.ID]())
	for i := 0; i < 10; i++ {
		rel.RecordAccepted()
	}
	return rel
}

// UnreliableSource creates a SourceReliability with high rejection rate.
func (f *fixtures) UnreliableSource(sourceID types.ID) *model.SourceReliability {
	rel := model.NewSourceReliability(sourceID, types.None[types.ID]())
	for i := 0; i < 10; i++ {
		rel.RecordRejected()
	}
	return rel
}

// ─────────────────────────────────────────────────────────────────
// ConfidenceScore
// ─────────────────────────────────────────────────────────────────

// HighConfidenceScore creates a high-confidence score (above accept threshold).
func (f *fixtures) HighConfidenceScore() model.ConfidenceScore {
	return model.NewConfidenceScore(0.9, 0.95, 0.95, 0.9, nil)
}

// MediumConfidenceScore creates a medium-confidence score (between accept and reject).
func (f *fixtures) MediumConfidenceScore() model.ConfidenceScore {
	return model.NewConfidenceScore(0.5, 0.6, 0.5, 0.5, nil)
}

// LowConfidenceScore creates a low-confidence score (below reject threshold).
func (f *fixtures) LowConfidenceScore() model.ConfidenceScore {
	return model.NewConfidenceScore(0.1, 0.2, 0.1, 0.1, nil)
}

// ZeroConfidenceScore creates a zero confidence score.
func (f *fixtures) ZeroConfidenceScore() model.ConfidenceScore {
	return model.ZeroConfidenceScore("test zero score")
}
