package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/domain/event"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/messaging"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/validation"
)

type processRawDataHandler struct {
	recordRepo       repository.IngestRecordRepository
	quarantineRepo   repository.QuarantinedRecordRepository
	reliabilityRepo  repository.SourceReliabilityRepository
	publisher        messaging.EventPublisher
	schemaValidator  validation.SchemaValidator
	anomalyDetector  validation.AnomalyDetector
	confidenceScorer validation.ConfidenceScorer
	verifier         provenance.Verifier
	signer           *provenance.EnvelopeBuilder
	thresholds       validation.ScoringThresholds
	quarantineExpiry time.Duration
}

type ProcessRawDataHandlerConfig struct {
	Thresholds       validation.ScoringThresholds
	QuarantineExpiry time.Duration
}

func DefaultProcessRawDataHandlerConfig() ProcessRawDataHandlerConfig {
	return ProcessRawDataHandlerConfig{
		Thresholds:       validation.DefaultScoringThresholds(),
		QuarantineExpiry: 72 * time.Hour,
	}
}

func NewProcessRawDataHandler(
	recordRepo repository.IngestRecordRepository,
	quarantineRepo repository.QuarantinedRecordRepository,
	reliabilityRepo repository.SourceReliabilityRepository,
	publisher messaging.EventPublisher,
	schemaValidator validation.SchemaValidator,
	anomalyDetector validation.AnomalyDetector,
	confidenceScorer validation.ConfidenceScorer,
	verifier provenance.Verifier,
	signer *provenance.EnvelopeBuilder,
	config ProcessRawDataHandlerConfig,
) command.ProcessRawDataHandler {
	return &processRawDataHandler{
		recordRepo:       recordRepo,
		quarantineRepo:   quarantineRepo,
		reliabilityRepo:  reliabilityRepo,
		publisher:        publisher,
		schemaValidator:  schemaValidator,
		anomalyDetector:  anomalyDetector,
		confidenceScorer: confidenceScorer,
		verifier:         verifier,
		signer:           signer,
		thresholds:       config.Thresholds,
		quarantineExpiry: config.QuarantineExpiry,
	}
}

