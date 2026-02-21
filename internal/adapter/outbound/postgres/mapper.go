package postgres

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/adapter/outbound/postgres/sqlc"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

// =============================================================================
// pgtype helpers
// =============================================================================

func ensureStringSlice(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func stringToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func pgTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func optionalIDToPgText(opt types.Optional[types.ID]) pgtype.Text {
	if opt.IsEmpty() {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: string(opt.MustGet()), Valid: true}
}

func pgTextToOptionalID(t pgtype.Text) types.Optional[types.ID] {
	if !t.Valid || t.String == "" {
		return types.None[types.ID]()
	}
	return types.Some(types.ID(t.String))
}

func optionalStringToPgText(opt types.Optional[string]) pgtype.Text {
	if opt.IsEmpty() {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: opt.MustGet(), Valid: true}
}

func pgTextToOptionalString(t pgtype.Text) types.Optional[string] {
	if !t.Valid {
		return types.None[string]()
	}
	return types.Some(t.String)
}

func boolToPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func optionalBoolToPgBool(opt types.Optional[bool]) pgtype.Bool {
	if opt.IsEmpty() {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: opt.MustGet(), Valid: true}
}

func pgBoolToOptionalBool(b pgtype.Bool) types.Optional[bool] {
	if !b.Valid {
		return types.None[bool]()
	}
	return types.Some(b.Bool)
}

func timeToPgTimestamptz(t types.Timestamp) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t.Time(), Valid: true}
}

func optionalTimeToPgTimestamptz(opt types.Optional[types.Timestamp]) pgtype.Timestamptz {
	if opt.IsEmpty() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: opt.MustGet().Time(), Valid: true}
}

func pgTimestamptzToOptionalTimestamp(t pgtype.Timestamptz) types.Optional[types.Timestamp] {
	if !t.Valid {
		return types.None[types.Timestamp]()
	}
	return types.Some(types.FromTime(t.Time))
}

// =============================================================================
// JSON helpers
// =============================================================================

func marshalJSON(v any) []byte {
	if v == nil {
		return []byte("null")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return []byte("null")
	}
	return data
}

func unmarshalJSON[T any](data []byte) T {
	var v T
	if len(data) == 0 {
		return v
	}
	_ = json.Unmarshal(data, &v)
	return v
}

// =============================================================================
// SignerInfo JSON serialization
// =============================================================================

func signerInfoToJSON(s *provenance.SignatureInfo) []byte {
	if s == nil {
		return nil
	}
	return marshalJSON(s)
}

func signerInfoFromJSON(data []byte) *provenance.SignatureInfo {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var s provenance.SignatureInfo
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	return &s
}

// =============================================================================
// Anomaly JSON serialization
// =============================================================================

type anomalyJSON struct {
	Field    string         `json:"field"`
	Type     string         `json:"type"`
	Severity string         `json:"severity"`
	Message  string         `json:"message"`
	Expected string         `json:"expected,omitempty"`
	Actual   string         `json:"actual,omitempty"`
	Context  map[string]any `json:"context,omitempty"`
}

func anomaliesToJSON(anomalies []model.Anomaly) []byte {
	if len(anomalies) == 0 {
		return []byte("[]")
	}
	jsonAnomalies := make([]anomalyJSON, len(anomalies))
	for i, a := range anomalies {
		jsonAnomalies[i] = anomalyJSON{
			Field:    a.Field(),
			Type:     a.Type().String(),
			Severity: a.Severity().String(),
			Message:  a.Message(),
			Expected: a.Expected(),
			Actual:   a.Actual(),
			Context:  a.Context(),
		}
	}
	return marshalJSON(jsonAnomalies)
}

func anomaliesFromJSON(data []byte) []model.Anomaly {
	if len(data) == 0 || string(data) == "[]" {
		return nil
	}
	var jsonAnomalies []anomalyJSON
	if err := json.Unmarshal(data, &jsonAnomalies); err != nil {
		return nil
	}
	anomalies := make([]model.Anomaly, len(jsonAnomalies))
	for i, a := range jsonAnomalies {
		typ, _ := model.ParseAnomalyType(a.Type)
		sev, _ := model.ParseAnomalySeverity(a.Severity)
		anomalies[i] = model.ReconstructAnomaly(
			a.Field,
			typ,
			sev,
			a.Message,
			a.Expected,
			a.Actual,
			a.Context,
		)
	}
	return anomalies
}

