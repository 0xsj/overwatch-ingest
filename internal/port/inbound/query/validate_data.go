package query

import (
	"context"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

// ValidateData is a preview request that runs the full validation pipeline
// without persisting any records. Used for testing source configs and data.
type ValidateData struct {
	SourceID   string
	SourceType string
	Payload    map[string]any
}

// ValidateDataResult holds the validation preview output.
type ValidateDataResult struct {
	Validation      model.ValidationResult
	Confidence      model.ConfidenceScore
	PredictedStatus model.IngestStatus
}

// ValidateDataHandler defines the interface for handling ValidateData queries.
type ValidateDataHandler interface {
	Handle(ctx context.Context, q ValidateData) (*ValidateDataResult, error)
}
