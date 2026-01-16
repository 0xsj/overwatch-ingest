package command

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type ResolveQuarantined struct {
	QuarantineID  types.ID
	Resolution    model.QuarantineResolution
	ResolvedBy    string
	ResolvedByDID string
	Notes         string
	ModifiedData  map[string]any
}

func (c ResolveQuarantined) CommandName() string {
	return "ingest.resolve_quarantined"
}

type ResolveQuarantinedResult struct {
	QuarantinedRecord *model.QuarantinedRecord
	IngestRecord      *model.IngestRecord
	EntityType        types.Optional[string]
	EntityID          types.Optional[string]
	EventIDs          []string
}

type ResolveQuarantinedHandler interface {
	Handle(ctx context.Context, cmd ResolveQuarantined) (ResolveQuarantinedResult, error)
}

type BulkResolveQuarantined struct {
	QuarantineIDs []types.ID
	Resolution    model.QuarantineResolution
	ResolvedBy    string
	ResolvedByDID string
	Notes         string
}

func (c BulkResolveQuarantined) CommandName() string {
	return "ingest.bulk_resolve_quarantined"
}

type BulkResolveQuarantinedResult struct {
	ResolvedCount int
	FailedIDs     []types.ID
	Errors        []string
}

type BulkResolveQuarantinedHandler interface {
	Handle(ctx context.Context, cmd BulkResolveQuarantined) (BulkResolveQuarantinedResult, error)
}
