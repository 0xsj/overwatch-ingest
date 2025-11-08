// platform/pkg/grpc/server/grpc_server.go
package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/0xsj/scout/platform/pkg/grpc/interceptors"
)

// grpcServer implements the Server interface using google.golang.org/grpc.
type grpcServer struct {
	config *Config
	server *grpc.Server
}

// New creates a new gRPC server with the given configuration.
func New(config *Config) (Server, error) {
	if config.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	// Server options
	opts := []grpc.ServerOption{
		// Connection parameters
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     config.MaxConnectionIdle,
			MaxConnectionAge:      config.MaxConnectionAge,
			MaxConnectionAgeGrace: config.MaxConnectionAgeGrace,
			Time:                  config.KeepAliveTime,
			Timeout:               config.KeepAliveTimeout,
		}),

		// Unary interceptors (middleware)
		grpc.ChainUnaryInterceptor(
			interceptors.UnaryLoggingInterceptor(config.Logger),
			interceptors.UnaryRecoveryInterceptor(config.Logger),
		),

		// Stream interceptors
		grpc.ChainStreamInterceptor(
			interceptors.StreamLoggingInterceptor(config.Logger),
			interceptors.StreamRecoveryInterceptor(config.Logger),
		),
	}

	server := grpc.NewServer(opts...)

	return &grpcServer{
		config: config,
		server: server,
	}, nil
}

// Start starts the gRPC server.
func (s *grpcServer) Start() error {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.config.Address, err)
	}

	s.config.Logger.Info("grpc server starting", "address", s.config.Address)

	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop gracefully stops the server.
func (s *grpcServer) Stop(ctx context.Context) error {
	s.config.Logger.Info("grpc server stopping")

	// Channel to signal when graceful stop completes
	done := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	// Wait for graceful stop or context timeout
	select {
	case <-done:
		s.config.Logger.Info("grpc server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.config.Logger.Warn("grpc server stop timeout, forcing shutdown")
		s.server.Stop()
		return ctx.Err()
	}
}

// RegisterService registers a gRPC service.
func (s *grpcServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
	s.config.Logger.Info("grpc service registered", "service", desc.ServiceName)
}

// GetServer returns the underlying grpc.Server.
func (s *grpcServer) GetServer() *grpc.Server {
	return s.server
}