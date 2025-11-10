// services/incidents/internal/domain/incident/aggregate.go
package incident

import (
	"time"

	"github.com/0xsj/scout/platform/pkg/errors"
)

// Incident is the aggregate root for the Incident bounded context.
// It represents an emergency incident that requires response.
// State is rebuilt from events (Event Sourcing).
type Incident struct {
	// Identity
	id IncidentID

	// Description
	title       string
	description string

	// Classification
	severity     Severity
	incidentType IncidentType

	// Location
	location Location
	address  string // Optional human-readable address

	// Status
	status IncidentStatus

	// Source/Reporter
	reportedBy   string
	reportSource ReportSource

	// Metadata (flexible for future expansion)
	metadata map[string]interface{}

	// Timestamps
	reportedAt   time.Time
	verifiedAt   *time.Time
	dispatchedAt *time.Time
	resolvedAt   *time.Time
	closedAt     *time.Time

	// Event Sourcing metadata
	version           int64
	uncommittedEvents []Event
}

// NewIncident creates a new incident aggregate (for new incidents).
func NewIncident(id IncidentID) *Incident {
	return &Incident{
		id:                id,
		status:            "", // Will be set when reported
		metadata:          make(map[string]interface{}),
		version:           0,
		uncommittedEvents: make([]Event, 0),
	}
}

// LoadFromHistory rebuilds an incident aggregate from historical events.
func LoadFromHistory(id IncidentID, events []Event) *Incident {
	incident := NewIncident(id)
	for _, event := range events {
		incident.apply(event, false) // false = don't add to uncommitted events
	}
	return incident
}

// ID returns the incident's ID.
func (i *Incident) ID() IncidentID {
	return i.id
}

// Version returns the current version.
func (i *Incident) Version() int64 {
	return i.version
}

// UncommittedEvents returns events that have been generated but not persisted.
func (i *Incident) UncommittedEvents() []Event {
	return i.uncommittedEvents
}

// ClearUncommittedEvents clears the uncommitted events (called after persistence).
func (i *Incident) ClearUncommittedEvents() {
	i.uncommittedEvents = make([]Event, 0)
}

// nextSequenceNumber returns the next sequence number for a new event.
func (i *Incident) nextSequenceNumber() int64 {
	return i.version + 1
}

// apply applies an event to the aggregate's state.
// If recordEvent is true, adds to uncommittedEvents.
func (i *Incident) apply(event Event, recordEvent bool) {
	// Apply state changes based on event type
	switch e := event.(type) {
	case *IncidentReported:
		i.title = e.Title
		i.description = e.Description
		i.severity = Severity(e.Severity)
		i.incidentType = IncidentType(e.IncidentType)
		i.location = Location{Latitude: e.Latitude, Longitude: e.Longitude}
		i.address = e.Address
		i.reportedBy = e.ReportedBy
		i.reportSource = ReportSource(e.ReportSource)
		i.metadata = e.Metadata
		i.reportedAt = e.ReportedAt
		i.status = StatusReported

	case *IncidentVerified:
		i.status = StatusVerified
		i.verifiedAt = &e.VerifiedAt

	case *IncidentDispatched:
		i.status = StatusDispatched
		i.dispatchedAt = &e.DispatchedAt

	case *IncidentUpdated:
		if e.Title != nil {
			i.title = *e.Title
		}
		if e.Description != nil {
			i.description = *e.Description
		}
		if e.Severity != nil {
			i.severity = Severity(*e.Severity)
		}

	case *IncidentResolved:
		i.status = StatusResolved
		i.resolvedAt = &e.ResolvedAt

	case *IncidentClosed:
		i.status = StatusClosed
		i.closedAt = &e.ClosedAt
	}

	// Increment version
	i.version++

	// Record uncommitted event if needed
	if recordEvent {
		i.uncommittedEvents = append(i.uncommittedEvents, event)
	}
}

// === COMMANDS (Business Operations) ===

// Report reports a new incident.
func (i *Incident) Report(
	title, description string,
	severity Severity,
	incidentType IncidentType,
	location Location,
	address string,
	reportedBy string,
	reportSource ReportSource,
	metadata map[string]interface{},
) error {
	// Business rules validation
	if i.status != "" {
		return ErrIncidentAlreadyReported
	}

	if title == "" {
		return errors.RequiredField("title")
	}

	if !severity.IsValid() {
		return NewInvalidSeverityError(severity.String())
	}

	if !incidentType.IsValid() {
		return NewInvalidIncidentTypeError(incidentType.String())
	}

	if !reportSource.IsValid() {
		return NewInvalidReportSourceError(reportSource.String())
	}

	if reportedBy == "" {
		return errors.RequiredField("reported_by")
	}

	// Generate and apply event
	event := NewIncidentReported(
		i.id.String(),
		title,
		description,
		severity.String(),
		incidentType.String(),
		location.Latitude,
		location.Longitude,
		address,
		reportedBy,
		reportSource.String(),
		metadata,
		i.nextSequenceNumber(),
	)

	i.apply(event, true)

	return nil
}

