// Package mocks provides mock implementations of ports for testing.
package mocks

import (
	"context"
	"sync"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

// IngestRecordRepository is a mock implementation of repository.IngestRecordRepository.
type IngestRecordRepository struct {
	mu sync.RWMutex

	// Storage
	records     map[string]*model.IngestRecord // by ID
	byRawDataID map[string]string              // rawDataID -> recordID

	// Call tracking
	Calls struct {
		Create            int
		Update            int
		FindByID          int
		FindByRawDataID   int
		List              int
		Count             int
		ExistsByRawDataID int
	}

	// Error injection
	Errors struct {
		Create            error
		Update            error
		FindByID          error
		FindByRawDataID   error
		List              error
		Count             error
		ExistsByRawDataID error
	}
}

// NewIngestRecordRepository creates a new mock IngestRecordRepository.
func NewIngestRecordRepository() *IngestRecordRepository {
	return &IngestRecordRepository{
		records:     make(map[string]*model.IngestRecord),
		byRawDataID: make(map[string]string),
	}
}

func (r *IngestRecordRepository) Create(_ context.Context, record *model.IngestRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls.Create++
	if r.Errors.Create != nil {
		return r.Errors.Create
	}
	id := record.ID().String()
	r.records[id] = record
	r.byRawDataID[record.RawDataID()] = id
	return nil
}

func (r *IngestRecordRepository) Update(_ context.Context, record *model.IngestRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls.Update++
	if r.Errors.Update != nil {
		return r.Errors.Update
	}
	r.records[record.ID().String()] = record
	return nil
}

func (r *IngestRecordRepository) FindByID(_ context.Context, id types.ID) (*model.IngestRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.FindByID++
	if r.Errors.FindByID != nil {
		return nil, r.Errors.FindByID
	}
	rec, ok := r.records[id.String()]
	if !ok {
		return nil, nil
	}
	return rec, nil
}

func (r *IngestRecordRepository) FindByRawDataID(_ context.Context, rawDataID string) (*model.IngestRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.FindByRawDataID++
	if r.Errors.FindByRawDataID != nil {
		return nil, r.Errors.FindByRawDataID
	}
	id, ok := r.byRawDataID[rawDataID]
	if !ok {
		return nil, nil
	}
	rec, ok := r.records[id]
	if !ok {
		return nil, nil
	}
	return rec, nil
}

func (r *IngestRecordRepository) List(_ context.Context, p repository.ListIngestRecordsParams) ([]*model.IngestRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.List++
	if r.Errors.List != nil {
		return nil, r.Errors.List
	}
	var res []*model.IngestRecord
	for _, rec := range r.records {
		if p.Status.IsPresent() && rec.Status() != p.Status.MustGet() {
			continue
		}
		if p.SourceID.IsPresent() && rec.SourceID() != p.SourceID.MustGet() {
			continue
		}
		if p.TenantID.IsPresent() && !rec.BelongsToTenant(p.TenantID.MustGet()) {
			continue
		}
		if p.SourceType.IsPresent() && rec.SourceType() != p.SourceType.MustGet() {
			continue
		}
		if p.EntityType.IsPresent() && rec.EntityType().OrElse("") != p.EntityType.MustGet() {
			continue
		}
		res = append(res, rec)
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

func (r *IngestRecordRepository) Count(_ context.Context, p repository.ListIngestRecordsParams) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.Count++
	if r.Errors.Count != nil {
		return 0, r.Errors.Count
	}
	var n int64
	for _, rec := range r.records {
		if p.Status.IsPresent() && rec.Status() != p.Status.MustGet() {
			continue
		}
		if p.SourceID.IsPresent() && rec.SourceID() != p.SourceID.MustGet() {
			continue
		}
		if p.TenantID.IsPresent() && !rec.BelongsToTenant(p.TenantID.MustGet()) {
			continue
		}
		if p.SourceType.IsPresent() && rec.SourceType() != p.SourceType.MustGet() {
			continue
		}
		if p.EntityType.IsPresent() && rec.EntityType().OrElse("") != p.EntityType.MustGet() {
			continue
		}
		n++
	}
	return n, nil
}

func (r *IngestRecordRepository) ExistsByRawDataID(_ context.Context, rawDataID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.ExistsByRawDataID++
	if r.Errors.ExistsByRawDataID != nil {
		return false, r.Errors.ExistsByRawDataID
	}
	_, exists := r.byRawDataID[rawDataID]
	return exists, nil
}

// Seed adds a record directly to storage for test setup.
func (r *IngestRecordRepository) Seed(rec *model.IngestRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := rec.ID().String()
	r.records[id] = rec
	r.byRawDataID[rec.RawDataID()] = id
}

// RecordCount returns the number of stored records.
func (r *IngestRecordRepository) RecordCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.records)
}
