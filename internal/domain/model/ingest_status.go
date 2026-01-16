package model

import (
	"fmt"
)

type IngestStatus int

const (
	IngestStatusUnspecified IngestStatus = iota
	IngestStatusPending
	IngestStatusAccepted
	IngestStatusQuarantined
	IngestStatusRejected
)

var ingestStatusNames = map[IngestStatus]string{
	IngestStatusUnspecified: "unspecified",
	IngestStatusPending:     "pending",
	IngestStatusAccepted:    "accepted",
	IngestStatusQuarantined: "quarantined",
	IngestStatusRejected:    "rejected",
}

var ingestStatusValues = map[string]IngestStatus{
	"unspecified": IngestStatusUnspecified,
	"pending":     IngestStatusPending,
	"accepted":    IngestStatusAccepted,
	"quarantined": IngestStatusQuarantined,
	"rejected":    IngestStatusRejected,
}

func (s IngestStatus) String() string {
	if name, ok := ingestStatusNames[s]; ok {
		return name
	}
	return fmt.Sprintf("IngestStatus(%d)", s)
}

func (s IngestStatus) IsValid() bool {
	_, ok := ingestStatusNames[s]
	return ok && s != IngestStatusUnspecified
}

func ParseIngestStatus(s string) (IngestStatus, error) {
	if status, ok := ingestStatusValues[s]; ok {
		return status, nil
	}
	return IngestStatusUnspecified, fmt.Errorf("invalid ingest status: %s", s)
}

func (s IngestStatus) IsTerminal() bool {
	return s == IngestStatusAccepted || s == IngestStatusRejected
}

func (s IngestStatus) IsPending() bool {
	return s == IngestStatusPending
}

func (s IngestStatus) IsQuarantined() bool {
	return s == IngestStatusQuarantined
}

func (s IngestStatus) IsAccepted() bool {
	return s == IngestStatusAccepted
}

func (s IngestStatus) IsRejected() bool {
	return s == IngestStatusRejected
}

func (s IngestStatus) CanTransitionTo(target IngestStatus) bool {
	switch s {
	case IngestStatusPending:
		return target == IngestStatusAccepted ||
			target == IngestStatusQuarantined ||
			target == IngestStatusRejected
	case IngestStatusQuarantined:
		return target == IngestStatusAccepted ||
			target == IngestStatusRejected
	case IngestStatusAccepted, IngestStatusRejected:
		return false
	default:
		return false
	}
}
