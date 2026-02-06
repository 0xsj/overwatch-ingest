package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"

	appcommand "github.com/0xsj/overwatch-ingest/internal/app/command"
	domainerror "github.com/0xsj/overwatch-ingest/internal/domain/error"
	"github.com/0xsj/overwatch-ingest/internal/domain/model"
	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
	"github.com/0xsj/overwatch-ingest/internal/port/outbound/validation"
	"github.com/0xsj/overwatch-ingest/tests/testutil"
	"github.com/0xsj/overwatch-ingest/tests/testutil/mocks"
)

type processRawDataDeps struct {
	recordRepo       *mocks.IngestRecordRepository
	quarantineRepo   *mocks.QuarantinedRecordRepository
	reliabilityRepo  *mocks.SourceReliabilityRepository
	publisher        *mocks.EventPublisher
	schemaValidator  *mocks.SchemaValidator
	anomalyDetector  *mocks.AnomalyDetector
	confidenceScorer *mocks.ConfidenceScorer
	verifier         *mocks.Verifier
}

func newProcessRawDataDeps() *processRawDataDeps {
	return &processRawDataDeps{
		recordRepo:       mocks.NewIngestRecordRepository(),
		quarantineRepo:   mocks.NewQuarantinedRecordRepository(),
		reliabilityRepo:  mocks.NewSourceReliabilityRepository(),
		publisher:        mocks.NewEventPublisher(),
		schemaValidator:  mocks.NewSchemaValidator(),
		anomalyDetector:  mocks.NewAnomalyDetector(),
		confidenceScorer: mocks.NewConfidenceScorer(),
		verifier:         mocks.NewVerifier(),
	}
}

func newProcessRawDataHandler(d *processRawDataDeps) command.ProcessRawDataHandler {
	return appcommand.NewProcessRawDataHandler(
		d.recordRepo,
		d.quarantineRepo,
		d.reliabilityRepo,
		d.publisher,
		d.schemaValidator,
		d.anomalyDetector,
		d.confidenceScorer,
		d.verifier,
		mocks.NilEnvelopeBuilder(),
		appcommand.DefaultProcessRawDataHandlerConfig(),
	)
}

func validProcessRawDataCmd() command.ProcessRawData {
	return command.ProcessRawData{
		TenantID:        types.None[types.ID](),
		SourceID:        types.NewID(),
		SourceType:      "ais",
		RawDataID:       testutil.Fake.RawDataID(),
		Payload:         testutil.Fake.VesselPayload(),
		Metadata:        testutil.Fake.Metadata(),
		SourceTimestamp: types.Some(types.Now()),
		CollectedAt:     types.Now(),
		SourceSigner:    nil,
		CollectorSigner: testutil.Fixtures.Signer(),
		JobID:           types.None[types.ID](),
		BatchID:         types.None[types.ID](),
		BatchIndex:      types.None[int32](),
	}
}

// =============================================================================
// Accept path
// =============================================================================

func TestProcessRawData_Accept_Success(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "accepted" {
		t.Errorf("Status = %q, want %q", result.Status, "accepted")
	}
	if result.IngestRecordID.IsEmpty() {
		t.Error("IngestRecordID should not be empty")
	}
	if result.EntityType.IsEmpty() {
		t.Error("EntityType should be present for accepted records")
	}
	if result.EntityID.IsEmpty() {
		t.Error("EntityID should be present for accepted records")
	}
	if result.QuarantineID.IsPresent() {
		t.Error("QuarantineID should not be present for accepted records")
	}
	if result.RejectionReason.IsPresent() {
		t.Error("RejectionReason should not be present for accepted records")
	}
	if d.recordRepo.Calls.Create != 1 {
		t.Errorf("RecordRepo.Create calls = %d, want 1", d.recordRepo.Calls.Create)
	}
	if d.recordRepo.RecordCount() != 1 {
		t.Errorf("RecordRepo count = %d, want 1", d.recordRepo.RecordCount())
	}
}

