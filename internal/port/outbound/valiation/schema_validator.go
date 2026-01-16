package validation

import (
	"context"

	"github.com/0xsj/overwatch-ingest/internal/domain/model"
)

type SchemaValidator interface {
	Validate(ctx context.Context, sourceType string, payload map[string]any) SchemaValidationResult
	ValidateWithSchema(ctx context.Context, schemaID string, payload map[string]any) SchemaValidationResult
	SupportsSourceType(sourceType string) bool
	Version() string
}

type SchemaValidationResult struct {
	Valid         bool
	FieldsPresent []string
	FieldsMissing []string
	Anomalies     []model.Anomaly
	SchemaID      string
	SchemaVersion string
}

func (r SchemaValidationResult) ToValidationResult(validatorVersion string) model.ValidationResult {
	return model.NewValidationResult(
		r.Valid,
		r.FieldsPresent,
		r.FieldsMissing,
		r.Anomalies,
		validatorVersion,
	)
}

func ValidSchemaResult(fieldsPresent []string, schemaID, schemaVersion string) SchemaValidationResult {
	return SchemaValidationResult{
		Valid:         true,
		FieldsPresent: fieldsPresent,
		FieldsMissing: nil,
		Anomalies:     nil,
		SchemaID:      schemaID,
		SchemaVersion: schemaVersion,
	}
}

func InvalidSchemaResult(fieldsPresent, fieldsMissing []string, anomalies []model.Anomaly, schemaID, schemaVersion string) SchemaValidationResult {
	return SchemaValidationResult{
		Valid:         false,
		FieldsPresent: fieldsPresent,
		FieldsMissing: fieldsMissing,
		Anomalies:     anomalies,
		SchemaID:      schemaID,
		SchemaVersion: schemaVersion,
	}
}
