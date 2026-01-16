package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/0xsj/overwatch-pkg/log"
	"github.com/0xsj/overwatch-pkg/provenance"
	"github.com/0xsj/overwatch-pkg/types"

	"github.com/0xsj/overwatch-ingest/internal/port/inbound/command"
)

type RawDataConsumer struct {
	conn                  *nats.Conn
	logger                log.Logger
	processRawDataHandler command.ProcessRawDataHandler
	subjectPattern        string
	queueGroup            string
	subscription          *nats.Subscription
	verifySignatures      bool
}

type RawDataConsumerConfig struct {
	SubjectPattern   string
	QueueGroup       string
	VerifySignatures bool
}

func DefaultRawDataConsumerConfig() RawDataConsumerConfig {
	return RawDataConsumerConfig{
		SubjectPattern:   "overwatch.ingest.raw.*",
		QueueGroup:       "ingest-service",
		VerifySignatures: true,
	}
}

func NewRawDataConsumer(
	conn *nats.Conn,
	logger log.Logger,
	processRawDataHandler command.ProcessRawDataHandler,
	config RawDataConsumerConfig,
) *RawDataConsumer {
	if config.SubjectPattern == "" {
		config.SubjectPattern = "overwatch.ingest.raw.*"
	}
	if config.QueueGroup == "" {
		config.QueueGroup = "ingest-service"
	}

	return &RawDataConsumer{
		conn:                  conn,
		logger:                logger,
		processRawDataHandler: processRawDataHandler,
		subjectPattern:        config.SubjectPattern,
		queueGroup:            config.QueueGroup,
		verifySignatures:      config.VerifySignatures,
	}
}

func (c *RawDataConsumer) Start(ctx context.Context) error {
	sub, err := c.conn.QueueSubscribe(c.subjectPattern, c.queueGroup, func(msg *nats.Msg) {
		c.handleMessage(ctx, msg)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", c.subjectPattern, err)
	}

	c.subscription = sub

	c.logger.Info("raw data consumer started",
		log.String("subject", c.subjectPattern),
		log.String("queue_group", c.queueGroup),
	)

	return nil
}

func (c *RawDataConsumer) Stop() error {
	if c.subscription == nil {
		return nil
	}

	if err := c.subscription.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	c.logger.Info("raw data consumer stopped")
	return nil
}

func (c *RawDataConsumer) handleMessage(ctx context.Context, msg *nats.Msg) {
	envelope, err := provenance.UnmarshalEnvelope(msg.Data)
	if err != nil {
		c.logger.Error("failed to unmarshal envelope",
			log.String("subject", msg.Subject),
			log.String("error", err.Error()),
		)
		return
	}

	if c.verifySignatures {
		if err := envelope.Verify(ctx); err != nil {
			c.logger.Error("envelope signature verification failed",
				log.String("subject", msg.Subject),
				log.String("signer_did", envelope.SignerDID),
				log.String("error", err.Error()),
			)
			return
		}
	}

	var payload rawDataPayload
	if err := envelope.UnmarshalPayload(&payload); err != nil {
		c.logger.Error("failed to unmarshal raw data payload",
			log.String("subject", msg.Subject),
			log.String("error", err.Error()),
		)
		return
	}

	cmd, err := c.buildCommand(payload, envelope)
	if err != nil {
		c.logger.Error("failed to build process command",
			log.String("raw_data_id", payload.ID),
			log.String("error", err.Error()),
		)
		return
	}

	result, err := c.processRawDataHandler.Handle(ctx, cmd)
	if err != nil {
		c.logger.Error("failed to process raw data",
			log.String("raw_data_id", payload.ID),
			log.String("source_id", payload.SourceID),
			log.String("error", err.Error()),
		)
		return
	}

	c.logger.Info("raw data processed",
		log.String("raw_data_id", payload.ID),
		log.String("ingest_record_id", result.IngestRecordID.String()),
		log.String("status", result.Status),
		log.Float64("confidence", result.ConfidenceScore),
	)
}

func (c *RawDataConsumer) buildCommand(payload rawDataPayload, envelope *provenance.SignedEnvelope) (command.ProcessRawData, error) {
	sourceID := types.ID(payload.SourceID)

	var tenantID types.Optional[types.ID]
	if payload.TenantID != nil {
		tenantID = types.Some(types.ID(*payload.TenantID))
	} else {
		tenantID = types.None[types.ID]()
	}

	var sourceTimestamp types.Optional[types.Timestamp]
	if payload.SourceTimestamp != nil {
		sourceTimestamp = types.Some(types.FromTime(time.Unix(*payload.SourceTimestamp, 0)))
	} else {
		sourceTimestamp = types.None[types.Timestamp]()
	}

	collectedAt := types.FromTime(time.Unix(payload.CollectedAt, 0))

	var jobID types.Optional[types.ID]
	if payload.JobID != nil {
		jobID = types.Some(types.ID(*payload.JobID))
	} else {
		jobID = types.None[types.ID]()
	}

	var batchID types.Optional[types.ID]
	if payload.BatchID != nil {
		batchID = types.Some(types.ID(*payload.BatchID))
	} else {
		batchID = types.None[types.ID]()
	}

	var batchIndex types.Optional[int32]
	if payload.BatchIndex != nil {
		batchIndex = types.Some(*payload.BatchIndex)
	} else {
		batchIndex = types.None[int32]()
	}

	collectorSigner := envelope.ToSignatureInfo()

	var sourceSigner *provenance.SignatureInfo
	if payload.SourceDID != nil && payload.SourceSignature != nil {
		sourceSigner = &provenance.SignatureInfo{
			DID:       *payload.SourceDID,
			Signature: *payload.SourceSignature,
		}
	}

	return command.ProcessRawData{
		TenantID:        tenantID,
		SourceID:        sourceID,
		SourceType:      payload.SourceType,
		RawDataID:       payload.ID,
		Payload:         payload.Payload,
		Metadata:        payload.Metadata,
		SourceTimestamp: sourceTimestamp,
		CollectedAt:     collectedAt,
		SourceSigner:    sourceSigner,
		CollectorSigner: collectorSigner,
		JobID:           jobID,
		BatchID:         batchID,
		BatchIndex:      batchIndex,
	}, nil
}

type rawDataPayload struct {
	ID              string            `json:"id"`
	TenantID        *string           `json:"tenant_id,omitempty"`
	SourceID        string            `json:"source_id"`
	SourceType      string            `json:"source_type"`
	Payload         map[string]any    `json:"payload"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	SourceTimestamp *int64            `json:"source_timestamp,omitempty"`
	CollectedAt     int64             `json:"collected_at"`
	SourceDID       *string           `json:"source_did,omitempty"`
	SourceSignature *string           `json:"source_signature,omitempty"`
	JobID           *string           `json:"job_id,omitempty"`
	BatchID         *string           `json:"batch_id,omitempty"`
	BatchIndex      *int32            `json:"batch_index,omitempty"`
}

var _ json.Unmarshaler = (*rawDataPayload)(nil)

func (p *rawDataPayload) UnmarshalJSON(data []byte) error {
	type alias rawDataPayload
	return json.Unmarshal(data, (*alias)(p))
}
