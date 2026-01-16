package validation

import (
	"context"

	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type AnomalyDetector interface {
	Detect(ctx context.Context, sourceType string, payload map[string]any, metadata DetectionMetadata) AnomalyDetectionResult
	SupportsSourceType(sourceType string) bool
	Version() string
}

type DetectionMetadata struct {
	SourceID        types.ID
	TenantID        types.Optional[types.ID]
	SourceTimestamp types.Optional[types.Timestamp]
	CollectedAt     types.Timestamp
	PreviousPayload map[string]any
}

type AnomalyDetectionResult struct {
	Anomalies        []model.Anomaly
	MaxSeverity      model.AnomalySeverity
	ShouldReject     bool
	ShouldQuarantine bool
	DetectorVersion  string
}

func (r AnomalyDetectionResult) HasAnomalies() bool {
	return len(r.Anomalies) > 0
}

func (r AnomalyDetectionResult) AnomalyCount() int {
	return len(r.Anomalies)
}

func (r AnomalyDetectionResult) CriticalCount() int {
	count := 0
	for _, a := range r.Anomalies {
		if a.Severity() == model.AnomalySeverityCritical {
			count++
		}
	}
	return count
}

func (r AnomalyDetectionResult) ErrorCount() int {
	count := 0
	for _, a := range r.Anomalies {
		if a.Severity() == model.AnomalySeverityError {
			count++
		}
	}
	return count
}

func (r AnomalyDetectionResult) WarningCount() int {
	count := 0
	for _, a := range r.Anomalies {
		if a.Severity() == model.AnomalySeverityWarning {
			count++
		}
	}
	return count
}

func NoAnomaliesResult(detectorVersion string) AnomalyDetectionResult {
	return AnomalyDetectionResult{
		Anomalies:        nil,
		MaxSeverity:      model.AnomalySeverityUnspecified,
		ShouldReject:     false,
		ShouldQuarantine: false,
		DetectorVersion:  detectorVersion,
	}
}

func NewAnomalyDetectionResult(anomalies []model.Anomaly, detectorVersion string) AnomalyDetectionResult {
	maxSeverity := model.AnomalySeverityUnspecified
	shouldReject := false
	shouldQuarantine := false

	for _, a := range anomalies {
		if a.Severity() > maxSeverity {
			maxSeverity = a.Severity()
		}
		if a.ShouldReject() {
			shouldReject = true
		}
		if a.ShouldQuarantine() {
			shouldQuarantine = true
		}
	}

	return AnomalyDetectionResult{
		Anomalies:        anomalies,
		MaxSeverity:      maxSeverity,
		ShouldReject:     shouldReject,
		ShouldQuarantine: shouldQuarantine && !shouldReject,
		DetectorVersion:  detectorVersion,
	}
}
