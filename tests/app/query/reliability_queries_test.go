package query_test

import (
	"context"
	"errors"
	"testing"

	"github.com/0xsj/overwatch-pkg/types"

	appquery "github.com/0xsj/overwatch-ingest/internal/app/query"
	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/query"
	"github.com/0xsj/overwatch-ingest/tests/testutil"
	"github.com/0xsj/overwatch-ingest/tests/testutil/mocks"
)

// =============================================================================
// GetSourceReliability
// =============================================================================

func TestGetSourceReliability_Found(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewGetSourceReliabilityHandler(repo)
	ctx := context.Background()

	sourceID := types.NewID()
	rel := testutil.Fixtures.ReliableSource(sourceID)
	repo.Seed(rel)

	result, err := handler.Handle(ctx, query.GetSourceReliability{
		SourceID: sourceID,
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Reliability == nil {
		t.Fatal("result.Reliability is nil")
	}
	if result.Reliability.SourceID() != sourceID {
		t.Errorf("Reliability.SourceID = %v, want %v", result.Reliability.SourceID(), sourceID)
	}
	if repo.Calls.FindBySourceID != 1 {
		t.Errorf("FindBySourceID calls = %d, want 1", repo.Calls.FindBySourceID)
	}
}

func TestGetSourceReliability_NotFound(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewGetSourceReliabilityHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetSourceReliability{
		SourceID: types.NewID(),
	})
	if err == nil {
		t.Fatal("Handle() expected error for not found source reliability")
	}
}

func TestGetSourceReliability_MissingSourceID(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewGetSourceReliabilityHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetSourceReliability{
		SourceID: "",
	})
	if err == nil {
		t.Fatal("Handle() expected error for missing source ID")
	}
	if !errors.Is(err, domainerror.ErrSourceIDRequired) {
		t.Errorf("error = %v, want ErrSourceIDRequired", err)
	}
}

func TestGetSourceReliability_RepoError(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	repo.Errors.FindBySourceID = errors.New("db error")
	handler := appquery.NewGetSourceReliabilityHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.GetSourceReliability{
		SourceID: types.NewID(),
	})
	if err == nil {
		t.Fatal("Handle() expected error from repository")
	}
}

func TestGetSourceReliability_ReturnsCorrectScore(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewGetSourceReliabilityHandler(repo)
	ctx := context.Background()

	sourceID := types.NewID()
	rel := testutil.Fixtures.ReliableSource(sourceID)
	repo.Seed(rel)

	result, err := handler.Handle(ctx, query.GetSourceReliability{
		SourceID: sourceID,
	})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Reliability.ReliabilityScore() <= 0 {
		t.Errorf("ReliabilityScore = %f, want > 0", result.Reliability.ReliabilityScore())
	}
}

// =============================================================================
// ListSourceReliability
// =============================================================================

func TestListSourceReliability_Default(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewListSourceReliabilityHandler(repo)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		sourceID := types.NewID()
		repo.Seed(testutil.Fixtures.ReliableSource(sourceID))
	}

	result, err := handler.Handle(ctx, query.DefaultListSourceReliability())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Reliabilities) != 3 {
		t.Errorf("len(Reliabilities) = %d, want 3", len(result.Reliabilities))
	}
	if result.TotalCount != 3 {
		t.Errorf("TotalCount = %d, want 3", result.TotalCount)
	}
}

func TestListSourceReliability_Empty(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewListSourceReliabilityHandler(repo)
	ctx := context.Background()

	result, err := handler.Handle(ctx, query.DefaultListSourceReliability())
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.TotalCount != 0 {
		t.Errorf("TotalCount = %d, want 0", result.TotalCount)
	}
}

func TestListSourceReliability_WithPagination(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewListSourceReliabilityHandler(repo)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		sourceID := types.NewID()
		repo.Seed(testutil.Fixtures.ReliableSource(sourceID))
	}

	qry := query.DefaultListSourceReliability().WithPagination(2, 0)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if len(result.Reliabilities) != 2 {
		t.Errorf("len(Reliabilities) = %d, want 2", len(result.Reliabilities))
	}
	if result.TotalCount != 5 {
		t.Errorf("TotalCount = %d, want 5", result.TotalCount)
	}
}

func TestListSourceReliability_RepoListError(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	repo.Errors.List = errors.New("db error")
	handler := appquery.NewListSourceReliabilityHandler(repo)
	ctx := context.Background()

	_, err := handler.Handle(ctx, query.DefaultListSourceReliability())
	if err == nil {
		t.Fatal("Handle() expected error from repository list")
	}
}

func TestListSourceReliability_RepoCountError(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	repo.Errors.Count = errors.New("db error")
	handler := appquery.NewListSourceReliabilityHandler(repo)
	ctx := context.Background()

	sourceID := types.NewID()
	repo.Seed(testutil.Fixtures.ReliableSource(sourceID))

	_, err := handler.Handle(ctx, query.DefaultListSourceReliability())
	if err == nil {
		t.Fatal("Handle() expected error from repository count")
	}
}

func TestListSourceReliability_WithMinScore(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewListSourceReliabilityHandler(repo)
	ctx := context.Background()

	reliableID := types.NewID()
	repo.Seed(testutil.Fixtures.ReliableSource(reliableID))

	unreliableID := types.NewID()
	repo.Seed(testutil.Fixtures.UnreliableSource(unreliableID))

	qry := query.DefaultListSourceReliability().WithMinScore(0.5)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	// The reliable source should have a high score, unreliable should be filtered
	if repo.Calls.List != 1 {
		t.Errorf("List calls = %d, want 1", repo.Calls.List)
	}
	_ = result
}

func TestListSourceReliability_WithMaxScore(t *testing.T) {
	repo := mocks.NewSourceReliabilityRepository()
	handler := appquery.NewListSourceReliabilityHandler(repo)
	ctx := context.Background()

	reliableID := types.NewID()
	repo.Seed(testutil.Fixtures.ReliableSource(reliableID))

	unreliableID := types.NewID()
	repo.Seed(testutil.Fixtures.UnreliableSource(unreliableID))

	qry := query.DefaultListSourceReliability().WithMaxScore(0.5)
	result, err := handler.Handle(ctx, qry)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if repo.Calls.List != 1 {
		t.Errorf("List calls = %d, want 1", repo.Calls.List)
	}
	_ = result
}
