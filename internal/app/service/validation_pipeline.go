package service

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/validation"
)

type ValidationPipeline interface {
	Validate(ctx context.Context, input ValidationInput) ValidationOutput
}

type ValidationInput struct {
	SourceID        types.ID
	SourceType      string
	TenantID        types.Optional[types.ID]
	Payload         map[string]any
	SourceTimestamp types.Optional[types.Timestamp]
	CollectedAt     types.Timestamp
	Reliability     *model.SourceReliability
}

type ValidationOutput struct {
	ValidationResult model.ValidationResult
	Anomalies        []model.Anomaly
	ConfidenceScore  model.ConfidenceScore
	ShouldAccept     bool
	ShouldReject     bool
	ShouldQuarantine bool
}

type validationPipeline struct {
	schemaValidator  validation.SchemaValidator
	anomalyDetector  validation.AnomalyDetector
	confidenceScorer validation.ConfidenceScorer
	thresholds       validation.ScoringThresholds
}

func NewValidationPipeline(
	schemaValidator validation.SchemaValidator,
	anomalyDetector validation.AnomalyDetector,
	confidenceScorer validation.ConfidenceScorer,
	thresholds validation.ScoringThresholds,
) ValidationPipeline {
	return &validationPipeline{
		schemaValidator:  schemaValidator,
		anomalyDetector:  anomalyDetector,
		confidenceScorer: confidenceScorer,
		thresholds:       thresholds,
	}
}

func (p *validationPipeline) Validate(ctx context.Context, input ValidationInput) ValidationOutput {
	schemaResult := p.schemaValidator.Validate(ctx, input.SourceType, input.Payload)
	validationResult := schemaResult.ToValidationResult(p.schemaValidator.Version())

	anomalyResult := p.anomalyDetector.Detect(ctx, input.SourceType, input.Payload, validation.DetectionMetadata{
		SourceID:        input.SourceID,
		TenantID:        input.TenantID,
		SourceTimestamp: input.SourceTimestamp,
		CollectedAt:     input.CollectedAt,
	})

	allAnomalies := mergeAnomalies(validationResult.Anomalies(), anomalyResult.Anomalies)

	validationResult = model.NewValidationResult(
		validationResult.SchemaValid(),
		validationResult.FieldsPresent(),
		validationResult.FieldsMissing(),
		allAnomalies,
		validationResult.ValidatorVersion(),
	)

	confidenceScore := p.confidenceScorer.Score(ctx, validation.ScoringInput{
		SourceID:          input.SourceID,
		SourceType:        input.SourceType,
		TenantID:          input.TenantID,
		Payload:           input.Payload,
		ValidationResult:  validationResult,
		Anomalies:         allAnomalies,
		SourceReliability: input.Reliability,
		SourceTimestamp:   input.SourceTimestamp,
		CollectedAt:       input.CollectedAt,
	})

	shouldReject := validationResult.ShouldReject() ||
		anomalyResult.ShouldReject ||
		p.thresholds.ShouldReject(confidenceScore)

	shouldQuarantine := !shouldReject && (validationResult.ShouldQuarantine() ||
		anomalyResult.ShouldQuarantine ||
		p.thresholds.ShouldQuarantine(confidenceScore))

	shouldAccept := !shouldReject && !shouldQuarantine

	return ValidationOutput{
		ValidationResult: validationResult,
		Anomalies:        allAnomalies,
		ConfidenceScore:  confidenceScore,
		ShouldAccept:     shouldAccept,
		ShouldReject:     shouldReject,
		ShouldQuarantine: shouldQuarantine,
	}
}

func mergeAnomalies(a, b []model.Anomaly) []model.Anomaly {
	result := make([]model.Anomaly, 0, len(a)+len(b))
	result = append(result, a...)
	result = append(result, b...)
	return result
}
