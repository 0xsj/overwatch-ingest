// gateway/config/config.go
package config

import (
	"time"

	pkgconfig "github.com/0xsj/scout/platform/pkg/config"
)

// Config defines the complete configuration for the Gateway service.
// It follows the interface-first design principle where the domain
// defines what configuration it needs.
type Config interface {
	// Application settings
	Port() int
	Environment() string
	LogLevel() string

	// Infrastructure
	Postgres() PostgresConfig
	Redis() RedisConfig
	NATS() NATSConfig
	RabbitMQ() RabbitMQConfig
	Observability() ObservabilityConfig
}

// PostgresConfig defines PostgreSQL database configuration.
type PostgresConfig interface {
	Host() string
	Port() int
	User() string
	Password() string
	Database() string
	MaxConnections() int
	MaxIdleConnections() int
	ConnectionTimeout() time.Duration
}

// RedisConfig defines Redis cache configuration.
type RedisConfig interface {
	Host() string
	Port() int
	MaxRetries() int
	PoolSize() int
}

// NATSConfig defines NATS messaging configuration.
type NATSConfig interface {
	URL() string
	MaxReconnects() int
	ReconnectWait() time.Duration
}

// RabbitMQConfig defines RabbitMQ queue configuration.
type RabbitMQConfig interface {
	URL() string
	MaxChannels() int
	PrefetchCount() int
}

// ObservabilityConfig defines observability configuration.
type ObservabilityConfig interface {
	OTELEndpoint() string
	JaegerEndpoint() string
	EnableTracing() bool
	EnableMetrics() bool
}

// envConfig is the concrete implementation that loads from environment variables.
type envConfig struct {
	port          int
	environment   string
	logLevel      string
	postgres      *envPostgresConfig
	redis         *envRedisConfig
	nats          *envNATSConfig
	rabbitmq      *envRabbitMQConfig
	observability *envObservabilityConfig
}

func (c *envConfig) Port() int                          { return c.port }
func (c *envConfig) Environment() string                { return c.environment }
func (c *envConfig) LogLevel() string                   { return c.logLevel }
func (c *envConfig) Postgres() PostgresConfig           { return c.postgres }
func (c *envConfig) Redis() RedisConfig                 { return c.redis }
func (c *envConfig) NATS() NATSConfig                   { return c.nats }
func (c *envConfig) RabbitMQ() RabbitMQConfig           { return c.rabbitmq }
func (c *envConfig) Observability() ObservabilityConfig { return c.observability }

// Postgres implementation
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

// Redis implementation
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

// NATS implementation
type envNATSConfig struct {
	url           string
	maxReconnects int
	reconnectWait time.Duration
}

func (c *envNATSConfig) URL() string                  { return c.url }
func (c *envNATSConfig) MaxReconnects() int           { return c.maxReconnects }
func (c *envNATSConfig) ReconnectWait() time.Duration { return c.reconnectWait }

// RabbitMQ implementation
type envRabbitMQConfig struct {
	url           string
	maxChannels   int
	prefetchCount int
}

func (c *envRabbitMQConfig) URL() string        { return c.url }
func (c *envRabbitMQConfig) MaxChannels() int   { return c.maxChannels }
func (c *envRabbitMQConfig) PrefetchCount() int { return c.prefetchCount }

// Observability implementation
type envObservabilityConfig struct {
	otelEndpoint   string
	jaegerEndpoint string
	enableTracing  bool
	enableMetrics  bool
}

func (c *envObservabilityConfig) OTELEndpoint() string   { return c.otelEndpoint }
func (c *envObservabilityConfig) JaegerEndpoint() string { return c.jaegerEndpoint }
func (c *envObservabilityConfig) EnableTracing() bool    { return c.enableTracing }
func (c *envObservabilityConfig) EnableMetrics() bool    { return c.enableMetrics }