// =============================================================================
// ConfidenceFactor JSON serialization
// =============================================================================

type confidenceFactorJSON struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
	Weight float64 `json:"weight"`
	Reason string  `json:"reason"`
}

func confidenceFactorsToJSON(factors []model.ConfidenceFactor) []byte {
	if len(factors) == 0 {
		return []byte("[]")
	}
	jsonFactors := make([]confidenceFactorJSON, len(factors))
	for i, f := range factors {
		jsonFactors[i] = confidenceFactorJSON{
			Name:   f.Name(),
			Score:  f.Score(),
			Weight: f.Weight(),
			Reason: f.Reason(),
		}
	}
	return marshalJSON(jsonFactors)
}

func confidenceFactorsFromJSON(data []byte) []model.ConfidenceFactor {
	if len(data) == 0 || string(data) == "[]" {
		return nil
	}
	var jsonFactors []confidenceFactorJSON
	if err := json.Unmarshal(data, &jsonFactors); err != nil {
		return nil
	}
	factors := make([]model.ConfidenceFactor, len(jsonFactors))
	for i, f := range jsonFactors {
		factors[i] = model.NewConfidenceFactor(f.Name, f.Score, f.Weight, f.Reason)
	}
	return factors
}

// =============================================================================
// IngestStatus mapping
// =============================================================================

func ingestStatusToString(s model.IngestStatus) string {
	return s.String()
}

func ingestStatusFromString(s string) model.IngestStatus {
	status, _ := model.ParseIngestStatus(s)
	return status
}

func float32ToPgFloat4(f float32) pgtype.Float4 {
	return pgtype.Float4{Float32: f, Valid: true}
}

// =============================================================================
// QuarantineReason mapping
// =============================================================================

func quarantineReasonToString(r model.QuarantineReason) string {
	return r.String()
}

func quarantineReasonFromString(s string) model.QuarantineReason {
	reason, _ := model.ParseQuarantineReason(s)
	return reason
}

// =============================================================================
// QuarantineResolution mapping
// =============================================================================

func quarantineResolutionToString(r model.QuarantineResolution) string {
	return r.String()
}

func quarantineResolutionFromString(s string) model.QuarantineResolution {
	resolution, _ := model.ParseQuarantineResolution(s)
	return resolution
}

// =============================================================================
// IngestRecord mapping
// =============================================================================

func toCreateIngestRecordParams(r *model.IngestRecord) sqlc.CreateIngestRecordParams {
	return sqlc.CreateIngestRecordParams{
		ID:                          string(r.ID()),
		TenantID:                    optionalIDToPgText(r.TenantID()),
		SourceID:                    string(r.SourceID()),
		SourceType:                  r.SourceType(),
		RawDataID:                   r.RawDataID(),
		Status:                      ingestStatusToString(r.Status()),
		ValidationValid:             r.Validation().Valid(),
		ValidationSchemaValid:       r.Validation().SchemaValid(),
		ValidationFieldsPresent:     ensureStringSlice(r.Validation().FieldsPresent()),
		ValidationFieldsMissing:     ensureStringSlice(r.Validation().FieldsMissing()),
		ValidationAnomalies:         anomaliesToJSON(r.Validation().Anomalies()),
		ValidationValidatorVersion:  stringToPgText(r.Validation().ValidatorVersion()),
		ConfidenceOverall:           float32(r.Confidence().Overall()),
		ConfidenceSourceReliability: float32(r.Confidence().SourceReliability()),
		ConfidenceDataCompleteness:  float32(r.Confidence().DataCompleteness()),
		ConfidenceTemporalFreshness: float32(r.Confidence().TemporalFreshness()),
		ConfidenceSignatureTrust:    float32(r.Confidence().SignatureTrust()),
		ConfidenceFactors:           confidenceFactorsToJSON(r.Confidence().Factors()),
		EntityType:                  optionalStringToPgText(r.EntityType()),
		EntityID:                    optionalStringToPgText(r.EntityID()),
		EventIds:                    ensureStringSlice(r.EventIDs()),
		RejectionReason:             optionalStringToPgText(r.RejectionReason()),
		QuarantineID:                optionalIDToPgText(r.QuarantineID()),
		SourceSigner:                signerInfoToJSON(r.SourceSigner()),
		CollectorSigner:             signerInfoToJSON(r.CollectorSigner()),
		IngestSigner:                signerInfoToJSON(r.IngestSigner()),
		SourceSignatureVerified:     optionalBoolToPgBool(r.SourceSignatureVerified()),
		CollectorSignatureVerified:  optionalBoolToPgBool(r.CollectorSignatureVerified()),
		ReceivedAt:                  r.ReceivedAt().Time(),
		ProcessedAt:                 r.ProcessedAt().Time(),
	}
}

