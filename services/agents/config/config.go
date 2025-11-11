// services/agents/config/config.go
package config

import (
	"fmt"
	"time"

	pkgconfig "github.com/0xsj/scout/platform/pkg/config"
)

// Config defines the complete configuration for the Agents service.
type Config interface {
	// Application settings
	Port() int
	Environment() string
	LogLevel() string

	// Infrastructure
	Postgres() PostgresConfig
	Redis() RedisConfig
	NATS() NATSConfig
	Observability() ObservabilityConfig
}

type PostgresConfig interface {
	Host() string
	Port() int
	User() string
	Password() string
	Database() string
	MaxConnections() int
	MaxIdleConnections() int
	ConnectionTimeout() time.Duration
	// DSN returns the full connection string
	DSN() string
}

type RedisConfig interface {
	Host() string
	Port() int
	MaxRetries() int
	PoolSize() int
}

type NATSConfig interface {
	URL() string
	MaxReconnects() int
	ReconnectWait() time.Duration
}

type ObservabilityConfig interface {
	OTELEndpoint() string
	EnableTracing() bool
	EnableMetrics() bool
}

// Implementation structs
type envConfig struct {
	port          int
	environment   string
	logLevel      string
	postgres      *envPostgresConfig
	redis         *envRedisConfig
	nats          *envNATSConfig
	observability *envObservabilityConfig
}

func (c *envConfig) Port() int                          { return c.port }
func (c *envConfig) Environment() string                { return c.environment }
func (c *envConfig) LogLevel() string                   { return c.logLevel }
func (c *envConfig) Postgres() PostgresConfig           { return c.postgres }
func (c *envConfig) Redis() RedisConfig                 { return c.redis }
func (c *envConfig) NATS() NATSConfig                   { return c.nats }
func (c *envConfig) Observability() ObservabilityConfig { return c.observability }

type envPostgresConfig struct {
	host               string
	port               int
	user               string
	password           string
	database           string
	maxConnections     int
	maxIdleConnections int
	connectionTimeout  time.Duration
}

func (c *envPostgresConfig) Host() string                     { return c.host }
func (c *envPostgresConfig) Port() int                        { return c.port }
func (c *envPostgresConfig) User() string                     { return c.user }
func (c *envPostgresConfig) Password() string                 { return c.password }
func (c *envPostgresConfig) Database() string                 { return c.database }
func (c *envPostgresConfig) MaxConnections() int              { return c.maxConnections }
func (c *envPostgresConfig) MaxIdleConnections() int          { return c.maxIdleConnections }
func (c *envPostgresConfig) ConnectionTimeout() time.Duration { return c.connectionTimeout }

// DSN returns the PostgreSQL connection string
func (c *envPostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.host, c.port, c.user, c.password, c.database,
	)
}

type envRedisConfig struct {
	host       string
	port       int
	maxRetries int
	poolSize   int
}

func (c *envRedisConfig) Host() string    { return c.host }
func (c *envRedisConfig) Port() int       { return c.port }
func (c *envRedisConfig) MaxRetries() int { return c.maxRetries }
func (c *envRedisConfig) PoolSize() int   { return c.poolSize }

type envNATSConfig struct {
	url           string
	maxReconnects int
	reconnectWait time.Duration
}

func (c *envNATSConfig) URL() string                  { return c.url }
func (c *envNATSConfig) MaxReconnects() int           { return c.maxReconnects }
func (c *envNATSConfig) ReconnectWait() time.Duration { return c.reconnectWait }

type envObservabilityConfig struct {
	otelEndpoint  string
	enableTracing bool
	enableMetrics bool
}

func (c *envObservabilityConfig) OTELEndpoint() string { return c.otelEndpoint }
func (c *envObservabilityConfig) EnableTracing() bool  { return c.enableTracing }
func (c *envObservabilityConfig) EnableMetrics() bool  { return c.enableMetrics }

