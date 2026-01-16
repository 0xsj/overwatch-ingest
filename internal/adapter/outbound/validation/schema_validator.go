package validation

import (
	"context"
	"fmt"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/validation"
)

const schemaValidatorVersion = "1.0.0"

type schemaValidator struct {
	registry *SourceTypeRegistry
}

func NewSchemaValidator(registry *SourceTypeRegistry) validation.SchemaValidator {
	if registry == nil {
		registry = NewSourceTypeRegistry()
	}
	return &schemaValidator{
		registry: registry,
	}
}

func (v *schemaValidator) Validate(ctx context.Context, sourceType string, payload map[string]any) validation.SchemaValidationResult {
	cfg := v.registry.GetOrDefault(sourceType)
	return v.validateAgainstConfig(payload, cfg)
}

func (v *schemaValidator) ValidateWithSchema(ctx context.Context, schemaID string, payload map[string]any) validation.SchemaValidationResult {
	for _, cfg := range v.registry.configs {
		if cfg.SchemaID == schemaID {
			return v.validateAgainstConfig(payload, cfg)
		}
	}

	return validation.InvalidSchemaResult(
		nil,
		nil,
		[]model.Anomaly{
			model.NewAnomaly(
				"schema",
				model.AnomalyTypeInvalidFormat,
				model.AnomalySeverityError,
				fmt.Sprintf("unknown schema ID: %s", schemaID),
			),
		},
		schemaID,
		"unknown",
	)
}

func (v *schemaValidator) SupportsSourceType(sourceType string) bool {
	return v.registry.Supports(sourceType)
}

func (v *schemaValidator) Version() string {
	return schemaValidatorVersion
}

func (v *schemaValidator) validateAgainstConfig(payload map[string]any, cfg SourceTypeConfig) validation.SchemaValidationResult {
	var fieldsPresent []string
	var fieldsMissing []string
	var anomalies []model.Anomaly

	allExpectedFields := append(cfg.RequiredFields, cfg.OptionalFields...)

	for _, field := range allExpectedFields {
		if _, ok := payload[field]; ok {
			fieldsPresent = append(fieldsPresent, field)
		}
	}

	for _, field := range cfg.RequiredFields {
		val, ok := payload[field]
		if !ok {
			fieldsMissing = append(fieldsMissing, field)
			anomalies = append(anomalies, model.NewAnomaly(
				field,
				model.AnomalyTypeMissingRequired,
				model.AnomalySeverityError,
				fmt.Sprintf("required field '%s' is missing", field),
			))
			continue
		}

		if isEmpty(val) {
			fieldsMissing = append(fieldsMissing, field)
			anomalies = append(anomalies, model.NewAnomalyWithExpected(
				field,
				model.AnomalyTypeMissingRequired,
				model.AnomalySeverityError,
				fmt.Sprintf("required field '%s' is empty", field),
				"non-empty value",
				"empty value",
			))
		}
	}

	for _, field := range cfg.OptionalFields {
		if val, ok := payload[field]; ok {
			if err := v.validateFieldType(field, val, cfg.SourceType); err != nil {
				anomalies = append(anomalies, model.NewAnomaly(
					field,
					model.AnomalyTypeInvalidFormat,
					model.AnomalySeverityWarning,
					err.Error(),
				))
			}
		}
	}

	valid := len(fieldsMissing) == 0 && !hasErrorOrCritical(anomalies)

	return validation.SchemaValidationResult{
		Valid:         valid,
		FieldsPresent: fieldsPresent,
		FieldsMissing: fieldsMissing,
		Anomalies:     anomalies,
		SchemaID:      cfg.SchemaID,
		SchemaVersion: cfg.SchemaVersion,
	}
}

func (v *schemaValidator) validateFieldType(field string, value any, sourceType string) error {
	switch field {
	case "latitude":
		lat, ok := toFloat64(value)
		if !ok {
			return fmt.Errorf("field '%s' must be a number", field)
		}
		if lat < -90 || lat > 90 {
			return fmt.Errorf("field '%s' must be between -90 and 90", field)
		}

	case "longitude":
		lon, ok := toFloat64(value)
		if !ok {
			return fmt.Errorf("field '%s' must be a number", field)
		}
		if lon < -180 || lon > 180 {
			return fmt.Errorf("field '%s' must be between -180 and 180", field)
		}

	case "mmsi":
		mmsi, ok := toString(value)
		if !ok {
			return fmt.Errorf("field '%s' must be a string", field)
		}
		if len(mmsi) != 9 {
			return fmt.Errorf("field '%s' must be 9 digits", field)
		}

	case "speed":
		speed, ok := toFloat64(value)
		if !ok {
			return fmt.Errorf("field '%s' must be a number", field)
		}
		if speed < 0 || speed > 102.2 {
			return fmt.Errorf("field '%s' must be between 0 and 102.2 knots", field)
		}

	case "course", "heading":
		angle, ok := toFloat64(value)
		if !ok {
			return fmt.Errorf("field '%s' must be a number", field)
		}
		if angle < 0 || angle >= 360 {
			return fmt.Errorf("field '%s' must be between 0 and 359.9", field)
		}

	case "frequency":
		freq, ok := toFloat64(value)
		if !ok {
			return fmt.Errorf("field '%s' must be a number", field)
		}
		if freq <= 0 {
			return fmt.Errorf("field '%s' must be positive", field)
		}

	case "cloud_cover":
		cc, ok := toFloat64(value)
		if !ok {
			return fmt.Errorf("field '%s' must be a number", field)
		}
		if cc < 0 || cc > 100 {
			return fmt.Errorf("field '%s' must be between 0 and 100", field)
		}

	case "resolution":
		res, ok := toFloat64(value)
		if !ok {
			return fmt.Errorf("field '%s' must be a number", field)
		}
		if res <= 0 {
			return fmt.Errorf("field '%s' must be positive", field)
		}
	}

	return nil
}

func isEmpty(val any) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case string:
		return v == ""
	case []any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	}

	return false
}

func hasErrorOrCritical(anomalies []model.Anomaly) bool {
	for _, a := range anomalies {
		if a.Severity() >= model.AnomalySeverityError {
			return true
		}
	}
	return false
}

func toFloat64(val any) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	}
	return 0, false
}

func toString(val any) (string, bool) {
	switch v := val.(type) {
	case string:
		return v, true
	case fmt.Stringer:
		return v.String(), true
	}
	return "", false
}
