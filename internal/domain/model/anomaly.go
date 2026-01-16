package model

import (
	"fmt"
)

type Anomaly struct {
	field    string
	typ      AnomalyType
	severity AnomalySeverity
	message  string
	expected string
	actual   string
	context  map[string]any
}

func NewAnomaly(
	field string,
	typ AnomalyType,
	severity AnomalySeverity,
	message string,
) Anomaly {
	return Anomaly{
		field:    field,
		typ:      typ,
		severity: severity,
		message:  message,
		context:  make(map[string]any),
	}
}

func NewAnomalyWithExpected(
	field string,
	typ AnomalyType,
	severity AnomalySeverity,
	message string,
	expected string,
	actual string,
) Anomaly {
	return Anomaly{
		field:    field,
		typ:      typ,
		severity: severity,
		message:  message,
		expected: expected,
		actual:   actual,
		context:  make(map[string]any),
	}
}

func ReconstructAnomaly(
	field string,
	typ AnomalyType,
	severity AnomalySeverity,
	message string,
	expected string,
	actual string,
	context map[string]any,
) Anomaly {
	if context == nil {
		context = make(map[string]any)
	}
	return Anomaly{
		field:    field,
		typ:      typ,
		severity: severity,
		message:  message,
		expected: expected,
		actual:   actual,
		context:  context,
	}
}

func (a Anomaly) Field() string             { return a.field }
func (a Anomaly) Type() AnomalyType         { return a.typ }
func (a Anomaly) Severity() AnomalySeverity { return a.severity }
func (a Anomaly) Message() string           { return a.message }
func (a Anomaly) Expected() string          { return a.expected }
func (a Anomaly) Actual() string            { return a.actual }
func (a Anomaly) Context() map[string]any   { return a.context }

func (a Anomaly) WithContext(key string, value any) Anomaly {
	newContext := make(map[string]any, len(a.context)+1)
	for k, v := range a.context {
		newContext[k] = v
	}
	newContext[key] = value

	return Anomaly{
		field:    a.field,
		typ:      a.typ,
		severity: a.severity,
		message:  a.message,
		expected: a.expected,
		actual:   a.actual,
		context:  newContext,
	}
}

func (a Anomaly) ShouldReject() bool {
	return a.severity.ShouldReject()
}

func (a Anomaly) ShouldQuarantine() bool {
	return a.severity.ShouldQuarantine()
}

func (a Anomaly) RequiresHumanReview() bool {
	return a.typ.RequiresHumanReview() || a.severity.ShouldQuarantine()
}

func OutOfRangeAnomaly(field string, severity AnomalySeverity, min, max, actual any) Anomaly {
	return Anomaly{
		field:    field,
		typ:      AnomalyTypeOutOfRange,
		severity: severity,
		message:  "value out of expected range",
		expected: formatRange(min, max),
		actual:   formatValue(actual),
		context:  map[string]any{"min": min, "max": max},
	}
}

func InvalidFormatAnomaly(field string, severity AnomalySeverity, expected, actual string) Anomaly {
	return Anomaly{
		field:    field,
		typ:      AnomalyTypeInvalidFormat,
		severity: severity,
		message:  "invalid format",
		expected: expected,
		actual:   actual,
		context:  make(map[string]any),
	}
}

func MissingRequiredAnomaly(field string, severity AnomalySeverity) Anomaly {
	return Anomaly{
		field:    field,
		typ:      AnomalyTypeMissingRequired,
		severity: severity,
		message:  "required field is missing",
		context:  make(map[string]any),
	}
}

func TemporalAnomaly(field string, severity AnomalySeverity, message string, timestamp any) Anomaly {
	return Anomaly{
		field:    field,
		typ:      AnomalyTypeTemporal,
		severity: severity,
		message:  message,
		actual:   formatValue(timestamp),
		context:  map[string]any{"timestamp": timestamp},
	}
}

func DuplicateAnomaly(field string, existingID string) Anomaly {
	return Anomaly{
		field:    field,
		typ:      AnomalyTypeDuplicate,
		severity: AnomalySeverityWarning,
		message:  "duplicate record detected",
		context:  map[string]any{"existing_id": existingID},
	}
}

func SuspiciousAnomaly(field string, severity AnomalySeverity, message string, context map[string]any) Anomaly {
	if context == nil {
		context = make(map[string]any)
	}
	return Anomaly{
		field:    field,
		typ:      AnomalyTypeSuspicious,
		severity: severity,
		message:  message,
		context:  context,
	}
}

func formatRange(min, max any) string {
	return formatValue(min) + " - " + formatValue(max)
}

func formatValue(v any) string {
	if v == nil {
		return "<nil>"
	}
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}
