package command

import (
	"context"

	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"
)

type ProcessRawData struct {
	TenantID        types.Optional[types.ID]
	SourceID        types.ID
	SourceType      string
	RawDataID       string
	Payload         map[string]any
	Metadata        map[string]string
	SourceTimestamp types.Optional[types.Timestamp]
	CollectedAt     types.Timestamp
	SourceSigner    *provenance.SignatureInfo
	CollectorSigner *provenance.SignatureInfo
	JobID           types.Optional[types.ID]
	BatchID         types.Optional[types.ID]
	BatchIndex      types.Optional[int32]
}

func (c ProcessRawData) CommandName() string {
	return "ingest.process_raw_data"
}

type ProcessRawDataResult struct {
	IngestRecordID  types.ID
	Status          string
	EntityType      types.Optional[string]
	EntityID        types.Optional[string]
	EventIDs        []string
	QuarantineID    types.Optional[types.ID]
	RejectionReason types.Optional[string]
	ConfidenceScore float64
}

type ProcessRawDataHandler interface {
	Handle(ctx context.Context, cmd ProcessRawData) (ProcessRawDataResult, error)
}