func (h *processRawDataHandler) Handle(ctx context.Context, cmd command.ProcessRawData) (command.ProcessRawDataResult, error) {
	if cmd.SourceID.IsEmpty() {
		return command.ProcessRawDataResult{}, domainerror.ErrSourceIDRequired
	}
	if cmd.SourceType == "" {
		return command.ProcessRawDataResult{}, domainerror.ErrSourceTypeRequired
	}
	if cmd.RawDataID == "" {
		return command.ProcessRawDataResult{}, domainerror.ErrRawDataIDRequired
	}
	if cmd.Payload == nil {
		return command.ProcessRawDataResult{}, domainerror.ErrPayloadRequired
	}

	exists, err := h.recordRepo.ExistsByRawDataID(ctx, cmd.RawDataID)
	if err != nil {
		return command.ProcessRawDataResult{}, err
	}
	if exists {
		return command.ProcessRawDataResult{}, domainerror.RecordAlreadyExists(cmd.RawDataID)
	}

	record, err := model.NewIngestRecord(
		cmd.TenantID,
		cmd.SourceID,
		cmd.SourceType,
		cmd.RawDataID,
		cmd.CollectedAt,
	)
	if err != nil {
		return command.ProcessRawDataResult{}, err
	}

	collectorVerified, err := h.verifyCollectorSignature(ctx, cmd, record)
	if err != nil {
		return command.ProcessRawDataResult{}, err
	}

	sourceVerified := h.verifySourceSignature(ctx, cmd, record)

	_ = h.publisher.PublishRecordReceived(ctx, event.NewRecordReceived(
		cmd.TenantID,
		cmd.SourceID,
		record.ID(),
		cmd.RawDataID,
		cmd.SourceType,
		h.extractDID(cmd.SourceSigner),
		h.extractDID(cmd.CollectorSigner),
	))

	schemaResult := h.schemaValidator.Validate(ctx, cmd.SourceType, cmd.Payload)
	validationResult := schemaResult.ToValidationResult(h.schemaValidator.Version())

	anomalyResult := h.anomalyDetector.Detect(ctx, cmd.SourceType, cmd.Payload, validation.DetectionMetadata{
		SourceID:        cmd.SourceID,
		TenantID:        cmd.TenantID,
		SourceTimestamp: cmd.SourceTimestamp,
		CollectedAt:     cmd.CollectedAt,
	})

	allAnomalies := h.mergeAnomalies(validationResult.Anomalies(), anomalyResult.Anomalies)
	validationResult = model.NewValidationResult(
		validationResult.SchemaValid(),
		validationResult.FieldsPresent(),
		validationResult.FieldsMissing(),
		allAnomalies,
		validationResult.ValidatorVersion(),
	)

	record.SetValidation(validationResult)

	reliability, _ := h.reliabilityRepo.FindBySourceID(ctx, cmd.SourceID)

	confidenceScore := h.confidenceScorer.Score(ctx, validation.ScoringInput{
		SourceID:                   cmd.SourceID,
		SourceType:                 cmd.SourceType,
		TenantID:                   cmd.TenantID,
		Payload:                    cmd.Payload,
		ValidationResult:           validationResult,
		Anomalies:                  allAnomalies,
		SourceReliability:          reliability,
		SourceTimestamp:            cmd.SourceTimestamp,
		CollectedAt:                cmd.CollectedAt,
		HasSourceSigner:            cmd.SourceSigner != nil,
		SourceSignatureVerified:    sourceVerified,
		HasCollectorSigner:         cmd.CollectorSigner != nil,
		CollectorSignatureVerified: collectorVerified,
	})

	record.SetConfidence(confidenceScore)

	_ = h.publisher.PublishRecordValidated(ctx, event.NewRecordValidated(
		cmd.TenantID,
		cmd.SourceID,
		record.ID(),
		cmd.RawDataID,
		cmd.SourceType,
		validationResult.Valid(),
		validationResult.AnomalyCount(),
		confidenceScore.Overall(),
	))

	result := h.routeRecord(ctx, cmd, record, validationResult, anomalyResult, confidenceScore, reliability)

	return result, nil
}

func (h *processRawDataHandler) verifyCollectorSignature(ctx context.Context, cmd command.ProcessRawData, record *model.IngestRecord) (bool, error) {
	if cmd.CollectorSigner == nil {
		return false, domainerror.ErrSignatureRequired
	}

	record.SetCollectorSigner(cmd.CollectorSigner, false)

	// If the envelope was already verified at the transport layer (NATS consumer),
	// trust that result instead of re-verifying against the inner payload
	// (which would fail because the collector signs the full envelope, not just the payload).
	if cmd.EnvelopeVerified {
		record.SetCollectorSigner(cmd.CollectorSigner, true)

		_ = h.publisher.PublishSignatureVerified(ctx, event.NewSignatureVerified(
			cmd.TenantID,
			cmd.SourceID,
			record.ID(),
			cmd.RawDataID,
			cmd.CollectorSigner.DID,
			"collector",
			types.FromTime(cmd.CollectorSigner.SignedAt),
		))

		return true, nil
	}

	// Fallback: verify against payload bytes for non-envelope sources (e.g., HTTP).
	payloadBytes, err := json.Marshal(cmd.Payload)
	if err != nil {
		return false, domainerror.ErrPayloadInvalid
	}

	err = h.verifier.VerifySignatureInfo(ctx, payloadBytes, cmd.CollectorSigner)
	if err != nil {
		_ = h.publisher.PublishSignatureFailed(ctx, event.NewSignatureFailed(
			cmd.TenantID,
			cmd.SourceID,
			record.ID(),
			cmd.RawDataID,
			cmd.CollectorSigner.DID,
			"collector",
			err.Error(),
		))
		return false, nil
	}

	record.SetCollectorSigner(cmd.CollectorSigner, true)

	_ = h.publisher.PublishSignatureVerified(ctx, event.NewSignatureVerified(
		cmd.TenantID,
		cmd.SourceID,
		record.ID(),
		cmd.RawDataID,
		cmd.CollectorSigner.DID,
		"collector",
		types.FromTime(cmd.CollectorSigner.SignedAt),
	))

	return true, nil
}

