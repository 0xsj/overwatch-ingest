package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/0xsj/overwatch-pkg/types"

	appcommand "github.com/0xsj/overwatch-ingest/internal/app/command"
	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
	"github.com/0xsj/overwatch-ingest/tests/testutil"
	"github.com/0xsj/overwatch-ingest/tests/testutil/mocks"
)

type resolveQuarantinedDeps struct {
	quarantineRepo  *mocks.QuarantinedRecordRepository
	recordRepo      *mocks.IngestRecordRepository
	reliabilityRepo *mocks.SourceReliabilityRepository
	publisher       *mocks.EventPublisher
}

func newResolveQuarantinedDeps() *resolveQuarantinedDeps {
	return &resolveQuarantinedDeps{
		quarantineRepo:  mocks.NewQuarantinedRecordRepository(),
		recordRepo:      mocks.NewIngestRecordRepository(),
		reliabilityRepo: mocks.NewSourceReliabilityRepository(),
		publisher:       mocks.NewEventPublisher(),
	}
}

func newResolveQuarantinedHandler(d *resolveQuarantinedDeps) command.ResolveQuarantinedHandler {
	return appcommand.NewResolveQuarantinedHandler(
		d.quarantineRepo,
		d.recordRepo,
		d.reliabilityRepo,
		d.publisher,
		mocks.NilEnvelopeBuilder(),
	)
}

// seedQuarantinedPair creates a matching quarantined record and ingest record pair,
// seeding them into the repositories and returning both.
func seedQuarantinedPair(d *resolveQuarantinedDeps) (*model.QuarantinedRecord, *model.IngestRecord) {
	sourceID := types.NewID()
	rawDataID := testutil.Fake.RawDataID()

	rec := testutil.Fixtures.IngestRecordBuilder().
		WithSourceID(sourceID).
		WithRawDataID(rawDataID).
		Build()
	rec.SetValidation(testutil.Fixtures.WarningValidationResult())
	rec.SetConfidence(testutil.Fixtures.MediumConfidenceScore())

	qr := testutil.Fixtures.QuarantinedRecordBuilder().
		WithSourceID(sourceID).
		WithIngestRecordID(rec.ID()).
		Build()

	if err := rec.Quarantine(qr.ID()); err != nil {
		panic("seedQuarantinedPair: failed to quarantine record: " + err.Error())
	}

	d.recordRepo.Seed(rec)
	d.quarantineRepo.Seed(qr)

	return qr, rec
}

// =============================================================================
// Approve
// =============================================================================

func TestResolveQuarantined_Approve_Success(t *testing.T) {
	d := newResolveQuarantinedDeps()
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	qr, _ := seedQuarantinedPair(d)

	result, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  qr.ID(),
		Resolution:    model.QuarantineResolutionApproved,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
		Notes:         "data verified manually",
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.QuarantinedRecord == nil {
		t.Fatal("result.QuarantinedRecord is nil")
	}
	if !result.QuarantinedRecord.IsResolved() {
		t.Error("quarantined record should be resolved")
	}
	if result.IngestRecord == nil {
		t.Fatal("result.IngestRecord is nil")
	}
	if !result.IngestRecord.IsAccepted() {
		t.Error("ingest record should be accepted after approval")
	}
	if result.EntityType.IsEmpty() {
		t.Error("EntityType should be present for approved resolution")
	}
	if d.quarantineRepo.Calls.Update != 1 {
		t.Errorf("QuarantineRepo.Update calls = %d, want 1", d.quarantineRepo.Calls.Update)
	}
	if d.recordRepo.Calls.Update != 1 {
		t.Errorf("RecordRepo.Update calls = %d, want 1", d.recordRepo.Calls.Update)
	}
	if d.publisher.Calls.PublishQuarantineResolved != 1 {
		t.Errorf("PublishQuarantineResolved calls = %d, want 1", d.publisher.Calls.PublishQuarantineResolved)
	}
}

// =============================================================================
// Approve with modifications
// =============================================================================

func TestResolveQuarantined_Modified_Success(t *testing.T) {
	d := newResolveQuarantinedDeps()
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	qr, _ := seedQuarantinedPair(d)

	modifiedData := map[string]any{
		"mmsi":      "123456789",
		"lat":       45.0,
		"lon":       -73.0,
		"corrected": true,
	}

	result, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  qr.ID(),
		Resolution:    model.QuarantineResolutionModified,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
		Notes:         "corrected coordinates",
		ModifiedData:  modifiedData,
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if !result.IngestRecord.IsAccepted() {
		t.Error("ingest record should be accepted after modified approval")
	}
	if result.EntityType.IsEmpty() {
		t.Error("EntityType should be present for modified resolution")
	}
}