func toUpdateIngestRecordParams(r *model.IngestRecord) sqlc.UpdateIngestRecordParams {
	return sqlc.UpdateIngestRecordParams{
		ID:                          string(r.ID()),
		Status:                      ingestStatusToString(r.Status()),
		ValidationValid:             r.Validation().Valid(),
		ValidationSchemaValid:       r.Validation().SchemaValid(),
		ValidationFieldsPresent:     ensureStringSlice(r.Validation().FieldsPresent()),
		ValidationFieldsMissing:     ensureStringSlice(r.Validation().FieldsMissing()),
		ValidationAnomalies:         anomaliesToJSON(r.Validation().Anomalies()),
		ValidationValidatorVersion:  stringToPgText(r.Validation().ValidatorVersion()),
		ConfidenceOverall:           float32(r.Confidence().Overall()),
		ConfidenceSourceReliability: float32(r.Confidence().SourceReliability()),
		ConfidenceDataCompleteness:  float32(r.Confidence().DataCompleteness()),
		ConfidenceTemporalFreshness: float32(r.Confidence().TemporalFreshness()),
		ConfidenceSignatureTrust:    float32(r.Confidence().SignatureTrust()),
		ConfidenceFactors:           confidenceFactorsToJSON(r.Confidence().Factors()),
		EntityType:                  optionalStringToPgText(r.EntityType()),
		EntityID:                    optionalStringToPgText(r.EntityID()),
		EventIds:                    ensureStringSlice(r.EventIDs()),
		RejectionReason:             optionalStringToPgText(r.RejectionReason()),
		QuarantineID:                optionalIDToPgText(r.QuarantineID()),
		SourceSigner:                signerInfoToJSON(r.SourceSigner()),
		CollectorSigner:             signerInfoToJSON(r.CollectorSigner()),
		IngestSigner:                signerInfoToJSON(r.IngestSigner()),
		SourceSignatureVerified:     optionalBoolToPgBool(r.SourceSignatureVerified()),
		CollectorSignatureVerified:  optionalBoolToPgBool(r.CollectorSignatureVerified()),
		ProcessedAt:                 r.ProcessedAt().Time(),
	}
}

