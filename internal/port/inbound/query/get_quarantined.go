package query

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type GetQuarantined struct {
	ID types.ID
}

func (q GetQuarantined) QueryName() string {
	return "ingest.get_quarantined"
}

type GetQuarantinedResult struct {
	Record *model.QuarantinedRecord
}

type GetQuarantinedHandler interface {
	Handle(ctx context.Context, qry GetQuarantined) (GetQuarantinedResult, error)
}

type GetQuarantinedByIngestRecord struct {
	IngestRecordID types.ID
}

func (q GetQuarantinedByIngestRecord) QueryName() string {
	return "ingest.get_quarantined_by_ingest_record"
}

type GetQuarantinedByIngestRecordResult struct {
	Record *model.QuarantinedRecord
}

type GetQuarantinedByIngestRecordHandler interface {
	Handle(ctx context.Context, qry GetQuarantinedByIngestRecord) (GetQuarantinedByIngestRecordResult, error)
}
