package query

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type GetRecord struct {
	ID types.ID
}

func (q GetRecord) QueryName() string {
	return "ingest.get_record"
}

type GetRecordResult struct {
	Record *model.IngestRecord
}

type GetRecordHandler interface {
	Handle(ctx context.Context, qry GetRecord) (GetRecordResult, error)
}

type GetRecordByRawData struct {
	RawDataID string
}

func (q GetRecordByRawData) QueryName() string {
	return "ingest.get_record_by_raw_data"
}

type GetRecordByRawDataResult struct {
	Record *model.IngestRecord
}

type GetRecordByRawDataHandler interface {
	Handle(ctx context.Context, qry GetRecordByRawData) (GetRecordByRawDataResult, error)
}