func toIngestRecordModel(row sqlc.IngestRecord) *model.IngestRecord {
	validation := model.ReconstructValidationResult(
		row.ValidationValid,
		row.ValidationSchemaValid,
		row.ValidationFieldsPresent,
		row.ValidationFieldsMissing,
		anomaliesFromJSON(row.ValidationAnomalies),
		pgTextToString(row.ValidationValidatorVersion),
	)

	confidence := model.ReconstructConfidenceScore(
		float64(row.ConfidenceOverall),
		float64(row.ConfidenceSourceReliability),
		float64(row.ConfidenceDataCompleteness),
		float64(row.ConfidenceTemporalFreshness),
		float64(row.ConfidenceSignatureTrust),
		confidenceFactorsFromJSON(row.ConfidenceFactors),
	)

	return model.ReconstructIngestRecord(
		types.ID(row.ID),
		pgTextToOptionalID(row.TenantID),
		types.ID(row.SourceID),
		row.SourceType,
		row.RawDataID,
		ingestStatusFromString(row.Status),
		validation,
		confidence,
		pgTextToOptionalString(row.EntityType),
		pgTextToOptionalString(row.EntityID),
		row.EventIds,
		pgTextToOptionalString(row.RejectionReason),
		pgTextToOptionalID(row.QuarantineID),
		signerInfoFromJSON(row.SourceSigner),
		signerInfoFromJSON(row.CollectorSigner),
		signerInfoFromJSON(row.IngestSigner),
		pgBoolToOptionalBool(row.SourceSignatureVerified),
		pgBoolToOptionalBool(row.CollectorSignatureVerified),
		types.FromTime(row.ReceivedAt),
		types.FromTime(row.ProcessedAt),
	)
}

// =============================================================================
// QuarantinedRecord mapping
// =============================================================================

func toCreateQuarantinedRecordParams(r *model.QuarantinedRecord) sqlc.CreateQuarantinedRecordParams {
	return sqlc.CreateQuarantinedRecordParams{
		ID:                          string(r.ID()),
		TenantID:                    optionalIDToPgText(r.TenantID()),
		SourceID:                    string(r.SourceID()),
		SourceType:                  r.SourceType(),
		RawDataID:                   r.RawDataID(),
		IngestRecordID:              string(r.IngestRecordID()),
		RawData:                     marshalJSON(r.RawData()),
		Reason:                      quarantineReasonToString(r.Reason()),
		ReasonDetail:                r.ReasonDetail(),
		Anomalies:                   anomaliesToJSON(r.Anomalies()),
		ConfidenceOverall:           float32(r.Confidence().Overall()),
		ConfidenceSourceReliability: float32(r.Confidence().SourceReliability()),
		ConfidenceDataCompleteness:  float32(r.Confidence().DataCompleteness()),
		ConfidenceTemporalFreshness: float32(r.Confidence().TemporalFreshness()),
		ConfidenceSignatureTrust:    float32(r.Confidence().SignatureTrust()),
		ConfidenceFactors:           confidenceFactorsToJSON(r.Confidence().Factors()),
		Resolution:                  quarantineResolutionToString(r.Resolution()),
		ResolvedBy:                  optionalStringToPgText(r.ResolvedBy()),
		ResolvedByDid:               optionalStringToPgText(r.ResolvedByDID()),
		ResolutionNotes:             optionalStringToPgText(r.ResolutionNotes()),
		ModifiedData:                marshalJSON(r.ModifiedData()),
		IngestSigner:                signerInfoToJSON(r.IngestSigner()),
		ResolverSigner:              signerInfoToJSON(r.ResolverSigner()),
		QuarantinedAt:               r.QuarantinedAt().Time(),
		ExpiresAt:                   optionalTimeToPgTimestamptz(r.ExpiresAt()),
		ResolvedAt:                  optionalTimeToPgTimestamptz(r.ResolvedAt()),
	}
}

func toUpdateQuarantinedRecordParams(r *model.QuarantinedRecord) sqlc.UpdateQuarantinedRecordParams {
	return sqlc.UpdateQuarantinedRecordParams{
		ID:              string(r.ID()),
		Resolution:      quarantineResolutionToString(r.Resolution()),
		ResolvedBy:      optionalStringToPgText(r.ResolvedBy()),
		ResolvedByDid:   optionalStringToPgText(r.ResolvedByDID()),
		ResolutionNotes: optionalStringToPgText(r.ResolutionNotes()),
		ModifiedData:    marshalJSON(r.ModifiedData()),
		ResolverSigner:  signerInfoToJSON(r.ResolverSigner()),
		ResolvedAt:      optionalTimeToPgTimestamptz(r.ResolvedAt()),
	}
}

