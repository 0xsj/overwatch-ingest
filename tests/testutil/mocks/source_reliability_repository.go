package mocks

import (
	"context"
	"sync"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

// SourceReliabilityRepository is a mock implementation of repository.SourceReliabilityRepository.
type SourceReliabilityRepository struct {
	mu sync.RWMutex

	// Storage
	reliabilities map[string]*model.SourceReliability // by sourceID

	// Call tracking
	Calls struct {
		Create         int
		Update         int
		Upsert         int
		FindBySourceID int
		List           int
		Count          int
		Delete         int
	}

	// Error injection
	Errors struct {
		Create         error
		Update         error
		Upsert         error
		FindBySourceID error
		List           error
		Count          error
		Delete         error
	}
}

// NewSourceReliabilityRepository creates a new mock SourceReliabilityRepository.
func NewSourceReliabilityRepository() *SourceReliabilityRepository {
	return &SourceReliabilityRepository{
		reliabilities: make(map[string]*model.SourceReliability),
	}
}

func (r *SourceReliabilityRepository) Create(_ context.Context, reliability *model.SourceReliability) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls.Create++
	if r.Errors.Create != nil {
		return r.Errors.Create
	}
	r.reliabilities[reliability.SourceID().String()] = reliability
	return nil
}

func (r *SourceReliabilityRepository) Update(_ context.Context, reliability *model.SourceReliability) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls.Update++
	if r.Errors.Update != nil {
		return r.Errors.Update
	}
	r.reliabilities[reliability.SourceID().String()] = reliability
	return nil
}

func (r *SourceReliabilityRepository) Upsert(_ context.Context, reliability *model.SourceReliability) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls.Upsert++
	if r.Errors.Upsert != nil {
		return r.Errors.Upsert
	}
	r.reliabilities[reliability.SourceID().String()] = reliability
	return nil
}

func (r *SourceReliabilityRepository) FindBySourceID(_ context.Context, sourceID types.ID) (*model.SourceReliability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.FindBySourceID++
	if r.Errors.FindBySourceID != nil {
		return nil, r.Errors.FindBySourceID
	}
	rel, ok := r.reliabilities[sourceID.String()]
	if !ok {
		return nil, nil
	}
	return rel, nil
}

func (r *SourceReliabilityRepository) List(_ context.Context, p repository.ListSourceReliabilityParams) ([]*model.SourceReliability, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.List++
	if r.Errors.List != nil {
		return nil, r.Errors.List
	}
	var res []*model.SourceReliability
	for _, rel := range r.reliabilities {
		if p.TenantID.IsPresent() {
			relTenant := rel.TenantID()
			if relTenant.IsEmpty() || relTenant.MustGet() != p.TenantID.MustGet() {
				continue
			}
		}
		if p.MinScore.IsPresent() && rel.ReliabilityScore() < p.MinScore.MustGet() {
			continue
		}
		if p.MaxScore.IsPresent() && rel.ReliabilityScore() > p.MaxScore.MustGet() {
			continue
		}
		res = append(res, rel)
	}
	if p.Offset >= len(res) {
		return nil, nil
	}
	end := p.Offset + p.Limit
	if end > len(res) {
		end = len(res)
	}
	return res[p.Offset:end], nil
}

func (r *SourceReliabilityRepository) Count(_ context.Context, p repository.ListSourceReliabilityParams) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.Count++
	if r.Errors.Count != nil {
		return 0, r.Errors.Count
	}
	var n int64
	for _, rel := range r.reliabilities {
		if p.TenantID.IsPresent() {
			relTenant := rel.TenantID()
			if relTenant.IsEmpty() || relTenant.MustGet() != p.TenantID.MustGet() {
				continue
			}
		}
		if p.MinScore.IsPresent() && rel.ReliabilityScore() < p.MinScore.MustGet() {
			continue
		}
		if p.MaxScore.IsPresent() && rel.ReliabilityScore() > p.MaxScore.MustGet() {
			continue
		}
		n++
	}
	return n, nil
}

func (r *SourceReliabilityRepository) Delete(_ context.Context, sourceID types.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls.Delete++
	if r.Errors.Delete != nil {
		return r.Errors.Delete
	}
	delete(r.reliabilities, sourceID.String())
	return nil
}

// Seed adds a reliability record directly to storage for test setup.
func (r *SourceReliabilityRepository) Seed(rel *model.SourceReliability) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reliabilities[rel.SourceID().String()] = rel
}
