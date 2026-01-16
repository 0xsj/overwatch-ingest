package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
)

func toGRPCError(err error) error {
	if err == nil {
		return nil
	}

	// Record errors
	if errors.Is(err, domainerror.ErrRecordNotFound) {
		return status.Error(codes.NotFound, "ingest record not found")
	}
	if errors.Is(err, domainerror.ErrRecordAlreadyExists) {
		return status.Error(codes.AlreadyExists, "ingest record already exists")
	}
	if errors.Is(err, domainerror.ErrRecordIDRequired) {
		return status.Error(codes.InvalidArgument, "record_id is required")
	}
	if errors.Is(err, domainerror.ErrRawDataIDRequired) {
		return status.Error(codes.InvalidArgument, "raw_data_id is required")
	}
	if errors.Is(err, domainerror.ErrSourceIDRequired) {
		return status.Error(codes.InvalidArgument, "source_id is required")
	}
	if errors.Is(err, domainerror.ErrSourceTypeRequired) {
		return status.Error(codes.InvalidArgument, "source_type is required")
	}
	if errors.Is(err, domainerror.ErrSourceTypeInvalid) {
		return status.Error(codes.InvalidArgument, "source_type is invalid")
	}
	if errors.Is(err, domainerror.ErrPayloadRequired) {
		return status.Error(codes.InvalidArgument, "payload is required")
	}
	if errors.Is(err, domainerror.ErrPayloadInvalid) {
		return status.Error(codes.InvalidArgument, "payload is invalid")
	}

	// Status errors
	if errors.Is(err, domainerror.ErrRecordAlreadyProcessed) {
		return status.Error(codes.FailedPrecondition, "record has already been processed")
	}
	if errors.Is(err, domainerror.ErrRecordNotPending) {
		return status.Error(codes.FailedPrecondition, "record is not in pending status")
	}
	if errors.Is(err, domainerror.ErrRecordNotQuarantined) {
		return status.Error(codes.FailedPrecondition, "record is not quarantined")
	}
	if errors.Is(err, domainerror.ErrInvalidStatusTransition) {
		return status.Error(codes.FailedPrecondition, "invalid status transition")
	}

	// Validation errors
	if errors.Is(err, domainerror.ErrValidationFailed) {
		return status.Error(codes.InvalidArgument, "validation failed")
	}
	if errors.Is(err, domainerror.ErrSchemaInvalid) {
		return status.Error(codes.InvalidArgument, "schema validation failed")
	}
	if errors.Is(err, domainerror.ErrAnomalyDetected) {
		return status.Error(codes.FailedPrecondition, "anomaly detected")
	}
	if errors.Is(err, domainerror.ErrFieldRequired) {
		return status.Error(codes.InvalidArgument, "required field missing")
	}
	if errors.Is(err, domainerror.ErrFieldInvalidFormat) {
		return status.Error(codes.InvalidArgument, "field has invalid format")
	}
	if errors.Is(err, domainerror.ErrFieldOutOfRange) {
		return status.Error(codes.InvalidArgument, "field value out of range")
	}
	if errors.Is(err, domainerror.ErrTemporalAnomaly) {
		return status.Error(codes.FailedPrecondition, "temporal anomaly detected")
	}
	if errors.Is(err, domainerror.ErrDuplicateDetected) {
		return status.Error(codes.AlreadyExists, "duplicate record detected")
	}

	// Confidence errors
	if errors.Is(err, domainerror.ErrConfidenceTooLow) {
		return status.Error(codes.FailedPrecondition, "confidence score too low")
	}
	if errors.Is(err, domainerror.ErrConfidenceInvalid) {
		return status.Error(codes.InvalidArgument, "invalid confidence score")
	}

	// Quarantine errors
	if errors.Is(err, domainerror.ErrQuarantineNotFound) {
		return status.Error(codes.NotFound, "quarantined record not found")
	}
	if errors.Is(err, domainerror.ErrQuarantineAlreadyResolved) {
		return status.Error(codes.FailedPrecondition, "quarantine has already been resolved")
	}
	if errors.Is(err, domainerror.ErrQuarantineExpired) {
		return status.Error(codes.FailedPrecondition, "quarantine has expired")
	}
	if errors.Is(err, domainerror.ErrResolutionInvalid) {
		return status.Error(codes.InvalidArgument, "invalid quarantine resolution")
	}

	// Provenance errors
	if errors.Is(err, domainerror.ErrSignatureRequired) {
		return status.Error(codes.InvalidArgument, "signature is required")
	}
	if errors.Is(err, domainerror.ErrSignatureInvalid) {
		return status.Error(codes.Unauthenticated, "signature is invalid")
	}
	if errors.Is(err, domainerror.ErrSignatureVerifyFailed) {
		return status.Error(codes.Unauthenticated, "signature verification failed")
	}
	if errors.Is(err, domainerror.ErrProvenanceChainBroken) {
		return status.Error(codes.FailedPrecondition, "provenance chain is broken")
	}
	if errors.Is(err, domainerror.ErrCollectorSignatureInvalid) {
		return status.Error(codes.Unauthenticated, "collector signature is invalid")
	}
	if errors.Is(err, domainerror.ErrSourceSignatureInvalid) {
		return status.Error(codes.Unauthenticated, "source signature is invalid")
	}

	// Source reliability errors
	if errors.Is(err, domainerror.ErrSourceReliabilityNotFound) {
		return status.Error(codes.NotFound, "source reliability not found")
	}

	// Processing errors
	if errors.Is(err, domainerror.ErrProcessingFailed) {
		return status.Error(codes.Internal, "processing failed")
	}
	if errors.Is(err, domainerror.ErrRoutingFailed) {
		return status.Error(codes.Internal, "routing failed")
	}
	if errors.Is(err, domainerror.ErrPublishFailed) {
		return status.Error(codes.Internal, "event publish failed")
	}

	// Default: internal error
	return status.Error(codes.Internal, "internal error")
}
