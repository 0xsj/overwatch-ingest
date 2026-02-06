package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/0xsj/overwatch-pkg/types"

	appquery "github.com/0xsj/overwatch-ingest/internal/app/query"
	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/query"
	"github.com/0xsj/overwatch-ingest/tests/testutil"
	"github.com/0xsj/overwatch-ingest/tests/testutil/mocks"
)

// =============================================================================
// GetQuarantined
// =============================================================================

func TestGetQuarantined_Success(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetQuarantinedHandler(repo)
	ctx := context.Background()

	qr := testutil.Fixtures.QuarantinedRecord()
	repo.Seed(qr)

	result, err := handler.Handle(ctx, query.GetQuarantined{ID: qr.ID()})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Record == nil {
		t.Fatal("result.Record is nil")
	}
	if result.Record.ID() != qr.ID() {
		t.Errorf("Record.ID = %v, want %v", result.Record.ID(), qr.ID())
	}
	if repo.Calls.FindByID != 1 {
		t.Errorf("Repo.FindByID calls = %d, want 1", repo.Calls.FindByID)
	}
}

func TestGetQuarantined_NotFound(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetQuarantinedHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetQuarantined{ID: types.NewID()})
	if err == nil {
		t.Fatal("Handle() expected error for not found record")
	}
}

func TestGetQuarantined_MissingID(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetQuarantinedHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetQuarantined{ID: ""})
	if err == nil {
		t.Fatal("Handle() expected error for missing ID")
	}
	if !errors.Is(err, domainerror.ErrRecordIDRequired) {
		t.Errorf("error = %v, want ErrRecordIDRequired", err)
	}
}

func TestGetQuarantined_RepoError(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	repo.Errors.FindByID = errors.New("db error")
	handler := appquery.NewGetQuarantinedHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetQuarantined{ID: types.NewID()})
	if err == nil {
		t.Fatal("Handle() expected error for repo error")
	}
}

// =============================================================================
// GetQuarantinedByIngestRecord
// =============================================================================

func TestGetQuarantinedByIngestRecord_Success(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetQuarantinedByIngestRecordHandler(repo)
	ctx := context.Background()

	qr := testutil.Fixtures.QuarantinedRecord()
	repo.Seed(qr)

	result, err := handler.Handle(ctx, query.GetQuarantinedByIngestRecord{
		IngestRecordID: qr.IngestRecordID(),
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Record == nil {
		t.Fatal("result.Record is nil")
	}
	if result.Record.IngestRecordID() != qr.IngestRecordID() {
		t.Errorf("Record.IngestRecordID = %v, want %v", result.Record.IngestRecordID(), qr.IngestRecordID())
	}
}

func TestGetQuarantinedByIngestRecord_MissingID(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetQuarantinedByIngestRecordHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetQuarantinedByIngestRecord{
		IngestRecordID: "",
	})
	if err == nil {
		t.Fatal("Handle() expected error for missing ingest record ID")
	}
	if !errors.Is(err, domainerror.ErrRecordIDRequired) {
		t.Errorf("error = %v, want ErrRecordIDRequired", err)
	}
}

func TestGetQuarantinedByIngestRecord_NotFound(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewGetQuarantinedByIngestRecordHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetQuarantinedByIngestRecord{
		IngestRecordID: types.NewID(),
	})
	if err == nil {
		t.Fatal("Handle() expected error for not found record")
	}
}

// =============================================================================
// ListQuarantined
// =============================================================================

func TestListQuarantined_Success(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		repo.Seed(testutil.Fixtures.QuarantinedRecord())
	}

	result, err := handler.Handle(ctx, query.DefaultListQuarantined())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 5 {
		t.Errorf("Records count = %d, want 5", len(result.Records))
	}
	if result.TotalCount != 5 {
		t.Errorf("TotalCount = %d, want 5", result.TotalCount)
	}
}

func TestListQuarantined_Empty(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	result, err := handler.Handle(ctx, query.DefaultListQuarantined())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.TotalCount != 0 {
		t.Errorf("TotalCount = %d, want 0", result.TotalCount)
	}
}

func TestListQuarantined_FilterByReason(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	anomalyQR := testutil.Fixtures.QuarantinedRecordBuilder().
		WithReason(model.QuarantineReasonAnomalyDetected).
		Build()
	lowConfQR := testutil.Fixtures.QuarantinedRecordBuilder().
		WithReason(model.QuarantineReasonLowConfidence).
		Build()
	repo.Seed(anomalyQR)
	repo.Seed(lowConfQR)

	qry := query.DefaultListQuarantined().WithReason(model.QuarantineReasonAnomalyDetected)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Errorf("Records count = %d, want 1", len(result.Records))
	}
}

func TestListQuarantined_FilterByResolution(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	// All newly created quarantined records have resolution = pending
	for i := 0; i < 3; i++ {
		repo.Seed(testutil.Fixtures.QuarantinedRecord())
	}

	qry := query.DefaultListQuarantined().WithPendingOnly()
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 3 {
		t.Errorf("Records count = %d, want 3", len(result.Records))
	}
}

func TestListQuarantined_FilterBySourceID(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	sourceID := types.NewID()
	qr := testutil.Fixtures.QuarantinedRecordBuilder().
		WithSourceID(sourceID).
		Build()
	repo.Seed(qr)
	repo.Seed(testutil.Fixtures.QuarantinedRecord()) // different source

	qry := query.DefaultListQuarantined().WithSourceID(sourceID)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Errorf("Records count = %d, want 1", len(result.Records))
	}
}

func TestListQuarantined_FilterBySourceType(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	aisQR := testutil.Fixtures.QuarantinedRecordBuilder().
		WithSourceType("ais").
		Build()
	adsbQR := testutil.Fixtures.QuarantinedRecordBuilder().
		WithSourceType("adsb").
		Build()
	repo.Seed(aisQR)
	repo.Seed(adsbQR)

	qry := query.DefaultListQuarantined().WithSourceType("adsb")
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Errorf("Records count = %d, want 1", len(result.Records))
	}
}

func TestListQuarantined_Pagination(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		repo.Seed(testutil.Fixtures.QuarantinedRecord())
	}

	qry := query.DefaultListQuarantined().WithPagination(3, 0)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 3 {
		t.Errorf("Records count = %d, want 3", len(result.Records))
	}
	if result.TotalCount != 10 {
		t.Errorf("TotalCount = %d, want 10", result.TotalCount)
	}
}

func TestListQuarantined_RepoListError(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	repo.Errors.List = errors.New("db error")
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.DefaultListQuarantined())
	if err == nil {
		t.Fatal("Handle() expected error for repo list error")
	}
}

func TestListQuarantined_RepoCountError(t *testing.T) {
	repo := mocks.NewQuarantinedRecordRepository()
	repo.Errors.Count = errors.New("db error")
	handler := appquery.NewListQuarantinedHandler(repo)
	ctx := context.Background()

	repo.Seed(testutil.Fixtures.QuarantinedRecord())

	_, err := handler.Handle(ctx, query.DefaultListQuarantined())
	if err == nil {
		t.Fatal("Handle() expected error for repo count error")
	}
}
