package query

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type GetSourceReliability struct {
	SourceID types.ID
}

func (q GetSourceReliability) QueryName() string {
	return "ingest.get_source_reliability"
}

type GetSourceReliabilityResult struct {
	Reliability *model.SourceReliability
}

type GetSourceReliabilityHandler interface {
	Handle(ctx context.Context, qry GetSourceReliability) (GetSourceReliabilityResult, error)
}

type ListSourceReliability struct {
	TenantID types.Optional[types.ID]
	MinScore types.Optional[float64]
	MaxScore types.Optional[float64]
	Limit    int
	Offset   int
}

func (q ListSourceReliability) QueryName() string {
	return "ingest.list_source_reliability"
}

func DefaultListSourceReliability() ListSourceReliability {
	return ListSourceReliability{
		TenantID: types.None[types.ID](),
		MinScore: types.None[float64](),
		MaxScore: types.None[float64](),
		Limit:    20,
		Offset:   0,
	}
}

func (q ListSourceReliability) WithTenantID(tenantID types.ID) ListSourceReliability {
	q.TenantID = types.Some(tenantID)
	return q
}

func (q ListSourceReliability) WithMinScore(minScore float64) ListSourceReliability {
	q.MinScore = types.Some(minScore)
	return q
}

func (q ListSourceReliability) WithMaxScore(maxScore float64) ListSourceReliability {
	q.MaxScore = types.Some(maxScore)
	return q
}

func (q ListSourceReliability) WithPagination(limit, offset int) ListSourceReliability {
	q.Limit = limit
	q.Offset = offset
	return q
}

type ListSourceReliabilityResult struct {
	Reliabilities []*model.SourceReliability
	TotalCount    int64
}

type ListSourceReliabilityHandler interface {
	Handle(ctx context.Context, qry ListSourceReliability) (ListSourceReliabilityResult, error)
}
