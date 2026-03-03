package grpc

import (
	"context"

	commonv1 "github.com/0xsj/overwatch-contracts/gen/go/common/v1"
	ingestv1 "github.com/0xsj/overwatch-contracts/gen/go/ingest/v1"
	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/query"
)

type Handler struct {
	ingestv1.UnimplementedIngestServiceServer

	// Command handlers
	resolveQuarantinedHandler     command.ResolveQuarantinedHandler
	bulkResolveQuarantinedHandler command.BulkResolveQuarantinedHandler
	reprocessRecordHandler        command.ReprocessRecordHandler
	reprocessBySourceHandler      command.ReprocessBySourceHandler

	// Query handlers
	getRecordHandler             query.GetRecordHandler
	getRecordByRawDataHandler    query.GetRecordByRawDataHandler
	listRecordsHandler           query.ListRecordsHandler
	getQuarantinedHandler        query.GetQuarantinedHandler
	listQuarantinedHandler       query.ListQuarantinedHandler
	getSourceReliabilityHandler  query.GetSourceReliabilityHandler
	listSourceReliabilityHandler query.ListSourceReliabilityHandler
	getIngestStatsHandler        query.GetIngestStatsHandler
	validateDataHandler          query.ValidateDataHandler
}

type HandlerConfig struct {
	ResolveQuarantinedHandler     command.ResolveQuarantinedHandler
	BulkResolveQuarantinedHandler command.BulkResolveQuarantinedHandler
	ReprocessRecordHandler        command.ReprocessRecordHandler
	ReprocessBySourceHandler      command.ReprocessBySourceHandler

	GetRecordHandler             query.GetRecordHandler
	GetRecordByRawDataHandler    query.GetRecordByRawDataHandler
	ListRecordsHandler           query.ListRecordsHandler
	GetQuarantinedHandler        query.GetQuarantinedHandler
	ListQuarantinedHandler       query.ListQuarantinedHandler
	GetSourceReliabilityHandler  query.GetSourceReliabilityHandler
	ListSourceReliabilityHandler query.ListSourceReliabilityHandler
	GetIngestStatsHandler        query.GetIngestStatsHandler
	ValidateDataHandler          query.ValidateDataHandler
}

func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		resolveQuarantinedHandler:     cfg.ResolveQuarantinedHandler,
		bulkResolveQuarantinedHandler: cfg.BulkResolveQuarantinedHandler,
		reprocessRecordHandler:        cfg.ReprocessRecordHandler,
		reprocessBySourceHandler:      cfg.ReprocessBySourceHandler,

		getRecordHandler:             cfg.GetRecordHandler,
		getRecordByRawDataHandler:    cfg.GetRecordByRawDataHandler,
		listRecordsHandler:           cfg.ListRecordsHandler,
		getQuarantinedHandler:        cfg.GetQuarantinedHandler,
		listQuarantinedHandler:       cfg.ListQuarantinedHandler,
		getSourceReliabilityHandler:  cfg.GetSourceReliabilityHandler,
		listSourceReliabilityHandler: cfg.ListSourceReliabilityHandler,
		getIngestStatsHandler:        cfg.GetIngestStatsHandler,
		validateDataHandler:          cfg.ValidateDataHandler,
	}
}

// =============================================================================
// Health
// =============================================================================

func (h *Handler) Ping(ctx context.Context, req *ingestv1.PingRequest) (*ingestv1.PingResponse, error) {
	return &ingestv1.PingResponse{
		Message: "pong",
	}, nil
}

// =============================================================================
// Record Access
// =============================================================================

func (h *Handler) GetRecord(ctx context.Context, req *ingestv1.GetRecordRequest) (*ingestv1.GetRecordResponse, error) {
	qry := query.GetRecord{
		ID: types.ID(req.GetId()),
	}

	result, err := h.getRecordHandler.Handle(ctx, qry)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.GetRecordResponse{
		Record: ingestRecordToProto(result.Record),
	}, nil
}