func TestProcessRawData_Accept_VesselEntity(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	cmd.Payload = testutil.Fake.VesselPayload()

	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "accepted" {
		t.Errorf("Status = %q, want %q", result.Status, "accepted")
	}
	if result.EntityType.OrElse("") != "vessel" {
		t.Errorf("EntityType = %q, want %q", result.EntityType.OrElse(""), "vessel")
	}
}

func TestProcessRawData_Accept_AircraftEntity(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	cmd.Payload = testutil.Fake.AircraftPayload()

	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "accepted" {
		t.Errorf("Status = %q, want %q", result.Status, "accepted")
	}
	if result.EntityType.OrElse("") != "aircraft" {
		t.Errorf("EntityType = %q, want %q", result.EntityType.OrElse(""), "aircraft")
	}
}

func TestProcessRawData_Accept_WithTenant(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	tenantID := types.NewID()
	cmd := validProcessRawDataCmd()
	cmd.TenantID = types.Some(tenantID)

	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "accepted" {
		t.Errorf("Status = %q, want %q", result.Status, "accepted")
	}
}

func TestProcessRawData_Accept_ReliabilityUpdated(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	_, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if d.reliabilityRepo.Calls.Upsert != 1 {
		t.Errorf("ReliabilityRepo.Upsert calls = %d, want 1", d.reliabilityRepo.Calls.Upsert)
	}
}

func TestProcessRawData_Accept_EventsPublished(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	_, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}

	// Expect: RecordReceived, SignatureVerified (collector), RecordValidated, RecordAccepted, SourceReliabilityUpdated
	if d.publisher.Calls.PublishRecordReceived != 1 {
		t.Errorf("PublishRecordReceived calls = %d, want 1", d.publisher.Calls.PublishRecordReceived)
	}
	if d.publisher.Calls.PublishRecordValidated != 1 {
		t.Errorf("PublishRecordValidated calls = %d, want 1", d.publisher.Calls.PublishRecordValidated)
	}
	if d.publisher.Calls.PublishRecordAccepted != 1 {
		t.Errorf("PublishRecordAccepted calls = %d, want 1", d.publisher.Calls.PublishRecordAccepted)
	}
}

// =============================================================================
// Reject path
// =============================================================================

func TestProcessRawData_Reject_ValidationFailed(t *testing.T) {
	d := newProcessRawDataDeps()
	d.schemaValidator.SetInvalid(
		[]string{"mmsi"},
		[]model.Anomaly{testutil.Fixtures.CriticalAnomaly()},
	)
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "rejected" {
		t.Errorf("Status = %q, want %q", result.Status, "rejected")
	}
	if result.RejectionReason.IsEmpty() {
		t.Error("RejectionReason should be present for rejected records")
	}
	if result.EntityType.IsPresent() {
		t.Error("EntityType should not be present for rejected records")
	}
	if d.publisher.Calls.PublishRecordRejected != 1 {
		t.Errorf("PublishRecordRejected calls = %d, want 1", d.publisher.Calls.PublishRecordRejected)
	}
}

func TestProcessRawData_Reject_AnomalyDetectorRejects(t *testing.T) {
	d := newProcessRawDataDeps()
	d.anomalyDetector.SetAnomalies([]model.Anomaly{testutil.Fixtures.CriticalAnomaly()})
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "rejected" {
		t.Errorf("Status = %q, want %q", result.Status, "rejected")
	}
}

func TestProcessRawData_Reject_LowConfidence(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetLowConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "rejected" {
		t.Errorf("Status = %q, want %q", result.Status, "rejected")
	}
	if result.RejectionReason.OrElse("") != "confidence too low" {
		t.Errorf("RejectionReason = %q, want %q", result.RejectionReason.OrElse(""), "confidence too low")
	}
}

func TestProcessRawData_Reject_RecordCreated(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetLowConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	_, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if d.recordRepo.Calls.Create != 1 {
		t.Errorf("RecordRepo.Create calls = %d, want 1", d.recordRepo.Calls.Create)
	}
	if d.recordRepo.RecordCount() != 1 {
		t.Errorf("RecordRepo count = %d, want 1", d.recordRepo.RecordCount())
	}
}