func (h *processRawDataHandler) verifySourceSignature(ctx context.Context, cmd command.ProcessRawData, record *model.IngestRecord) bool {
	if cmd.SourceSigner == nil {
		return false
	}

	record.SetSourceSigner(cmd.SourceSigner, false)

	payloadBytes, err := json.Marshal(cmd.Payload)
	if err != nil {
		return false
	}

	err = h.verifier.VerifySignatureInfo(ctx, payloadBytes, cmd.SourceSigner)
	if err != nil {
		_ = h.publisher.PublishSignatureFailed(ctx, event.NewSignatureFailed(
			cmd.TenantID,
			cmd.SourceID,
			record.ID(),
			cmd.RawDataID,
			cmd.SourceSigner.DID,
			"source",
			err.Error(),
		))
		return false
	}

	record.SetSourceSigner(cmd.SourceSigner, true)

	_ = h.publisher.PublishSignatureVerified(ctx, event.NewSignatureVerified(
		cmd.TenantID,
		cmd.SourceID,
		record.ID(),
		cmd.RawDataID,
		cmd.SourceSigner.DID,
		"source",
		types.FromTime(cmd.SourceSigner.SignedAt),
	))

	return true
}

func (h *processRawDataHandler) routeRecord(
	ctx context.Context,
	cmd command.ProcessRawData,
	record *model.IngestRecord,
	validationResult model.ValidationResult,
	anomalyResult validation.AnomalyDetectionResult,
	confidenceScore model.ConfidenceScore,
	reliability *model.SourceReliability,
) command.ProcessRawDataResult {
	if validationResult.ShouldReject() || anomalyResult.ShouldReject {
		return h.rejectRecord(ctx, cmd, record, "validation failed", reliability)
	}

	if h.thresholds.ShouldReject(confidenceScore) {
		return h.rejectRecord(ctx, cmd, record, "confidence too low", reliability)
	}

	if validationResult.ShouldQuarantine() || anomalyResult.ShouldQuarantine || h.thresholds.ShouldQuarantine(confidenceScore) {
		return h.quarantineRecord(ctx, cmd, record, validationResult, confidenceScore, reliability)
	}

	return h.acceptRecord(ctx, cmd, record, confidenceScore, reliability)
}

func (h *processRawDataHandler) acceptRecord(
	ctx context.Context,
	cmd command.ProcessRawData,
	record *model.IngestRecord,
	confidenceScore model.ConfidenceScore,
	reliability *model.SourceReliability,
) command.ProcessRawDataResult {
	entityType := h.inferEntityType(cmd.SourceType, cmd.Payload)
	entityID := h.inferEntityID(cmd.SourceType, cmd.Payload)
	eventIDs := []string{}

	if err := record.Accept(entityType, entityID, eventIDs); err != nil {
		return command.ProcessRawDataResult{}
	}

	h.signRecord(record)

	if err := h.recordRepo.Create(ctx, record); err != nil {
		return command.ProcessRawDataResult{}
	}

	h.updateReliability(ctx, cmd.SourceID, cmd.TenantID, reliability, model.IngestStatusAccepted)

	_ = h.publisher.PublishRecordAccepted(ctx, event.NewRecordAccepted(
		cmd.TenantID,
		cmd.SourceID,
		record.ID(),
		cmd.RawDataID,
		cmd.SourceType,
		entityType,
		entityID,
		eventIDs,
		confidenceScore.Overall(),
		record.SourceSigner(),
		record.CollectorSigner(),
		record.IngestSigner(),
	))

	return command.ProcessRawDataResult{
		IngestRecordID:  record.ID(),
		Status:          model.IngestStatusAccepted.String(),
		EntityType:      types.Some(entityType),
		EntityID:        types.Some(entityID),
		EventIDs:        eventIDs,
		QuarantineID:    types.None[types.ID](),
		RejectionReason: types.None[string](),
		ConfidenceScore: confidenceScore.Overall(),
	}
}