func TestResolveQuarantined_Modified_MissingData(t *testing.T) {
	d := newResolveQuarantinedDeps()
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	qr, _ := seedQuarantinedPair(d)

	_, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  qr.ID(),
		Resolution:    model.QuarantineResolutionModified,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
		Notes:         "corrected",
		ModifiedData:  nil, // missing required data
	})
	if err == nil {
		t.Fatal("Handle() expected error for modified resolution without data")
	}
	if !errors.Is(err, domainerror.ErrPayloadRequired) {
		t.Errorf("error = %v, want ErrPayloadRequired", err)
	}
}

// =============================================================================
// Reject
// =============================================================================

func TestResolveQuarantined_Reject_Success(t *testing.T) {
	d := newResolveQuarantinedDeps()
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	qr, _ := seedQuarantinedPair(d)

	result, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  qr.ID(),
		Resolution:    model.QuarantineResolutionRejected,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
		Notes:         "data confirmed invalid",
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if !result.IngestRecord.IsRejected() {
		t.Error("ingest record should be rejected")
	}
	if result.EntityType.IsPresent() {
		t.Error("EntityType should not be present for rejected resolution")
	}
	if result.EntityID.IsPresent() {
		t.Error("EntityID should not be present for rejected resolution")
	}
}

// =============================================================================
// Validation errors
// =============================================================================

func TestResolveQuarantined_MissingQuarantineID(t *testing.T) {
	d := newResolveQuarantinedDeps()
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  "",
		Resolution:    model.QuarantineResolutionApproved,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
	})
	if err == nil {
		t.Fatal("Handle() expected error for missing quarantine ID")
	}
	if !errors.Is(err, domainerror.ErrRecordIDRequired) {
		t.Errorf("error = %v, want ErrRecordIDRequired", err)
	}
}

func TestResolveQuarantined_InvalidResolution(t *testing.T) {
	d := newResolveQuarantinedDeps()
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  types.NewID(),
		Resolution:    model.QuarantineResolutionUnspecified,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
	})
	if err == nil {
		t.Fatal("Handle() expected error for invalid resolution")
	}
	if !errors.Is(err, domainerror.ErrResolutionInvalid) {
		t.Errorf("error = %v, want ErrResolutionInvalid", err)
	}
}

func TestResolveQuarantined_PendingResolution(t *testing.T) {
	d := newResolveQuarantinedDeps()
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  types.NewID(),
		Resolution:    model.QuarantineResolutionPending,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
	})
	if err == nil {
		t.Fatal("Handle() expected error for pending resolution")
	}
	if !errors.Is(err, domainerror.ErrResolutionInvalid) {
		t.Errorf("error = %v, want ErrResolutionInvalid", err)
	}
}

func TestResolveQuarantined_NotFound(t *testing.T) {
	d := newResolveQuarantinedDeps()
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  types.NewID(),
		Resolution:    model.QuarantineResolutionApproved,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
	})
	if err == nil {
		t.Fatal("Handle() expected error for not found quarantine")
	}
}

// =============================================================================
// Repository errors
// =============================================================================

func TestResolveQuarantined_QuarantineRepoFindError(t *testing.T) {
	d := newResolveQuarantinedDeps()
	d.quarantineRepo.Errors.FindByID = errors.New("db error")
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  types.NewID(),
		Resolution:    model.QuarantineResolutionApproved,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
	})
	if err == nil {
		t.Fatal("Handle() expected error for quarantine repo find error")
	}
}

func TestResolveQuarantined_QuarantineRepoUpdateError(t *testing.T) {
	d := newResolveQuarantinedDeps()
	d.quarantineRepo.Errors.Update = errors.New("db error")
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	qr, _ := seedQuarantinedPair(d)

	_, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  qr.ID(),
		Resolution:    model.QuarantineResolutionApproved,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
	})
	if err == nil {
		t.Fatal("Handle() expected error for quarantine repo update error")
	}
}

func TestResolveQuarantined_RecordRepoUpdateError(t *testing.T) {
	d := newResolveQuarantinedDeps()
	d.recordRepo.Errors.Update = errors.New("db error")
	handler := newResolveQuarantinedHandler(d)
	ctx := context.Background()

	qr, _ := seedQuarantinedPair(d)

	_, err := handler.Handle(ctx, command.ResolveQuarantined{
		QuarantineID:  qr.ID(),
		Resolution:    model.QuarantineResolutionApproved,
		ResolvedBy:    "analyst@example.com",
		ResolvedByDID: "did:key:z6MkAnalyst",
	})
	if err == nil {
		t.Fatal("Handle() expected error for record repo update error")
	}
}
