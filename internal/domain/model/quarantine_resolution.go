package model

import (
	"fmt"
)

type QuarantineResolution int

const (
	QuarantineResolutionUnspecified QuarantineResolution = iota
	QuarantineResolutionPending
	QuarantineResolutionApproved
	QuarantineResolutionModified
	QuarantineResolutionRejected
	QuarantineResolutionExpired
)

var quarantineResolutionNames = map[QuarantineResolution]string{
	QuarantineResolutionUnspecified: "unspecified",
	QuarantineResolutionPending:     "pending",
	QuarantineResolutionApproved:    "approved",
	QuarantineResolutionModified:    "modified",
	QuarantineResolutionRejected:    "rejected",
	QuarantineResolutionExpired:     "expired",
}

var quarantineResolutionValues = map[string]QuarantineResolution{
	"unspecified": QuarantineResolutionUnspecified,
	"pending":     QuarantineResolutionPending,
	"approved":    QuarantineResolutionApproved,
	"modified":    QuarantineResolutionModified,
	"rejected":    QuarantineResolutionRejected,
	"expired":     QuarantineResolutionExpired,
}

func (r QuarantineResolution) String() string {
	if name, ok := quarantineResolutionNames[r]; ok {
		return name
	}
	return fmt.Sprintf("QuarantineResolution(%d)", r)
}

func (r QuarantineResolution) IsValid() bool {
	_, ok := quarantineResolutionNames[r]
	return ok && r != QuarantineResolutionUnspecified
}

func ParseQuarantineResolution(s string) (QuarantineResolution, error) {
	if resolution, ok := quarantineResolutionValues[s]; ok {
		return resolution, nil
	}
	return QuarantineResolutionUnspecified, fmt.Errorf("invalid quarantine resolution: %s", s)
}

func (r QuarantineResolution) IsPending() bool {
	return r == QuarantineResolutionPending
}

func (r QuarantineResolution) IsResolved() bool {
	return r == QuarantineResolutionApproved ||
		r == QuarantineResolutionModified ||
		r == QuarantineResolutionRejected ||
		r == QuarantineResolutionExpired
}

func (r QuarantineResolution) IsAccepted() bool {
	return r == QuarantineResolutionApproved || r == QuarantineResolutionModified
}

func (r QuarantineResolution) IsRejected() bool {
	return r == QuarantineResolutionRejected || r == QuarantineResolutionExpired
}

func (r QuarantineResolution) IsManual() bool {
	return r == QuarantineResolutionApproved ||
		r == QuarantineResolutionModified ||
		r == QuarantineResolutionRejected
}

func (r QuarantineResolution) IsAutomatic() bool {
	return r == QuarantineResolutionExpired
}

func (r QuarantineResolution) RequiresModifiedData() bool {
	return r == QuarantineResolutionModified
}

func (r QuarantineResolution) ToIngestStatus() IngestStatus {
	switch r {
	case QuarantineResolutionApproved, QuarantineResolutionModified:
		return IngestStatusAccepted
	case QuarantineResolutionRejected, QuarantineResolutionExpired:
		return IngestStatusRejected
	default:
		return IngestStatusQuarantined
	}
}
