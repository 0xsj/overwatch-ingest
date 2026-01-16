package model

type ValidationResult struct {
	valid            bool
	schemaValid      bool
	fieldsPresent    []string
	fieldsMissing    []string
	anomalies        []Anomaly
	validatorVersion string
}

func NewValidationResult(
	schemaValid bool,
	fieldsPresent []string,
	fieldsMissing []string,
	anomalies []Anomaly,
	validatorVersion string,
) ValidationResult {
	valid := schemaValid && !hasBlockingAnomalies(anomalies)

	return ValidationResult{
		valid:            valid,
		schemaValid:      schemaValid,
		fieldsPresent:    fieldsPresent,
		fieldsMissing:    fieldsMissing,
		anomalies:        anomalies,
		validatorVersion: validatorVersion,
	}
}

func ReconstructValidationResult(
	valid bool,
	schemaValid bool,
	fieldsPresent []string,
	fieldsMissing []string,
	anomalies []Anomaly,
	validatorVersion string,
) ValidationResult {
	return ValidationResult{
		valid:            valid,
		schemaValid:      schemaValid,
		fieldsPresent:    fieldsPresent,
		fieldsMissing:    fieldsMissing,
		anomalies:        anomalies,
		validatorVersion: validatorVersion,
	}
}

func (v ValidationResult) Valid() bool              { return v.valid }
func (v ValidationResult) SchemaValid() bool        { return v.schemaValid }
func (v ValidationResult) FieldsPresent() []string  { return v.fieldsPresent }
func (v ValidationResult) FieldsMissing() []string  { return v.fieldsMissing }
func (v ValidationResult) Anomalies() []Anomaly     { return v.anomalies }
func (v ValidationResult) ValidatorVersion() string { return v.validatorVersion }

func (v ValidationResult) AnomalyCount() int {
	return len(v.anomalies)
}

func (v ValidationResult) HasAnomalies() bool {
	return len(v.anomalies) > 0
}

func (v ValidationResult) HasCriticalAnomalies() bool {
	for _, a := range v.anomalies {
		if a.Severity() == AnomalySeverityCritical {
			return true
		}
	}
	return false
}

func (v ValidationResult) HasErrorAnomalies() bool {
	for _, a := range v.anomalies {
		if a.Severity() == AnomalySeverityError || a.Severity() == AnomalySeverityCritical {
			return true
		}
	}
	return false
}

func (v ValidationResult) MaxSeverity() AnomalySeverity {
	max := AnomalySeverityUnspecified
	for _, a := range v.anomalies {
		if a.Severity() > max {
			max = a.Severity()
		}
	}
	return max
}

func (v ValidationResult) AnomaliesBySeverity(severity AnomalySeverity) []Anomaly {
	var result []Anomaly
	for _, a := range v.anomalies {
		if a.Severity() == severity {
			result = append(result, a)
		}
	}
	return result
}

func (v ValidationResult) AnomaliesByType(typ AnomalyType) []Anomaly {
	var result []Anomaly
	for _, a := range v.anomalies {
		if a.Type() == typ {
			result = append(result, a)
		}
	}
	return result
}

func (v ValidationResult) AnomaliesForField(field string) []Anomaly {
	var result []Anomaly
	for _, a := range v.anomalies {
		if a.Field() == field {
			result = append(result, a)
		}
	}
	return result
}

func (v ValidationResult) ShouldReject() bool {
	if !v.schemaValid {
		return true
	}
	for _, a := range v.anomalies {
		if a.ShouldReject() {
			return true
		}
	}
	return false
}

func (v ValidationResult) ShouldQuarantine() bool {
	if v.ShouldReject() {
		return false
	}
	for _, a := range v.anomalies {
		if a.ShouldQuarantine() {
			return true
		}
	}
	return false
}

func (v ValidationResult) CompletenessScore() float64 {
	totalFields := len(v.fieldsPresent) + len(v.fieldsMissing)
	if totalFields == 0 {
		return 1.0
	}
	return float64(len(v.fieldsPresent)) / float64(totalFields)
}

func hasBlockingAnomalies(anomalies []Anomaly) bool {
	for _, a := range anomalies {
		if a.Severity() == AnomalySeverityError || a.Severity() == AnomalySeverityCritical {
			return true
		}
	}
	return false
}

func ValidResult(fieldsPresent []string, validatorVersion string) ValidationResult {
	return ValidationResult{
		valid:            true,
		schemaValid:      true,
		fieldsPresent:    fieldsPresent,
		fieldsMissing:    nil,
		anomalies:        nil,
		validatorVersion: validatorVersion,
	}
}

func InvalidSchemaResult(fieldsMissing []string, anomalies []Anomaly, validatorVersion string) ValidationResult {
	return ValidationResult{
		valid:            false,
		schemaValid:      false,
		fieldsPresent:    nil,
		fieldsMissing:    fieldsMissing,
		anomalies:        anomalies,
		validatorVersion: validatorVersion,
	}
}
