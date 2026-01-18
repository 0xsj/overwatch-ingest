package config

import (
	"fmt"
	"time"

	pkgconfig "github.com/0xsj/overwatch-pkg/config"
)

type Config struct {
	Server          ServerConfig
	Database        DatabaseConfig
	Redis           RedisConfig
	NATS            NATSConfig
	ServiceIdentity ServiceIdentityConfig
	Ingest          IngestConfig
}

type ServerConfig struct {
	Host              string        `env:"SERVER_HOST" default:"0.0.0.0"`
	Port              int           `env:"SERVER_PORT" default:"50055"`
	EnableReflection  bool          `env:"SERVER_ENABLE_REFLECTION" default:"true"`
	EnableHealthCheck bool          `env:"SERVER_ENABLE_HEALTH_CHECK" default:"true"`
	ShutdownTimeout   time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" default:"30s"`
}

type DatabaseConfig struct {
	Host              string        `env:"DATABASE_HOST" default:"localhost"`
	Port              int           `env:"DATABASE_PORT" default:"5450"`
	User              string        `env:"DATABASE_USER" default:"overwatch"`
	Password          string        `env:"DATABASE_PASSWORD" default:"overwatch" sensitive:"true"`
	Database          string        `env:"DATABASE_NAME" default:"overwatch_ingest"`
	SSLMode           string        `env:"DATABASE_SSL_MODE" default:"disable"`
	MaxConns          int           `env:"DATABASE_MAX_CONNS" default:"25"`
	MinConns          int           `env:"DATABASE_MIN_CONNS" default:"5"`
	MaxConnLifetime   time.Duration `env:"DATABASE_MAX_CONN_LIFETIME" default:"1h"`
	MaxConnIdleTime   time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME" default:"30m"`
	HealthCheckPeriod time.Duration `env:"DATABASE_HEALTH_CHECK_PERIOD" default:"1m"`
}

type RedisConfig struct {
	Host         string        `env:"REDIS_HOST" default:"localhost"`
	Port         int           `env:"REDIS_PORT" default:"6390"`
	Password     string        `env:"REDIS_PASSWORD" default:"" sensitive:"true"`
	DB           int           `env:"REDIS_DB" default:"3"`
	PoolSize     int           `env:"REDIS_POOL_SIZE" default:"10"`
	MinIdleConns int           `env:"REDIS_MIN_IDLE_CONNS" default:"5"`
	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT" default:"5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT" default:"3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT" default:"3s"`
}

type NATSConfig struct {
	URL           string        `env:"NATS_URL" default:"nats://localhost:4230"`
	SubjectPrefix string        `env:"NATS_SUBJECT_PREFIX" default:"overwatch.ingest"`
	MaxReconnects int           `env:"NATS_MAX_RECONNECTS" default:"10"`
	ReconnectWait time.Duration `env:"NATS_RECONNECT_WAIT" default:"2s"`
}

type ServiceIdentityConfig struct {
	ID                string `env:"SERVICE_IDENTITY_ID" default:"ingest-service"`
	Name              string `env:"SERVICE_IDENTITY_NAME" default:"ingest"`
	PrivateKeyPath    string `env:"SERVICE_IDENTITY_PRIVATE_KEY_PATH" default:""`
	PrivateKeyBase64  string `env:"SERVICE_IDENTITY_PRIVATE_KEY" default:"" sensitive:"true"`
	GenerateIfMissing bool   `env:"SERVICE_IDENTITY_GENERATE_IF_MISSING" default:"true"`
}

type IngestConfig struct {
	AcceptThreshold                float64       `env:"ACCEPT_THRESHOLD" default:"0.7"`
	RejectThreshold                float64       `env:"REJECT_THRESHOLD" default:"0.3"`
	QuarantineExpiry               time.Duration `env:"QUARANTINE_EXPIRY" default:"72h"`
	ReliabilityQuarantineThreshold float64       `env:"RELIABILITY_QUARANTINE_THRESHOLD" default:"0.4"`
	MinReliabilityRecords          int64         `env:"MIN_RELIABILITY_RECORDS" default:"10"`
	RequireCollectorSignature      bool          `env:"REQUIRE_COLLECTOR_SIGNATURE" default:"true"`
	RequireSourceSignature         bool          `env:"REQUIRE_SOURCE_SIGNATURE" default:"false"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := pkgconfig.Load(cfg, pkgconfig.WithPrefix("INGEST_")); err != nil {
		return nil, err
	}
	return cfg, nil
}

func MustLoad() *Config {
	cfg := &Config{}
	pkgconfig.MustLoad(cfg, pkgconfig.WithPrefix("INGEST_"))
	return cfg
}

func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

func (c *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *ServiceIdentityConfig) HasPrivateKey() bool {
	return c.PrivateKeyBase64 != "" || c.PrivateKeyPath != ""
}
