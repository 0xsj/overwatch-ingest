package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/0xsj/overwatch-pkg/types"

	appcommand "github.com/0xsj/overwatch-ingest/internal/app/command"
	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
	"github.com/0xsj/overwatch-ingest/tests/testutil"
	"github.com/0xsj/overwatch-ingest/tests/testutil/mocks"
)

// =============================================================================
// ReprocessRecord
// =============================================================================

func TestReprocessRecord_SuccessByID(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	// ProcessRawDataHandler is a dependency but the current implementation
	// only looks up the record, so we pass nil (it's not called in the lookup path).
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	rec := testutil.Fixtures.IngestRecord()
	recordRepo.Seed(rec)

	result, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.Some(rec.ID()),
		RawDataID:      types.None[string](),
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.IngestRecord == nil {
		t.Fatal("result.IngestRecord is nil")
	}
	if result.IngestRecord.ID() != rec.ID() {
		t.Errorf("IngestRecord.ID = %v, want %v", result.IngestRecord.ID(), rec.ID())
	}
	if recordRepo.Calls.FindByID != 1 {
		t.Errorf("FindByID calls = %d, want 1", recordRepo.Calls.FindByID)
	}
}

func TestReprocessRecord_SuccessByRawDataID(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	rec := testutil.Fixtures.IngestRecord()
	recordRepo.Seed(rec)

	result, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.None[types.ID](),
		RawDataID:      types.Some(rec.RawDataID()),
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.IngestRecord == nil {
		t.Fatal("result.IngestRecord is nil")
	}
	if result.IngestRecord.RawDataID() != rec.RawDataID() {
		t.Errorf("IngestRecord.RawDataID = %v, want %v", result.IngestRecord.RawDataID(), rec.RawDataID())
	}
	if recordRepo.Calls.FindByRawDataID != 1 {
		t.Errorf("FindByRawDataID calls = %d, want 1", recordRepo.Calls.FindByRawDataID)
	}
}

func TestReprocessRecord_NotFound_ByID(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.Some(types.NewID()),
		RawDataID:      types.None[string](),
	})
	if err == nil {
		t.Fatal("Handle() expected error for not found record")
	}
}

func TestReprocessRecord_NotFound_ByRawDataID(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.None[types.ID](),
		RawDataID:      types.Some("nonexistent-raw-data-id"),
	})
	if err == nil {
		t.Fatal("Handle() expected error for not found record")
	}
}

func TestReprocessRecord_MissingBothIDs(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.None[types.ID](),
		RawDataID:      types.None[string](),
	})
	if err == nil {
		t.Fatal("Handle() expected error for missing both IDs")
	}
	if !errors.Is(err, domainerror.ErrRecordIDRequired) {
		t.Errorf("error = %v, want ErrRecordIDRequired", err)
	}
}

func TestReprocessRecord_RepoFindByIDError(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	recordRepo.Errors.FindByID = errors.New("db error")
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.Some(types.NewID()),
		RawDataID:      types.None[string](),
	})
	if err == nil {
		t.Fatal("Handle() expected error from repository")
	}
}

func TestReprocessRecord_RepoFindByRawDataIDError(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	recordRepo.Errors.FindByRawDataID = errors.New("db error")
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	_, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.None[types.ID](),
		RawDataID:      types.Some(testutil.Fake.RawDataID()),
	})
	if err == nil {
		t.Fatal("Handle() expected error from repository")
	}
}

func TestReprocessRecord_PrefersIDOverRawDataID(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	rec := testutil.Fixtures.IngestRecord()
	recordRepo.Seed(rec)

	// When both are present, IngestRecordID takes precedence
	result, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.Some(rec.ID()),
		RawDataID:      types.Some(rec.RawDataID()),
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.IngestRecord == nil {
		t.Fatal("result.IngestRecord is nil")
	}
	// FindByID should be called (not FindByRawDataID)
	if recordRepo.Calls.FindByID != 1 {
		t.Errorf("FindByID calls = %d, want 1", recordRepo.Calls.FindByID)
	}
	if recordRepo.Calls.FindByRawDataID != 0 {
		t.Errorf("FindByRawDataID calls = %d, want 0", recordRepo.Calls.FindByRawDataID)
	}
}

func TestReprocessRecord_AcceptedRecord(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	rec := testutil.Fixtures.AcceptedIngestRecord()
	recordRepo.Seed(rec)

	result, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.Some(rec.ID()),
		RawDataID:      types.None[string](),
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.IngestRecord == nil {
		t.Fatal("result.IngestRecord is nil")
	}
	if result.IngestRecord.ID() != rec.ID() {
		t.Errorf("IngestRecord.ID = %v, want %v", result.IngestRecord.ID(), rec.ID())
	}
}

func TestReprocessRecord_EmptyIngestRecordID(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	handler := appcommand.NewReprocessRecordHandler(recordRepo, nil)
	ctx := context.Background()

	// Empty ID (types.Some("")) is not the same as None -- types.ID("").IsEmpty() is true
	_, err := handler.Handle(ctx, command.ReprocessRecord{
		IngestRecordID: types.Some[types.ID](""),
		RawDataID:      types.None[string](),
	})
	// When IngestRecordID.IsPresent() is true but the value is empty,
	// it will try FindByID with empty string. The handler checks
	// cmd.IngestRecordID.IsEmpty() && cmd.RawDataID.IsEmpty() first.
	// types.Some[types.ID]("").IsEmpty() returns false (Optional is present),
	// so it proceeds to FindByID with empty ID. The mock returns nil.
	if err == nil {
		t.Fatal("Handle() expected error for empty ID")
	}
}
