package command

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

type reprocessRecordHandler struct {
	recordRepo            repository.IngestRecordRepository
	processRawDataHandler command.ProcessRawDataHandler
}

func NewReprocessRecordHandler(
	recordRepo repository.IngestRecordRepository,
	processRawDataHandler command.ProcessRawDataHandler,
) command.ReprocessRecordHandler {
	return &reprocessRecordHandler{
		recordRepo:            recordRepo,
		processRawDataHandler: processRawDataHandler,
	}
}

func (h *reprocessRecordHandler) Handle(ctx context.Context, cmd command.ReprocessRecord) (command.ReprocessRecordResult, error) {
	if cmd.IngestRecordID.IsEmpty() && cmd.RawDataID.IsEmpty() {
		return command.ReprocessRecordResult{}, domainerror.ErrRecordIDRequired
	}

	var record *model.IngestRecord
	var err error

	if cmd.IngestRecordID.IsPresent() {
		record, err = h.recordRepo.FindByID(ctx, cmd.IngestRecordID.MustGet())
	} else {
		record, err = h.recordRepo.FindByRawDataID(ctx, cmd.RawDataID.MustGet())
	}

	if err != nil {
		return command.ReprocessRecordResult{}, err
	}
	if record == nil {
		if cmd.IngestRecordID.IsPresent() {
			return command.ReprocessRecordResult{}, domainerror.RecordNotFound(cmd.IngestRecordID.MustGet().String())
		}
		return command.ReprocessRecordResult{}, domainerror.RecordNotFound(cmd.RawDataID.MustGet())
	}

	return command.ReprocessRecordResult{
		IngestRecord: record,
	}, nil
}

type reprocessBySourceHandler struct {
	recordRepo repository.IngestRecordRepository
}

func NewReprocessBySourceHandler(
	recordRepo repository.IngestRecordRepository,
) command.ReprocessBySourceHandler {
	return &reprocessBySourceHandler{
		recordRepo: recordRepo,
	}
}

func (h *reprocessBySourceHandler) Handle(ctx context.Context, cmd command.ReprocessBySource) (command.ReprocessBySourceResult, error) {
	if cmd.SourceID.IsEmpty() {
		return command.ReprocessBySourceResult{}, domainerror.ErrSourceIDRequired
	}

	params := repository.DefaultListIngestRecordsParams().
		WithSourceID(cmd.SourceID)

	if cmd.StatusFilter.IsPresent() {
		params = params.WithStatus(cmd.StatusFilter.MustGet())
	}

	if cmd.TimeRange.IsPresent() {
		timeRange := cmd.TimeRange.MustGet()
		if timeRange.Start.IsPresent() && timeRange.End.IsPresent() {
			params = params.WithTimeRange(timeRange.Start.MustGet(), timeRange.End.MustGet())
		}
	}

	count, err := h.recordRepo.Count(ctx, params)
	if err != nil {
		return command.ReprocessBySourceResult{}, err
	}

	batchID := types.NewID().String()

	return command.ReprocessBySourceResult{
		RecordsQueued: int(count),
		BatchID:       batchID,
	}, nil
}
