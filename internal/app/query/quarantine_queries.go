package query

import (
	"context"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/query"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

type getQuarantinedHandler struct {
	quarantineRepo repository.QuarantinedRecordRepository
}

func NewGetQuarantinedHandler(quarantineRepo repository.QuarantinedRecordRepository) query.GetQuarantinedHandler {
	return &getQuarantinedHandler{
		quarantineRepo: quarantineRepo,
	}
}

func (h *getQuarantinedHandler) Handle(ctx context.Context, qry query.GetQuarantined) (query.GetQuarantinedResult, error) {
	if qry.ID.IsEmpty() {
		return query.GetQuarantinedResult{}, domainerror.ErrRecordIDRequired
	}

	record, err := h.quarantineRepo.FindByID(ctx, qry.ID)
	if err != nil {
		return query.GetQuarantinedResult{}, err
	}
	if record == nil {
		return query.GetQuarantinedResult{}, domainerror.QuarantineNotFound(qry.ID.String())
	}

	return query.GetQuarantinedResult{Record: record}, nil
}

type getQuarantinedByIngestRecordHandler struct {
	quarantineRepo repository.QuarantinedRecordRepository
}

func NewGetQuarantinedByIngestRecordHandler(quarantineRepo repository.QuarantinedRecordRepository) query.GetQuarantinedByIngestRecordHandler {
	return &getQuarantinedByIngestRecordHandler{
		quarantineRepo: quarantineRepo,
	}
}

func (h *getQuarantinedByIngestRecordHandler) Handle(ctx context.Context, qry query.GetQuarantinedByIngestRecord) (query.GetQuarantinedByIngestRecordResult, error) {
	if qry.IngestRecordID.IsEmpty() {
		return query.GetQuarantinedByIngestRecordResult{}, domainerror.ErrRecordIDRequired
	}

	record, err := h.quarantineRepo.FindByIngestRecordID(ctx, qry.IngestRecordID)
	if err != nil {
		return query.GetQuarantinedByIngestRecordResult{}, err
	}
	if record == nil {
		return query.GetQuarantinedByIngestRecordResult{}, domainerror.QuarantineNotFound(qry.IngestRecordID.String())
	}

	return query.GetQuarantinedByIngestRecordResult{Record: record}, nil
}

type listQuarantinedHandler struct {
	quarantineRepo repository.QuarantinedRecordRepository
}

func NewListQuarantinedHandler(quarantineRepo repository.QuarantinedRecordRepository) query.ListQuarantinedHandler {
	return &listQuarantinedHandler{
		quarantineRepo: quarantineRepo,
	}
}

func (h *listQuarantinedHandler) Handle(ctx context.Context, qry query.ListQuarantined) (query.ListQuarantinedResult, error) {
	params := repository.DefaultListQuarantinedRecordsParams().
		WithPagination(qry.Limit, qry.Offset)

	if qry.TenantID.IsPresent() {
		params = params.WithTenantID(qry.TenantID.MustGet())
	}

	if qry.SourceID.IsPresent() {
		params = params.WithSourceID(qry.SourceID.MustGet())
	}

	if qry.SourceType.IsPresent() {
		params = params.WithSourceType(qry.SourceType.MustGet())
	}

	if qry.Reason.IsPresent() {
		params = params.WithReason(qry.Reason.MustGet())
	}

	if qry.Resolution.IsPresent() {
		params = params.WithResolution(qry.Resolution.MustGet())
	}

	if qry.TimeRange.IsPresent() {
		timeRange := qry.TimeRange.MustGet()
		if timeRange.Start.IsPresent() && timeRange.End.IsPresent() {
			params = params.WithTimeRange(timeRange.Start.MustGet(), timeRange.End.MustGet())
		}
	}

	records, err := h.quarantineRepo.List(ctx, params)
	if err != nil {
		return query.ListQuarantinedResult{}, err
	}

	count, err := h.quarantineRepo.Count(ctx, params)
	if err != nil {
		return query.ListQuarantinedResult{}, err
	}

	return query.ListQuarantinedResult{
		Records:    records,
		TotalCount: count,
	}, nil
}
