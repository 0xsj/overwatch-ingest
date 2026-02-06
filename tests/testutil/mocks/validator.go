package mocks

import (
	"context"
	"sync"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/validation"
)

// SchemaValidator is a mock implementation of validation.SchemaValidator.
type SchemaValidator struct {
	mu sync.RWMutex

	// Configurable result returned by Validate.
	Result validation.SchemaValidationResult

	// Version string returned by Version().
	VersionStr string

	Calls struct {
		Validate           int
		ValidateWithSchema int
		SupportsSourceType int
		Version            int
	}

	Errors struct {
		Validate error
	}
}

func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		Result: validation.ValidSchemaResult(
			[]string{"mmsi", "lat", "lon", "speed", "course", "timestamp"},
			"test-schema",
			"v1",
		),
		VersionStr: "v1.0.0-test",
	}
}

func (v *SchemaValidator) Validate(_ context.Context, _ string, _ map[string]any) validation.SchemaValidationResult {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Calls.Validate++
	return v.Result
}

func (v *SchemaValidator) ValidateWithSchema(_ context.Context, _ string, _ map[string]any) validation.SchemaValidationResult {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Calls.ValidateWithSchema++
	return v.Result
}

func (v *SchemaValidator) SupportsSourceType(_ string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Calls.SupportsSourceType++
	return true
}

func (v *SchemaValidator) Version() string {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Calls.Version++
	return v.VersionStr
}

// SetInvalid configures the validator to return an invalid schema result.
func (v *SchemaValidator) SetInvalid(fieldsMissing []string, anomalies []model.Anomaly) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Result = validation.InvalidSchemaResult(nil, fieldsMissing, anomalies, "test-schema", "v1")
}

// SetValid configures the validator to return a valid schema result.
func (v *SchemaValidator) SetValid(fieldsPresent []string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Result = validation.ValidSchemaResult(fieldsPresent, "test-schema", "v1")
}

// ─────────────────────────────────────────────────────────────────
// AnomalyDetector
// ─────────────────────────────────────────────────────────────────

// AnomalyDetector is a mock implementation of validation.AnomalyDetector.
type AnomalyDetector struct {
	mu sync.RWMutex

	// Configurable result returned by Detect.
	Result validation.AnomalyDetectionResult

	// Version string returned by Version().
	VersionStr string

	Calls struct {
		Detect             int
		SupportsSourceType int
		Version            int
	}
}

func NewAnomalyDetector() *AnomalyDetector {
	return &AnomalyDetector{
		Result:     validation.NoAnomaliesResult("v1.0.0-test"),
		VersionStr: "v1.0.0-test",
	}
}

func (d *AnomalyDetector) Detect(_ context.Context, _ string, _ map[string]any, _ validation.DetectionMetadata) validation.AnomalyDetectionResult {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Calls.Detect++
	return d.Result
}

func (d *AnomalyDetector) SupportsSourceType(_ string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Calls.SupportsSourceType++
	return true
}

func (d *AnomalyDetector) Version() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Calls.Version++
	return d.VersionStr
}

// SetAnomalies configures the detector to return specific anomalies.
func (d *AnomalyDetector) SetAnomalies(anomalies []model.Anomaly) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Result = validation.NewAnomalyDetectionResult(anomalies, d.VersionStr)
}

// SetNoAnomalies configures the detector to return no anomalies.
func (d *AnomalyDetector) SetNoAnomalies() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Result = validation.NoAnomaliesResult(d.VersionStr)
}

// ─────────────────────────────────────────────────────────────────
// ConfidenceScorer
// ─────────────────────────────────────────────────────────────────

// ConfidenceScorer is a mock implementation of validation.ConfidenceScorer.
type ConfidenceScorer struct {
	mu sync.RWMutex

	// Configurable result returned by Score.
	Result model.ConfidenceScore

	// Version string returned by Version().
	VersionStr string

	Calls struct {
		Score              int
		SupportsSourceType int
		Version            int
	}
}

func NewConfidenceScorer() *ConfidenceScorer {
	return &ConfidenceScorer{
		Result:     model.DefaultConfidenceScore(),
		VersionStr: "v1.0.0-test",
	}
}

func (s *ConfidenceScorer) Score(_ context.Context, _ validation.ScoringInput) model.ConfidenceScore {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Calls.Score++
	return s.Result
}

func (s *ConfidenceScorer) SupportsSourceType(_ string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Calls.SupportsSourceType++
	return true
}

func (s *ConfidenceScorer) Version() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Calls.Version++
	return s.VersionStr
}

// SetScore configures the scorer to return a specific confidence score.
func (s *ConfidenceScorer) SetScore(score model.ConfidenceScore) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Result = score
}

// SetHighConfidence configures the scorer to return a high confidence score.
func (s *ConfidenceScorer) SetHighConfidence() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Result = model.NewConfidenceScore(0.9, 0.95, 0.95, 0.9, nil)
}

// SetLowConfidence configures the scorer to return a low confidence score.
func (s *ConfidenceScorer) SetLowConfidence() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Result = model.NewConfidenceScore(0.1, 0.2, 0.1, 0.1, nil)
}

// SetMediumConfidence configures the scorer to return a medium confidence score (quarantine range).
func (s *ConfidenceScorer) SetMediumConfidence() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Result = model.NewConfidenceScore(0.5, 0.6, 0.5, 0.5, nil)
}
