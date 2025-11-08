"""Configuration for Analytics service."""

from dataclasses import dataclass
from datetime import timedelta

from scout_common.config import (
    load_port_optional,
    load_string_optional,
    load_string_required,
    load_int_optional,
    load_duration_optional,
    load_bool_optional,
)


@dataclass
class PostgresConfig:
    """PostgreSQL configuration."""
    
    host: str
    port: int
    user: str
    password: str
    database: str
    max_connections: int
    max_idle_connections: int
    connection_timeout: timedelta


@dataclass
class RedisConfig:
    """Redis configuration."""
    
    host: str
    port: int
    max_retries: int
    pool_size: int


@dataclass
class NATSConfig:
    """NATS configuration."""
    
    url: str
    max_reconnects: int
    reconnect_wait: timedelta


@dataclass
class ObservabilityConfig:
    """Observability configuration."""
    
    otel_endpoint: str
    enable_tracing: bool
    enable_metrics: bool


@dataclass
class Config:
    """Complete configuration for Analytics service."""
    
    port: int
    environment: str
    log_level: str
    postgres: PostgresConfig
    redis: RedisConfig
    nats: NATSConfig
    observability: ObservabilityConfig


def load(use_prefix: bool = False) -> Config:
    """
    Load Analytics configuration from environment variables.
    
    Args:
        use_prefix: If True, use ANALYTICS_ prefix. If False, no prefix.
        
    Returns:
        Loaded configuration
        
    Raises:
        Error: If required configuration is missing or invalid
    """
    prefix = "ANALYTICS_" if use_prefix else ""
    
    def key(name: str) -> str:
        return prefix + name
    
    # Application config
    port = load_port_optional(key("PORT"), 8084)
    environment = load_string_optional(key("ENV"), "development")
    log_level = load_string_optional(key("LOG_LEVEL"), "info")
    
    # Postgres config
    postgres = PostgresConfig(
        host=load_string_required(key("POSTGRES_HOST")),
        port=load_port_optional(key("POSTGRES_PORT"), 5432),
        user=load_string_required(key("POSTGRES_USER")),
        password=load_string_required(key("POSTGRES_PASSWORD")),
        database=load_string_required(key("POSTGRES_DB")),
        max_connections=load_int_optional(key("POSTGRES_MAX_CONNECTIONS"), 25),
        max_idle_connections=load_int_optional(key("POSTGRES_MAX_IDLE_CONNECTIONS"), 5),
        connection_timeout=load_duration_optional(
            key("POSTGRES_CONNECTION_TIMEOUT"),
            timedelta(seconds=10),
        ),
    )
    
    # Redis config
    redis = RedisConfig(
        host=load_string_required(key("REDIS_HOST")),
        port=load_port_optional(key("REDIS_PORT"), 6379),
        max_retries=load_int_optional(key("REDIS_MAX_RETRIES"), 3),
        pool_size=load_int_optional(key("REDIS_POOL_SIZE"), 10),
    )
    
    # NATS config
    nats = NATSConfig(
        url=load_string_required(key("NATS_URL")),
        max_reconnects=load_int_optional(key("NATS_MAX_RECONNECTS"), 10),
        reconnect_wait=load_duration_optional(
            key("NATS_RECONNECT_WAIT"),
            timedelta(seconds=2),
        ),
    )
    
    # Observability config
    observability = ObservabilityConfig(
        otel_endpoint=load_string_optional(
            key("OTEL_EXPORTER_OTLP_ENDPOINT"),
            "http://localhost:4317",
        ),
        enable_tracing=load_bool_optional(key("ENABLE_TRACING"), True),
        enable_metrics=load_bool_optional(key("ENABLE_METRICS"), True),
    )
    
    return Config(
        port=port,
        environment=environment,
        log_level=log_level,
        postgres=postgres,
        redis=redis,
        nats=nats,
        observability=observability,
    )