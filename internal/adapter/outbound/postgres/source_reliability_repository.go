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

type sourceReliabilityRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewSourceReliabilityRepository(pool *pgxpool.Pool) repository.SourceReliabilityRepository {
	return &sourceReliabilityRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *sourceReliabilityRepository) Create(ctx context.Context, reliability *model.SourceReliability) error {
	params := toUpsertSourceReliabilityParams(reliability)
	return r.queries.UpsertSourceReliability(ctx, params)
}

func (r *sourceReliabilityRepository) Update(ctx context.Context, reliability *model.SourceReliability) error {
	params := toUpdateSourceReliabilityParams(reliability)
	return r.queries.UpdateSourceReliability(ctx, params)
}

func (r *sourceReliabilityRepository) Upsert(ctx context.Context, reliability *model.SourceReliability) error {
	params := toUpsertSourceReliabilityParams(reliability)
	return r.queries.UpsertSourceReliability(ctx, params)
}

func (r *sourceReliabilityRepository) FindBySourceID(ctx context.Context, sourceID types.ID) (*model.SourceReliability, error) {
	row, err := r.queries.FindSourceReliabilityByID(ctx, string(sourceID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return toSourceReliabilityModel(row), nil
}

func (r *sourceReliabilityRepository) List(ctx context.Context, params repository.ListSourceReliabilityParams) ([]*model.SourceReliability, error) {
	sqlcParams := sqlc.ListSourceReliabilityParams{
		Limit:    int32(params.Limit),
		Offset:   int32(params.Offset),
		TenantID: optionalIDToPgText(params.TenantID),
	}

	if params.MinScore.IsPresent() {
		sqlcParams.MinScore = float32ToPgFloat4(float32(params.MinScore.MustGet()))
	}

	if params.MaxScore.IsPresent() {
		sqlcParams.MaxScore = float32ToPgFloat4(float32(params.MaxScore.MustGet()))
	}

	rows, err := r.queries.ListSourceReliability(ctx, sqlcParams)
	if err != nil {
		return nil, err
	}

	reliabilities := make([]*model.SourceReliability, len(rows))
	for i, row := range rows {
		reliabilities[i] = toSourceReliabilityModel(row)
	}
	return reliabilities, nil
}

func (r *sourceReliabilityRepository) Count(ctx context.Context, params repository.ListSourceReliabilityParams) (int64, error) {
	sqlcParams := sqlc.CountSourceReliabilityParams{
		TenantID: optionalIDToPgText(params.TenantID),
	}

	if params.MinScore.IsPresent() {
		sqlcParams.MinScore = float32ToPgFloat4(float32(params.MinScore.MustGet()))
	}

	if params.MaxScore.IsPresent() {
		sqlcParams.MaxScore = float32ToPgFloat4(float32(params.MaxScore.MustGet()))
	}

	return r.queries.CountSourceReliability(ctx, sqlcParams)
}

func (r *sourceReliabilityRepository) Delete(ctx context.Context, sourceID types.ID) error {
	return r.queries.DeleteSourceReliability(ctx, string(sourceID))
}
