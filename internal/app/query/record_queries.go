package query

import (
	"context"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/query"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

type getRecordHandler struct {
	recordRepo repository.IngestRecordRepository
}

func NewGetRecordHandler(recordRepo repository.IngestRecordRepository) query.GetRecordHandler {
	return &getRecordHandler{
		recordRepo: recordRepo,
	}
}

func (h *getRecordHandler) Handle(ctx context.Context, qry query.GetRecord) (query.GetRecordResult, error) {
	if qry.ID.IsEmpty() {
		return query.GetRecordResult{}, domainerror.ErrRecordIDRequired
	}

	record, err := h.recordRepo.FindByID(ctx, qry.ID)
	if err != nil {
		return query.GetRecordResult{}, err
	}
	if record == nil {
		return query.GetRecordResult{}, domainerror.RecordNotFound(qry.ID.String())
	}

	return query.GetRecordResult{Record: record}, nil
}

type getRecordByRawDataHandler struct {
	recordRepo repository.IngestRecordRepository
}

func NewGetRecordByRawDataHandler(recordRepo repository.IngestRecordRepository) query.GetRecordByRawDataHandler {
	return &getRecordByRawDataHandler{
		recordRepo: recordRepo,
	}
}

func (h *getRecordByRawDataHandler) Handle(ctx context.Context, qry query.GetRecordByRawData) (query.GetRecordByRawDataResult, error) {
	if qry.RawDataID == "" {
		return query.GetRecordByRawDataResult{}, domainerror.ErrRawDataIDRequired
	}

	record, err := h.recordRepo.FindByRawDataID(ctx, qry.RawDataID)
	if err != nil {
		return query.GetRecordByRawDataResult{}, err
	}
	if record == nil {
		return query.GetRecordByRawDataResult{}, domainerror.RecordNotFound(qry.RawDataID)
	}

	return query.GetRecordByRawDataResult{Record: record}, nil
}

type listRecordsHandler struct {
	recordRepo repository.IngestRecordRepository
}

func NewListRecordsHandler(recordRepo repository.IngestRecordRepository) query.ListRecordsHandler {
	return &listRecordsHandler{
		recordRepo: recordRepo,
	}
}

func (h *listRecordsHandler) Handle(ctx context.Context, qry query.ListRecords) (query.ListRecordsResult, error) {
	params := repository.DefaultListIngestRecordsParams().
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

	if qry.Status.IsPresent() {
		params = params.WithStatus(qry.Status.MustGet())
	}

	if qry.EntityType.IsPresent() {
		params = params.WithEntityType(qry.EntityType.MustGet())
	}

	if qry.TimeRange.IsPresent() {
		timeRange := qry.TimeRange.MustGet()
		if timeRange.Start.IsPresent() && timeRange.End.IsPresent() {
			params = params.WithTimeRange(timeRange.Start.MustGet(), timeRange.End.MustGet())
		}
	}

	records, err := h.recordRepo.List(ctx, params)
	if err != nil {
		return query.ListRecordsResult{}, err
	}

	count, err := h.recordRepo.Count(ctx, params)
	if err != nil {
		return query.ListRecordsResult{}, err
	}

	return query.ListRecordsResult{
		Records:    records,
		TotalCount: count,
	}, nil
}
