package grpc

import (
	"time"

	commonv1 "github.com/0xsj/overwatch-contracts/gen/go/common/v1"
	ingestv1 "github.com/0xsj/overwatch-contracts/gen/go/ingest/v1"
	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

// =============================================================================
// IngestStatus
// =============================================================================

func ingestStatusToProto(s model.IngestStatus) ingestv1.IngestStatus {
	switch s {
	case model.IngestStatusPending:
		return ingestv1.IngestStatus_INGEST_STATUS_PENDING
	case model.IngestStatusAccepted:
		return ingestv1.IngestStatus_INGEST_STATUS_ACCEPTED
	case model.IngestStatusQuarantined:
		return ingestv1.IngestStatus_INGEST_STATUS_QUARANTINED
	case model.IngestStatusRejected:
		return ingestv1.IngestStatus_INGEST_STATUS_REJECTED
	default:
		return ingestv1.IngestStatus_INGEST_STATUS_UNSPECIFIED
	}
}

func ingestStatusFromProto(s ingestv1.IngestStatus) model.IngestStatus {
	switch s {
	case ingestv1.IngestStatus_INGEST_STATUS_PENDING:
		return model.IngestStatusPending
	case ingestv1.IngestStatus_INGEST_STATUS_ACCEPTED:
		return model.IngestStatusAccepted
	case ingestv1.IngestStatus_INGEST_STATUS_QUARANTINED:
		return model.IngestStatusQuarantined
	case ingestv1.IngestStatus_INGEST_STATUS_REJECTED:
		return model.IngestStatusRejected
	default:
		return model.IngestStatusUnspecified
	}
}

// =============================================================================
// AnomalyType
// =============================================================================

func anomalyTypeToProto(t model.AnomalyType) ingestv1.AnomalyType {
	switch t {
	case model.AnomalyTypeOutOfRange:
		return ingestv1.AnomalyType_ANOMALY_TYPE_OUT_OF_RANGE
	case model.AnomalyTypeInvalidFormat:
		return ingestv1.AnomalyType_ANOMALY_TYPE_INVALID_FORMAT
	case model.AnomalyTypeMissingRequired:
		return ingestv1.AnomalyType_ANOMALY_TYPE_MISSING_REQUIRED
	case model.AnomalyTypeUnexpectedValue:
		return ingestv1.AnomalyType_ANOMALY_TYPE_UNEXPECTED_VALUE
	case model.AnomalyTypeTemporal:
		return ingestv1.AnomalyType_ANOMALY_TYPE_TEMPORAL
	case model.AnomalyTypeStatistical:
		return ingestv1.AnomalyType_ANOMALY_TYPE_STATISTICAL
	case model.AnomalyTypeDuplicate:
		return ingestv1.AnomalyType_ANOMALY_TYPE_DUPLICATE
	case model.AnomalyTypeSuspicious:
		return ingestv1.AnomalyType_ANOMALY_TYPE_SUSPICIOUS
	default:
		return ingestv1.AnomalyType_ANOMALY_TYPE_UNSPECIFIED
	}
}

func anomalyTypeFromProto(t ingestv1.AnomalyType) model.AnomalyType {
	switch t {
	case ingestv1.AnomalyType_ANOMALY_TYPE_OUT_OF_RANGE:
		return model.AnomalyTypeOutOfRange
	case ingestv1.AnomalyType_ANOMALY_TYPE_INVALID_FORMAT:
		return model.AnomalyTypeInvalidFormat
	case ingestv1.AnomalyType_ANOMALY_TYPE_MISSING_REQUIRED:
		return model.AnomalyTypeMissingRequired
	case ingestv1.AnomalyType_ANOMALY_TYPE_UNEXPECTED_VALUE:
		return model.AnomalyTypeUnexpectedValue
	case ingestv1.AnomalyType_ANOMALY_TYPE_TEMPORAL:
		return model.AnomalyTypeTemporal
	case ingestv1.AnomalyType_ANOMALY_TYPE_STATISTICAL:
		return model.AnomalyTypeStatistical
	case ingestv1.AnomalyType_ANOMALY_TYPE_DUPLICATE:
		return model.AnomalyTypeDuplicate
	case ingestv1.AnomalyType_ANOMALY_TYPE_SUSPICIOUS:
		return model.AnomalyTypeSuspicious
	default:
		return model.AnomalyTypeUnspecified
	}
}

// =============================================================================
// AnomalySeverity
// =============================================================================

func anomalySeverityToProto(s model.AnomalySeverity) ingestv1.AnomalySeverity {
	switch s {
	case model.AnomalySeverityInfo:
		return ingestv1.AnomalySeverity_ANOMALY_SEVERITY_INFO
	case model.AnomalySeverityWarning:
		return ingestv1.AnomalySeverity_ANOMALY_SEVERITY_WARNING
	case model.AnomalySeverityError:
		return ingestv1.AnomalySeverity_ANOMALY_SEVERITY_ERROR
	case model.AnomalySeverityCritical:
		return ingestv1.AnomalySeverity_ANOMALY_SEVERITY_CRITICAL
	default:
		return ingestv1.AnomalySeverity_ANOMALY_SEVERITY_UNSPECIFIED
	}
}

func anomalySeverityFromProto(s ingestv1.AnomalySeverity) model.AnomalySeverity {
	switch s {
	case ingestv1.AnomalySeverity_ANOMALY_SEVERITY_INFO:
		return model.AnomalySeverityInfo
	case ingestv1.AnomalySeverity_ANOMALY_SEVERITY_WARNING:
		return model.AnomalySeverityWarning
	case ingestv1.AnomalySeverity_ANOMALY_SEVERITY_ERROR:
		return model.AnomalySeverityError
	case ingestv1.AnomalySeverity_ANOMALY_SEVERITY_CRITICAL:
		return model.AnomalySeverityCritical
	default:
		return model.AnomalySeverityUnspecified
	}
}

// =============================================================================
// QuarantineReason
// =============================================================================

func quarantineReasonToProto(r model.QuarantineReason) ingestv1.QuarantineReason {
	switch r {
	case model.QuarantineReasonValidationFailed:
		return ingestv1.QuarantineReason_QUARANTINE_REASON_VALIDATION_FAILED
	case model.QuarantineReasonLowConfidence:
		return ingestv1.QuarantineReason_QUARANTINE_REASON_LOW_CONFIDENCE
	case model.QuarantineReasonAnomalyDetected:
		return ingestv1.QuarantineReason_QUARANTINE_REASON_ANOMALY_DETECTED
	case model.QuarantineReasonSignatureInvalid:
		return ingestv1.QuarantineReason_QUARANTINE_REASON_SIGNATURE_INVALID
	case model.QuarantineReasonDuplicateSuspected:
		return ingestv1.QuarantineReason_QUARANTINE_REASON_DUPLICATE_SUSPECTED
	case model.QuarantineReasonManualReview:
		return ingestv1.QuarantineReason_QUARANTINE_REASON_MANUAL_REVIEW
	default:
		return ingestv1.QuarantineReason_QUARANTINE_REASON_UNSPECIFIED
	}
}

func quarantineReasonFromProto(r ingestv1.QuarantineReason) model.QuarantineReason {
	switch r {
	case ingestv1.QuarantineReason_QUARANTINE_REASON_VALIDATION_FAILED:
		return model.QuarantineReasonValidationFailed
	case ingestv1.QuarantineReason_QUARANTINE_REASON_LOW_CONFIDENCE:
		return model.QuarantineReasonLowConfidence
	case ingestv1.QuarantineReason_QUARANTINE_REASON_ANOMALY_DETECTED:
		return model.QuarantineReasonAnomalyDetected
	case ingestv1.QuarantineReason_QUARANTINE_REASON_SIGNATURE_INVALID:
		return model.QuarantineReasonSignatureInvalid
	case ingestv1.QuarantineReason_QUARANTINE_REASON_DUPLICATE_SUSPECTED:
		return model.QuarantineReasonDuplicateSuspected
	case ingestv1.QuarantineReason_QUARANTINE_REASON_MANUAL_REVIEW:
		return model.QuarantineReasonManualReview
	default:
		return model.QuarantineReasonUnspecified
	}
}

// =============================================================================
// QuarantineResolution
// =============================================================================

func quarantineResolutionToProto(r model.QuarantineResolution) ingestv1.QuarantineResolution {
	switch r {
	case model.QuarantineResolutionPending:
		return ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_PENDING
	case model.QuarantineResolutionApproved:
		return ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_APPROVED
	case model.QuarantineResolutionModified:
		return ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_MODIFIED
	case model.QuarantineResolutionRejected:
		return ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_REJECTED
	case model.QuarantineResolutionExpired:
		return ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_EXPIRED
	default:
		return ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_UNSPECIFIED
	}
}

func quarantineResolutionFromProto(r ingestv1.QuarantineResolution) model.QuarantineResolution {
	switch r {
	case ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_PENDING:
		return model.QuarantineResolutionPending
	case ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_APPROVED:
		return model.QuarantineResolutionApproved
	case ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_MODIFIED:
		return model.QuarantineResolutionModified
	case ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_REJECTED:
		return model.QuarantineResolutionRejected
	case ingestv1.QuarantineResolution_QUARANTINE_RESOLUTION_EXPIRED:
		return model.QuarantineResolutionExpired
	default:
		return model.QuarantineResolutionUnspecified
	}
}

// =============================================================================
// Anomaly
// =============================================================================

func anomalyToProto(a model.Anomaly) *ingestv1.Anomaly {
	pb := &ingestv1.Anomaly{
		Field:    a.Field(),
		Type:     anomalyTypeToProto(a.Type()),
		Severity: anomalySeverityToProto(a.Severity()),
		Message:  a.Message(),
	}

	if a.Expected() != "" {
		pb.Expected = &[]string{a.Expected()}[0]
	}
	if a.Actual() != "" {
		pb.Actual = &[]string{a.Actual()}[0]
	}
	if len(a.Context()) > 0 {
		pb.Context, _ = structpb.NewStruct(a.Context())
	}

	return pb
}

func anomalyFromProto(pb *ingestv1.Anomaly) model.Anomaly {
	if pb == nil {
		return model.Anomaly{}
	}

	expected := ""
	if pb.Expected != nil {
		expected = *pb.Expected
	}
	actual := ""
	if pb.Actual != nil {
		actual = *pb.Actual
	}

	var context map[string]any
	if pb.Context != nil {
		context = pb.Context.AsMap()
	}

	return model.ReconstructAnomaly(
		pb.Field,
		anomalyTypeFromProto(pb.Type),
		anomalySeverityFromProto(pb.Severity),
		pb.Message,
		expected,
		actual,
		context,
	)
}

func anomaliesToProto(anomalies []model.Anomaly) []*ingestv1.Anomaly {
	if anomalies == nil {
		return nil
	}
	result := make([]*ingestv1.Anomaly, len(anomalies))
	for i, a := range anomalies {
		result[i] = anomalyToProto(a)
	}
	return result
}

func anomaliesFromProto(pbs []*ingestv1.Anomaly) []model.Anomaly {
	if pbs == nil {
		return nil
	}
	result := make([]model.Anomaly, len(pbs))
	for i, pb := range pbs {
		result[i] = anomalyFromProto(pb)
	}
	return result
}

// =============================================================================
// ValidationResult
// =============================================================================

func validationResultToProto(v model.ValidationResult) *ingestv1.ValidationResult {
	pb := &ingestv1.ValidationResult{
		Valid:         v.Valid(),
		SchemaValid:   v.SchemaValid(),
		FieldsPresent: v.FieldsPresent(),
		FieldsMissing: v.FieldsMissing(),
		Anomalies:     anomaliesToProto(v.Anomalies()),
	}

	if ver := v.ValidatorVersion(); ver != "" {
		pb.ValidatorVersion = &ver
	}

	return pb
}

func validationResultFromProto(pb *ingestv1.ValidationResult) model.ValidationResult {
	if pb == nil {
		return model.ValidationResult{}
	}

	validatorVersion := ""
	if pb.ValidatorVersion != nil {
		validatorVersion = *pb.ValidatorVersion
	}

	return model.ReconstructValidationResult(
		pb.Valid,
		pb.SchemaValid,
		pb.FieldsPresent,
		pb.FieldsMissing,
		anomaliesFromProto(pb.Anomalies),
		validatorVersion,
	)
}

// =============================================================================
// ConfidenceFactor
// =============================================================================

func confidenceFactorToProto(f model.ConfidenceFactor) *ingestv1.ConfidenceFactor {
	return &ingestv1.ConfidenceFactor{
		Name:   f.Name(),
		Score:  float32(f.Score()),
		Weight: float32(f.Weight()),
		Reason: f.Reason(),
	}
}

func confidenceFactorFromProto(pb *ingestv1.ConfidenceFactor) model.ConfidenceFactor {
	if pb == nil {
		return model.ConfidenceFactor{}
	}
	return model.NewConfidenceFactor(pb.Name, float64(pb.Score), float64(pb.Weight), pb.Reason)
}

func confidenceFactorsToProto(factors []model.ConfidenceFactor) []*ingestv1.ConfidenceFactor {
	if factors == nil {
		return nil
	}
	result := make([]*ingestv1.ConfidenceFactor, len(factors))
	for i, f := range factors {
		result[i] = confidenceFactorToProto(f)
	}
	return result
}

func confidenceFactorsFromProto(pbs []*ingestv1.ConfidenceFactor) []model.ConfidenceFactor {
	if pbs == nil {
		return nil
	}
	result := make([]model.ConfidenceFactor, len(pbs))
	for i, pb := range pbs {
		result[i] = confidenceFactorFromProto(pb)
	}
	return result
}

// =============================================================================
// ConfidenceScore
// =============================================================================

func confidenceScoreToProto(c model.ConfidenceScore) *ingestv1.ConfidenceScore {
	return &ingestv1.ConfidenceScore{
		Overall:           float32(c.Overall()),
		SourceReliability: float32(c.SourceReliability()),
		DataCompleteness:  float32(c.DataCompleteness()),
		TemporalFreshness: float32(c.TemporalFreshness()),
		SignatureTrust:    float32(c.SignatureTrust()),
		Factors:           confidenceFactorsToProto(c.Factors()),
	}
}

func confidenceScoreFromProto(pb *ingestv1.ConfidenceScore) model.ConfidenceScore {
	if pb == nil {
		return model.ConfidenceScore{}
	}
	return model.ReconstructConfidenceScore(
		float64(pb.Overall),
		float64(pb.SourceReliability),
		float64(pb.DataCompleteness),
		float64(pb.TemporalFreshness),
		float64(pb.SignatureTrust),
		confidenceFactorsFromProto(pb.Factors),
	)
}

// =============================================================================
// SignerInfo (uses common/v1 and provenance package)
// =============================================================================

func signerInfoToProto(s *provenance.SignatureInfo) *commonv1.SignerInfo {
	if s == nil {
		return nil
	}
	return &commonv1.SignerInfo{
		Did:       s.DID,
		Signature: s.Signature,
		SignedAt:  timestamppb.New(s.SignedAt),
	}
}

func signerInfoFromProto(pb *commonv1.SignerInfo) *provenance.SignatureInfo {
	if pb == nil {
		return nil
	}
	var signedAt time.Time
	if pb.SignedAt != nil {
		signedAt = pb.SignedAt.AsTime()
	}
	return &provenance.SignatureInfo{
		DID:       pb.Did,
		Signature: pb.Signature,
		SignedAt:  signedAt,
	}
}

// =============================================================================
// IngestRecord
// =============================================================================

func ingestRecordToProto(r *model.IngestRecord) *ingestv1.IngestRecord {
	if r == nil {
		return nil
	}

	pb := &ingestv1.IngestRecord{
		Id:          string(r.ID()),
		SourceId:    string(r.SourceID()),
		SourceType:  r.SourceType(),
		RawDataId:   r.RawDataID(),
		Status:      ingestStatusToProto(r.Status()),
		Validation:  validationResultToProto(r.Validation()),
		Confidence:  confidenceScoreToProto(r.Confidence()),
		EventIds:    r.EventIDs(),
		ReceivedAt:  timestamppb.New(r.ReceivedAt().Time()),
		ProcessedAt: timestamppb.New(r.ProcessedAt().Time()),
	}

	if r.TenantID().IsPresent() {
		tid := string(r.TenantID().MustGet())
		pb.TenantId = &tid
	}
	if r.EntityType().IsPresent() {
		et := r.EntityType().MustGet()
		pb.EntityType = &et
	}
	if r.EntityID().IsPresent() {
		eid := r.EntityID().MustGet()
		pb.EntityId = &eid
	}
	if r.RejectionReason().IsPresent() {
		rr := r.RejectionReason().MustGet()
		pb.RejectionReason = &rr
	}
	if r.QuarantineID().IsPresent() {
		qid := string(r.QuarantineID().MustGet())
		pb.QuarantineId = &qid
	}

	pb.SourceSigner = signerInfoToProto(r.SourceSigner())
	pb.CollectorSigner = signerInfoToProto(r.CollectorSigner())
	pb.IngestSigner = signerInfoToProto(r.IngestSigner())

	if r.SourceSignatureVerified().IsPresent() {
		v := r.SourceSignatureVerified().MustGet()
		pb.SourceSignatureVerified = &v
	}
	if r.CollectorSignatureVerified().IsPresent() {
		v := r.CollectorSignatureVerified().MustGet()
		pb.CollectorSignatureVerified = &v
	}

	return pb
}

func ingestRecordsToProto(records []*model.IngestRecord) []*ingestv1.IngestRecord {
	if records == nil {
		return nil
	}
	result := make([]*ingestv1.IngestRecord, len(records))
	for i, r := range records {
		result[i] = ingestRecordToProto(r)
	}
	return result
}

// =============================================================================
// QuarantinedRecord
// =============================================================================

func quarantinedRecordToProto(r *model.QuarantinedRecord) *ingestv1.QuarantinedRecord {
	if r == nil {
		return nil
	}

	pb := &ingestv1.QuarantinedRecord{
		Id:             string(r.ID()),
		SourceId:       string(r.SourceID()),
		SourceType:     r.SourceType(),
		RawDataId:      r.RawDataID(),
		IngestRecordId: string(r.IngestRecordID()),
		Reason:         quarantineReasonToProto(r.Reason()),
		ReasonDetail:   r.ReasonDetail(),
		Anomalies:      anomaliesToProto(r.Anomalies()),
		Confidence:     confidenceScoreToProto(r.Confidence()),
		Resolution:     quarantineResolutionToProto(r.Resolution()),
		QuarantinedAt:  timestamppb.New(r.QuarantinedAt().Time()),
		IngestSigner:   signerInfoToProto(r.IngestSigner()),
	}

	if r.TenantID().IsPresent() {
		tid := string(r.TenantID().MustGet())
		pb.TenantId = &tid
	}
	if r.RawData() != nil {
		pb.RawData, _ = structpb.NewStruct(r.RawData())
	}
	if r.ResolvedBy().IsPresent() {
		rb := r.ResolvedBy().MustGet()
		pb.ResolvedBy = &rb
	}
	if r.ResolvedByDID().IsPresent() {
		did := r.ResolvedByDID().MustGet()
		pb.ResolvedByDid = &did
	}
	if r.ResolutionNotes().IsPresent() {
		n := r.ResolutionNotes().MustGet()
		pb.ResolutionNotes = &n
	}
	if r.ModifiedData() != nil {
		pb.ModifiedData, _ = structpb.NewStruct(r.ModifiedData())
	}
	if r.ExpiresAt().IsPresent() {
		pb.ExpiresAt = timestamppb.New(r.ExpiresAt().MustGet().Time())
	}
	if r.ResolvedAt().IsPresent() {
		pb.ResolvedAt = timestamppb.New(r.ResolvedAt().MustGet().Time())
	}

	pb.ResolverSigner = signerInfoToProto(r.ResolverSigner())

	return pb
}

func quarantinedRecordsToProto(records []*model.QuarantinedRecord) []*ingestv1.QuarantinedRecord {
	if records == nil {
		return nil
	}
	result := make([]*ingestv1.QuarantinedRecord, len(records))
	for i, r := range records {
		result[i] = quarantinedRecordToProto(r)
	}
	return result
}

// =============================================================================
// SourceReliability
// =============================================================================

func sourceReliabilityToProto(r *model.SourceReliability) *ingestv1.SourceReliability {
	if r == nil {
		return nil
	}

	pb := &ingestv1.SourceReliability{
		SourceId:            string(r.SourceID()),
		ReliabilityScore:    float32(r.ReliabilityScore()),
		TotalRecords:        r.TotalRecords(),
		AcceptedRecords:     r.AcceptedRecords(),
		RejectedRecords:     r.RejectedRecords(),
		QuarantinedRecords:  r.QuarantinedRecords(),
		CorroboratedRecords: r.CorroboratedRecords(),
		DisputedRecords:     r.DisputedRecords(),
		CalculatedAt:        timestamppb.New(r.CalculatedAt().Time()),
		WindowStart:         timestamppb.New(r.WindowStart().Time()),
		WindowEnd:           timestamppb.New(r.WindowEnd().Time()),
	}

	if r.TenantID().IsPresent() {
		tid := string(r.TenantID().MustGet())
		pb.TenantId = &tid
	}

	return pb
}

func sourceReliabilitiesToProto(reliabilities []*model.SourceReliability) []*ingestv1.SourceReliability {
	if reliabilities == nil {
		return nil
	}
	result := make([]*ingestv1.SourceReliability, len(reliabilities))
	for i, r := range reliabilities {
		result[i] = sourceReliabilityToProto(r)
	}
	return result
}

// =============================================================================
// TimeRange
// =============================================================================

func timeRangeFromProto(pb *ingestv1.TimeRange) (start, end types.Optional[time.Time]) {
	if pb == nil {
		return types.None[time.Time](), types.None[time.Time]()
	}
	if pb.Start != nil {
		start = types.Some(pb.Start.AsTime())
	}
	if pb.End != nil {
		end = types.Some(pb.End.AsTime())
	}
	return start, end
}

// =============================================================================
// Helpers
// =============================================================================

func optionalString(s *string) types.Optional[string] {
	if s == nil {
		return types.None[string]()
	}
	return types.Some(*s)
}

func optionalID(s *string) types.Optional[types.ID] {
	if s == nil {
		return types.None[types.ID]()
	}
	return types.Some(types.ID(*s))
}
