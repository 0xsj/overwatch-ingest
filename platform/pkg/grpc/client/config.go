// platform/pkg/grpc/client/config.go
package client

import (
	"time"

	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// Config holds gRPC client configuration.
type Config struct {
	// Target address (e.g., "localhost:50051", "agents:50051")
	Target string

	// Logger for client operations
	Logger logger.Logger

	// Timeout for connection establishment
	ConnectTimeout time.Duration

	// MaxRetries for connection attempts
	MaxRetries int

	// Insecure disables TLS (for development only)
	Insecure bool

	// KeepAliveTime is the time after which a keepalive ping is sent
	KeepAliveTime time.Duration

	// KeepAliveTimeout is the time to wait for keepalive ping ack
	KeepAliveTimeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Target:           "localhost:50051",
		Logger:           nil, // Must be provided
		ConnectTimeout:   10 * time.Second,
		MaxRetries:       3,
		Insecure:         true, // Default to insecure for local dev
		KeepAliveTime:    30 * time.Second,
		KeepAliveTimeout: 10 * time.Second,
	}
}

// WithTarget sets the target address.
func (c *Config) WithTarget(target string) *Config {
	c.Target = target
	return c
}

// WithLogger sets the logger.
func (c *Config) WithLogger(logger logger.Logger) *Config {
	c.Logger = logger
	return c
}

// WithInsecure sets whether to use insecure connections.
func (c *Config) WithInsecure(insecure bool) *Config {
	c.Insecure = insecure
	return c
}