// Load loads the Agents configuration from environment variables.
func Load(usePrefix bool) (Config, error) {
	prefix := ""
	if usePrefix {
		prefix = "AGENTS_"
	}

	key := func(name string) string {
		return prefix + name
	}

	// Application config
	port, err := pkgconfig.LoadPortOptional(key("PORT"), 8081)
	if err != nil {
		return nil, err
	}

	environment := pkgconfig.LoadStringOptional(key("ENV"), "development")
	logLevel := pkgconfig.LoadStringOptional(key("LOG_LEVEL"), "info")

	// Postgres config
	postgresHost := pkgconfig.LoadStringOptional(key("POSTGRES_HOST"), "localhost")
	
	postgresPort, err := pkgconfig.LoadPortOptional(key("POSTGRES_PORT"), 5432)
	if err != nil {
		return nil, err
	}

	postgresUser := pkgconfig.LoadStringOptional(key("POSTGRES_USER"), "scout")
	postgresPassword := pkgconfig.LoadStringOptional(key("POSTGRES_PASSWORD"), "scout_dev_password")
	postgresDB := pkgconfig.LoadStringOptional(key("POSTGRES_DB"), "agents_db")

	postgresMaxConns, err := pkgconfig.LoadIntOptional(key("POSTGRES_MAX_CONNECTIONS"), 25)
	if err != nil {
		return nil, err
	}

	postgresMaxIdleConns, err := pkgconfig.LoadIntOptional(key("POSTGRES_MAX_IDLE_CONNECTIONS"), 5)
	if err != nil {
		return nil, err
	}

	postgresConnTimeout, err := pkgconfig.LoadDurationOptional(key("POSTGRES_CONNECTION_TIMEOUT"), 10*time.Second)
	if err != nil {
		return nil, err
	}

	// Redis config
	redisHost := pkgconfig.LoadStringOptional(key("REDIS_HOST"), "localhost")
	
	redisPort, err := pkgconfig.LoadPortOptional(key("REDIS_PORT"), 6379)
	if err != nil {
		return nil, err
	}

	redisMaxRetries, err := pkgconfig.LoadIntOptional(key("REDIS_MAX_RETRIES"), 3)
	if err != nil {
		return nil, err
	}

	redisPoolSize, err := pkgconfig.LoadIntOptional(key("REDIS_POOL_SIZE"), 10)
	if err != nil {
		return nil, err
	}

	// NATS config
	natsURL := pkgconfig.LoadStringOptional(key("NATS_URL"), "nats://localhost:4222")
	
	natsMaxReconnects, err := pkgconfig.LoadIntOptional(key("NATS_MAX_RECONNECTS"), 10)
	if err != nil {
		return nil, err
	}

	natsReconnectWait, err := pkgconfig.LoadDurationOptional(key("NATS_RECONNECT_WAIT"), 2*time.Second)
	if err != nil {
		return nil, err
	}

	// Observability config
	otelEndpoint := pkgconfig.LoadStringOptional(key("OTEL_EXPORTER_OTLP_ENDPOINT"), "http://localhost:4317")

	enableTracing, err := pkgconfig.LoadBoolOptional(key("ENABLE_TRACING"), true)
	if err != nil {
		return nil, err
	}

	enableMetrics, err := pkgconfig.LoadBoolOptional(key("ENABLE_METRICS"), true)
	if err != nil {
		return nil, err
	}

	return &envConfig{
		port:        port,
		environment: environment,
		logLevel:    logLevel,
		postgres: &envPostgresConfig{
			host:               postgresHost,
			port:               postgresPort,
			user:               postgresUser,
			password:           postgresPassword,
			database:           postgresDB,
			maxConnections:     postgresMaxConns,
			maxIdleConnections: postgresMaxIdleConns,
			connectionTimeout:  postgresConnTimeout,
		},
		redis: &envRedisConfig{
			host:       redisHost,
			port:       redisPort,
			maxRetries: redisMaxRetries,
			poolSize:   redisPoolSize,
		},
		nats: &envNATSConfig{
			url:           natsURL,
			maxReconnects: natsMaxReconnects,
			reconnectWait: natsReconnectWait,
		},
		observability: &envObservabilityConfig{
			otelEndpoint:  otelEndpoint,
			enableTracing: enableTracing,
			enableMetrics: enableMetrics,
		},
	}, nil
}