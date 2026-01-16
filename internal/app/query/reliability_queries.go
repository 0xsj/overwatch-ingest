package query

import (
	"context"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/query"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

type getSourceReliabilityHandler struct {
	reliabilityRepo repository.SourceReliabilityRepository
}

func NewGetSourceReliabilityHandler(reliabilityRepo repository.SourceReliabilityRepository) query.GetSourceReliabilityHandler {
	return &getSourceReliabilityHandler{
		reliabilityRepo: reliabilityRepo,
	}
}

func (h *getSourceReliabilityHandler) Handle(ctx context.Context, qry query.GetSourceReliability) (query.GetSourceReliabilityResult, error) {
	if qry.SourceID.IsEmpty() {
		return query.GetSourceReliabilityResult{}, domainerror.ErrSourceIDRequired
	}

	reliability, err := h.reliabilityRepo.FindBySourceID(ctx, qry.SourceID)
	if err != nil {
		return query.GetSourceReliabilityResult{}, err
	}
	if reliability == nil {
		return query.GetSourceReliabilityResult{}, domainerror.SourceReliabilityNotFound(qry.SourceID.String())
	}

	return query.GetSourceReliabilityResult{Reliability: reliability}, nil
}

type listSourceReliabilityHandler struct {
	reliabilityRepo repository.SourceReliabilityRepository
}

func NewListSourceReliabilityHandler(reliabilityRepo repository.SourceReliabilityRepository) query.ListSourceReliabilityHandler {
	return &listSourceReliabilityHandler{
		reliabilityRepo: reliabilityRepo,
	}
}

func (h *listSourceReliabilityHandler) Handle(ctx context.Context, qry query.ListSourceReliability) (query.ListSourceReliabilityResult, error) {
	params := repository.DefaultListSourceReliabilityParams().
		WithPagination(qry.Limit, qry.Offset)

	if qry.TenantID.IsPresent() {
		params = params.WithTenantID(qry.TenantID.MustGet())
	}

	if qry.MinScore.IsPresent() {
		params = params.WithMinScore(qry.MinScore.MustGet())
	}

	if qry.MaxScore.IsPresent() {
		params = params.WithMaxScore(qry.MaxScore.MustGet())
	}

	reliabilities, err := h.reliabilityRepo.List(ctx, params)
	if err != nil {
		return query.ListSourceReliabilityResult{}, err
	}

	count, err := h.reliabilityRepo.Count(ctx, params)
	if err != nil {
		return query.ListSourceReliabilityResult{}, err
	}

	return query.ListSourceReliabilityResult{
		Reliabilities: reliabilities,
		TotalCount:    count,
	}, nil
}

type getIngestStatsHandler struct {
	recordRepo     repository.IngestRecordRepository
	quarantineRepo repository.QuarantinedRecordRepository
}

func NewGetIngestStatsHandler(
	recordRepo repository.IngestRecordRepository,
	quarantineRepo repository.QuarantinedRecordRepository,
) query.GetIngestStatsHandler {
	return &getIngestStatsHandler{
		recordRepo:     recordRepo,
		quarantineRepo: quarantineRepo,
	}
}

func (h *getIngestStatsHandler) Handle(ctx context.Context, qry query.GetIngestStats) (query.GetIngestStatsResult, error) {
	baseParams := repository.DefaultListIngestRecordsParams()

	if qry.TenantID.IsPresent() {
		baseParams = baseParams.WithTenantID(qry.TenantID.MustGet())
	}

	if qry.SourceID.IsPresent() {
		baseParams = baseParams.WithSourceID(qry.SourceID.MustGet())
	}

	if qry.SourceType.IsPresent() {
		baseParams = baseParams.WithSourceType(qry.SourceType.MustGet())
	}

	if qry.TimeRange.IsPresent() {
		timeRange := qry.TimeRange.MustGet()
		if timeRange.Start.IsPresent() && timeRange.End.IsPresent() {
			baseParams = baseParams.WithTimeRange(timeRange.Start.MustGet(), timeRange.End.MustGet())
		}
	}

	totalCount, err := h.recordRepo.Count(ctx, baseParams)
	if err != nil {
		return query.GetIngestStatsResult{}, err
	}

	acceptedParams := baseParams.WithStatus(model.IngestStatusAccepted)
	acceptedCount, err := h.recordRepo.Count(ctx, acceptedParams)
	if err != nil {
		return query.GetIngestStatsResult{}, err
	}

	rejectedParams := baseParams.WithStatus(model.IngestStatusRejected)
	rejectedCount, err := h.recordRepo.Count(ctx, rejectedParams)
	if err != nil {
		return query.GetIngestStatsResult{}, err
	}

	quarantinedParams := baseParams.WithStatus(model.IngestStatusQuarantined)
	quarantinedCount, err := h.recordRepo.Count(ctx, quarantinedParams)
	if err != nil {
		return query.GetIngestStatsResult{}, err
	}

	pendingParams := baseParams.WithStatus(model.IngestStatusPending)
	pendingCount, err := h.recordRepo.Count(ctx, pendingParams)
	if err != nil {
		return query.GetIngestStatsResult{}, err
	}

	stats := query.IngestStats{
		TotalRecords:       totalCount,
		AcceptedRecords:    acceptedCount,
		RejectedRecords:    rejectedCount,
		QuarantinedRecords: quarantinedCount,
		PendingRecords:     pendingCount,

		RecordsBySource:     make(map[string]int64),
		RecordsBySourceType: make(map[string]int64),
		RecordsByStatus:     make(map[string]int64),
		RecordsByEntityType: make(map[string]int64),
		AnomaliesByType:     make(map[string]int64),
		QuarantineByReason:  make(map[string]int64),
	}

	stats.RecordsByStatus["accepted"] = acceptedCount
	stats.RecordsByStatus["rejected"] = rejectedCount
	stats.RecordsByStatus["quarantined"] = quarantinedCount
	stats.RecordsByStatus["pending"] = pendingCount

	return query.GetIngestStatsResult{Stats: stats}, nil
}
