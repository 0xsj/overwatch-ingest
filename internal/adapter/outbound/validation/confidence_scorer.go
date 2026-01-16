package validation

import (
	"context"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/validation"
)

const confidenceScorerVersion = "1.0.0"

type confidenceScorer struct {
	registry *SourceTypeRegistry
}

func NewConfidenceScorer(registry *SourceTypeRegistry) validation.ConfidenceScorer {
	if registry == nil {
		registry = NewSourceTypeRegistry()
	}
	return &confidenceScorer{
		registry: registry,
	}
}

func (s *confidenceScorer) Score(ctx context.Context, input validation.ScoringInput) model.ConfidenceScore {
	factors := make([]model.ConfidenceFactor, 0, 5)

	sourceReliability := s.scoreSourceReliability(input)
	factors = append(factors, sourceReliability)

	dataCompleteness := s.scoreDataCompleteness(input)
	factors = append(factors, dataCompleteness)

	temporalFreshness := s.scoreTemporalFreshness(input)
	factors = append(factors, temporalFreshness)

	signatureTrust := s.scoreSignatureTrust(input)
	factors = append(factors, signatureTrust)

	anomalyPenalty := s.scoreAnomalyPenalty(input)
	factors = append(factors, anomalyPenalty)

	return model.NewConfidenceScore(
		sourceReliability.Score(),
		dataCompleteness.Score(),
		temporalFreshness.Score(),
		signatureTrust.Score(),
		factors,
	)
}

func (s *confidenceScorer) SupportsSourceType(sourceType string) bool {
	return s.registry.Supports(sourceType)
}

func (s *confidenceScorer) Version() string {
	return confidenceScorerVersion
}

func (s *confidenceScorer) scoreSourceReliability(input validation.ScoringInput) model.ConfidenceFactor {
	score := 0.5
	reason := "default source reliability"

	if input.SourceReliability != nil {
		score = input.SourceReliability.ReliabilityScore()
		reason = "based on historical source performance"
	} else {
		cfg := s.registry.GetOrDefault(input.SourceType)
		score = cfg.BaseReliability
		reason = "based on source type default reliability"
	}

	return model.NewConfidenceFactor(
		"source_reliability",
		score,
		0.25,
		reason,
	)
}

func (s *confidenceScorer) scoreDataCompleteness(input validation.ScoringInput) model.ConfidenceFactor {
	cfg := s.registry.GetOrDefault(input.SourceType)

	totalRequired := len(cfg.RequiredFields)
	totalOptional := len(cfg.OptionalFields)

	if totalRequired == 0 && totalOptional == 0 {
		return model.NewConfidenceFactor(
			"data_completeness",
			0.5,
			0.25,
			"no schema defined for source type",
		)
	}

	presentRequired := 0
	presentOptional := 0

	presentFields := make(map[string]bool)
	for _, f := range input.ValidationResult.FieldsPresent() {
		presentFields[f] = true
	}

	for _, f := range cfg.RequiredFields {
		if presentFields[f] {
			presentRequired++
		}
	}

	for _, f := range cfg.OptionalFields {
		if presentFields[f] {
			presentOptional++
		}
	}

	requiredScore := 1.0
	if totalRequired > 0 {
		requiredScore = float64(presentRequired) / float64(totalRequired)
	}

	optionalScore := 0.5
	if totalOptional > 0 {
		optionalScore = float64(presentOptional) / float64(totalOptional)
	}

	score := (requiredScore * 0.7) + (optionalScore * 0.3)

	return model.NewConfidenceFactor(
		"data_completeness",
		score,
		0.25,
		"based on required and optional field presence",
	)
}

func (s *confidenceScorer) scoreTemporalFreshness(input validation.ScoringInput) model.ConfidenceFactor {
	score := 0.5
	reason := "no timestamp information available"

	if input.SourceTimestamp.IsPresent() {
		sourceTs := input.SourceTimestamp.MustGet()
		collectedTs := input.CollectedAt

		diff := collectedTs.Time().Sub(sourceTs.Time())
		if diff < 0 {
			diff = -diff
		}

		switch {
		case diff.Minutes() < 5:
			score = 1.0
			reason = "data is very fresh (< 5 minutes)"
		case diff.Hours() < 1:
			score = 0.9
			reason = "data is fresh (< 1 hour)"
		case diff.Hours() < 6:
			score = 0.7
			reason = "data is moderately fresh (< 6 hours)"
		case diff.Hours() < 24:
			score = 0.5
			reason = "data is aging (< 24 hours)"
		case diff.Hours() < 168:
			score = 0.3
			reason = "data is stale (< 1 week)"
		default:
			score = 0.1
			reason = "data is very old (> 1 week)"
		}
	}

	return model.NewConfidenceFactor(
		"temporal_freshness",
		score,
		0.20,
		reason,
	)
}

func (s *confidenceScorer) scoreSignatureTrust(input validation.ScoringInput) model.ConfidenceFactor {
	score := 0.5
	reason := "no cryptographic signatures present"

	hasSignatures := input.HasSourceSigner || input.HasCollectorSigner
	verifiedCount := 0
	totalSignatures := 0

	if input.HasSourceSigner {
		totalSignatures++
		if input.SourceSignatureVerified {
			verifiedCount++
		}
	}

	if input.HasCollectorSigner {
		totalSignatures++
		if input.CollectorSignatureVerified {
			verifiedCount++
		}
	}

	if !hasSignatures {
		score = 0.5
		reason = "no cryptographic signatures present"
	} else if verifiedCount == totalSignatures {
		score = 1.0
		reason = "all signatures verified"
	} else if verifiedCount > 0 {
		score = 0.7
		reason = "partial signature verification"
	} else {
		score = 0.2
		reason = "signature verification failed"
	}

	return model.NewConfidenceFactor(
		"signature_trust",
		score,
		0.15,
		reason,
	)
}

func (s *confidenceScorer) scoreAnomalyPenalty(input validation.ScoringInput) model.ConfidenceFactor {
	if len(input.Anomalies) == 0 {
		return model.NewConfidenceFactor(
			"anomaly_penalty",
			1.0,
			0.15,
			"no anomalies detected",
		)
	}

	penalty := 0.0

	for _, anomaly := range input.Anomalies {
		switch anomaly.Severity() {
		case model.AnomalySeverityCritical:
			penalty += 0.4
		case model.AnomalySeverityError:
			penalty += 0.25
		case model.AnomalySeverityWarning:
			penalty += 0.1
		case model.AnomalySeverityInfo:
			penalty += 0.02
		}
	}

	if penalty > 0.9 {
		penalty = 0.9
	}

	score := 1.0 - penalty

	return model.NewConfidenceFactor(
		"anomaly_penalty",
		score,
		0.15,
		"penalty applied for detected anomalies",
	)
}