func (h *Handler) ListRecords(ctx context.Context, req *ingestv1.ListRecordsRequest) (*ingestv1.ListRecordsResponse, error) {
	qry := query.DefaultListRecords()

	if req.TenantId != nil {
		qry = qry.WithTenantID(types.ID(*req.TenantId))
	}
	if req.SourceId != nil {
		qry = qry.WithSourceID(types.ID(*req.SourceId))
	}
	if req.SourceType != nil {
		qry = qry.WithSourceType(*req.SourceType)
	}
	if req.Status != nil {
		qry = qry.WithStatus(ingestStatusFromProto(*req.Status))
	}
	if req.EntityType != nil {
		qry = qry.WithEntityType(*req.EntityType)
	}
	if req.TimeRange != nil {
		qry = qry.WithTimeRange(timeRangeFromProtoToCommand(req.TimeRange))
	}
	if req.Pagination != nil {
		limit, offset := paginationFromProto(req.Pagination)
		qry = qry.WithPagination(limit, offset)
	}

	result, err := h.listRecordsHandler.Handle(ctx, qry)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.ListRecordsResponse{
		Records:    ingestRecordsToProto(result.Records),
		Pagination: paginationToProto(result.TotalCount, qry.Limit, qry.Offset),
	}, nil
}

func (h *Handler) GetRecordByRawData(ctx context.Context, req *ingestv1.GetRecordByRawDataRequest) (*ingestv1.GetRecordByRawDataResponse, error) {
	qry := query.GetRecordByRawData{
		RawDataID: req.GetRawDataId(),
	}

	result, err := h.getRecordByRawDataHandler.Handle(ctx, qry)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.GetRecordByRawDataResponse{
		Record: ingestRecordToProto(result.Record),
	}, nil
}

// =============================================================================
// Quarantine Management
// =============================================================================

func (h *Handler) GetQuarantined(ctx context.Context, req *ingestv1.GetQuarantinedRequest) (*ingestv1.GetQuarantinedResponse, error) {
	qry := query.GetQuarantined{
		ID: types.ID(req.GetId()),
	}

	result, err := h.getQuarantinedHandler.Handle(ctx, qry)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.GetQuarantinedResponse{
		Record: quarantinedRecordToProto(result.Record),
	}, nil
}

func (h *Handler) ListQuarantined(ctx context.Context, req *ingestv1.ListQuarantinedRequest) (*ingestv1.ListQuarantinedResponse, error) {
	qry := query.DefaultListQuarantined()

	if req.TenantId != nil {
		qry = qry.WithTenantID(types.ID(*req.TenantId))
	}
	if req.SourceId != nil {
		qry = qry.WithSourceID(types.ID(*req.SourceId))
	}
	if req.SourceType != nil {
		qry = qry.WithSourceType(*req.SourceType)
	}
	if req.Reason != nil {
		qry = qry.WithReason(quarantineReasonFromProto(*req.Reason))
	}
	if req.Resolution != nil {
		qry = qry.WithResolution(quarantineResolutionFromProto(*req.Resolution))
	}
	if req.TimeRange != nil {
		qry = qry.WithTimeRange(timeRangeFromProtoToCommand(req.TimeRange))
	}
	if req.Pagination != nil {
		limit, offset := paginationFromProto(req.Pagination)
		qry = qry.WithPagination(limit, offset)
	}

	result, err := h.listQuarantinedHandler.Handle(ctx, qry)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.ListQuarantinedResponse{
		Records:    quarantinedRecordsToProto(result.Records),
		Pagination: paginationToProto(result.TotalCount, qry.Limit, qry.Offset),
	}, nil
}

