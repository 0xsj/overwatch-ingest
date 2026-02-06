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
// GetRecord
// =============================================================================

func TestGetRecord_Success(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewGetRecordHandler(repo)
	ctx := context.Background()

	rec := testutil.Fixtures.IngestRecord()
	repo.Seed(rec)

	result, err := handler.Handle(ctx, query.GetRecord{ID: rec.ID()})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Record == nil {
		t.Fatal("result.Record is nil")
	}
	if result.Record.ID() != rec.ID() {
		t.Errorf("Record.ID = %v, want %v", result.Record.ID(), rec.ID())
	}
	if repo.Calls.FindByID != 1 {
		t.Errorf("Repo.FindByID calls = %d, want 1", repo.Calls.FindByID)
	}
}

func TestGetRecord_NotFound(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewGetRecordHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetRecord{ID: types.NewID()})
	if err == nil {
		t.Fatal("Handle() expected error for not found record")
	}
}

func TestGetRecord_MissingID(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewGetRecordHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetRecord{ID: ""})
	if err == nil {
		t.Fatal("Handle() expected error for missing ID")
	}
	if !errors.Is(err, domainerror.ErrRecordIDRequired) {
		t.Errorf("error = %v, want ErrRecordIDRequired", err)
	}
}

func TestGetRecord_RepoError(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	repo.Errors.FindByID = errors.New("db error")
	handler := appquery.NewGetRecordHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetRecord{ID: types.NewID()})
	if err == nil {
		t.Fatal("Handle() expected error for repo error")
	}
}

// =============================================================================
// GetRecordByRawData
// =============================================================================

func TestGetRecordByRawData_Success(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewGetRecordByRawDataHandler(repo)
	ctx := context.Background()

	rec := testutil.Fixtures.IngestRecord()
	repo.Seed(rec)

	result, err := handler.Handle(ctx, query.GetRecordByRawData{RawDataID: rec.RawDataID()})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Record == nil {
		t.Fatal("result.Record is nil")
	}
	if result.Record.RawDataID() != rec.RawDataID() {
		t.Errorf("Record.RawDataID = %v, want %v", result.Record.RawDataID(), rec.RawDataID())
	}
}

func TestGetRecordByRawData_NotFound(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewGetRecordByRawDataHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetRecordByRawData{RawDataID: "nonexistent"})
	if err == nil {
		t.Fatal("Handle() expected error for not found record")
	}
}

func TestGetRecordByRawData_MissingRawDataID(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewGetRecordByRawDataHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetRecordByRawData{RawDataID: ""})
	if err == nil {
		t.Fatal("Handle() expected error for missing raw data ID")
	}
	if !errors.Is(err, domainerror.ErrRawDataIDRequired) {
		t.Errorf("error = %v, want ErrRawDataIDRequired", err)
	}
}

// =============================================================================
// ListRecords
// =============================================================================

func TestListRecords_Success(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewListRecordsHandler(repo)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		repo.Seed(testutil.Fixtures.IngestRecord())
	}

	result, err := handler.Handle(ctx, query.DefaultListRecords())
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

func TestListRecords_Empty(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewListRecordsHandler(repo)
	ctx := context.Background()

	result, err := handler.Handle(ctx, query.DefaultListRecords())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.TotalCount != 0 {
		t.Errorf("TotalCount = %d, want 0", result.TotalCount)
	}
}

func TestListRecords_FilterByStatus(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewListRecordsHandler(repo)
	ctx := context.Background()

	// Seed accepted records
	for i := 0; i < 3; i++ {
		repo.Seed(testutil.Fixtures.AcceptedIngestRecord())
	}
	// Seed rejected records
	for i := 0; i < 2; i++ {
		repo.Seed(testutil.Fixtures.RejectedIngestRecord())
	}

	qry := query.DefaultListRecords().WithStatus(model.IngestStatusAccepted)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 3 {
		t.Errorf("Records count = %d, want 3", len(result.Records))
	}
	if result.TotalCount != 3 {
		t.Errorf("TotalCount = %d, want 3", result.TotalCount)
	}
}

func TestListRecords_FilterBySourceID(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewListRecordsHandler(repo)
	ctx := context.Background()

	sourceID := types.NewID()
	rec := testutil.Fixtures.IngestRecordBuilder().
		WithSourceID(sourceID).
		Build()
	repo.Seed(rec)
	repo.Seed(testutil.Fixtures.IngestRecord()) // different source

	qry := query.DefaultListRecords().WithSourceID(sourceID)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Errorf("Records count = %d, want 1", len(result.Records))
	}
}

func TestListRecords_FilterBySourceType(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewListRecordsHandler(repo)
	ctx := context.Background()

	aisRec := testutil.Fixtures.IngestRecordBuilder().
		WithSourceType("ais").
		Build()
	adsbRec := testutil.Fixtures.IngestRecordBuilder().
		WithSourceType("adsb").
		Build()
	repo.Seed(aisRec)
	repo.Seed(adsbRec)

	qry := query.DefaultListRecords().WithSourceType("ais")
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Records) != 1 {
		t.Errorf("Records count = %d, want 1", len(result.Records))
	}
}

func TestListRecords_Pagination(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	handler := appquery.NewListRecordsHandler(repo)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		repo.Seed(testutil.Fixtures.IngestRecord())
	}

	qry := query.DefaultListRecords().WithPagination(3, 0)
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

func TestListRecords_RepoError(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	repo.Errors.List = errors.New("db error")
	handler := appquery.NewListRecordsHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.DefaultListRecords())
	if err == nil {
		t.Fatal("Handle() expected error for repo error")
	}
}

func TestListRecords_CountError(t *testing.T) {
	repo := mocks.NewIngestRecordRepository()
	repo.Errors.Count = errors.New("db error")
	handler := appquery.NewListRecordsHandler(repo)
	ctx := context.Background()

	repo.Seed(testutil.Fixtures.IngestRecord())

	_, err := handler.Handle(ctx, query.DefaultListRecords())
	if err == nil {
		t.Fatal("Handle() expected error for count error")
	}
}
