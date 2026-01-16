package repository

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type QuarantinedRecordRepository interface {
	Create(ctx context.Context, record *model.QuarantinedRecord) error
	Update(ctx context.Context, record *model.QuarantinedRecord) error
	FindByID(ctx context.Context, id types.ID) (*model.QuarantinedRecord, error)
	FindByIngestRecordID(ctx context.Context, ingestRecordID types.ID) (*model.QuarantinedRecord, error)
	List(ctx context.Context, params ListQuarantinedRecordsParams) ([]*model.QuarantinedRecord, error)
	Count(ctx context.Context, params ListQuarantinedRecordsParams) (int64, error)
	FindExpired(ctx context.Context, limit int) ([]*model.QuarantinedRecord, error)
	CountPending(ctx context.Context, tenantID types.Optional[types.ID]) (int64, error)
}

type ListQuarantinedRecordsParams struct {
	Limit      int
	Offset     int
	TenantID   types.Optional[types.ID]
	SourceID   types.Optional[types.ID]
	SourceType types.Optional[string]
	Reason     types.Optional[model.QuarantineReason]
	Resolution types.Optional[model.QuarantineResolution]
	StartTime  types.Optional[types.Timestamp]
	EndTime    types.Optional[types.Timestamp]
	SortBy     QuarantinedRecordSortField
	SortOrder  SortOrder
}

type QuarantinedRecordSortField string

const (
	QuarantinedRecordSortFieldQuarantinedAt QuarantinedRecordSortField = "quarantined_at"
	QuarantinedRecordSortFieldExpiresAt     QuarantinedRecordSortField = "expires_at"
	QuarantinedRecordSortFieldPriority      QuarantinedRecordSortField = "priority"
)

func DefaultListQuarantinedRecordsParams() ListQuarantinedRecordsParams {
	return ListQuarantinedRecordsParams{
		Limit:      20,
		Offset:     0,
		TenantID:   types.None[types.ID](),
		SourceID:   types.None[types.ID](),
		SourceType: types.None[string](),
		Reason:     types.None[model.QuarantineReason](),
		Resolution: types.None[model.QuarantineResolution](),
		StartTime:  types.None[types.Timestamp](),
		EndTime:    types.None[types.Timestamp](),
		SortBy:     QuarantinedRecordSortFieldQuarantinedAt,
		SortOrder:  SortOrderDesc,
	}
}

func (p ListQuarantinedRecordsParams) WithTenantID(tenantID types.ID) ListQuarantinedRecordsParams {
	p.TenantID = types.Some(tenantID)
	return p
}

func (p ListQuarantinedRecordsParams) WithSourceID(sourceID types.ID) ListQuarantinedRecordsParams {
	p.SourceID = types.Some(sourceID)
	return p
}

func (p ListQuarantinedRecordsParams) WithSourceType(sourceType string) ListQuarantinedRecordsParams {
	p.SourceType = types.Some(sourceType)
	return p
}

func (p ListQuarantinedRecordsParams) WithReason(reason model.QuarantineReason) ListQuarantinedRecordsParams {
	p.Reason = types.Some(reason)
	return p
}

func (p ListQuarantinedRecordsParams) WithResolution(resolution model.QuarantineResolution) ListQuarantinedRecordsParams {
	p.Resolution = types.Some(resolution)
	return p
}

func (p ListQuarantinedRecordsParams) WithPendingOnly() ListQuarantinedRecordsParams {
	p.Resolution = types.Some(model.QuarantineResolutionPending)
	return p
}

func (p ListQuarantinedRecordsParams) WithTimeRange(start, end types.Timestamp) ListQuarantinedRecordsParams {
	p.StartTime = types.Some(start)
	p.EndTime = types.Some(end)
	return p
}

func (p ListQuarantinedRecordsParams) WithPagination(limit, offset int) ListQuarantinedRecordsParams {
	p.Limit = limit
	p.Offset = offset
	return p
}

func (p ListQuarantinedRecordsParams) WithSort(field QuarantinedRecordSortField, order SortOrder) ListQuarantinedRecordsParams {
	p.SortBy = field
	p.SortOrder = order
	return p
}
