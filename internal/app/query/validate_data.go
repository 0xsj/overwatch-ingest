package query

import (
	"context"
	"fmt"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/app/service"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/query"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

// validateDataHandler runs the validation pipeline without persisting.
type validateDataHandler struct {
	pipeline       service.ValidationPipeline
	reliabilityRepo repository.SourceReliabilityRepository
}

func NewValidateDataHandler(
	pipeline service.ValidationPipeline,
	reliabilityRepo repository.SourceReliabilityRepository,
) query.ValidateDataHandler {
	return &validateDataHandler{
		pipeline:        pipeline,
		reliabilityRepo: reliabilityRepo,
	}
}

func (h *validateDataHandler) Handle(ctx context.Context, q query.ValidateData) (*query.ValidateDataResult, error) {
	if q.SourceType == "" {
		return nil, fmt.Errorf("source_type is required")
	}
	if q.Payload == nil {
		return nil, fmt.Errorf("data payload is required")
	}

	sourceID := types.ID(q.SourceID)

	// Look up source reliability (optional — may not exist for new sources).
	var reliability *model.SourceReliability
	if q.SourceID != "" {
		reliability, _ = h.reliabilityRepo.FindBySourceID(ctx, sourceID)
	}

	output := h.pipeline.Validate(ctx, service.ValidationInput{
		SourceID:   sourceID,
		SourceType: q.SourceType,
		Payload:    q.Payload,
		CollectedAt: types.Now(),
		Reliability: reliability,
	})

	var predictedStatus model.IngestStatus
	switch {
	case output.ShouldReject:
		predictedStatus = model.IngestStatusRejected
	case output.ShouldQuarantine:
		predictedStatus = model.IngestStatusQuarantined
	default:
		predictedStatus = model.IngestStatusAccepted
	}

	return &query.ValidateDataResult{
		Validation:      output.ValidationResult,
		Confidence:      output.ConfidenceScore,
		PredictedStatus: predictedStatus,
	}, nil
}
