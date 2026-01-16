package query

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
)

type GetIngestStats struct {
	TenantID   types.Optional[types.ID]
	SourceID   types.Optional[types.ID]
	SourceType types.Optional[string]
	TimeRange  types.Optional[command.TimeRange]
}

func (q GetIngestStats) QueryName() string {
	return "ingest.get_ingest_stats"
}

func DefaultGetIngestStats() GetIngestStats {
	return GetIngestStats{
		TenantID:   types.None[types.ID](),
		SourceID:   types.None[types.ID](),
		SourceType: types.None[string](),
		TimeRange:  types.None[command.TimeRange](),
	}
}

func (q GetIngestStats) WithTenantID(tenantID types.ID) GetIngestStats {
	q.TenantID = types.Some(tenantID)
	return q
}

func (q GetIngestStats) WithSourceID(sourceID types.ID) GetIngestStats {
	q.SourceID = types.Some(sourceID)
	return q
}

func (q GetIngestStats) WithSourceType(sourceType string) GetIngestStats {
	q.SourceType = types.Some(sourceType)
	return q
}

func (q GetIngestStats) WithTimeRange(timeRange command.TimeRange) GetIngestStats {
	q.TimeRange = types.Some(timeRange)
	return q
}

type IngestStats struct {
	TotalRecords       int64
	AcceptedRecords    int64
	RejectedRecords    int64
	QuarantinedRecords int64
	PendingRecords     int64

	AverageConfidence       float64
	AverageProcessingTimeMs float64

	RecordsBySource     map[string]int64
	RecordsBySourceType map[string]int64
	RecordsByStatus     map[string]int64
	RecordsByEntityType map[string]int64
	AnomaliesByType     map[string]int64
	QuarantineByReason  map[string]int64

	SourceSignaturesVerified    int64
	SourceSignaturesFailed      int64
	CollectorSignaturesVerified int64
	CollectorSignaturesFailed   int64
}

type GetIngestStatsResult struct {
	Stats IngestStats
}

type GetIngestStatsHandler interface {
	Handle(ctx context.Context, qry GetIngestStats) (GetIngestStatsResult, error)
}
