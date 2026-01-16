package query

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
)

type ListQuarantined struct {
	TenantID   types.Optional[types.ID]
	SourceID   types.Optional[types.ID]
	SourceType types.Optional[string]
	Reason     types.Optional[model.QuarantineReason]
	Resolution types.Optional[model.QuarantineResolution]
	TimeRange  types.Optional[command.TimeRange]
	Limit      int
	Offset     int
}

func (q ListQuarantined) QueryName() string {
	return "ingest.list_quarantined"
}

func DefaultListQuarantined() ListQuarantined {
	return ListQuarantined{
		TenantID:   types.None[types.ID](),
		SourceID:   types.None[types.ID](),
		SourceType: types.None[string](),
		Reason:     types.None[model.QuarantineReason](),
		Resolution: types.None[model.QuarantineResolution](),
		TimeRange:  types.None[command.TimeRange](),
		Limit:      20,
		Offset:     0,
	}
}

func (q ListQuarantined) WithTenantID(tenantID types.ID) ListQuarantined {
	q.TenantID = types.Some(tenantID)
	return q
}

func (q ListQuarantined) WithSourceID(sourceID types.ID) ListQuarantined {
	q.SourceID = types.Some(sourceID)
	return q
}

func (q ListQuarantined) WithSourceType(sourceType string) ListQuarantined {
	q.SourceType = types.Some(sourceType)
	return q
}

func (q ListQuarantined) WithReason(reason model.QuarantineReason) ListQuarantined {
	q.Reason = types.Some(reason)
	return q
}

func (q ListQuarantined) WithResolution(resolution model.QuarantineResolution) ListQuarantined {
	q.Resolution = types.Some(resolution)
	return q
}

func (q ListQuarantined) WithTimeRange(timeRange command.TimeRange) ListQuarantined {
	q.TimeRange = types.Some(timeRange)
	return q
}

func (q ListQuarantined) WithPagination(limit, offset int) ListQuarantined {
	q.Limit = limit
	q.Offset = offset
	return q
}

func (q ListQuarantined) WithPendingOnly() ListQuarantined {
	q.Resolution = types.Some(model.QuarantineResolutionPending)
	return q
}

type ListQuarantinedResult struct {
	Records    []*model.QuarantinedRecord
	TotalCount int64
}

type ListQuarantinedHandler interface {
	Handle(ctx context.Context, qry ListQuarantined) (ListQuarantinedResult, error)
}
