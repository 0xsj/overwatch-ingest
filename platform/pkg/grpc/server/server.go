// platform/pkg/grpc/server/server.go
package server

import (
	"context"

	"google.golang.org/grpc"
)

// Server defines the gRPC server interface.
// Implementations handle server lifecycle (start, stop, graceful shutdown).
type Server interface {
	// Start starts the gRPC server on the configured address.
	// This is a blocking call that returns when the server stops.
	Start() error

	// Stop gracefully stops the server.
	// It stops accepting new connections and waits for existing RPCs to complete.
	Stop(ctx context.Context) error

	// RegisterService registers a gRPC service implementation.
	// Must be called before Start().
	RegisterService(desc *grpc.ServiceDesc, impl interface{})

	// GetServer returns the underlying grpc.Server for advanced use cases.
	GetServer() *grpc.Server
}