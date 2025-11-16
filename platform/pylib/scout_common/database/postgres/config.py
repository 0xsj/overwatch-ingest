"""PostgreSQL client configuration."""

from dataclasses import dataclass, field
from datetime import timedelta
from typing import Optional
from urllib.parse import quote_plus

from scout_common.errors import Error, required_field, invalid_field, validation_error
from scout_common.logging import Logger


@dataclass
class Config:
    """
    PostgreSQL connection pool configuration.
    
    This configuration is used to create a connection pool with psycopg.
    
    Attributes:
        dsn: PostgreSQL connection string
        max_connections: Maximum number of connections in the pool
        min_connections: Minimum number of idle connections to maintain
        max_connection_lifetime: Maximum duration a connection can be reused
        max_connection_idle_time: Maximum time a connection can be idle
        connect_timeout: Timeout for establishing new connections
        logger: Optional logger for connection pool events
    """
    
    dsn: str = ""
    max_connections: int = 25
    min_connections: int = 5
    max_connection_lifetime: timedelta = field(default_factory=lambda: timedelta(hours=1))
    max_connection_idle_time: timedelta = field(default_factory=lambda: timedelta(minutes=30))
    connect_timeout: timedelta = field(default_factory=lambda: timedelta(seconds=10))
    logger: Optional[Logger] = None
    
    def with_dsn(self, dsn: str) -> "Config":
        """Set the connection string."""
        self.dsn = dsn
        return self
    
    def with_max_connections(self, max_conns: int) -> "Config":
        """Set the maximum number of connections."""
        self.max_connections = max_conns
        return self
    
    def with_min_connections(self, min_conns: int) -> "Config":
        """Set the minimum number of idle connections."""
        self.min_connections = min_conns
        return self
    
    def with_max_connection_lifetime(self, lifetime: timedelta) -> "Config":
        """Set the maximum connection lifetime."""
        self.max_connection_lifetime = lifetime
        return self
    
    def with_max_connection_idle_time(self, idle_time: timedelta) -> "Config":
        """Set the maximum connection idle time."""
        self.max_connection_idle_time = idle_time
        return self
    
    def with_connect_timeout(self, timeout: timedelta) -> "Config":
        """Set the connection timeout."""
        self.connect_timeout = timeout
        return self
    
    def with_logger(self, logger: Logger) -> "Config":
        """Set the logger for connection pool events."""
        self.logger = logger
        return self
    
    def validate(self) -> Optional[Error]:
        """
        Validate the configuration.
        
        Returns:
            Error if validation fails, None otherwise
        """
        if not self.dsn:
            return required_field("dsn")
        
        if self.max_connections <= 0:
            return invalid_field("max_connections", "must be greater than 0")
        
        if self.min_connections < 0:
            return invalid_field("min_connections", "must be non-negative")
        
        if self.min_connections > self.max_connections:
            return validation_error(
                f"min_connections ({self.min_connections}) cannot exceed "
                f"max_connections ({self.max_connections})"
            )
        
        if self.max_connection_lifetime.total_seconds() <= 0:
            return invalid_field("max_connection_lifetime", "must be greater than 0")
        
        if self.max_connection_idle_time.total_seconds() <= 0:
            return invalid_field("max_connection_idle_time", "must be greater than 0")
        
        if self.connect_timeout.total_seconds() <= 0:
            return invalid_field("connect_timeout", "must be greater than 0")
        
        return None


def default_config() -> Config:
    """
    Return a Config with production-ready defaults.
    
    Returns:
        Config with default values
    """
    return Config()


def from_service_config(
    host: str,
    port: int,
    user: str,
    password: str,
    database: str,
    max_connections: int = 25,
    max_idle_connections: int = 5,
    connection_timeout: timedelta = timedelta(seconds=10),
    logger: Optional[Logger] = None,
) -> Config:
    """
    Create a database Config from service configuration parameters.
    
    This allows seamless integration with existing service configuration.
    
    Args:
        host: PostgreSQL host
        port: PostgreSQL port
        user: Database user
        password: Database password
        database: Database name
        max_connections: Maximum connections in pool
        max_idle_connections: Minimum idle connections
        connection_timeout: Connection timeout
        logger: Optional logger
        
    Returns:
        Configured Config instance
        
    Example:
        from services.agents.config import load
        
        service_cfg = load(use_prefix=False)
        pg_cfg = service_cfg.postgres
        
        db_cfg = from_service_config(
            host=pg_cfg.host,
            port=pg_cfg.port,
            user=pg_cfg.user,
            password=pg_cfg.password,
            database=pg_cfg.database,
            max_connections=pg_cfg.max_connections,
            max_idle_connections=pg_cfg.max_idle_connections,
            connection_timeout=pg_cfg.connection_timeout,
            logger=app_logger,
        )
    """
    # Build DSN
    dsn = _build_dsn(host, port, user, password, database)
    
    # Create config with defaults
    config = default_config()
    config.dsn = dsn
    config.max_connections = max_connections
    config.min_connections = max_idle_connections
    config.connect_timeout = connection_timeout
    config.logger = logger
    
    return config


def _build_dsn(
    host: str,
    port: int,
    user: str,
    password: str,
    database: str,
    sslmode: str = "disable",
) -> str:
    """
    Build a PostgreSQL connection string (DSN).
    
    Args:
        host: Database host
        port: Database port
        user: Database user
        password: Database password
        database: Database name
        sslmode: SSL mode (default: disable for local dev)
        
    Returns:
        PostgreSQL connection string
    """
    # URL-encode password to handle special characters
    encoded_password = quote_plus(password)
    
    return (
        f"postgresql://{user}:{encoded_password}@{host}:{port}/{database}"
        f"?sslmode={sslmode}"
    )


def mask_password(dsn: str) -> str:
    """
    Mask the password in a DSN for safe logging.
    
    Args:
        dsn: PostgreSQL connection string
        
    Returns:
        DSN with masked password
        
    Example:
        >>> mask_password("postgresql://user:secret@localhost:5432/db")
        "postgresql://user:****@localhost:5432/db"
    """
    # Find password between : and @
    try:
        # Split by :// to get past the protocol
        if "://" not in dsn:
            return dsn
        
        protocol, rest = dsn.split("://", 1)
        
        # Find @ to locate the password section
        if "@" not in rest:
            return dsn
        
        credentials, host_part = rest.split("@", 1)
        
        # Find : in credentials to separate user and password
        if ":" not in credentials:
            return dsn
        
        user, _ = credentials.split(":", 1)
        
        # Rebuild with masked password
        return f"{protocol}://{user}:****@{host_part}"
    
    except Exception:
        # If anything fails, return original (better safe than sorry)
        return dsn