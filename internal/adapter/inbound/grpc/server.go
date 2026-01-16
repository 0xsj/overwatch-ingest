package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	ingestv1 "github.com/0xsj/overwatch-contracts/gen/go/ingest/v1"
	pkggrpc "github.com/0xsj/overwatch-pkg/grpc"
	"github.com/0xsj/overwatch-pkg/log"
)

type ServerConfig struct {
	Host              string
	Port              int
	EnableReflection  bool
	EnableHealthCheck bool
}

func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Server struct {
	server  *pkggrpc.Server
	handler *Handler
	logger  log.Logger
}

func NewServer(
	cfg ServerConfig,
	handler *Handler,
	logger log.Logger,
	interceptors ...grpc.UnaryServerInterceptor,
) (*Server, error) {
	opts := []pkggrpc.ServerOption{
		pkggrpc.WithServerAddress(cfg.Address()),
		pkggrpc.WithServerLogger(logger),
		pkggrpc.WithServerReflection(cfg.EnableReflection),
		pkggrpc.WithServerHealthCheck(cfg.EnableHealthCheck),
	}

	if len(interceptors) > 0 {
		opts = append(opts, pkggrpc.WithUnaryInterceptors(interceptors...))
	}

	server, err := pkggrpc.NewServer(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %w", err)
	}

	return &Server{
		server:  server,
		handler: handler,
		logger:  logger,
	}, nil
}

func (s *Server) RegisterServices() {
	s.server.RegisterService(
		&ingestv1.IngestService_ServiceDesc,
		s.handler,
	)
}

func (s *Server) Start(ctx context.Context) error {
	s.RegisterServices()
	return s.server.Start(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Stop(ctx)
}

func (s *Server) Run() error {
	s.RegisterServices()
	return s.server.Run()
}

func (s *Server) GRPCServer() *grpc.Server {
	return s.server.Server()
}

func (s *Server) Address() string {
	return s.server.Address()
}

func (s *Server) IsRunning() bool {
	return s.server.IsRunning()
}
