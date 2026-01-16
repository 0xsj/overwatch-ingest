package error

import (
	"github.com/0xsj/overwatch-pkg/errors"
)

const (
	// Record errors
	CodeRecordNotFound      = "INGEST_RECORD_NOT_FOUND"
	CodeRecordAlreadyExists = "INGEST_RECORD_ALREADY_EXISTS"
	CodeRecordIDRequired    = "INGEST_RECORD_ID_REQUIRED"
	CodeRawDataIDRequired   = "RAW_DATA_ID_REQUIRED"
	CodeSourceIDRequired    = "SOURCE_ID_REQUIRED"
	CodeSourceTypeRequired  = "SOURCE_TYPE_REQUIRED"
	CodeSourceTypeInvalid   = "SOURCE_TYPE_INVALID"
	CodePayloadRequired     = "PAYLOAD_REQUIRED"
	CodePayloadInvalid      = "PAYLOAD_INVALID"

	// Status errors
	CodeRecordAlreadyProcessed  = "RECORD_ALREADY_PROCESSED"
	CodeRecordNotPending        = "RECORD_NOT_PENDING"
	CodeRecordNotQuarantined    = "RECORD_NOT_QUARANTINED"
	CodeInvalidStatusTransition = "INVALID_STATUS_TRANSITION"

	// Validation errors
	CodeValidationFailed   = "VALIDATION_FAILED"
	CodeSchemaInvalid      = "SCHEMA_INVALID"
	CodeAnomalyDetected    = "ANOMALY_DETECTED"
	CodeFieldRequired      = "FIELD_REQUIRED"
	CodeFieldInvalidFormat = "FIELD_INVALID_FORMAT"
	CodeFieldOutOfRange    = "FIELD_OUT_OF_RANGE"
	CodeTemporalAnomaly    = "TEMPORAL_ANOMALY"
	CodeDuplicateDetected  = "DUPLICATE_DETECTED"

	// Confidence errors
	CodeConfidenceTooLow  = "CONFIDENCE_TOO_LOW"
	CodeConfidenceInvalid = "CONFIDENCE_INVALID"

	// Quarantine errors
	CodeQuarantineNotFound        = "QUARANTINE_NOT_FOUND"
	CodeQuarantineAlreadyResolved = "QUARANTINE_ALREADY_RESOLVED"
	CodeQuarantineExpired         = "QUARANTINE_EXPIRED"
	CodeResolutionInvalid         = "RESOLUTION_INVALID"

	// Provenance errors
	CodeSignatureRequired         = "SIGNATURE_REQUIRED"
	CodeSignatureInvalid          = "SIGNATURE_INVALID"
	CodeSignatureVerifyFailed     = "SIGNATURE_VERIFY_FAILED"
	CodeProvenanceChainBroken     = "PROVENANCE_CHAIN_BROKEN"
	CodeCollectorSignatureInvalid = "COLLECTOR_SIGNATURE_INVALID"
	CodeSourceSignatureInvalid    = "SOURCE_SIGNATURE_INVALID"

	// Source reliability errors
	CodeSourceReliabilityNotFound = "SOURCE_RELIABILITY_NOT_FOUND"

	// Processing errors
	CodeProcessingFailed = "PROCESSING_FAILED"
	CodeRoutingFailed    = "ROUTING_FAILED"
	CodePublishFailed    = "PUBLISH_FAILED"
)

var (
	ErrRecordNotFound      = errors.New(errors.KindNotFound, CodeRecordNotFound, "ingest record not found")
	ErrRecordAlreadyExists = errors.New(errors.KindConflict, CodeRecordAlreadyExists, "ingest record already exists")
	ErrRecordIDRequired    = errors.New(errors.KindValidation, CodeRecordIDRequired, "ingest record id is required")
	ErrRawDataIDRequired   = errors.New(errors.KindValidation, CodeRawDataIDRequired, "raw data id is required")
	ErrSourceIDRequired    = errors.New(errors.KindValidation, CodeSourceIDRequired, "source id is required")
	ErrSourceTypeRequired  = errors.New(errors.KindValidation, CodeSourceTypeRequired, "source type is required")
	ErrSourceTypeInvalid   = errors.New(errors.KindValidation, CodeSourceTypeInvalid, "invalid source type")
	ErrPayloadRequired     = errors.New(errors.KindValidation, CodePayloadRequired, "payload is required")
	ErrPayloadInvalid      = errors.New(errors.KindValidation, CodePayloadInvalid, "payload is invalid")
)

var (
	ErrRecordAlreadyProcessed  = errors.New(errors.KindConflict, CodeRecordAlreadyProcessed, "record has already been processed")
	ErrRecordNotPending        = errors.New(errors.KindDomain, CodeRecordNotPending, "record is not in pending status")
	ErrRecordNotQuarantined    = errors.New(errors.KindDomain, CodeRecordNotQuarantined, "record is not quarantined")
	ErrInvalidStatusTransition = errors.New(errors.KindDomain, CodeInvalidStatusTransition, "invalid status transition")
)

