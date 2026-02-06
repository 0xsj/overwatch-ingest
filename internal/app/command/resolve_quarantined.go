package command

import (
	"context"

	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/domain/event"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/messaging"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

type resolveQuarantinedHandler struct {
	quarantineRepo  repository.QuarantinedRecordRepository
	recordRepo      repository.IngestRecordRepository
	reliabilityRepo repository.SourceReliabilityRepository
	publisher       messaging.EventPublisher
	signer          *provenance.EnvelopeBuilder
}

func NewResolveQuarantinedHandler(
	quarantineRepo repository.QuarantinedRecordRepository,
	recordRepo repository.IngestRecordRepository,
	reliabilityRepo repository.SourceReliabilityRepository,
	publisher messaging.EventPublisher,
	signer *provenance.EnvelopeBuilder,
) command.ResolveQuarantinedHandler {
	return &resolveQuarantinedHandler{
		quarantineRepo:  quarantineRepo,
		recordRepo:      recordRepo,
		reliabilityRepo: reliabilityRepo,
		publisher:       publisher,
		signer:          signer,
	}
}

func (h *resolveQuarantinedHandler) Handle(ctx context.Context, cmd command.ResolveQuarantined) (command.ResolveQuarantinedResult, error) {
	if cmd.QuarantineID.IsEmpty() {
		return command.ResolveQuarantinedResult{}, domainerror.ErrRecordIDRequired
	}
	if !cmd.Resolution.IsValid() {
		return command.ResolveQuarantinedResult{}, domainerror.ErrResolutionInvalid
	}
	if cmd.Resolution == model.QuarantineResolutionPending {
		return command.ResolveQuarantinedResult{}, domainerror.ErrResolutionInvalid
	}
	if cmd.Resolution.RequiresModifiedData() && cmd.ModifiedData == nil {
		return command.ResolveQuarantinedResult{}, domainerror.ErrPayloadRequired
	}

	quarantined, err := h.quarantineRepo.FindByID(ctx, cmd.QuarantineID)
	if err != nil {
		return command.ResolveQuarantinedResult{}, err
	}
	if quarantined == nil {
		return command.ResolveQuarantinedResult{}, domainerror.QuarantineNotFound(cmd.QuarantineID.String())
	}

	if quarantined.IsExpired() {
		return command.ResolveQuarantinedResult{}, domainerror.ErrQuarantineExpired
	}

	record, err := h.recordRepo.FindByID(ctx, quarantined.IngestRecordID())
	if err != nil {
		return command.ResolveQuarantinedResult{}, err
	}
	if record == nil {
		return command.ResolveQuarantinedResult{}, domainerror.RecordNotFound(quarantined.IngestRecordID().String())
	}

	resolverSigner := h.signResolution(cmd)

	switch cmd.Resolution {
	case model.QuarantineResolutionApproved:
		err = quarantined.Approve(cmd.ResolvedBy, cmd.ResolvedByDID, cmd.Notes, resolverSigner)
	case model.QuarantineResolutionModified:
		err = quarantined.ApproveWithModifications(cmd.ResolvedBy, cmd.ResolvedByDID, cmd.Notes, cmd.ModifiedData, resolverSigner)
	case model.QuarantineResolutionRejected:
		err = quarantined.Reject(cmd.ResolvedBy, cmd.ResolvedByDID, cmd.Notes, resolverSigner)
	default:
		return command.ResolveQuarantinedResult{}, domainerror.ErrResolutionInvalid
	}

	if err != nil {
		return command.ResolveQuarantinedResult{}, err
	}

	var entityType, entityID string
	var eventIDs []string

	if cmd.Resolution.IsAccepted() {
		entityType = h.inferEntityType(quarantined.SourceType(), quarantined.GetEffectiveData())
		entityID = h.inferEntityID(quarantined.SourceType(), quarantined.GetEffectiveData())
		eventIDs = []string{}

		err = record.ResolveFromQuarantine(cmd.Resolution, entityType, entityID, eventIDs)
		if err != nil {
			return command.ResolveQuarantinedResult{}, err
		}
	} else {
		err = record.ResolveFromQuarantine(cmd.Resolution, "", "", nil)
		if err != nil {
			return command.ResolveQuarantinedResult{}, err
		}
	}

	if err := h.quarantineRepo.Update(ctx, quarantined); err != nil {
		return command.ResolveQuarantinedResult{}, err
	}

	if err := h.recordRepo.Update(ctx, record); err != nil {
		return command.ResolveQuarantinedResult{}, err
	}

	h.updateReliability(ctx, quarantined.SourceID(), quarantined.TenantID(), record.Status())

	_ = h.publisher.PublishQuarantineResolved(ctx, event.NewQuarantineResolved(
		quarantined.TenantID(),
		quarantined.SourceID(),
		quarantined.ID(),
		record.ID(),
		quarantined.SourceType(),
		cmd.Resolution,
		cmd.ResolvedBy,
		types.Some(cmd.ResolvedByDID),
		h.optionalString(cmd.Notes),
		h.optionalString(entityType),
		h.optionalString(entityID),
		eventIDs,
	))

	result := command.ResolveQuarantinedResult{
		QuarantinedRecord: quarantined,
		IngestRecord:      record,
	}

	if cmd.Resolution.IsAccepted() {
		result.EntityType = types.Some(entityType)
		result.EntityID = types.Some(entityID)
		result.EventIDs = eventIDs
	} else {
		result.EntityType = types.None[string]()
		result.EntityID = types.None[string]()
		result.EventIDs = nil
	}

	return result, nil
}

func (h *resolveQuarantinedHandler) signResolution(cmd command.ResolveQuarantined) *provenance.SignatureInfo {
	if h.signer == nil {
		return nil
	}

	envelope, err := h.signer.Build("quarantine.resolved", map[string]any{
		"quarantine_id": cmd.QuarantineID.String(),
		"resolution":    cmd.Resolution.String(),
		"resolved_by":   cmd.ResolvedBy,
	})
	if err != nil {
		return nil
	}

	return envelope.ToSignatureInfo()
}

func (h *resolveQuarantinedHandler) updateReliability(
	ctx context.Context,
	sourceID types.ID,
	tenantID types.Optional[types.ID],
	status model.IngestStatus,
) {
	reliability, _ := h.reliabilityRepo.FindBySourceID(ctx, sourceID)
	if reliability == nil {
		return
	}

	switch status {
	case model.IngestStatusAccepted:
		reliability.RecordAccepted()
	case model.IngestStatusRejected:
		reliability.RecordRejected()
	}

	_ = h.reliabilityRepo.Update(ctx, reliability)
}

func (h *resolveQuarantinedHandler) inferEntityType(sourceType string, payload map[string]any) string {
	if _, ok := payload["mmsi"]; ok {
		return "vessel"
	}
	if _, ok := payload["icao"]; ok {
		return "aircraft"
	}
	if _, ok := payload["lat"]; ok {
		if _, ok := payload["lon"]; ok {
			return "location"
		}
	}
	return "unknown"
}

func (h *resolveQuarantinedHandler) inferEntityID(sourceType string, payload map[string]any) string {
	if mmsi, ok := payload["mmsi"]; ok {
		return "mmsi:" + toString(mmsi)
	}
	if icao, ok := payload["icao"]; ok {
		return "icao:" + toString(icao)
	}
	return "unknown"
}

func (h *resolveQuarantinedHandler) optionalString(s string) types.Optional[string] {
	if s == "" {
		return types.None[string]()
	}
	return types.Some(s)
}

// =============================================================================
// Bulk Resolve
// =============================================================================

type bulkResolveQuarantinedHandler struct {
	resolveHandler command.ResolveQuarantinedHandler
}

func NewBulkResolveQuarantinedHandler(
	resolveHandler command.ResolveQuarantinedHandler,
) command.BulkResolveQuarantinedHandler {
	return &bulkResolveQuarantinedHandler{
		resolveHandler: resolveHandler,
	}
}

func (h *bulkResolveQuarantinedHandler) Handle(ctx context.Context, cmd command.BulkResolveQuarantined) (command.BulkResolveQuarantinedResult, error) {
	if len(cmd.QuarantineIDs) == 0 {
		return command.BulkResolveQuarantinedResult{}, domainerror.ErrRecordIDRequired
	}

	var resolvedCount int
	var failedIDs []types.ID
	var errors []string

	for _, id := range cmd.QuarantineIDs {
		resolveCmd := command.ResolveQuarantined{
			QuarantineID:  id,
			Resolution:    cmd.Resolution,
			ResolvedBy:    cmd.ResolvedBy,
			ResolvedByDID: cmd.ResolvedByDID,
			Notes:         cmd.Notes,
		}

		_, err := h.resolveHandler.Handle(ctx, resolveCmd)
		if err != nil {
			failedIDs = append(failedIDs, id)
			errors = append(errors, err.Error())
			continue
		}

		resolvedCount++
	}

	return command.BulkResolveQuarantinedResult{
		ResolvedCount: resolvedCount,
		FailedIDs:     failedIDs,
		Errors:        errors,
	}, nil
}
