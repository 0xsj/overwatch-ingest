package repository

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type IngestRecordRepository interface {
	Create(ctx context.Context, record *model.IngestRecord) error
	Update(ctx context.Context, record *model.IngestRecord) error
	FindByID(ctx context.Context, id types.ID) (*model.IngestRecord, error)
	FindByRawDataID(ctx context.Context, rawDataID string) (*model.IngestRecord, error)
	List(ctx context.Context, params ListIngestRecordsParams) ([]*model.IngestRecord, error)
	Count(ctx context.Context, params ListIngestRecordsParams) (int64, error)
	ExistsByRawDataID(ctx context.Context, rawDataID string) (bool, error)
}

type ListIngestRecordsParams struct {
	Limit      int
	Offset     int
	TenantID   types.Optional[types.ID]
	SourceID   types.Optional[types.ID]
	SourceType types.Optional[string]
	Status     types.Optional[model.IngestStatus]
	EntityType types.Optional[string]
	StartTime  types.Optional[types.Timestamp]
	EndTime    types.Optional[types.Timestamp]
	SortBy     IngestRecordSortField
	SortOrder  SortOrder
}

type IngestRecordSortField string

const (
	IngestRecordSortFieldReceivedAt  IngestRecordSortField = "received_at"
	IngestRecordSortFieldProcessedAt IngestRecordSortField = "processed_at"
	IngestRecordSortFieldConfidence  IngestRecordSortField = "confidence"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func DefaultListIngestRecordsParams() ListIngestRecordsParams {
	return ListIngestRecordsParams{
		Limit:      20,
		Offset:     0,
		TenantID:   types.None[types.ID](),
		SourceID:   types.None[types.ID](),
		SourceType: types.None[string](),
		Status:     types.None[model.IngestStatus](),
		EntityType: types.None[string](),
		StartTime:  types.None[types.Timestamp](),
		EndTime:    types.None[types.Timestamp](),
		SortBy:     IngestRecordSortFieldReceivedAt,
		SortOrder:  SortOrderDesc,
	}
}

func (p ListIngestRecordsParams) WithTenantID(tenantID types.ID) ListIngestRecordsParams {
	p.TenantID = types.Some(tenantID)
	return p
}

func (p ListIngestRecordsParams) WithSourceID(sourceID types.ID) ListIngestRecordsParams {
	p.SourceID = types.Some(sourceID)
	return p
}

func (p ListIngestRecordsParams) WithSourceType(sourceType string) ListIngestRecordsParams {
	p.SourceType = types.Some(sourceType)
	return p
}

func (p ListIngestRecordsParams) WithStatus(status model.IngestStatus) ListIngestRecordsParams {
	p.Status = types.Some(status)
	return p
}

func (p ListIngestRecordsParams) WithEntityType(entityType string) ListIngestRecordsParams {
	p.EntityType = types.Some(entityType)
	return p
}

func (p ListIngestRecordsParams) WithTimeRange(start, end types.Timestamp) ListIngestRecordsParams {
	p.StartTime = types.Some(start)
	p.EndTime = types.Some(end)
	return p
}

func (p ListIngestRecordsParams) WithPagination(limit, offset int) ListIngestRecordsParams {
	p.Limit = limit
	p.Offset = offset
	return p
}

func (p ListIngestRecordsParams) WithSort(field IngestRecordSortField, order SortOrder) ListIngestRecordsParams {
	p.SortBy = field
	p.SortOrder = order
	return p
}
