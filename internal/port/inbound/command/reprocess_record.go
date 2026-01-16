package command

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type ReprocessRecord struct {
	IngestRecordID types.Optional[types.ID]
	RawDataID      types.Optional[string]
}

func (c ReprocessRecord) CommandName() string {
	return "ingest.reprocess_record"
}

type ReprocessRecordResult struct {
	IngestRecord *model.IngestRecord
}

type ReprocessRecordHandler interface {
	Handle(ctx context.Context, cmd ReprocessRecord) (ReprocessRecordResult, error)
}

type ReprocessBySource struct {
	SourceID     types.ID
	TimeRange    types.Optional[TimeRange]
	StatusFilter types.Optional[model.IngestStatus]
}

func (c ReprocessBySource) CommandName() string {
	return "ingest.reprocess_by_source"
}

type TimeRange struct {
	Start types.Optional[types.Timestamp]
	End   types.Optional[types.Timestamp]
}

type ReprocessBySourceResult struct {
	RecordsQueued int
	BatchID       string
}

type ReprocessBySourceHandler interface {
	Handle(ctx context.Context, cmd ReprocessBySource) (ReprocessBySourceResult, error)
}