func toQuarantinedRecordModel(row sqlc.QuarantinedRecord) *model.QuarantinedRecord {
	confidence := model.ReconstructConfidenceScore(
		float64(row.ConfidenceOverall),
		float64(row.ConfidenceSourceReliability),
		float64(row.ConfidenceDataCompleteness),
		float64(row.ConfidenceTemporalFreshness),
		float64(row.ConfidenceSignatureTrust),
		confidenceFactorsFromJSON(row.ConfidenceFactors),
	)

	var rawData map[string]any
	if len(row.RawData) > 0 {
		_ = json.Unmarshal(row.RawData, &rawData)
	}

	var modifiedData map[string]any
	if len(row.ModifiedData) > 0 {
		_ = json.Unmarshal(row.ModifiedData, &modifiedData)
	}

	return model.ReconstructQuarantinedRecord(
		types.ID(row.ID),
		pgTextToOptionalID(row.TenantID),
		types.ID(row.SourceID),
		row.SourceType,
		row.RawDataID,
		types.ID(row.IngestRecordID),
		rawData,
		quarantineReasonFromString(row.Reason),
		row.ReasonDetail,
		anomaliesFromJSON(row.Anomalies),
		confidence,
		quarantineResolutionFromString(row.Resolution),
		pgTextToOptionalString(row.ResolvedBy),
		pgTextToOptionalString(row.ResolvedByDid),
		pgTextToOptionalString(row.ResolutionNotes),
		modifiedData,
		types.FromTime(row.QuarantinedAt),
		pgTimestamptzToOptionalTimestamp(row.ExpiresAt),
		pgTimestamptzToOptionalTimestamp(row.ResolvedAt),
		signerInfoFromJSON(row.IngestSigner),
		signerInfoFromJSON(row.ResolverSigner),
	)
}

// =============================================================================
// SourceReliability mapping
// =============================================================================

func toUpsertSourceReliabilityParams(r *model.SourceReliability) sqlc.UpsertSourceReliabilityParams {
	return sqlc.UpsertSourceReliabilityParams{
		SourceID:            string(r.SourceID()),
		TenantID:            optionalIDToPgText(r.TenantID()),
		ReliabilityScore:    float32(r.ReliabilityScore()),
		TotalRecords:        r.TotalRecords(),
		AcceptedRecords:     r.AcceptedRecords(),
		RejectedRecords:     r.RejectedRecords(),
		QuarantinedRecords:  r.QuarantinedRecords(),
		CorroboratedRecords: r.CorroboratedRecords(),
		DisputedRecords:     r.DisputedRecords(),
		CalculatedAt:        r.CalculatedAt().Time(),
		WindowStart:         r.WindowStart().Time(),
		WindowEnd:           r.WindowEnd().Time(),
	}
}

func toUpdateSourceReliabilityParams(r *model.SourceReliability) sqlc.UpdateSourceReliabilityParams {
	return sqlc.UpdateSourceReliabilityParams{
		SourceID:            string(r.SourceID()),
		ReliabilityScore:    float32(r.ReliabilityScore()),
		TotalRecords:        r.TotalRecords(),
		AcceptedRecords:     r.AcceptedRecords(),
		RejectedRecords:     r.RejectedRecords(),
		QuarantinedRecords:  r.QuarantinedRecords(),
		CorroboratedRecords: r.CorroboratedRecords(),
		DisputedRecords:     r.DisputedRecords(),
		CalculatedAt:        r.CalculatedAt().Time(),
		WindowStart:         r.WindowStart().Time(),
		WindowEnd:           r.WindowEnd().Time(),
	}
}

func toSourceReliabilityModel(row sqlc.SourceReliability) *model.SourceReliability {
	return model.ReconstructSourceReliability(
		types.ID(row.SourceID),
		pgTextToOptionalID(row.TenantID),
		float64(row.ReliabilityScore),
		row.TotalRecords,
		row.AcceptedRecords,
		row.RejectedRecords,
		row.QuarantinedRecords,
		row.CorroboratedRecords,
		row.DisputedRecords,
		types.FromTime(row.CalculatedAt),
		types.FromTime(row.WindowStart),
		types.FromTime(row.WindowEnd),
	)
}
