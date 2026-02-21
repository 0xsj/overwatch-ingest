package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	natsclient "github.com/nats-io/nats.go"

	"github.com/0xsj/overwatch-pkg/log"
	"github.com/0xsj/overwatch-pkg/provenance"

	ingestgrpc "github.com/0xsj/overwatch-ingest/internal/adapter/inbound/grpc"
	natsadapter "github.com/0xsj/overwatch-ingest/internal/adapter/inbound/nats"
	natspublisher "github.com/0xsj/overwatch-ingest/internal/adapter/outbound/nats"
	"github.com/0xsj/overwatch-ingest/internal/adapter/outbound/postgres"
	"github.com/0xsj/overwatch-ingest/internal/adapter/outbound/validation"
	"github.com/0xsj/overwatch-ingest/internal/app/command"
	"github.com/0xsj/overwatch-ingest/internal/app/query"
	"github.com/0xsj/overwatch-ingest/internal/config"
	portvalidation "github.com/0xsj/overwatch-ingest/internal/port/outbound/validation"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := log.NewPretty(log.DefaultConfig())

	logger.Info("starting ingest service",
		log.String("version", "1.0.0"),
		log.String("address", cfg.Server.Address()),
	)

	identity, signer, verifier, err := initializeProvenance(cfg.ServiceIdentity, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize provenance: %w", err)
	}

	logger.Info("service identity initialized",
		log.String("service_id", cfg.ServiceIdentity.ID),
		log.String("service_name", cfg.ServiceIdentity.Name),
		log.String("did", identity.DID()),
	)

	pool, err := connectPostgres(ctx, cfg.Database, logger)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer pool.Close()

	natsConn, err := connectNATS(cfg.NATS, logger)
	if err != nil {
		return fmt.Errorf("failed to connect to nats: %w", err)
	}
	defer natsConn.Close()

	recordRepo := postgres.NewIngestRecordRepository(pool)
	quarantineRepo := postgres.NewQuarantinedRecordRepository(pool)
	reliabilityRepo := postgres.NewSourceReliabilityRepository(pool)

	eventPublisher, err := natspublisher.NewSignedEventPublisher(
		natsConn,
		cfg.NATS.SubjectPrefix,
		identity,
	)
	if err != nil {
		return fmt.Errorf("failed to create event publisher: %w", err)
	}

	registry := validation.NewSourceTypeRegistry()
	schemaValidator := validation.NewSchemaValidator(registry)
	anomalyDetector := validation.NewAnomalyDetector(registry)
	confidenceScorer := validation.NewConfidenceScorer(registry)

	processRawDataHandler := command.NewProcessRawDataHandler(
		recordRepo,
		quarantineRepo,
		reliabilityRepo,
		eventPublisher,
		schemaValidator,
		anomalyDetector,
		confidenceScorer,
		verifier,
		signer,
		command.ProcessRawDataHandlerConfig{
			Thresholds: portvalidation.ScoringThresholds{
				AcceptThreshold: cfg.Ingest.AcceptThreshold,
				RejectThreshold: cfg.Ingest.RejectThreshold,
			},
			QuarantineExpiry: cfg.Ingest.QuarantineExpiry,
		},
	)

	resolveQuarantinedHandler := command.NewResolveQuarantinedHandler(
		quarantineRepo,
		recordRepo,
		reliabilityRepo,
		eventPublisher,
		signer,
	)

	bulkResolveQuarantinedHandler := command.NewBulkResolveQuarantinedHandler(
		resolveQuarantinedHandler,
	)

	reprocessRecordHandler := command.NewReprocessRecordHandler(
		recordRepo,
		processRawDataHandler,
	)

	reprocessBySourceHandler := command.NewReprocessBySourceHandler(
		recordRepo,
	)

	getRecordHandler := query.NewGetRecordHandler(recordRepo)
	getRecordByRawDataHandler := query.NewGetRecordByRawDataHandler(recordRepo)
	listRecordsHandler := query.NewListRecordsHandler(recordRepo)
	getQuarantinedHandler := query.NewGetQuarantinedHandler(quarantineRepo)
	getQuarantinedByIngestRecordHandler := query.NewGetQuarantinedByIngestRecordHandler(quarantineRepo)
	listQuarantinedHandler := query.NewListQuarantinedHandler(quarantineRepo)
	getSourceReliabilityHandler := query.NewGetSourceReliabilityHandler(reliabilityRepo)
	listSourceReliabilityHandler := query.NewListSourceReliabilityHandler(reliabilityRepo)
	getIngestStatsHandler := query.NewGetIngestStatsHandler(recordRepo, quarantineRepo)

	_ = getQuarantinedByIngestRecordHandler

	rawDataConsumer := natsadapter.NewRawDataConsumer(
		natsConn,
		logger,
		processRawDataHandler,
		natsadapter.RawDataConsumerConfig{
			SubjectPattern:   "overwatch.ingest.raw.*",
			QueueGroup:       "ingest-service",
			VerifySignatures: cfg.Ingest.RequireCollectorSignature,
		},
	)

	if err := rawDataConsumer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start raw data consumer: %w", err)
	}
	defer rawDataConsumer.Stop()

	logger.Info("raw data consumer started",
		log.String("subject_pattern", "overwatch.ingest.raw.*"),
		log.String("queue_group", "ingest-service"),
	)

	handler := ingestgrpc.NewHandler(ingestgrpc.HandlerConfig{
		ResolveQuarantinedHandler:     resolveQuarantinedHandler,
		BulkResolveQuarantinedHandler: bulkResolveQuarantinedHandler,
		ReprocessRecordHandler:        reprocessRecordHandler,
		ReprocessBySourceHandler:      reprocessBySourceHandler,

		GetRecordHandler:             getRecordHandler,
		GetRecordByRawDataHandler:    getRecordByRawDataHandler,
		ListRecordsHandler:           listRecordsHandler,
		GetQuarantinedHandler:        getQuarantinedHandler,
		ListQuarantinedHandler:       listQuarantinedHandler,
		GetSourceReliabilityHandler:  getSourceReliabilityHandler,
		ListSourceReliabilityHandler: listSourceReliabilityHandler,
		GetIngestStatsHandler:        getIngestStatsHandler,
	})

	loggingInterceptor := ingestgrpc.NewLoggingInterceptor(newGRPCLogger(logger))
	recoveryInterceptor := ingestgrpc.NewRecoveryInterceptor(newGRPCLogger(logger))

	serverCfg := ingestgrpc.ServerConfig{
		Host:              cfg.Server.Host,
		Port:              cfg.Server.Port,
		EnableReflection:  cfg.Server.EnableReflection,
		EnableHealthCheck: cfg.Server.EnableHealthCheck,
	}

	server, err := ingestgrpc.NewServer(
		serverCfg,
		handler,
		logger,
		recoveryInterceptor.Unary(),
		loggingInterceptor.Unary(),
	)
	if err != nil {
		return fmt.Errorf("failed to create grpc server: %w", err)
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("ingest service started", log.String("address", serverCfg.Address()))

	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-sigChan:
		logger.Info("received shutdown signal", log.String("signal", sig.String()))
		cancel()

		if err := rawDataConsumer.Stop(); err != nil {
			logger.Error("failed to stop raw data consumer", log.String("error", err.Error()))
		}

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer shutdownCancel()

		if err := server.Stop(shutdownCtx); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}

		logger.Info("ingest service stopped gracefully")
		return nil
	}
}

