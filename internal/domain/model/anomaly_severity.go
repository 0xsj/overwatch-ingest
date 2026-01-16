package model

import (
	"fmt"
)

type AnomalySeverity int

const (
	AnomalySeverityUnspecified AnomalySeverity = iota
	AnomalySeverityInfo
	AnomalySeverityWarning
	AnomalySeverityError
	AnomalySeverityCritical
)

var anomalySeverityNames = map[AnomalySeverity]string{
	AnomalySeverityUnspecified: "unspecified",
	AnomalySeverityInfo:        "info",
	AnomalySeverityWarning:     "warning",
	AnomalySeverityError:       "error",
	AnomalySeverityCritical:    "critical",
}

var anomalySeverityValues = map[string]AnomalySeverity{
	"unspecified": AnomalySeverityUnspecified,
	"info":        AnomalySeverityInfo,
	"warning":     AnomalySeverityWarning,
	"error":       AnomalySeverityError,
	"critical":    AnomalySeverityCritical,
}

func (s AnomalySeverity) String() string {
	if name, ok := anomalySeverityNames[s]; ok {
		return name
	}
	return fmt.Sprintf("AnomalySeverity(%d)", s)
}

func (s AnomalySeverity) IsValid() bool {
	_, ok := anomalySeverityNames[s]
	return ok && s != AnomalySeverityUnspecified
}

func ParseAnomalySeverity(str string) (AnomalySeverity, error) {
	if severity, ok := anomalySeverityValues[str]; ok {
		return severity, nil
	}
	return AnomalySeverityUnspecified, fmt.Errorf("invalid anomaly severity: %s", str)
}

func (s AnomalySeverity) Weight() float64 {
	switch s {
	case AnomalySeverityInfo:
		return 0.1
	case AnomalySeverityWarning:
		return 0.3
	case AnomalySeverityError:
		return 0.6
	case AnomalySeverityCritical:
		return 1.0
	default:
		return 0.0
	}
}

func (s AnomalySeverity) ShouldReject() bool {
	return s == AnomalySeverityCritical
}

func (s AnomalySeverity) ShouldQuarantine() bool {
	return s == AnomalySeverityError || s == AnomalySeverityCritical
}

func (s AnomalySeverity) IsHigherThan(other AnomalySeverity) bool {
	return s > other
}

func MaxSeverity(severities ...AnomalySeverity) AnomalySeverity {
	max := AnomalySeverityUnspecified
	for _, s := range severities {
		if s > max {
			max = s
		}
	}
	return max
}