func (h *processRawDataHandler) rejectRecord(
	ctx context.Context,
	cmd command.ProcessRawData,
	record *model.IngestRecord,
	reason string,
	reliability *model.SourceReliability,
) command.ProcessRawDataResult {
	if err := record.Reject(reason); err != nil {
		return command.ProcessRawDataResult{}
	}

	h.signRecord(record)

	if err := h.recordRepo.Create(ctx, record); err != nil {
		return command.ProcessRawDataResult{}
	}

	h.updateReliability(ctx, cmd.SourceID, cmd.TenantID, reliability, model.IngestStatusRejected)

	_ = h.publisher.PublishRecordRejected(ctx, event.NewRecordRejected(
		cmd.TenantID,
		cmd.SourceID,
		record.ID(),
		cmd.RawDataID,
		cmd.SourceType,
		reason,
		record.Validation().Anomalies(),
		record.Confidence().Overall(),
	))

	return command.ProcessRawDataResult{
		IngestRecordID:  record.ID(),
		Status:          model.IngestStatusRejected.String(),
		EntityType:      types.None[string](),
		EntityID:        types.None[string](),
		EventIDs:        nil,
		QuarantineID:    types.None[types.ID](),
		RejectionReason: types.Some(reason),
		ConfidenceScore: record.Confidence().Overall(),
	}
}

func (h *processRawDataHandler) quarantineRecord(
	ctx context.Context,
	cmd command.ProcessRawData,
	record *model.IngestRecord,
	validationResult model.ValidationResult,
	confidenceScore model.ConfidenceScore,
	reliability *model.SourceReliability,
) command.ProcessRawDataResult {
	reason := h.determineQuarantineReason(validationResult, confidenceScore)
	reasonDetail := h.buildQuarantineReasonDetail(validationResult, confidenceScore)

	expiresAt := types.Some(types.FromTime(time.Now().Add(h.quarantineExpiry)))

	h.signRecord(record)

	quarantined, err := model.NewQuarantinedRecord(
		cmd.TenantID,
		cmd.SourceID,
		cmd.SourceType,
		cmd.RawDataID,
		record.ID(),
		cmd.Payload,
		reason,
		reasonDetail,
		validationResult.Anomalies(),
		confidenceScore,
		expiresAt,
		record.IngestSigner(),
	)
	if err != nil {
		return command.ProcessRawDataResult{}
	}

	if err := record.Quarantine(quarantined.ID()); err != nil {
		return command.ProcessRawDataResult{}
	}

	if err := h.recordRepo.Create(ctx, record); err != nil {
		return command.ProcessRawDataResult{}
	}

	if err := h.quarantineRepo.Create(ctx, quarantined); err != nil {
		return command.ProcessRawDataResult{}
	}

	h.updateReliability(ctx, cmd.SourceID, cmd.TenantID, reliability, model.IngestStatusQuarantined)

	_ = h.publisher.PublishRecordQuarantined(ctx, event.NewRecordQuarantined(
		cmd.TenantID,
		cmd.SourceID,
		record.ID(),
		cmd.RawDataID,
		quarantined.ID(),
		cmd.SourceType,
		reason,
		reasonDetail,
		validationResult.Anomalies(),
		confidenceScore.Overall(),
		expiresAt,
	))

	return command.ProcessRawDataResult{
		IngestRecordID:  record.ID(),
		Status:          model.IngestStatusQuarantined.String(),
		EntityType:      types.None[string](),
		EntityID:        types.None[string](),
		EventIDs:        nil,
		QuarantineID:    types.Some(quarantined.ID()),
		RejectionReason: types.None[string](),
		ConfidenceScore: confidenceScore.Overall(),
	}
}

