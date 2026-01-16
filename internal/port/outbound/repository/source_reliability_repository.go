package repository

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type SourceReliabilityRepository interface {
	Create(ctx context.Context, reliability *model.SourceReliability) error
	Update(ctx context.Context, reliability *model.SourceReliability) error
	Upsert(ctx context.Context, reliability *model.SourceReliability) error
	FindBySourceID(ctx context.Context, sourceID types.ID) (*model.SourceReliability, error)
	List(ctx context.Context, params ListSourceReliabilityParams) ([]*model.SourceReliability, error)
	Count(ctx context.Context, params ListSourceReliabilityParams) (int64, error)
	Delete(ctx context.Context, sourceID types.ID) error
}

type ListSourceReliabilityParams struct {
	Limit     int
	Offset    int
	TenantID  types.Optional[types.ID]
	MinScore  types.Optional[float64]
	MaxScore  types.Optional[float64]
	SortBy    SourceReliabilitySortField
	SortOrder SortOrder
}

type SourceReliabilitySortField string

const (
	SourceReliabilitySortFieldScore        SourceReliabilitySortField = "reliability_score"
	SourceReliabilitySortFieldTotalRecords SourceReliabilitySortField = "total_records"
	SourceReliabilitySortFieldCalculatedAt SourceReliabilitySortField = "calculated_at"
)

func DefaultListSourceReliabilityParams() ListSourceReliabilityParams {
	return ListSourceReliabilityParams{
		Limit:     20,
		Offset:    0,
		TenantID:  types.None[types.ID](),
		MinScore:  types.None[float64](),
		MaxScore:  types.None[float64](),
		SortBy:    SourceReliabilitySortFieldScore,
		SortOrder: SortOrderDesc,
	}
}

func (p ListSourceReliabilityParams) WithTenantID(tenantID types.ID) ListSourceReliabilityParams {
	p.TenantID = types.Some(tenantID)
	return p
}

func (p ListSourceReliabilityParams) WithMinScore(minScore float64) ListSourceReliabilityParams {
	p.MinScore = types.Some(minScore)
	return p
}

func (p ListSourceReliabilityParams) WithMaxScore(maxScore float64) ListSourceReliabilityParams {
	p.MaxScore = types.Some(maxScore)
	return p
}

func (p ListSourceReliabilityParams) WithScoreRange(min, max float64) ListSourceReliabilityParams {
	p.MinScore = types.Some(min)
	p.MaxScore = types.Some(max)
	return p
}

func (p ListSourceReliabilityParams) WithPagination(limit, offset int) ListSourceReliabilityParams {
	p.Limit = limit
	p.Offset = offset
	return p
}

func (p ListSourceReliabilityParams) WithSort(field SourceReliabilitySortField, order SortOrder) ListSourceReliabilityParams {
	p.SortBy = field
	p.SortOrder = order
	return p
}
