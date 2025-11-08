// platform/pkg/grpc/client/grpc_client.go
package client

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/0xsj/scout/platform/pkg/grpc/interceptors"
)

// grpcClient implements the Client interface using google.golang.org/grpc.
type grpcClient struct {
	config *Config
	conn   *grpc.ClientConn
	mu     sync.RWMutex
}

// New creates a new gRPC client with the given configuration.
func New(config *Config) (Client, error) {
	if config.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &grpcClient{
		config: config,
	}, nil
}

// Connect establishes a connection to the gRPC server.
func (c *grpcClient) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return fmt.Errorf("client already connected")
	}

	c.config.Logger.Info("grpc client connecting", "target", c.config.Target)

	// Dial options
	opts := []grpc.DialOption{
		// Keepalive parameters
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                c.config.KeepAliveTime,
			Timeout:             c.config.KeepAliveTimeout,
			PermitWithoutStream: true,
		}),

		// Unary interceptors
		grpc.WithChainUnaryInterceptor(
			interceptors.UnaryClientLoggingInterceptor(c.config.Logger),
		),

		// Stream interceptors
		grpc.WithChainStreamInterceptor(
			interceptors.StreamClientLoggingInterceptor(c.config.Logger),
		),
	}

	// TLS or insecure
	if c.config.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Create connection with timeout
	dialCtx, cancel := context.WithTimeout(ctx, c.config.ConnectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, c.config.Target, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", c.config.Target, err)
	}

	c.conn = conn

	c.config.Logger.Info("grpc client connected", "target", c.config.Target)

	return nil
}

// Close closes the connection.
func (c *grpcClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return nil
	}

	c.config.Logger.Info("grpc client closing", "target", c.config.Target)

	err := c.conn.Close()
	c.conn = nil

	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	c.config.Logger.Info("grpc client closed")

	return nil
}

// GetConnection returns the underlying grpc.ClientConn.
func (c *grpcClient) GetConnection() *grpc.ClientConn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

// IsConnected returns true if the client is connected.
func (c *grpcClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return false
	}

	state := c.conn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}