func (h *processRawDataHandler) signRecord(record *model.IngestRecord) {
	if h.signer == nil {
		return
	}

	envelope, err := h.signer.Build(record.Status().String(), map[string]any{
		"record_id":   record.ID().String(),
		"raw_data_id": record.RawDataID(),
		"source_id":   record.SourceID().String(),
		"source_type": record.SourceType(),
	})
	if err != nil {
		return
	}

	record.SetIngestSigner(envelope.ToSignatureInfo())
}

func (h *processRawDataHandler) updateReliability(
	ctx context.Context,
	sourceID types.ID,
	tenantID types.Optional[types.ID],
	reliability *model.SourceReliability,
	status model.IngestStatus,
) {
	if reliability == nil {
		reliability = model.NewSourceReliability(sourceID, tenantID)
	}

	previousScore := reliability.ReliabilityScore()

	switch status {
	case model.IngestStatusAccepted:
		reliability.RecordAccepted()
	case model.IngestStatusRejected:
		reliability.RecordRejected()
	case model.IngestStatusQuarantined:
		reliability.RecordQuarantined()
	}

	_ = h.reliabilityRepo.Upsert(ctx, reliability)

	if previousScore != reliability.ReliabilityScore() {
		_ = h.publisher.PublishSourceReliabilityUpdated(ctx, event.NewSourceReliabilityUpdated(
			tenantID,
			sourceID,
			previousScore,
			reliability.ReliabilityScore(),
			reliability.TotalRecords(),
			reliability.AcceptedRecords(),
			reliability.RejectedRecords(),
			"record_processed",
		))
	}
}

func (h *processRawDataHandler) determineQuarantineReason(validationResult model.ValidationResult, confidenceScore model.ConfidenceScore) model.QuarantineReason {
	if !validationResult.SchemaValid() {
		return model.QuarantineReasonValidationFailed
	}

	if validationResult.HasErrorAnomalies() {
		return model.QuarantineReasonAnomalyDetected
	}

	return model.QuarantineReasonLowConfidence
}

func (h *processRawDataHandler) buildQuarantineReasonDetail(validationResult model.ValidationResult, confidenceScore model.ConfidenceScore) string {
	if !validationResult.SchemaValid() {
		return "schema validation failed"
	}

	if validationResult.HasErrorAnomalies() {
		return "anomalies detected in data"
	}

	return "confidence score below acceptance threshold"
}

func (h *processRawDataHandler) inferEntityType(sourceType string, payload map[string]any) string {
	if _, ok := payload["mmsi"]; ok {
		return "vessel"
	}
	if _, ok := payload["icao"]; ok {
		return "aircraft"
	}
	if _, ok := payload["mag"]; ok {
		return "earthquake"
	}
	if _, ok := payload["lat"]; ok {
		if _, ok := payload["lon"]; ok {
			return "location"
		}
	}
	return "unknown"
}

func (h *processRawDataHandler) inferEntityID(sourceType string, payload map[string]any) string {
	if mmsi, ok := payload["mmsi"]; ok {
		return "mmsi:" + toString(mmsi)
	}
	if icao, ok := payload["icao"]; ok {
		return "icao:" + toString(icao)
	}
	if eqID, ok := payload["earthquake_id"]; ok {
		return "earthquake:" + toString(eqID)
	}
	return "unknown"
}

func (h *processRawDataHandler) mergeAnomalies(a, b []model.Anomaly) []model.Anomaly {
	result := make([]model.Anomaly, 0, len(a)+len(b))
	result = append(result, a...)
	result = append(result, b...)
	return result
}

func (h *processRawDataHandler) extractDID(signer *provenance.SignatureInfo) types.Optional[string] {
	if signer == nil {
		return types.None[string]()
	}
	return types.Some(signer.DID)
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", v)
	}
}
