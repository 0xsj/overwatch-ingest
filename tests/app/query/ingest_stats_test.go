package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/0xsj/overwatch-pkg/types"

	appquery "github.com/0xsj/overwatch-ingest/internal/app/query"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/query"
	"github.com/0xsj/overwatch-ingest/tests/testutil"
	"github.com/0xsj/overwatch-ingest/tests/testutil/mocks"
)

// =============================================================================
// GetIngestStats
// =============================================================================

func TestGetIngestStats_Success(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	quarantineRepo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetIngestStatsHandler(recordRepo, quarantineRepo)
	ctx := context.Background()

	// Seed accepted, rejected, and quarantined records
	for i := 0; i < 3; i++ {
		recordRepo.Seed(testutil.Fixtures.AcceptedIngestRecord())
	}
	for i := 0; i < 2; i++ {
		recordRepo.Seed(testutil.Fixtures.RejectedIngestRecord())
	}
	recordRepo.Seed(testutil.Fixtures.QuarantinedIngestRecord())

	result, err := handler.Handle(ctx, query.DefaultGetIngestStats())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Stats.TotalRecords != 6 {
		t.Errorf("TotalRecords = %d, want 6", result.Stats.TotalRecords)
	}
	if result.Stats.AcceptedRecords != 3 {
		t.Errorf("AcceptedRecords = %d, want 3", result.Stats.AcceptedRecords)
	}
	if result.Stats.RejectedRecords != 2 {
		t.Errorf("RejectedRecords = %d, want 2", result.Stats.RejectedRecords)
	}
	if result.Stats.QuarantinedRecords != 1 {
		t.Errorf("QuarantinedRecords = %d, want 1", result.Stats.QuarantinedRecords)
	}

	// Verify RecordsByStatus map
	if result.Stats.RecordsByStatus["accepted"] != 3 {
		t.Errorf("RecordsByStatus[accepted] = %d, want 3", result.Stats.RecordsByStatus["accepted"])
	}
	if result.Stats.RecordsByStatus["rejected"] != 2 {
		t.Errorf("RecordsByStatus[rejected] = %d, want 2", result.Stats.RecordsByStatus["rejected"])
	}
	if result.Stats.RecordsByStatus["quarantined"] != 1 {
		t.Errorf("RecordsByStatus[quarantined] = %d, want 1", result.Stats.RecordsByStatus["quarantined"])
	}
}

func TestGetIngestStats_Empty(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	quarantineRepo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetIngestStatsHandler(recordRepo, quarantineRepo)
	ctx := context.Background()

	result, err := handler.Handle(ctx, query.DefaultGetIngestStats())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Stats.TotalRecords != 0 {
		t.Errorf("TotalRecords = %d, want 0", result.Stats.TotalRecords)
	}
	if result.Stats.AcceptedRecords != 0 {
		t.Errorf("AcceptedRecords = %d, want 0", result.Stats.AcceptedRecords)
	}
	if result.Stats.RejectedRecords != 0 {
		t.Errorf("RejectedRecords = %d, want 0", result.Stats.RejectedRecords)
	}
	if result.Stats.QuarantinedRecords != 0 {
		t.Errorf("QuarantinedRecords = %d, want 0", result.Stats.QuarantinedRecords)
	}
	if result.Stats.PendingRecords != 0 {
		t.Errorf("PendingRecords = %d, want 0", result.Stats.PendingRecords)
	}
}

func TestGetIngestStats_WithSourceFilter(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	quarantineRepo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetIngestStatsHandler(recordRepo, quarantineRepo)
	ctx := context.Background()

	sourceID := types.NewID()
	rec := testutil.Fixtures.IngestRecordBuilder().
		WithSourceID(sourceID).
		Build()
	recordRepo.Seed(rec)
	recordRepo.Seed(testutil.Fixtures.IngestRecord()) // different source

	qry := query.DefaultGetIngestStats().WithSourceID(sourceID)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	// The mock Count filters by sourceID when provided
	if result.Stats.TotalRecords != 1 {
		t.Errorf("TotalRecords = %d, want 1", result.Stats.TotalRecords)
	}
}

func TestGetIngestStats_WithSourceTypeFilter(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	quarantineRepo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetIngestStatsHandler(recordRepo, quarantineRepo)
	ctx := context.Background()

	aisRec := testutil.Fixtures.IngestRecordBuilder().
		WithSourceType("ais").
		Build()
	adsbRec := testutil.Fixtures.IngestRecordBuilder().
		WithSourceType("adsb").
		Build()
	recordRepo.Seed(aisRec)
	recordRepo.Seed(adsbRec)

	qry := query.DefaultGetIngestStats().WithSourceType("ais")
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Stats.TotalRecords != 1 {
		t.Errorf("TotalRecords = %d, want 1", result.Stats.TotalRecords)
	}
}

func TestGetIngestStats_RepoCountError(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	recordRepo.Errors.Count = errors.New("db error")
	quarantineRepo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetIngestStatsHandler(recordRepo, quarantineRepo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.DefaultGetIngestStats())
	if err == nil {
		t.Fatal("Handle() expected error from repository count")
	}
}

func TestGetIngestStats_PendingRecords(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	quarantineRepo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetIngestStatsHandler(recordRepo, quarantineRepo)
	ctx := context.Background()

	// IngestRecord() creates a pending record by default
	for i := 0; i < 4; i++ {
		recordRepo.Seed(testutil.Fixtures.IngestRecord())
	}

	result, err := handler.Handle(ctx, query.DefaultGetIngestStats())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Stats.TotalRecords != 4 {
		t.Errorf("TotalRecords = %d, want 4", result.Stats.TotalRecords)
	}
	if result.Stats.PendingRecords != 4 {
		t.Errorf("PendingRecords = %d, want 4", result.Stats.PendingRecords)
	}
	if result.Stats.RecordsByStatus["pending"] != 4 {
		t.Errorf("RecordsByStatus[pending] = %d, want 4", result.Stats.RecordsByStatus["pending"])
	}
}

func TestGetIngestStats_MapsInitialized(t *testing.T) {
	recordRepo := mocks.NewIngestRecordRepository()
	quarantineRepo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetIngestStatsHandler(recordRepo, quarantineRepo)
	ctx := context.Background()

	result, err := handler.Handle(ctx, query.DefaultGetIngestStats())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Stats.RecordsBySource == nil {
		t.Error("RecordsBySource should not be nil")
	}
	if result.Stats.RecordsBySourceType == nil {
		t.Error("RecordsBySourceType should not be nil")
	}
	if result.Stats.RecordsByStatus == nil {
		t.Error("RecordsByStatus should not be nil")
	}
	if result.Stats.RecordsByEntityType == nil {
		t.Error("RecordsByEntityType should not be nil")
	}
	if result.Stats.AnomaliesByType == nil {
		t.Error("AnomaliesByType should not be nil")
	}
	if result.Stats.QuarantineByReason == nil {
		t.Error("QuarantineByReason should not be nil")
	}
}
