package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/0xsj/scout/platform/pkg/errors"
	pkgerrors "github.com/0xsj/scout/platform/pkg/errors"
	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// Config holds PostgreSQL connection pool configuration.
type Config struct {
    // DSN is the PostgreSQL connection string
    // Format: postgres://user:password@host:port/dbname?sslmode=disable
    dsn string
    
    // MaxConnections is the maximum number of connections in the pool
    maxConnections int32
    
    // MinConnections is the minimum number of idle connections to maintain
    minConnections int32
    
    // MaxConnectionLifetime is the maximum duration a connection can be reused
    // Connections older than this are closed before reuse
    maxConnectionLifetime time.Duration
    
    // MaxConnectionIdleTime is the maximum time a connection can be idle
    // Idle connections older than this may be closed
    maxConnectionIdleTime time.Duration
    
    // HealthCheckPeriod is how often to perform background health checks
    healthCheckPeriod time.Duration
    
    // ConnectTimeout is the timeout for establishing new connections
    connectTimeout time.Duration
    
    // Logger for connection pool events (optional)
    logger logger.Logger
}

// DefaultConfig returns a Config with production-ready defaults.
func DefaultConfig() *Config {
    return &Config{
        maxConnections:        25,
        minConnections:        5,
        maxConnectionLifetime: 1 * time.Hour,
        maxConnectionIdleTime: 30 * time.Minute,
        healthCheckPeriod:     1 * time.Minute,
        connectTimeout:        10 * time.Second,
        logger:                nil,
    }
}

// Builder methods

// WithDSN sets the connection string.
func (c *Config) WithDSN(dsn string) *Config {
    c.dsn = dsn
    return c
}

// WithMaxConnections sets the maximum number of connections.
func (c *Config) WithMaxConnections(max int32) *Config {
    c.maxConnections = max
    return c
}

// WithMinConnections sets the minimum number of idle connections.
func (c *Config) WithMinConnections(min int32) *Config {
    c.minConnections = min
    return c
}

// WithMaxConnectionLifetime sets the maximum connection lifetime.
func (c *Config) WithMaxConnectionLifetime(d time.Duration) *Config {
    c.maxConnectionLifetime = d
    return c
}

// WithMaxConnectionIdleTime sets the maximum connection idle time.
func (c *Config) WithMaxConnectionIdleTime(d time.Duration) *Config {
    c.maxConnectionIdleTime = d
    return c
}

// WithHealthCheckPeriod sets the health check period.
func (c *Config) WithHealthCheckPeriod(d time.Duration) *Config {
    c.healthCheckPeriod = d
    return c
}

// WithConnectTimeout sets the connection timeout.
func (c *Config) WithConnectTimeout(d time.Duration) *Config {
    c.connectTimeout = d
    return c
}

// WithLogger sets the logger for connection pool events.
func (c *Config) WithLogger(l logger.Logger) *Config {
    c.logger = l
    return c
}

// Getters (for reading values)

func (c *Config) DSN() string                            { return c.dsn }
func (c *Config) MaxConnections() int32                  { return c.maxConnections }
func (c *Config) MinConnections() int32                  { return c.minConnections }
func (c *Config) MaxConnectionLifetime() time.Duration   { return c.maxConnectionLifetime }
func (c *Config) MaxConnectionIdleTime() time.Duration   { return c.maxConnectionIdleTime }
func (c *Config) HealthCheckPeriod() time.Duration       { return c.healthCheckPeriod }
func (c *Config) ConnectTimeout() time.Duration          { return c.connectTimeout }
func (c *Config) Logger() logger.Logger                  { return c.logger }

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
    if c.dsn == "" {
        return errors.RequiredField("dsn")
    }
    
    if c.maxConnections <= 0 {
        return errors.InvalidField("max_connections", "must be greater than 0")
    }
    
    if c.minConnections < 0 {
        return errors.InvalidField("min_connections", "must be non-negative")
    }
    
    if c.minConnections > c.maxConnections {
        return errors.Validation(
            fmt.Sprintf("min_connections (%d) cannot exceed max_connections (%d)", 
                c.minConnections, c.maxConnections),
        )
    }
    
    if c.maxConnectionLifetime <= 0 {
        return errors.InvalidField("max_connection_lifetime", "must be greater than 0")
    }
    
    if c.maxConnectionIdleTime <= 0 {
        return errors.InvalidField("max_connection_idle_time", "must be greater than 0")
    }
    
    if c.connectTimeout <= 0 {
        return errors.InvalidField("connect_timeout", "must be greater than 0")
    }
    
    return nil
}

// ToPgxConfig converts our Config to pgxpool.Config.
// This is where we map our abstraction to the underlying pgx library.
func (c *Config) ToPgxConfig(ctx context.Context) (*pgxpool.Config, error) {
    // Parse the DSN into pgxpool config
    poolConfig, err := pgxpool.ParseConfig(c.dsn)
    if err != nil {
        return nil, pkgerrors.New(
            pkgerrors.ErrorTypeValidation,
            "INVALID_DSN",
            "failed to parse database connection string",
        ).WithCause(err).
            WithDetail("dsn", maskPassword(c.dsn))
    }
    
    // Set connection pool parameters
    poolConfig.MaxConns = c.maxConnections
    poolConfig.MinConns = c.minConnections
    poolConfig.MaxConnLifetime = c.maxConnectionLifetime
    poolConfig.MaxConnIdleTime = c.maxConnectionIdleTime
    poolConfig.HealthCheckPeriod = c.healthCheckPeriod
    
    // Set connection timeout
    poolConfig.ConnConfig.ConnectTimeout = c.connectTimeout
    
    // TODO: Add pgx logger adapter when we implement it
    // if c.logger != nil {
    //     poolConfig.ConnConfig.Logger = newPgxLoggerAdapter(c.logger)
    // }
    
    return poolConfig, nil
}

// maskPassword masks the password in a DSN for safe logging.
// Format: postgres://user:password@host:port/dbname
func maskPassword(dsn string) string {
    // Simple masking - find password between : and @
    start := -1
    end := -1
    
    for i := 0; i < len(dsn); i++ {
        if dsn[i] == ':' && start == -1 {
            // Skip the first : (after postgres://)
            if i > 0 && dsn[i-1] == '/' {
                continue
            }
            start = i + 1
        }
        if dsn[i] == '@' && start != -1 {
            end = i
            break
        }
    }
    
    if start == -1 || end == -1 {
        return dsn // No password found, return as-is
    }
    
    return dsn[:start] + "****" + dsn[end:]
}