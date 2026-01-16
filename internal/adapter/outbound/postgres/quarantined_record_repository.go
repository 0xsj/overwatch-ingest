package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/adapter/outbound/postgres/sqlc"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

type quarantinedRecordRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewQuarantinedRecordRepository(pool *pgxpool.Pool) repository.QuarantinedRecordRepository {
	return &quarantinedRecordRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *quarantinedRecordRepository) Create(ctx context.Context, record *model.QuarantinedRecord) error {
	params := toCreateQuarantinedRecordParams(record)
	return r.queries.CreateQuarantinedRecord(ctx, params)
}

func (r *quarantinedRecordRepository) Update(ctx context.Context, record *model.QuarantinedRecord) error {
	params := toUpdateQuarantinedRecordParams(record)
	return r.queries.UpdateQuarantinedRecord(ctx, params)
}

func (r *quarantinedRecordRepository) FindByID(ctx context.Context, id types.ID) (*model.QuarantinedRecord, error) {
	row, err := r.queries.FindQuarantinedRecordByID(ctx, string(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return toQuarantinedRecordModel(row), nil
}

func (r *quarantinedRecordRepository) FindByIngestRecordID(ctx context.Context, ingestRecordID types.ID) (*model.QuarantinedRecord, error) {
	row, err := r.queries.FindQuarantinedRecordByIngestRecordID(ctx, string(ingestRecordID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return toQuarantinedRecordModel(row), nil
}

func (r *quarantinedRecordRepository) List(ctx context.Context, params repository.ListQuarantinedRecordsParams) ([]*model.QuarantinedRecord, error) {
	sqlcParams := sqlc.ListQuarantinedRecordsParams{
		Limit:      int32(params.Limit),
		Offset:     int32(params.Offset),
		TenantID:   optionalIDToPgText(params.TenantID),
		SourceID:   optionalIDToPgText(params.SourceID),
		SourceType: optionalStringToPgText(params.SourceType),
	}

	if params.Reason.IsPresent() {
		sqlcParams.Reason = sqlc.NullQuarantineReason{
			QuarantineReason: sqlc.QuarantineReason(quarantineReasonToString(params.Reason.MustGet())),
			Valid:            true,
		}
	}

	if params.Resolution.IsPresent() {
		sqlcParams.Resolution = sqlc.NullQuarantineResolution{
			QuarantineResolution: sqlc.QuarantineResolution(quarantineResolutionToString(params.Resolution.MustGet())),
			Valid:                true,
		}
	}

	if params.StartTime.IsPresent() {
		sqlcParams.StartTime = timeToPgTimestamptz(params.StartTime.MustGet())
	}

	if params.EndTime.IsPresent() {
		sqlcParams.EndTime = timeToPgTimestamptz(params.EndTime.MustGet())
	}

	rows, err := r.queries.ListQuarantinedRecords(ctx, sqlcParams)
	if err != nil {
		return nil, err
	}

	records := make([]*model.QuarantinedRecord, len(rows))
	for i, row := range rows {
		records[i] = toQuarantinedRecordModel(row)
	}
	return records, nil
}

func (r *quarantinedRecordRepository) Count(ctx context.Context, params repository.ListQuarantinedRecordsParams) (int64, error) {
	sqlcParams := sqlc.CountQuarantinedRecordsParams{
		TenantID:   optionalIDToPgText(params.TenantID),
		SourceID:   optionalIDToPgText(params.SourceID),
		SourceType: optionalStringToPgText(params.SourceType),
	}

	if params.Reason.IsPresent() {
		sqlcParams.Reason = sqlc.NullQuarantineReason{
			QuarantineReason: sqlc.QuarantineReason(quarantineReasonToString(params.Reason.MustGet())),
			Valid:            true,
		}
	}

	if params.Resolution.IsPresent() {
		sqlcParams.Resolution = sqlc.NullQuarantineResolution{
			QuarantineResolution: sqlc.QuarantineResolution(quarantineResolutionToString(params.Resolution.MustGet())),
			Valid:                true,
		}
	}

	if params.StartTime.IsPresent() {
		sqlcParams.StartTime = timeToPgTimestamptz(params.StartTime.MustGet())
	}

	if params.EndTime.IsPresent() {
		sqlcParams.EndTime = timeToPgTimestamptz(params.EndTime.MustGet())
	}

	return r.queries.CountQuarantinedRecords(ctx, sqlcParams)
}

func (r *quarantinedRecordRepository) FindExpired(ctx context.Context, limit int) ([]*model.QuarantinedRecord, error) {
	rows, err := r.queries.ListExpiredQuarantinedRecords(ctx, sqlc.ListExpiredQuarantinedRecordsParams{
		ExpiresAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	records := make([]*model.QuarantinedRecord, len(rows))
	for i, row := range rows {
		records[i] = toQuarantinedRecordModel(row)
	}
	return records, nil
}

func (r *quarantinedRecordRepository) CountPending(ctx context.Context, tenantID types.Optional[types.ID]) (int64, error) {
	return r.queries.CountPendingByTenant(ctx, optionalIDToPgText(tenantID))
}
