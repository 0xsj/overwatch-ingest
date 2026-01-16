package model

import (
	"fmt"
)

type QuarantineReason int

const (
	QuarantineReasonUnspecified QuarantineReason = iota
	QuarantineReasonValidationFailed
	QuarantineReasonLowConfidence
	QuarantineReasonAnomalyDetected
	QuarantineReasonSignatureInvalid
	QuarantineReasonDuplicateSuspected
	QuarantineReasonManualReview
)

var quarantineReasonNames = map[QuarantineReason]string{
	QuarantineReasonUnspecified:        "unspecified",
	QuarantineReasonValidationFailed:   "validation_failed",
	QuarantineReasonLowConfidence:      "low_confidence",
	QuarantineReasonAnomalyDetected:    "anomaly_detected",
	QuarantineReasonSignatureInvalid:   "signature_invalid",
	QuarantineReasonDuplicateSuspected: "duplicate_suspected",
	QuarantineReasonManualReview:       "manual_review",
}

var quarantineReasonValues = map[string]QuarantineReason{
	"unspecified":         QuarantineReasonUnspecified,
	"validation_failed":   QuarantineReasonValidationFailed,
	"low_confidence":      QuarantineReasonLowConfidence,
	"anomaly_detected":    QuarantineReasonAnomalyDetected,
	"signature_invalid":   QuarantineReasonSignatureInvalid,
	"duplicate_suspected": QuarantineReasonDuplicateSuspected,
	"manual_review":       QuarantineReasonManualReview,
}

func (r QuarantineReason) String() string {
	if name, ok := quarantineReasonNames[r]; ok {
		return name
	}
	return fmt.Sprintf("QuarantineReason(%d)", r)
}

func (r QuarantineReason) IsValid() bool {
	_, ok := quarantineReasonNames[r]
	return ok && r != QuarantineReasonUnspecified
}

func ParseQuarantineReason(s string) (QuarantineReason, error) {
	if reason, ok := quarantineReasonValues[s]; ok {
		return reason, nil
	}
	return QuarantineReasonUnspecified, fmt.Errorf("invalid quarantine reason: %s", s)
}

func (r QuarantineReason) IsAutomatic() bool {
	return r == QuarantineReasonValidationFailed ||
		r == QuarantineReasonLowConfidence ||
		r == QuarantineReasonAnomalyDetected ||
		r == QuarantineReasonSignatureInvalid ||
		r == QuarantineReasonDuplicateSuspected
}

func (r QuarantineReason) IsManual() bool {
	return r == QuarantineReasonManualReview
}

func (r QuarantineReason) RequiresSecurityReview() bool {
	return r == QuarantineReasonSignatureInvalid
}

func (r QuarantineReason) Priority() int {
	switch r {
	case QuarantineReasonSignatureInvalid:
		return 1
	case QuarantineReasonValidationFailed:
		return 2
	case QuarantineReasonAnomalyDetected:
		return 3
	case QuarantineReasonLowConfidence:
		return 4
	case QuarantineReasonDuplicateSuspected:
		return 5
	case QuarantineReasonManualReview:
		return 6
	default:
		return 99
	}
}
