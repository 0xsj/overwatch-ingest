package query

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
)

type ListRecords struct {
	TenantID   types.Optional[types.ID]
	SourceID   types.Optional[types.ID]
	SourceType types.Optional[string]
	Status     types.Optional[model.IngestStatus]
	EntityType types.Optional[string]
	TimeRange  types.Optional[command.TimeRange]
	Limit      int
	Offset     int
}

func (q ListRecords) QueryName() string {
	return "ingest.list_records"
}

func DefaultListRecords() ListRecords {
	return ListRecords{
		TenantID:   types.None[types.ID](),
		SourceID:   types.None[types.ID](),
		SourceType: types.None[string](),
		Status:     types.None[model.IngestStatus](),
		EntityType: types.None[string](),
		TimeRange:  types.None[command.TimeRange](),
		Limit:      20,
		Offset:     0,
	}
}

func (q ListRecords) WithTenantID(tenantID types.ID) ListRecords {
	q.TenantID = types.Some(tenantID)
	return q
}

func (q ListRecords) WithSourceID(sourceID types.ID) ListRecords {
	q.SourceID = types.Some(sourceID)
	return q
}

func (q ListRecords) WithSourceType(sourceType string) ListRecords {
	q.SourceType = types.Some(sourceType)
	return q
}

func (q ListRecords) WithStatus(status model.IngestStatus) ListRecords {
	q.Status = types.Some(status)
	return q
}

func (q ListRecords) WithEntityType(entityType string) ListRecords {
	q.EntityType = types.Some(entityType)
	return q
}

func (q ListRecords) WithTimeRange(timeRange command.TimeRange) ListRecords {
	q.TimeRange = types.Some(timeRange)
	return q
}

func (q ListRecords) WithPagination(limit, offset int) ListRecords {
	q.Limit = limit
	q.Offset = offset
	return q
}

type ListRecordsResult struct {
	Records    []*model.IngestRecord
	TotalCount int64
}

type ListRecordsHandler interface {
	Handle(ctx context.Context, qry ListRecords) (ListRecordsResult, error)
}
