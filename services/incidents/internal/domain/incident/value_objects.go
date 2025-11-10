// services/incidents/internal/domain/incident/value_objects.go
package incident

import (
	"fmt"

	"github.com/google/uuid"
)

// IncidentID is a unique identifier for an incident.
type IncidentID struct {
	value string
}

// NewIncidentID creates a new IncidentID.
func NewIncidentID() IncidentID {
	return IncidentID{value: uuid.New().String()}
}

// ParseIncidentID parses a string into an IncidentID.
func ParseIncidentID(id string) (IncidentID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return IncidentID{}, fmt.Errorf("invalid incident id: %w", err)
	}
	return IncidentID{value: id}, nil
}

func (id IncidentID) String() string {
	return id.value
}

func (id IncidentID) Equals(other IncidentID) bool {
	return id.value == other.value
}

// IncidentStatus represents an incident's lifecycle status.
type IncidentStatus string

const (
	StatusReported    IncidentStatus = "reported"    // Initial report received
	StatusVerified    IncidentStatus = "verified"    // Confirmed as real
	StatusDispatched  IncidentStatus = "dispatched"  // Responders assigned
	StatusInProgress  IncidentStatus = "in_progress" // Actively being handled
	StatusResolved    IncidentStatus = "resolved"    // Emergency handled
	StatusClosed      IncidentStatus = "closed"      // Fully closed
)

// String returns the string representation.
func (s IncidentStatus) String() string {
	return string(s)
}

// IsValid checks if the status is valid.
func (s IncidentStatus) IsValid() bool {
	switch s {
	case StatusReported, StatusVerified, StatusDispatched, StatusInProgress, StatusResolved, StatusClosed:
		return true
	}
	return false
}

// ParseIncidentStatus parses a string into an IncidentStatus.
func ParseIncidentStatus(s string) (IncidentStatus, error) {
	status := IncidentStatus(s)
	if !status.IsValid() {
		return "", fmt.Errorf("invalid incident status: %s", s)
	}
	return status, nil
}

// CanTransitionTo checks if a status transition is valid.
func (s IncidentStatus) CanTransitionTo(newStatus IncidentStatus) bool {
	transitions := map[IncidentStatus][]IncidentStatus{
		StatusReported:   {StatusVerified, StatusClosed}, // Can verify or close (false alarm)
		StatusVerified:   {StatusDispatched, StatusClosed},
		StatusDispatched: {StatusInProgress, StatusResolved, StatusClosed},
		StatusInProgress: {StatusResolved, StatusClosed},
		StatusResolved:   {StatusClosed},
		StatusClosed:     {}, // Terminal state
	}

	validTransitions, ok := transitions[s]
	if !ok {
		return false
	}

	for _, validStatus := range validTransitions {
		if validStatus == newStatus {
			return true
		}
	}

	return false
}

// Severity represents the urgency/severity of an incident.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// String returns the string representation.
func (s Severity) String() string {
	return string(s)
}

// IsValid checks if the severity is valid.
func (s Severity) IsValid() bool {
	switch s {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
		return true
	}
	return false
}

// ParseSeverity parses a string into a Severity.
func ParseSeverity(s string) (Severity, error) {
	severity := Severity(s)
	if !severity.IsValid() {
		return "", fmt.Errorf("invalid severity: %s", s)
	}
	return severity, nil
}

// IncidentType represents the type of incident.
type IncidentType string

const (
	TypeFire            IncidentType = "fire"
	TypeMedical         IncidentType = "medical"
	TypeAccident        IncidentType = "accident"
	TypeCrime           IncidentType = "crime"
	TypeNaturalDisaster IncidentType = "natural_disaster"
	TypeHazmat          IncidentType = "hazmat"
	TypeOther           IncidentType = "other"
)

// String returns the string representation.
func (t IncidentType) String() string {
	return string(t)
}

// IsValid checks if the incident type is valid.
func (t IncidentType) IsValid() bool {
	switch t {
	case TypeFire, TypeMedical, TypeAccident, TypeCrime, TypeNaturalDisaster, TypeHazmat, TypeOther:
		return true
	}
	return false
}

// ParseIncidentType parses a string into an IncidentType.
func ParseIncidentType(s string) (IncidentType, error) {
	incidentType := IncidentType(s)
	if !incidentType.IsValid() {
		return "", fmt.Errorf("invalid incident type: %s", s)
	}
	return incidentType, nil
}

// Location represents a geographic location.
type Location struct {
	Latitude  float64
	Longitude float64
}

// NewLocation creates a new Location.
func NewLocation(lat, lon float64) (Location, error) {
	if lat < -90 || lat > 90 {
		return Location{}, fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", lat)
	}
	if lon < -180 || lon > 180 {
		return Location{}, fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", lon)
	}
	return Location{Latitude: lat, Longitude: lon}, nil
}

func (l Location) Equals(other Location) bool {
	return l.Latitude == other.Latitude && l.Longitude == other.Longitude
}

// ReportSource represents how the incident was reported.
type ReportSource string

const (
	SourceAgent     ReportSource = "agent"      // AI agent created it
	SourceHuman     ReportSource = "human"      // Human operator
	SourceAPI       ReportSource = "api"        // External API
	SourcePhone     ReportSource = "phone"      // Phone call (911)
	SourceMobileApp ReportSource = "mobile_app" // Mobile app (Citizen-style)
	SourceSystem    ReportSource = "system"     // Automated system
)

// String returns the string representation.
func (s ReportSource) String() string {
	return string(s)
}

// IsValid checks if the report source is valid.
func (s ReportSource) IsValid() bool {
	switch s {
	case SourceAgent, SourceHuman, SourceAPI, SourcePhone, SourceMobileApp, SourceSystem:
		return true
	}
	return false
}

// ParseReportSource parses a string into a ReportSource.
func ParseReportSource(s string) (ReportSource, error) {
	source := ReportSource(s)
	if !source.IsValid() {
		return "", fmt.Errorf("invalid report source: %s", s)
	}
	return source, nil
}