// Verify verifies the incident as real (not a false alarm).
func (i *Incident) Verify(verifiedBy, notes string) error {
	// Business rules validation
	if i.status == "" {
		return ErrIncidentNotReported
	}

	if i.status == StatusClosed {
		return ErrIncidentAlreadyClosed
	}

	if i.status != StatusReported {
		return ErrIncidentAlreadyVerified
	}

	// Generate and apply event
	event := NewIncidentVerified(
		i.id.String(),
		verifiedBy,
		notes,
		i.nextSequenceNumber(),
	)

	i.apply(event, true)

	return nil
}

// Dispatch dispatches responders to the incident.
func (i *Incident) Dispatch(responderIDs []string, dispatchedBy, notes string) error {
	// Business rules validation
	if i.status == "" {
		return ErrIncidentNotReported
	}

	if i.status == StatusClosed {
		return ErrIncidentAlreadyClosed
	}

	if len(responderIDs) == 0 {
		return ErrNoRespondersProvided
	}

	// Generate and apply event
	event := NewIncidentDispatched(
		i.id.String(),
		responderIDs,
		dispatchedBy,
		notes,
		i.nextSequenceNumber(),
	)

	i.apply(event, true)

	return nil
}

// UpdateDetails updates incident details.
func (i *Incident) UpdateDetails(title, description *string, severity *Severity, updatedBy, reason string) error {
	// Business rules validation
	if i.status == "" {
		return ErrIncidentNotReported
	}

	if i.status == StatusClosed {
		return ErrIncidentAlreadyClosed
	}

	// At least one field must be updated
	if title == nil && description == nil && severity == nil {
		return errors.Validation("at least one field must be updated")
	}

	// Validate severity if provided
	if severity != nil && !severity.IsValid() {
		return NewInvalidSeverityError(severity.String())
	}

	// Convert severity to string pointer for event
	var severityStr *string
	if severity != nil {
		s := severity.String()
		severityStr = &s
	}

	// Generate and apply event
	event := NewIncidentUpdated(
		i.id.String(),
		title,
		description,
		severityStr,
		updatedBy,
		reason,
		i.nextSequenceNumber(),
	)

	i.apply(event, true)

	return nil
}

// Resolve marks the incident as resolved.
func (i *Incident) Resolve(resolution, resolvedBy string) error {
	// Business rules validation
	if i.status == "" {
		return ErrIncidentNotReported
	}

	if i.status == StatusClosed {
		return ErrIncidentAlreadyClosed
	}

	if resolution == "" {
		return errors.RequiredField("resolution")
	}

	// Generate and apply event
	event := NewIncidentResolved(
		i.id.String(),
		resolution,
		resolvedBy,
		i.nextSequenceNumber(),
	)

	i.apply(event, true)

	return nil
}

// Close closes the incident.
func (i *Incident) Close(closureNotes, closedBy string) error {
	// Business rules validation
	if i.status == "" {
		return ErrIncidentNotReported
	}

	if i.status == StatusClosed {
		return ErrIncidentAlreadyClosed
	}

	// Generate and apply event
	event := NewIncidentClosed(
		i.id.String(),
		closureNotes,
		closedBy,
		i.nextSequenceNumber(),
	)

	i.apply(event, true)

	return nil
}

// === GETTERS (Query State) ===

func (i *Incident) Title() string {
	return i.title
}

func (i *Incident) Description() string {
	return i.description
}

func (i *Incident) Severity() Severity {
	return i.severity
}

func (i *Incident) Type() IncidentType {
	return i.incidentType
}

func (i *Incident) Location() Location {
	return i.location
}

func (i *Incident) Address() string {
	return i.address
}

func (i *Incident) Status() IncidentStatus {
	return i.status
}

func (i *Incident) ReportedBy() string {
	return i.reportedBy
}

func (i *Incident) ReportSource() ReportSource {
	return i.reportSource
}

func (i *Incident) Metadata() map[string]interface{} {
	// Return a copy to prevent external modification
	metadata := make(map[string]interface{}, len(i.metadata))
	for k, v := range i.metadata {
		metadata[k] = v
	}
	return metadata
}

func (i *Incident) ReportedAt() time.Time {
	return i.reportedAt
}

func (i *Incident) VerifiedAt() *time.Time {
	return i.verifiedAt
}

func (i *Incident) DispatchedAt() *time.Time {
	return i.dispatchedAt
}

func (i *Incident) ResolvedAt() *time.Time {
	return i.resolvedAt
}

func (i *Incident) ClosedAt() *time.Time {
	return i.closedAt
}

func (i *Incident) IsReported() bool {
	return i.status != ""
}

func (i *Incident) IsVerified() bool {
	return i.status == StatusVerified || i.status == StatusDispatched || i.status == StatusInProgress || i.status == StatusResolved || i.status == StatusClosed
}

func (i *Incident) IsDispatched() bool {
	return i.status == StatusDispatched || i.status == StatusInProgress
}

func (i *Incident) IsResolved() bool {
	return i.status == StatusResolved
}

func (i *Incident) IsClosed() bool {
	return i.status == StatusClosed
}