package model

import (
	"fmt"
)

type AnomalyType int

const (
	AnomalyTypeUnspecified AnomalyType = iota
	AnomalyTypeOutOfRange
	AnomalyTypeInvalidFormat
	AnomalyTypeMissingRequired
	AnomalyTypeUnexpectedValue
	AnomalyTypeTemporal
	AnomalyTypeStatistical
	AnomalyTypeDuplicate
	AnomalyTypeSuspicious
)

var anomalyTypeNames = map[AnomalyType]string{
	AnomalyTypeUnspecified:     "unspecified",
	AnomalyTypeOutOfRange:      "out_of_range",
	AnomalyTypeInvalidFormat:   "invalid_format",
	AnomalyTypeMissingRequired: "missing_required",
	AnomalyTypeUnexpectedValue: "unexpected_value",
	AnomalyTypeTemporal:        "temporal",
	AnomalyTypeStatistical:     "statistical",
	AnomalyTypeDuplicate:       "duplicate",
	AnomalyTypeSuspicious:      "suspicious",
}

var anomalyTypeValues = map[string]AnomalyType{
	"unspecified":      AnomalyTypeUnspecified,
	"out_of_range":     AnomalyTypeOutOfRange,
	"invalid_format":   AnomalyTypeInvalidFormat,
	"missing_required": AnomalyTypeMissingRequired,
	"unexpected_value": AnomalyTypeUnexpectedValue,
	"temporal":         AnomalyTypeTemporal,
	"statistical":      AnomalyTypeStatistical,
	"duplicate":        AnomalyTypeDuplicate,
	"suspicious":       AnomalyTypeSuspicious,
}

func (t AnomalyType) String() string {
	if name, ok := anomalyTypeNames[t]; ok {
		return name
	}
	return fmt.Sprintf("AnomalyType(%d)", t)
}

func (t AnomalyType) IsValid() bool {
	_, ok := anomalyTypeNames[t]
	return ok && t != AnomalyTypeUnspecified
}

func ParseAnomalyType(s string) (AnomalyType, error) {
	if typ, ok := anomalyTypeValues[s]; ok {
		return typ, nil
	}
	return AnomalyTypeUnspecified, fmt.Errorf("invalid anomaly type: %s", s)
}

func (t AnomalyType) IsDataQuality() bool {
	return t == AnomalyTypeOutOfRange ||
		t == AnomalyTypeInvalidFormat ||
		t == AnomalyTypeMissingRequired ||
		t == AnomalyTypeUnexpectedValue
}

func (t AnomalyType) IsBehavioral() bool {
	return t == AnomalyTypeTemporal ||
		t == AnomalyTypeStatistical ||
		t == AnomalyTypeSuspicious
}

func (t AnomalyType) RequiresHumanReview() bool {
	return t == AnomalyTypeSuspicious ||
		t == AnomalyTypeStatistical
}
