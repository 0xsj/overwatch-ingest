// platform/pkg/grpc/server/config.go
package server

import (
	"time"

	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// Config holds gRPC server configuration.
type Config struct {
	// Address to listen on (e.g., ":50051", "0.0.0.0:50051")
	Address string

	// Logger for server operations
	Logger logger.Logger

	// MaxConnectionIdle is the maximum time a connection can be idle
	MaxConnectionIdle time.Duration

	// MaxConnectionAge is the maximum time a connection can exist
	MaxConnectionAge time.Duration

	// MaxConnectionAgeGrace is additional time for graceful connection closure
	MaxConnectionAgeGrace time.Duration

	// KeepAliveTime is the time after which a keepalive ping is sent
	KeepAliveTime time.Duration

	// KeepAliveTimeout is the time to wait for keepalive ping ack
	KeepAliveTimeout time.Duration
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Address:               ":50051",
		Logger:                nil, // Must be provided
		MaxConnectionIdle:     15 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Second,
		KeepAliveTime:         5 * time.Minute,
		KeepAliveTimeout:      20 * time.Second,
	}
}

// WithAddress sets the server address.
func (c *Config) WithAddress(address string) *Config {
	c.Address = address
	return c
}

// WithLogger sets the logger.
func (c *Config) WithLogger(logger logger.Logger) *Config {
	c.Logger = logger
	return c
}