// =============================================================================
// Quarantine path
// =============================================================================

func TestProcessRawData_Quarantine_MediumConfidence(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetMediumConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "quarantined" {
		t.Errorf("Status = %q, want %q", result.Status, "quarantined")
	}
	if result.QuarantineID.IsEmpty() {
		t.Error("QuarantineID should be present for quarantined records")
	}
	if result.EntityType.IsPresent() {
		t.Error("EntityType should not be present for quarantined records")
	}
	if result.RejectionReason.IsPresent() {
		t.Error("RejectionReason should not be present for quarantined records")
	}
}

func TestProcessRawData_Quarantine_AnomalyDetectorQuarantines(t *testing.T) {
	d := newProcessRawDataDeps()
	d.anomalyDetector.SetAnomalies([]model.Anomaly{testutil.Fixtures.ErrorAnomaly()})
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "quarantined" {
		t.Errorf("Status = %q, want %q", result.Status, "quarantined")
	}
}

func TestProcessRawData_Quarantine_ValidationQuarantines(t *testing.T) {
	d := newProcessRawDataDeps()
	// Schema valid, but with error-level anomalies that trigger quarantine
	d.schemaValidator.Result = validation.SchemaValidationResult{
		Valid:         true,
		FieldsPresent: []string{"mmsi", "lat", "lon"},
		FieldsMissing: []string{"speed"},
		Anomalies:     []model.Anomaly{testutil.Fixtures.ErrorAnomaly()},
		SchemaID:      "test-schema",
		SchemaVersion: "v1",
	}
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.Status != "quarantined" {
		t.Errorf("Status = %q, want %q", result.Status, "quarantined")
	}
}

func TestProcessRawData_Quarantine_BothRecordAndQuarantineCreated(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetMediumConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	_, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if d.recordRepo.Calls.Create != 1 {
		t.Errorf("RecordRepo.Create calls = %d, want 1", d.recordRepo.Calls.Create)
	}
	if d.quarantineRepo.Calls.Create != 1 {
		t.Errorf("QuarantineRepo.Create calls = %d, want 1", d.quarantineRepo.Calls.Create)
	}
}

func TestProcessRawData_Quarantine_EventPublished(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetMediumConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	_, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if d.publisher.Calls.PublishRecordQuarantined != 1 {
		t.Errorf("PublishRecordQuarantined calls = %d, want 1", d.publisher.Calls.PublishRecordQuarantined)
	}
}

// =============================================================================
// Validation errors
// =============================================================================

func TestProcessRawData_MissingSourceID(t *testing.T) {
	d := newProcessRawDataDeps()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	cmd.SourceID = ""

	_, err := handler.Handle(ctx, cmd)
	if err == nil {
		t.Fatal("Handle() expected error for missing source ID")
	}
	if !errors.Is(err, domainerror.ErrSourceIDRequired) {
		t.Errorf("error = %v, want ErrSourceIDRequired", err)
	}
}

func TestProcessRawData_MissingSourceType(t *testing.T) {
	d := newProcessRawDataDeps()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	cmd.SourceType = ""

	_, err := handler.Handle(ctx, cmd)
	if err == nil {
		t.Fatal("Handle() expected error for missing source type")
	}
	if !errors.Is(err, domainerror.ErrSourceTypeRequired) {
		t.Errorf("error = %v, want ErrSourceTypeRequired", err)
	}
}

func TestProcessRawData_MissingRawDataID(t *testing.T) {
	d := newProcessRawDataDeps()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	cmd.RawDataID = ""

	_, err := handler.Handle(ctx, cmd)
	if err == nil {
		t.Fatal("Handle() expected error for missing raw data ID")
	}
	if !errors.Is(err, domainerror.ErrRawDataIDRequired) {
		t.Errorf("error = %v, want ErrRawDataIDRequired", err)
	}
}