// Load loads the Gateway configuration from environment variables with GATEWAY_ prefix.
func Load() (Config, error) {
	const prefix = "GATEWAY_"

	// Application config
	port, err := pkgconfig.LoadPortOptional(pkgconfig.WithPrefix(prefix, "PORT"), 8080)
	if err != nil {
		return nil, err
	}

	environment := pkgconfig.LoadStringOptional(
		pkgconfig.WithPrefix(prefix, "ENV"),
		"development",
	)

	logLevel := pkgconfig.LoadStringOptional(
		pkgconfig.WithPrefix(prefix, "LOG_LEVEL"),
		"info",
	)

	// Postgres config
	postgresHost, err := pkgconfig.LoadStringRequired(pkgconfig.WithPrefix(prefix, "POSTGRES_HOST"))
	if err != nil {
		return nil, err
	}

	postgresPort, err := pkgconfig.LoadPortOptional(pkgconfig.WithPrefix(prefix, "POSTGRES_PORT"), 5432)
	if err != nil {
		return nil, err
	}

	postgresUser, err := pkgconfig.LoadStringRequired(pkgconfig.WithPrefix(prefix, "POSTGRES_USER"))
	if err != nil {
		return nil, err
	}

	postgresPassword, err := pkgconfig.LoadStringRequired(pkgconfig.WithPrefix(prefix, "POSTGRES_PASSWORD"))
	if err != nil {
		return nil, err
	}

	postgresDB, err := pkgconfig.LoadStringRequired(pkgconfig.WithPrefix(prefix, "POSTGRES_DB"))
	if err != nil {
		return nil, err
	}

	postgresMaxConns, err := pkgconfig.LoadIntOptional(
		pkgconfig.WithPrefix(prefix, "POSTGRES_MAX_CONNECTIONS"),
		25,
	)
	if err != nil {
		return nil, err
	}

	postgresMaxIdleConns, err := pkgconfig.LoadIntOptional(
		pkgconfig.WithPrefix(prefix, "POSTGRES_MAX_IDLE_CONNECTIONS"),
		5,
	)
	if err != nil {
		return nil, err
	}

	postgresConnTimeout, err := pkgconfig.LoadDurationOptional(
		pkgconfig.WithPrefix(prefix, "POSTGRES_CONNECTION_TIMEOUT"),
		10*time.Second,
	)
	if err != nil {
		return nil, err
	}

	// Redis config
	redisHost, err := pkgconfig.LoadStringRequired(pkgconfig.WithPrefix(prefix, "REDIS_HOST"))
	if err != nil {
		return nil, err
	}

	redisPort, err := pkgconfig.LoadPortOptional(pkgconfig.WithPrefix(prefix, "REDIS_PORT"), 6379)
	if err != nil {
		return nil, err
	}

	redisMaxRetries, err := pkgconfig.LoadIntOptional(
		pkgconfig.WithPrefix(prefix, "REDIS_MAX_RETRIES"),
		3,
	)
	if err != nil {
		return nil, err
	}

	redisPoolSize, err := pkgconfig.LoadIntOptional(
		pkgconfig.WithPrefix(prefix, "REDIS_POOL_SIZE"),
		10,
	)
	if err != nil {
		return nil, err
	}

	// NATS config
	natsURL, err := pkgconfig.LoadStringRequired(pkgconfig.WithPrefix(prefix, "NATS_URL"))
	if err != nil {
		return nil, err
	}

	natsMaxReconnects, err := pkgconfig.LoadIntOptional(
		pkgconfig.WithPrefix(prefix, "NATS_MAX_RECONNECTS"),
		10,
	)
	if err != nil {
		return nil, err
	}

	natsReconnectWait, err := pkgconfig.LoadDurationOptional(
		pkgconfig.WithPrefix(prefix, "NATS_RECONNECT_WAIT"),
		2*time.Second,
	)
	if err != nil {
		return nil, err
	}

	// RabbitMQ config
	rabbitmqURL, err := pkgconfig.LoadStringRequired(pkgconfig.WithPrefix(prefix, "RABBITMQ_URL"))
	if err != nil {
		return nil, err
	}

	rabbitmqMaxChannels, err := pkgconfig.LoadIntOptional(
		pkgconfig.WithPrefix(prefix, "RABBITMQ_MAX_CHANNELS"),
		100,
	)
	if err != nil {
		return nil, err
	}

	rabbitmqPrefetch, err := pkgconfig.LoadIntOptional(
		pkgconfig.WithPrefix(prefix, "RABBITMQ_PREFETCH_COUNT"),
		10,
	)
	if err != nil {
		return nil, err
	}

	// Observability config
	otelEndpoint := pkgconfig.LoadStringOptional(
		pkgconfig.WithPrefix(prefix, "OTEL_EXPORTER_OTLP_ENDPOINT"),
		"http://localhost:4317",
	)

	jaegerEndpoint := pkgconfig.LoadStringOptional(
		pkgconfig.WithPrefix(prefix, "JAEGER_ENDPOINT"),
		"http://localhost:14268/api/traces",
	)

	enableTracing, err := pkgconfig.LoadBoolOptional(
		pkgconfig.WithPrefix(prefix, "ENABLE_TRACING"),
		true,
	)
	if err != nil {
		return nil, err
	}

	enableMetrics, err := pkgconfig.LoadBoolOptional(
		pkgconfig.WithPrefix(prefix, "ENABLE_METRICS"),
		true,
	)
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
		rabbitmq: &envRabbitMQConfig{
			url:           rabbitmqURL,
			maxChannels:   rabbitmqMaxChannels,
			prefetchCount: rabbitmqPrefetch,
		},
		observability: &envObservabilityConfig{
			otelEndpoint:   otelEndpoint,
			jaegerEndpoint: jaegerEndpoint,
			enableTracing:  enableTracing,
			enableMetrics:  enableMetrics,
		},
	}, nil
}
