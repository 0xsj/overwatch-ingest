package postgres

import (
	"fmt"
	"time"

	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// ServicePostgresConfig is an interface that service configs must implement
// to be compatible with our database package.
// This matches the existing PostgresConfig interface in your services.
type ServicePostgresConfig interface {
    Host() string
    Port() int
    User() string
    Password() string
    Database() string
    MaxConnections() int
    MaxIdleConnections() int
    ConnectionTimeout() time.Duration
}

// FromServiceConfig creates a database Config from a service's PostgresConfig.
// This allows seamless integration with existing service configuration.
//
// Example:
//   serviceCfg, _ := config.Load(false)
//   dbCfg := postgres.FromServiceConfig(serviceCfg.Postgres())
//   client, _ := postgres.NewClient(ctx, dbCfg)
func FromServiceConfig(cfg ServicePostgresConfig, opts ...ConfigOption) *Config {
    // Build DSN from service config
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=disable",
        cfg.User(),
        cfg.Password(),
        cfg.Host(),
        cfg.Port(),
        cfg.Database(),
    )
    
    // Create config with defaults
    dbCfg := DefaultConfig().
        WithDSN(dsn).
        WithMaxConnections(int32(cfg.MaxConnections())).
        WithMinConnections(int32(cfg.MaxIdleConnections())).
        WithConnectTimeout(cfg.ConnectionTimeout())
    
    // Apply any additional options
    for _, opt := range opts {
        opt(dbCfg)
    }
    
    return dbCfg
}

// ConfigOption is a function that modifies a Config.
// This allows optional overrides when using FromServiceConfig.
type ConfigOption func(*Config)

// WithLoggerOption returns a ConfigOption that sets the logger.
func WithLoggerOption(l logger.Logger) ConfigOption {
    return func(c *Config) {
        c.logger = l
    }
}

// WithMaxConnectionLifetimeOption returns a ConfigOption that sets max lifetime.
func WithMaxConnectionLifetimeOption(d time.Duration) ConfigOption {
    return func(c *Config) {
        c.maxConnectionLifetime = d
    }
}

// WithMaxConnectionIdleTimeOption returns a ConfigOption that sets max idle time.
func WithMaxConnectionIdleTimeOption(d time.Duration) ConfigOption {
    return func(c *Config) {
        c.maxConnectionIdleTime = d
    }
}

// WithHealthCheckPeriodOption returns a ConfigOption that sets health check period.
func WithHealthCheckPeriodOption(d time.Duration) ConfigOption {
    return func(c *Config) {
        c.healthCheckPeriod = d
    }
}