func (h *Handler) ResolveQuarantined(ctx context.Context, req *ingestv1.ResolveQuarantinedRequest) (*ingestv1.ResolveQuarantinedResponse, error) {
	actorID, actorDID, err := getActorFromContext(ctx)
	if err != nil {
		return nil, toGRPCError(err)
	}

	var modifiedData map[string]any
	if req.ModifiedData != nil {
		modifiedData = req.ModifiedData.AsMap()
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	cmd := command.ResolveQuarantined{
		QuarantineID:  types.ID(req.GetId()),
		Resolution:    quarantineResolutionFromProto(req.GetResolution()),
		ResolvedBy:    actorID,
		ResolvedByDID: actorDID,
		Notes:         notes,
		ModifiedData:  modifiedData,
	}

	result, err := h.resolveQuarantinedHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	resp := &ingestv1.ResolveQuarantinedResponse{
		Record: quarantinedRecordToProto(result.QuarantinedRecord),
	}

	if result.IngestRecord != nil {
		resp.IngestRecord = ingestRecordToProto(result.IngestRecord)
	}

	return resp, nil
}

func (h *Handler) BulkResolveQuarantined(ctx context.Context, req *ingestv1.BulkResolveQuarantinedRequest) (*ingestv1.BulkResolveQuarantinedResponse, error) {
	actorID, actorDID, err := getActorFromContext(ctx)
	if err != nil {
		return nil, toGRPCError(err)
	}

	ids := make([]types.ID, len(req.GetIds()))
	for i, id := range req.GetIds() {
		ids[i] = types.ID(id)
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	cmd := command.BulkResolveQuarantined{
		QuarantineIDs: ids,
		Resolution:    quarantineResolutionFromProto(req.GetResolution()),
		ResolvedBy:    actorID,
		ResolvedByDID: actorDID,
		Notes:         notes,
	}

	result, err := h.bulkResolveQuarantinedHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	failedIDs := make([]string, len(result.FailedIDs))
	for i, id := range result.FailedIDs {
		failedIDs[i] = string(id)
	}

	return &ingestv1.BulkResolveQuarantinedResponse{
		ResolvedCount: int32(result.ResolvedCount),
		FailedIds:     failedIDs,
		Errors:        result.Errors,
	}, nil
}

// =============================================================================
// Reprocessing
// =============================================================================

func (h *Handler) ReprocessRecord(ctx context.Context, req *ingestv1.ReprocessRecordRequest) (*ingestv1.ReprocessRecordResponse, error) {
	cmd := command.ReprocessRecord{
		IngestRecordID: types.Some(types.ID(req.GetId())),
		RawDataID:      optionalString(req.RawDataId),
	}

	result, err := h.reprocessRecordHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.ReprocessRecordResponse{
		Record: ingestRecordToProto(result.IngestRecord),
	}, nil
}

func (h *Handler) ReprocessBySource(ctx context.Context, req *ingestv1.ReprocessBySourceRequest) (*ingestv1.ReprocessBySourceResponse, error) {
	cmd := command.ReprocessBySource{
		SourceID: types.ID(req.GetSourceId()),
	}

	if req.TimeRange != nil {
		cmd.TimeRange = types.Some(timeRangeFromProtoToCommand(req.TimeRange))
	}
	if req.StatusFilter != nil {
		cmd.StatusFilter = types.Some(ingestStatusFromProto(*req.StatusFilter))
	}

	result, err := h.reprocessBySourceHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.ReprocessBySourceResponse{
		RecordsQueued: int32(result.RecordsQueued),
		BatchId:       result.BatchID,
	}, nil
}

// =============================================================================
// Validation (preview)
// =============================================================================

func (h *Handler) ValidateData(ctx context.Context, req *ingestv1.ValidateDataRequest) (*ingestv1.ValidateDataResponse, error) {
	if h.validateDataHandler == nil {
		return &ingestv1.ValidateDataResponse{
			Validation:      validationResultToProto(model.ValidationResult{}),
			Confidence:      confidenceScoreToProto(model.ConfidenceScore{}),
			PredictedStatus: ingestv1.IngestStatus_INGEST_STATUS_UNSPECIFIED,
		}, nil
	}

	var payload map[string]any
	if req.GetData() != nil {
		payload = req.GetData().AsMap()
	}

	result, err := h.validateDataHandler.Handle(ctx, query.ValidateData{
		SourceID:   req.GetSourceId(),
		SourceType: req.GetSourceType(),
		Payload:    payload,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.ValidateDataResponse{
		Validation:      validationResultToProto(result.Validation),
		Confidence:      confidenceScoreToProto(result.Confidence),
		PredictedStatus: ingestStatusToProto(result.PredictedStatus),
	}, nil
}

// =============================================================================
// Source Reliability
// =============================================================================

func (h *Handler) GetSourceReliability(ctx context.Context, req *ingestv1.GetSourceReliabilityRequest) (*ingestv1.GetSourceReliabilityResponse, error) {
	qry := query.GetSourceReliability{
		SourceID: types.ID(req.GetSourceId()),
	}

	result, err := h.getSourceReliabilityHandler.Handle(ctx, qry)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.GetSourceReliabilityResponse{
		Reliability: sourceReliabilityToProto(result.Reliability),
	}, nil
}

func (h *Handler) ListSourceReliability(ctx context.Context, req *ingestv1.ListSourceReliabilityRequest) (*ingestv1.ListSourceReliabilityResponse, error) {
	qry := query.DefaultListSourceReliability()

	if req.TenantId != nil {
		qry = qry.WithTenantID(types.ID(*req.TenantId))
	}
	if req.MinScore != nil {
		qry = qry.WithMinScore(float64(*req.MinScore))
	}
	if req.MaxScore != nil {
		qry = qry.WithMaxScore(float64(*req.MaxScore))
	}
	if req.Pagination != nil {
		limit, offset := paginationFromProto(req.Pagination)
		qry = qry.WithPagination(limit, offset)
	}

	result, err := h.listSourceReliabilityHandler.Handle(ctx, qry)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.ListSourceReliabilityResponse{
		Reliabilities: sourceReliabilitiesToProto(result.Reliabilities),
		Pagination:    paginationToProto(result.TotalCount, qry.Limit, qry.Offset),
	}, nil
}

// =============================================================================
// Stats
// =============================================================================

func (h *Handler) GetIngestStats(ctx context.Context, req *ingestv1.GetIngestStatsRequest) (*ingestv1.GetIngestStatsResponse, error) {
	qry := query.DefaultGetIngestStats()

	if req.TenantId != nil {
		qry = qry.WithTenantID(types.ID(*req.TenantId))
	}
	if req.SourceId != nil {
		qry = qry.WithSourceID(types.ID(*req.SourceId))
	}
	if req.SourceType != nil {
		qry = qry.WithSourceType(*req.SourceType)
	}
	if req.TimeRange != nil {
		qry = qry.WithTimeRange(timeRangeFromProtoToCommand(req.TimeRange))
	}

	result, err := h.getIngestStatsHandler.Handle(ctx, qry)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &ingestv1.GetIngestStatsResponse{
		TotalRecords:                result.Stats.TotalRecords,
		AcceptedRecords:             result.Stats.AcceptedRecords,
		RejectedRecords:             result.Stats.RejectedRecords,
		QuarantinedRecords:          result.Stats.QuarantinedRecords,
		PendingRecords:              result.Stats.PendingRecords,
		AverageConfidence:           float32(result.Stats.AverageConfidence),
		AverageProcessingTimeMs:     float32(result.Stats.AverageProcessingTimeMs),
		RecordsBySource:             result.Stats.RecordsBySource,
		RecordsBySourceType:         result.Stats.RecordsBySourceType,
		RecordsByStatus:             result.Stats.RecordsByStatus,
		RecordsByEntityType:         result.Stats.RecordsByEntityType,
		AnomaliesByType:             result.Stats.AnomaliesByType,
		QuarantineByReason:          result.Stats.QuarantineByReason,
		SourceSignaturesVerified:    result.Stats.SourceSignaturesVerified,
		SourceSignaturesFailed:      result.Stats.SourceSignaturesFailed,
		CollectorSignaturesVerified: result.Stats.CollectorSignaturesVerified,
		CollectorSignaturesFailed:   result.Stats.CollectorSignaturesFailed,
	}, nil
}

// =============================================================================
// Context helpers
// =============================================================================

func getActorFromContext(ctx context.Context) (actorID string, actorDID string, err error) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return "", "", err
	}

	userDID, err := getUserDIDFromContext(ctx)
	if err != nil {
		return "", "", err
	}

	return string(userID), userDID, nil
}

// =============================================================================
// Pagination helpers
// =============================================================================

func paginationFromProto(pb *commonv1.PageRequest) (limit, offset int) {
	if pb == nil {
		return 20, 0
	}
	limit = int(pb.PageSize)
	if limit <= 0 {
		limit = 20
	}
	page := int(pb.Page)
	if page <= 0 {
		page = 1
	}
	offset = (page - 1) * limit
	return limit, offset
}

func paginationToProto(totalCount int64, limit, offset int) *commonv1.PageResponse {
	page := int32(1)
	if limit > 0 {
		page = int32(offset/limit) + 1
	}
	pageSize := int32(limit)
	totalPages := int32(0)
	if limit > 0 {
		totalPages = int32((totalCount + int64(limit) - 1) / int64(limit))
	}

	return &commonv1.PageResponse{
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		TotalItems:  totalCount,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}

// =============================================================================
// TimeRange helper
// =============================================================================

func timeRangeFromProtoToCommand(pb *ingestv1.TimeRange) command.TimeRange {
	if pb == nil {
		return command.TimeRange{}
	}

	tr := command.TimeRange{}
	if pb.Start != nil {
		tr.Start = types.Some(types.FromTime(pb.Start.AsTime()))
	}
	if pb.End != nil {
		tr.End = types.Some(types.FromTime(pb.End.AsTime()))
	}
	return tr
}
