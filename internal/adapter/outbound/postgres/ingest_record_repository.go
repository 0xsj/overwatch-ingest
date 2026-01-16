package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/adapter/outbound/postgres/sqlc"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

type ingestRecordRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewIngestRecordRepository(pool *pgxpool.Pool) repository.IngestRecordRepository {
	return &ingestRecordRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *ingestRecordRepository) Create(ctx context.Context, record *model.IngestRecord) error {
	params := toCreateIngestRecordParams(record)
	return r.queries.CreateIngestRecord(ctx, params)
}

func (r *ingestRecordRepository) Update(ctx context.Context, record *model.IngestRecord) error {
	params := toUpdateIngestRecordParams(record)
	return r.queries.UpdateIngestRecord(ctx, params)
}

func (r *ingestRecordRepository) FindByID(ctx context.Context, id types.ID) (*model.IngestRecord, error) {
	row, err := r.queries.FindIngestRecordByID(ctx, string(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return toIngestRecordModel(row), nil
}

func (r *ingestRecordRepository) FindByRawDataID(ctx context.Context, rawDataID string) (*model.IngestRecord, error) {
	row, err := r.queries.FindIngestRecordByRawDataID(ctx, rawDataID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return toIngestRecordModel(row), nil
}

func (r *ingestRecordRepository) ExistsByRawDataID(ctx context.Context, rawDataID string) (bool, error) {
	_, err := r.queries.FindIngestRecordByRawDataID(ctx, rawDataID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *ingestRecordRepository) List(ctx context.Context, params repository.ListIngestRecordsParams) ([]*model.IngestRecord, error) {
	sqlcParams := sqlc.ListIngestRecordsParams{
		Limit:      int32(params.Limit),
		Offset:     int32(params.Offset),
		TenantID:   optionalIDToPgText(params.TenantID),
		SourceID:   optionalIDToPgText(params.SourceID),
		SourceType: optionalStringToPgText(params.SourceType),
		EntityType: optionalStringToPgText(params.EntityType),
	}

	if params.Status.IsPresent() {
		sqlcParams.Status = sqlc.NullIngestStatus{
			IngestStatus: sqlc.IngestStatus(ingestStatusToString(params.Status.MustGet())),
			Valid:        true,
		}
	}

	if params.StartTime.IsPresent() {
		sqlcParams.StartTime = timeToPgTimestamptz(params.StartTime.MustGet())
	}

	if params.EndTime.IsPresent() {
		sqlcParams.EndTime = timeToPgTimestamptz(params.EndTime.MustGet())
	}

	rows, err := r.queries.ListIngestRecords(ctx, sqlcParams)
	if err != nil {
		return nil, err
	}

	records := make([]*model.IngestRecord, len(rows))
	for i, row := range rows {
		records[i] = toIngestRecordModel(row)
	}
	return records, nil
}

func (r *ingestRecordRepository) Count(ctx context.Context, params repository.ListIngestRecordsParams) (int64, error) {
	sqlcParams := sqlc.CountIngestRecordsParams{
		TenantID:   optionalIDToPgText(params.TenantID),
		SourceID:   optionalIDToPgText(params.SourceID),
		SourceType: optionalStringToPgText(params.SourceType),
		EntityType: optionalStringToPgText(params.EntityType),
	}

	if params.Status.IsPresent() {
		sqlcParams.Status = sqlc.NullIngestStatus{
			IngestStatus: sqlc.IngestStatus(ingestStatusToString(params.Status.MustGet())),
			Valid:        true,
		}
	}

	if params.StartTime.IsPresent() {
		sqlcParams.StartTime = timeToPgTimestamptz(params.StartTime.MustGet())
	}

	if params.EndTime.IsPresent() {
		sqlcParams.EndTime = timeToPgTimestamptz(params.EndTime.MustGet())
	}

	return r.queries.CountIngestRecords(ctx, sqlcParams)
}