func initializeProvenance(cfg config.ServiceIdentityConfig, logger log.Logger) (provenance.Identity, *provenance.EnvelopeBuilder, provenance.Verifier, error) {
	var identity *provenance.ServiceIdentity
	var err error

	if cfg.HasPrivateKey() {
		logger.Warn("private key loading not yet implemented, generating new identity")
		identity, err = provenance.GenerateServiceIdentity(cfg.ID, cfg.Name)
	} else if cfg.GenerateIfMissing {
		identity, err = provenance.GenerateServiceIdentity(cfg.ID, cfg.Name)
		if err == nil {
			logger.Warn("generated new service identity - consider persisting the private key",
				log.String("did", identity.DID()),
			)
		}
	} else {
		return nil, nil, nil, fmt.Errorf("no private key configured and generation disabled")
	}

	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create identity: %w", err)
	}

	signer, err := provenance.NewEnvelopeBuilder(identity)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create envelope builder: %w", err)
	}

	verifier := provenance.NewVerifier()

	return identity, signer, verifier, nil
}

func connectPostgres(ctx context.Context, cfg config.DatabaseConfig, logger log.Logger) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxConns)
	poolCfg.MinConns = int32(cfg.MinConns)
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("connected to postgres",
		log.String("host", cfg.Host),
		log.String("database", cfg.Database),
	)

	return pool, nil
}

func connectNATS(cfg config.NATSConfig, logger log.Logger) (*natsclient.Conn, error) {
	opts := []natsclient.Option{
		natsclient.MaxReconnects(cfg.MaxReconnects),
		natsclient.ReconnectWait(cfg.ReconnectWait),
		natsclient.DisconnectErrHandler(func(nc *natsclient.Conn, err error) {
			if err != nil {
				logger.Warn("nats disconnected", log.String("error", err.Error()))
			}
		}),
		natsclient.ReconnectHandler(func(nc *natsclient.Conn) {
			logger.Info("nats reconnected", log.String("url", nc.ConnectedUrl()))
		}),
	}

	conn, err := natsclient.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	logger.Info("connected to nats",
		log.String("url", conn.ConnectedUrl()),
	)

	return conn, nil
}

type grpcLogger struct {
	logger log.Logger
}

func newGRPCLogger(logger log.Logger) *grpcLogger {
	return &grpcLogger{logger: logger}
}

func (l *grpcLogger) Info(msg string, fields ...interface{}) {
	l.logger.Info(msg, toLogFields(fields)...)
}

func (l *grpcLogger) Error(msg string, fields ...interface{}) {
	l.logger.Error(msg, toLogFields(fields)...)
}

func toLogFields(fields []interface{}) []log.Field {
	if len(fields) == 0 {
		return nil
	}

	result := make([]log.Field, 0, len(fields)/2)
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		result = append(result, log.Any(key, fields[i+1]))
	}
	return result
}
