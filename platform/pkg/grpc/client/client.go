// platform/pkg/grpc/client/client.go
package client

import (
	"context"

	"google.golang.org/grpc"
)

// Client defines the gRPC client interface.
// Implementations manage connection lifecycle and provide access to the underlying connection.
type Client interface {
	// Connect establishes a connection to the gRPC server.
	Connect(ctx context.Context) error

	// Close closes the connection gracefully.
	Close() error

	// GetConnection returns the underlying grpc.ClientConn.
	GetConnection() *grpc.ClientConn

	// IsConnected returns true if the client is connected.
	IsConnected() bool
}