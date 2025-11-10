// services/incidents/internal/domain/incident/domain_events.go
package incident

import (
	"time"
)

// Event type constants
const (
	EventTypeIncidentReported   = "IncidentReported"
	EventTypeIncidentVerified   = "IncidentVerified"
	EventTypeIncidentDispatched = "IncidentDispatched"
	EventTypeIncidentUpdated    = "IncidentUpdated"
	EventTypeIncidentResolved   = "IncidentResolved"
	EventTypeIncidentClosed     = "IncidentClosed"
)

// IncidentReported is emitted when a new incident is reported.
type IncidentReported struct {
	BaseEvent
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Severity      string                 `json:"severity"`
	IncidentType  string                 `json:"incident_type"`
	Latitude      float64                `json:"latitude"`
	Longitude     float64                `json:"longitude"`
	Address       string                 `json:"address,omitempty"`
	ReportedBy    string                 `json:"reported_by"`
	ReportSource  string                 `json:"report_source"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	ReportedAt    time.Time              `json:"reported_at"`
}

// NewIncidentReported creates a new IncidentReported event.
func NewIncidentReported(
	aggregateID, title, description, severity, incidentType string,
	lat, lon float64, address, reportedBy, reportSource string,
	metadata map[string]interface{}, sequence int64,
) *IncidentReported {
	return &IncidentReported{
		BaseEvent:     NewBaseEvent(aggregateID, EventTypeIncidentReported, sequence),
		Title:         title,
		Description:   description,
		Severity:      severity,
		IncidentType:  incidentType,
		Latitude:      lat,
		Longitude:     lon,
		Address:       address,
		ReportedBy:    reportedBy,
		ReportSource:  reportSource,
		Metadata:      metadata,
		ReportedAt:    time.Now(),
	}
}

// IncidentVerified is emitted when an incident is verified as real.
type IncidentVerified struct {
	BaseEvent
	VerifiedBy string    `json:"verified_by,omitempty"`
	VerifiedAt time.Time `json:"verified_at"`
	Notes      string    `json:"notes,omitempty"`
}

// NewIncidentVerified creates a new IncidentVerified event.
func NewIncidentVerified(aggregateID, verifiedBy, notes string, sequence int64) *IncidentVerified {
	return &IncidentVerified{
		BaseEvent:  NewBaseEvent(aggregateID, EventTypeIncidentVerified, sequence),
		VerifiedBy: verifiedBy,
		VerifiedAt: time.Now(),
		Notes:      notes,
	}
}

// IncidentDispatched is emitted when responders are dispatched to the incident.
type IncidentDispatched struct {
	BaseEvent
	ResponderIDs []string  `json:"responder_ids"`
	DispatchedBy string    `json:"dispatched_by,omitempty"`
	DispatchedAt time.Time `json:"dispatched_at"`
	Notes        string    `json:"notes,omitempty"`
}

// NewIncidentDispatched creates a new IncidentDispatched event.
func NewIncidentDispatched(aggregateID string, responderIDs []string, dispatchedBy, notes string, sequence int64) *IncidentDispatched {
	return &IncidentDispatched{
		BaseEvent:    NewBaseEvent(aggregateID, EventTypeIncidentDispatched, sequence),
		ResponderIDs: responderIDs,
		DispatchedBy: dispatchedBy,
		DispatchedAt: time.Now(),
		Notes:        notes,
	}
}

// IncidentUpdated is emitted when incident details are updated.
type IncidentUpdated struct {
	BaseEvent
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	Severity     *string `json:"severity,omitempty"`
	UpdatedBy    string  `json:"updated_by,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdateReason string  `json:"update_reason,omitempty"`
}

// NewIncidentUpdated creates a new IncidentUpdated event.
func NewIncidentUpdated(aggregateID string, title, description, severity *string, updatedBy, reason string, sequence int64) *IncidentUpdated {
	return &IncidentUpdated{
		BaseEvent:    NewBaseEvent(aggregateID, EventTypeIncidentUpdated, sequence),
		Title:        title,
		Description:  description,
		Severity:     severity,
		UpdatedBy:    updatedBy,
		UpdatedAt:    time.Now(),
		UpdateReason: reason,
	}
}

// IncidentResolved is emitted when the incident is resolved.
type IncidentResolved struct {
	BaseEvent
	Resolution string    `json:"resolution"`
	ResolvedBy string    `json:"resolved_by,omitempty"`
	ResolvedAt time.Time `json:"resolved_at"`
}

// NewIncidentResolved creates a new IncidentResolved event.
func NewIncidentResolved(aggregateID, resolution, resolvedBy string, sequence int64) *IncidentResolved {
	return &IncidentResolved{
		BaseEvent:  NewBaseEvent(aggregateID, EventTypeIncidentResolved, sequence),
		Resolution: resolution,
		ResolvedBy: resolvedBy,
		ResolvedAt: time.Now(),
	}
}

// IncidentClosed is emitted when the incident is fully closed.
type IncidentClosed struct {
	BaseEvent
	ClosureNotes string    `json:"closure_notes,omitempty"`
	ClosedBy     string    `json:"closed_by,omitempty"`
	ClosedAt     time.Time `json:"closed_at"`
}

// NewIncidentClosed creates a new IncidentClosed event.
func NewIncidentClosed(aggregateID, closureNotes, closedBy string, sequence int64) *IncidentClosed {
	return &IncidentClosed{
		BaseEvent:    NewBaseEvent(aggregateID, EventTypeIncidentClosed, sequence),
		ClosureNotes: closureNotes,
		ClosedBy:     closedBy,
		ClosedAt:     time.Now(),
	}
}