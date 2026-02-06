package mocks

import (
	"context"
	"sync"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/repository"
)

// QuarantinedRecordRepository is a mock implementation of repository.QuarantinedRecordRepository.
type QuarantinedRecordRepository struct {
	mu sync.RWMutex

	// Storage
	records          map[string]*model.QuarantinedRecord // by ID
	byIngestRecordID map[string]string                   // ingestRecordID -> quarantineID

	// Call tracking
	Calls struct {
		Create               int
		Update               int
		FindByID             int
		FindByIngestRecordID int
		List                 int
		Count                int
		FindExpired          int
		CountPending         int
	}

	// Error injection
	Errors struct {
		Create               error
		Update               error
		FindByID             error
		FindByIngestRecordID error
		List                 error
		Count                error
		FindExpired          error
		CountPending         error
	}
}

// NewQuarantinedRecordRepository creates a new mock QuarantinedRecordRepository.
func NewQuarantinedRecordRepository() *QuarantinedRecordRepository {
	return &QuarantinedRecordRepository{
		records:          make(map[string]*model.QuarantinedRecord),
		byIngestRecordID: make(map[string]string),
	}
}

func (r *QuarantinedRecordRepository) Create(_ context.Context, record *model.QuarantinedRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls.Create++
	if r.Errors.Create != nil {
		return r.Errors.Create
	}
	id := record.ID().String()
	r.records[id] = record
	r.byIngestRecordID[record.IngestRecordID().String()] = id
	return nil
}

func (r *QuarantinedRecordRepository) Update(_ context.Context, record *model.QuarantinedRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls.Update++
	if r.Errors.Update != nil {
		return r.Errors.Update
	}
	r.records[record.ID().String()] = record
	return nil
}

func (r *QuarantinedRecordRepository) FindByID(_ context.Context, id types.ID) (*model.QuarantinedRecord, error) {
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

func (r *QuarantinedRecordRepository) FindByIngestRecordID(_ context.Context, ingestRecordID types.ID) (*model.QuarantinedRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.FindByIngestRecordID++
	if r.Errors.FindByIngestRecordID != nil {
		return nil, r.Errors.FindByIngestRecordID
	}
	id, ok := r.byIngestRecordID[ingestRecordID.String()]
	if !ok {
		return nil, nil
	}
	rec, ok := r.records[id]
	if !ok {
		return nil, nil
	}
	return rec, nil
}

func (r *QuarantinedRecordRepository) List(_ context.Context, p repository.ListQuarantinedRecordsParams) ([]*model.QuarantinedRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.List++
	if r.Errors.List != nil {
		return nil, r.Errors.List
	}
	var res []*model.QuarantinedRecord
	for _, rec := range r.records {
		if p.TenantID.IsPresent() && !rec.BelongsToTenant(p.TenantID.MustGet()) {
			continue
		}
		if p.SourceID.IsPresent() && rec.SourceID() != p.SourceID.MustGet() {
			continue
		}
		if p.SourceType.IsPresent() && rec.SourceType() != p.SourceType.MustGet() {
			continue
		}
		if p.Reason.IsPresent() && rec.Reason() != p.Reason.MustGet() {
			continue
		}
		if p.Resolution.IsPresent() && rec.Resolution() != p.Resolution.MustGet() {
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

func (r *QuarantinedRecordRepository) Count(_ context.Context, p repository.ListQuarantinedRecordsParams) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.Count++
	if r.Errors.Count != nil {
		return 0, r.Errors.Count
	}
	var n int64
	for _, rec := range r.records {
		if p.TenantID.IsPresent() && !rec.BelongsToTenant(p.TenantID.MustGet()) {
			continue
		}
		if p.SourceID.IsPresent() && rec.SourceID() != p.SourceID.MustGet() {
			continue
		}
		if p.SourceType.IsPresent() && rec.SourceType() != p.SourceType.MustGet() {
			continue
		}
		if p.Reason.IsPresent() && rec.Reason() != p.Reason.MustGet() {
			continue
		}
		if p.Resolution.IsPresent() && rec.Resolution() != p.Resolution.MustGet() {
			continue
		}
		n++
	}
	return n, nil
}

func (r *QuarantinedRecordRepository) FindExpired(_ context.Context, limit int) ([]*model.QuarantinedRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.FindExpired++
	if r.Errors.FindExpired != nil {
		return nil, r.Errors.FindExpired
	}
	var res []*model.QuarantinedRecord
	for _, rec := range r.records {
		if rec.IsExpired() {
			res = append(res, rec)
			if len(res) >= limit {
				break
			}
		}
	}
	return res, nil
}

func (r *QuarantinedRecordRepository) CountPending(_ context.Context, tenantID types.Optional[types.ID]) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.Calls.CountPending++
	if r.Errors.CountPending != nil {
		return 0, r.Errors.CountPending
	}
	var n int64
	for _, rec := range r.records {
		if !rec.IsPending() {
			continue
		}
		if tenantID.IsPresent() && !rec.BelongsToTenant(tenantID.MustGet()) {
			continue
		}
		n++
	}
	return n, nil
}

// Seed adds a quarantined record directly to storage for test setup.
func (r *QuarantinedRecordRepository) Seed(rec *model.QuarantinedRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := rec.ID().String()
	r.records[id] = rec
	r.byIngestRecordID[rec.IngestRecordID().String()] = id
}

// RecordCount returns the number of stored records.
func (r *QuarantinedRecordRepository) RecordCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.records)
}
