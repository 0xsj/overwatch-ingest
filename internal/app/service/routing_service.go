package service

import (
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type RoutingDecision string

const (
	RoutingDecisionAccept     RoutingDecision = "accept"
	RoutingDecisionReject     RoutingDecision = "reject"
	RoutingDecisionQuarantine RoutingDecision = "quarantine"
)

type RoutingService interface {
	Route(input RoutingInput) RoutingOutput
}

type RoutingInput struct {
	ValidationResult  model.ValidationResult
	ConfidenceScore   model.ConfidenceScore
	Anomalies         []model.Anomaly
	SourceReliability *model.SourceReliability
}

type RoutingOutput struct {
	Decision         RoutingDecision
	Reason           string
	QuarantineReason model.QuarantineReason
}

type routingService struct {
	acceptThreshold                float64
	rejectThreshold                float64
	reliabilityQuarantineThreshold float64
	minReliabilityRecords          int64
}

type RoutingConfig struct {
	AcceptThreshold                float64
	RejectThreshold                float64
	ReliabilityQuarantineThreshold float64
	MinReliabilityRecords          int64
}

func DefaultRoutingConfig() RoutingConfig {
	return RoutingConfig{
		AcceptThreshold:                0.7,
		RejectThreshold:                0.3,
		ReliabilityQuarantineThreshold: 0.4,
		MinReliabilityRecords:          10,
	}
}

func NewRoutingService(config RoutingConfig) RoutingService {
	return &routingService{
		acceptThreshold:                config.AcceptThreshold,
		rejectThreshold:                config.RejectThreshold,
		reliabilityQuarantineThreshold: config.ReliabilityQuarantineThreshold,
		minReliabilityRecords:          config.MinReliabilityRecords,
	}
}

func (s *routingService) Route(input RoutingInput) RoutingOutput {
	if s.shouldRejectValidation(input) {
		return RoutingOutput{
			Decision: RoutingDecisionReject,
			Reason:   "validation failed with critical errors",
		}
	}

	if s.shouldRejectConfidence(input) {
		return RoutingOutput{
			Decision: RoutingDecisionReject,
			Reason:   "confidence score below rejection threshold",
		}
	}

	if s.shouldQuarantineValidation(input) {
		return RoutingOutput{
			Decision:         RoutingDecisionQuarantine,
			Reason:           "validation requires review",
			QuarantineReason: s.determineQuarantineReason(input),
		}
	}

	if s.shouldQuarantineConfidence(input) {
		return RoutingOutput{
			Decision:         RoutingDecisionQuarantine,
			Reason:           "confidence score requires review",
			QuarantineReason: model.QuarantineReasonLowConfidence,
		}
	}

	if s.shouldQuarantineReliability(input) {
		return RoutingOutput{
			Decision:         RoutingDecisionQuarantine,
			Reason:           "source reliability requires review",
			QuarantineReason: model.QuarantineReasonManualReview,
		}
	}

	return RoutingOutput{
		Decision: RoutingDecisionAccept,
		Reason:   "validation passed",
	}
}

func (s *routingService) shouldRejectValidation(input RoutingInput) bool {
	if input.ValidationResult.ShouldReject() {
		return true
	}

	for _, anomaly := range input.Anomalies {
		if anomaly.ShouldReject() {
			return true
		}
	}

	return false
}

func (s *routingService) shouldRejectConfidence(input RoutingInput) bool {
	return input.ConfidenceScore.Overall() < s.rejectThreshold
}

func (s *routingService) shouldQuarantineValidation(input RoutingInput) bool {
	if input.ValidationResult.ShouldQuarantine() {
		return true
	}

	for _, anomaly := range input.Anomalies {
		if anomaly.ShouldQuarantine() {
			return true
		}
	}

	return false
}

func (s *routingService) shouldQuarantineConfidence(input RoutingInput) bool {
	score := input.ConfidenceScore.Overall()
	return score >= s.rejectThreshold && score < s.acceptThreshold
}

func (s *routingService) shouldQuarantineReliability(input RoutingInput) bool {
	if input.SourceReliability == nil {
		return false
	}

	if !input.SourceReliability.HasSufficientData(s.minReliabilityRecords) {
		return false
	}

	return input.SourceReliability.IsUnreliable(s.reliabilityQuarantineThreshold)
}

func (s *routingService) determineQuarantineReason(input RoutingInput) model.QuarantineReason {
	if !input.ValidationResult.SchemaValid() {
		return model.QuarantineReasonValidationFailed
	}

	if input.ValidationResult.HasErrorAnomalies() {
		return model.QuarantineReasonAnomalyDetected
	}

	for _, anomaly := range input.Anomalies {
		if anomaly.Type() == model.AnomalyTypeDuplicate {
			return model.QuarantineReasonDuplicateSuspected
		}
		if anomaly.Type() == model.AnomalyTypeSuspicious {
			return model.QuarantineReasonAnomalyDetected
		}
	}

	return model.QuarantineReasonManualReview
}
