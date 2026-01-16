package validation

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type ConfidenceScorer interface {
	Score(ctx context.Context, input ScoringInput) model.ConfidenceScore
	SupportsSourceType(sourceType string) bool
	Version() string
}

type ScoringInput struct {
	SourceID                   types.ID
	SourceType                 string
	TenantID                   types.Optional[types.ID]
	Payload                    map[string]any
	ValidationResult           model.ValidationResult
	Anomalies                  []model.Anomaly
	SourceReliability          *model.SourceReliability
	SourceTimestamp            types.Optional[types.Timestamp]
	CollectedAt                types.Timestamp
	HasSourceSigner            bool
	SourceSignatureVerified    bool
	HasCollectorSigner         bool
	CollectorSignatureVerified bool
}

type ScoringThresholds struct {
	AcceptThreshold float64
	RejectThreshold float64
}

func DefaultScoringThresholds() ScoringThresholds {
	return ScoringThresholds{
		AcceptThreshold: 0.7,
		RejectThreshold: 0.3,
	}
}

func (t ScoringThresholds) Evaluate(score model.ConfidenceScore) model.IngestStatus {
	if score.Overall() >= t.AcceptThreshold {
		return model.IngestStatusAccepted
	}
	if score.Overall() < t.RejectThreshold {
		return model.IngestStatusRejected
	}
	return model.IngestStatusQuarantined
}

func (t ScoringThresholds) ShouldAccept(score model.ConfidenceScore) bool {
	return score.Overall() >= t.AcceptThreshold
}

func (t ScoringThresholds) ShouldReject(score model.ConfidenceScore) bool {
	return score.Overall() < t.RejectThreshold
}

func (t ScoringThresholds) ShouldQuarantine(score model.ConfidenceScore) bool {
	return score.Overall() >= t.RejectThreshold && score.Overall() < t.AcceptThreshold
}