func TestProcessRawData_MissingPayload(t *testing.T) {
	d := newProcessRawDataDeps()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	cmd.Payload = nil

	_, err := handler.Handle(ctx, cmd)
	if err == nil {
		t.Fatal("Handle() expected error for missing payload")
	}
	if !errors.Is(err, domainerror.ErrPayloadRequired) {
		t.Errorf("error = %v, want ErrPayloadRequired", err)
	}
}

func TestProcessRawData_MissingCollectorSigner(t *testing.T) {
	d := newProcessRawDataDeps()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	cmd.CollectorSigner = nil

	_, err := handler.Handle(ctx, cmd)
	if err == nil {
		t.Fatal("Handle() expected error for missing collector signer")
	}
	if !errors.Is(err, domainerror.ErrSignatureRequired) {
		t.Errorf("error = %v, want ErrSignatureRequired", err)
	}
}

// =============================================================================
// Duplicate detection
// =============================================================================

func TestProcessRawData_DuplicateRawDataID(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()

	// First call should succeed
	_, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("First Handle() unexpected error: %v", err)
	}

	// Second call with same rawDataID should fail
	cmd2 := validProcessRawDataCmd()
	cmd2.RawDataID = cmd.RawDataID
	_, err = handler.Handle(ctx, cmd2)
	if err == nil {
		t.Fatal("Second Handle() expected error for duplicate raw data ID")
	}
}

// =============================================================================
// Signature verification
// =============================================================================

func TestProcessRawData_CollectorSignatureVerified(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	_, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if d.verifier.Calls.VerifySignatureInfo < 1 {
		t.Errorf("Verifier.VerifySignatureInfo calls = %d, want >= 1", d.verifier.Calls.VerifySignatureInfo)
	}
	if d.publisher.Calls.PublishSignatureVerified < 1 {
		t.Errorf("PublishSignatureVerified calls = %d, want >= 1", d.publisher.Calls.PublishSignatureVerified)
	}
}

func TestProcessRawData_CollectorSignatureFailedContinues(t *testing.T) {
	d := newProcessRawDataDeps()
	d.verifier.SetVerifyError(errors.New("invalid signature"))
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	// Handler should still process (collector sig failure publishes event but continues)
	if result.IngestRecordID.IsEmpty() {
		t.Error("IngestRecordID should not be empty even with failed signature")
	}
	if d.publisher.Calls.PublishSignatureFailed < 1 {
		t.Errorf("PublishSignatureFailed calls = %d, want >= 1", d.publisher.Calls.PublishSignatureFailed)
	}
}

func TestProcessRawData_WithSourceSigner(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	cmd.SourceSigner = &provenance.SignatureInfo{
		DID:        "did:key:z6MkTestSource",
		SignerType: provenance.SignerTypeSource,
		Signature:  testutil.Fake.Hex(64),
		SignedAt:   time.Now(),
	}

	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.IngestRecordID.IsEmpty() {
		t.Error("IngestRecordID should not be empty")
	}
	// Both collector and source signatures should be verified
	if d.verifier.Calls.VerifySignatureInfo < 2 {
		t.Errorf("Verifier.VerifySignatureInfo calls = %d, want >= 2", d.verifier.Calls.VerifySignatureInfo)
	}
}

// =============================================================================
// Repository error handling
// =============================================================================

func TestProcessRawData_ExistsByRawDataID_Error(t *testing.T) {
	d := newProcessRawDataDeps()
	d.recordRepo.Errors.ExistsByRawDataID = errors.New("db error")
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	_, err := handler.Handle(ctx, cmd)
	if err == nil {
		t.Fatal("Handle() expected error when ExistsByRawDataID fails")
	}
}

// =============================================================================
// Confidence score in result
// =============================================================================

func TestProcessRawData_ConfidenceScoreInResult(t *testing.T) {
	d := newProcessRawDataDeps()
	d.confidenceScorer.SetHighConfidence()
	handler := newProcessRawDataHandler(d)
	ctx := context.Background()

	cmd := validProcessRawDataCmd()
	result, err := handler.Handle(ctx, cmd)
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if result.ConfidenceScore <= 0 {
		t.Errorf("ConfidenceScore = %f, want > 0", result.ConfidenceScore)
	}
}