var (
	ErrValidationFailed   = errors.New(errors.KindValidation, CodeValidationFailed, "validation failed")
	ErrSchemaInvalid      = errors.New(errors.KindValidation, CodeSchemaInvalid, "schema validation failed")
	ErrAnomalyDetected    = errors.New(errors.KindValidation, CodeAnomalyDetected, "anomaly detected")
	ErrFieldRequired      = errors.New(errors.KindValidation, CodeFieldRequired, "required field missing")
	ErrFieldInvalidFormat = errors.New(errors.KindValidation, CodeFieldInvalidFormat, "field has invalid format")
	ErrFieldOutOfRange    = errors.New(errors.KindValidation, CodeFieldOutOfRange, "field value out of range")
	ErrTemporalAnomaly    = errors.New(errors.KindValidation, CodeTemporalAnomaly, "temporal anomaly detected")
	ErrDuplicateDetected  = errors.New(errors.KindConflict, CodeDuplicateDetected, "duplicate record detected")
)

var (
	ErrConfidenceTooLow  = errors.New(errors.KindDomain, CodeConfidenceTooLow, "confidence score too low")
	ErrConfidenceInvalid = errors.New(errors.KindValidation, CodeConfidenceInvalid, "invalid confidence score")
)

var (
	ErrQuarantineNotFound        = errors.New(errors.KindNotFound, CodeQuarantineNotFound, "quarantined record not found")
	ErrQuarantineAlreadyResolved = errors.New(errors.KindConflict, CodeQuarantineAlreadyResolved, "quarantine has already been resolved")
	ErrQuarantineExpired         = errors.New(errors.KindDomain, CodeQuarantineExpired, "quarantine has expired")
	ErrResolutionInvalid         = errors.New(errors.KindValidation, CodeResolutionInvalid, "invalid quarantine resolution")
)

var (
	ErrSignatureRequired         = errors.New(errors.KindValidation, CodeSignatureRequired, "signature is required")
	ErrSignatureInvalid          = errors.New(errors.KindValidation, CodeSignatureInvalid, "signature is invalid")
	ErrSignatureVerifyFailed     = errors.New(errors.KindDomain, CodeSignatureVerifyFailed, "signature verification failed")
	ErrProvenanceChainBroken     = errors.New(errors.KindDomain, CodeProvenanceChainBroken, "provenance chain is broken")
	ErrCollectorSignatureInvalid = errors.New(errors.KindDomain, CodeCollectorSignatureInvalid, "collector signature is invalid")
	ErrSourceSignatureInvalid    = errors.New(errors.KindDomain, CodeSourceSignatureInvalid, "source signature is invalid")
)

var ErrSourceReliabilityNotFound = errors.New(errors.KindNotFound, CodeSourceReliabilityNotFound, "source reliability not found")

var (
	ErrProcessingFailed = errors.New(errors.KindInternal, CodeProcessingFailed, "processing failed")
	ErrRoutingFailed    = errors.New(errors.KindInternal, CodeRoutingFailed, "routing failed")
	ErrPublishFailed    = errors.New(errors.KindInternal, CodePublishFailed, "event publish failed")
)

func RecordNotFound(id string) *errors.Error {
	return ErrRecordNotFound.WithEntity("ingest_record", id)
}

func RecordAlreadyExists(rawDataID string) *errors.Error {
	return ErrRecordAlreadyExists.WithMeta("raw_data_id", rawDataID)
}

func QuarantineNotFound(id string) *errors.Error {
	return ErrQuarantineNotFound.WithEntity("quarantined_record", id)
}

func SourceReliabilityNotFound(sourceID string) *errors.Error {
	return ErrSourceReliabilityNotFound.WithEntity("source", sourceID)
}

func ValidationFailed(reason string) *errors.Error {
	return ErrValidationFailed.WithMeta("reason", reason)
}

func SchemaInvalid(reason string) *errors.Error {
	return ErrSchemaInvalid.WithMeta("reason", reason)
}

func FieldRequired(field string) *errors.Error {
	return ErrFieldRequired.WithMeta("field", field)
}

func FieldInvalidFormat(field, expected string) *errors.Error {
	return ErrFieldInvalidFormat.WithMeta("field", field).WithMeta("expected", expected)
}

func FieldOutOfRange(field string, min, max any) *errors.Error {
	return ErrFieldOutOfRange.
		WithMeta("field", field).
		WithMeta("min", min).
		WithMeta("max", max)
}

func TemporalAnomaly(field, reason string) *errors.Error {
	return ErrTemporalAnomaly.WithMeta("field", field).WithMeta("reason", reason)
}

func SignatureVerifyFailed(signerDID, reason string) *errors.Error {
	return ErrSignatureVerifyFailed.
		WithMeta("signer_did", signerDID).
		WithMeta("reason", reason)
}

func ProcessingFailed(reason string) *errors.Error {
	return ErrProcessingFailed.WithMeta("reason", reason)
}

func RoutingFailed(reason string) *errors.Error {
	return ErrRoutingFailed.WithMeta("reason", reason)
}

func PublishFailed(subject, reason string) *errors.Error {
	return ErrPublishFailed.WithMeta("subject", subject).WithMeta("reason", reason)
}
