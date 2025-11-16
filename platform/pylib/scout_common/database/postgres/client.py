"""PostgreSQL client with connection pooling."""

from contextlib import contextmanager
from datetime import timedelta
from typing import Any, Generator, Optional

import psycopg
from psycopg import Connection, Cursor
from psycopg.rows import dict_row
from psycopg_pool import ConnectionPool

from scout_common.database.postgres.config import Config, mask_password
from scout_common.database.postgres.errors import (
    connection_error,
    map_error,
)
from scout_common.errors import Error, required_field
from scout_common.logging import Logger


class PoolStats:
    """
    Connection pool statistics.
    
    This is a simplified wrapper around psycopg_pool statistics
    for easier consumption and monitoring.
    """
    
    def __init__(self, pool: ConnectionPool):
        """
        Initialize pool stats from a ConnectionPool.
        
        Args:
            pool: The connection pool to get stats from
        """
        self._pool = pool
    
    @property
    def size(self) -> int:
        """Current pool size (total connections)."""
        return self._pool.get_stats().get("pool_size", 0)
    
    @property
    def available(self) -> int:
        """Number of available (idle) connections."""
        return self._pool.get_stats().get("pool_available", 0)
    
    @property
    def max_size(self) -> int:
        """Maximum pool size."""
        return self._pool.max_size
    
    @property
    def min_size(self) -> int:
        """Minimum pool size."""
        return self._pool.min_size
    
    def __str__(self) -> str:
        """Return a human-readable representation of pool statistics."""
        return (
            f"total={self.size} available={self.available} "
            f"max={self.max_size} min={self.min_size}"
        )
    
    def health_status(self) -> str:
        """
        Return the health status of the pool.
        
        Returns:
            Health status string: "healthy", "warning", or "unhealthy"
        """
        if self.size == 0:
            return "unhealthy: no connections"
        
        if self.available == 0 and self.size >= self.max_size:
            return "warning: pool exhausted"
        
        if self.available == 0:
            return "warning: no available connections"
        
        return "healthy"


class Client:
    """
    PostgreSQL client with connection pooling.
    
    Provides a clean interface for database operations with automatic
    connection management and error handling.
    
    Example:
        config = default_config().with_dsn("postgresql://...")
        client = Client(config)
        client.open()
        
        try:
            with client.connection() as conn:
                with conn.cursor() as cur:
                    cur.execute("SELECT * FROM users")
                    users = cur.fetchall()
        finally:
            client.close()
    """
    
    def __init__(self, config: Config):
        """
        Initialize the PostgreSQL client.
        
        Args:
            config: Database configuration
            
        Raises:
            Error: If configuration is invalid
        """
        if config is None:
            raise required_field("config")
        
        # Validate configuration
        if err := config.validate():
            raise err
        
        self._config = config
        self._pool: Optional[ConnectionPool] = None
        self._logger = config.logger
    
    def open(self) -> None:
        """
        Open the connection pool.
        
        This establishes the initial connections and verifies connectivity.
        
        Raises:
            Error: If connection fails
            
        Example:
            client = Client(config)
            client.open()
        """
        if self._pool is not None:
            raise Error(
                error_type="INTERNAL",
                code="POOL_ALREADY_OPEN",
                message="connection pool already open",
            )
        
        try:
            # Create connection pool
            # Convert timedelta to seconds for psycopg_pool
            max_lifetime = int(self._config.max_connection_lifetime.total_seconds())
            max_idle = int(self._config.max_connection_idle_time.total_seconds())
            timeout = int(self._config.connect_timeout.total_seconds())
            
            self._pool = ConnectionPool(
                conninfo=self._config.dsn,
                min_size=self._config.min_connections,
                max_size=self._config.max_connections,
                max_lifetime=max_lifetime,
                max_idle=max_idle,
                timeout=timeout,
                open=True,  # Open pool immediately
            )
            
            # Verify connection with a ping
            with self._pool.connection() as conn:
                conn.execute("SELECT 1")
            
            if self._logger:
                self._logger.info(
                    "postgres client connected",
                    max_connections=self._config.max_connections,
                    min_connections=self._config.min_connections,
                )
        
        except Exception as e:
            raise connection_error("open pool", e)
    
    def close(self) -> None:
        """
        Close the connection pool.
        
        This closes all connections and releases resources.
        After calling close, the client should not be used.
        
        Example:
            client.close()
        """
        if self._logger:
            self._logger.info("postgres client closing")
        
        if self._pool is not None:
            self._pool.close()
            self._pool = None
        
        if self._logger:
            self._logger.info("postgres client closed")
    
    def ping(self) -> None:
        """
        Verify the connection to the database.
        
        Raises:
            Error: If ping fails
            
        Example:
            try:
                client.ping()
                print("Database is alive")
            except Error:
                print("Database is down")
        """
        if self._pool is None:
            raise connection_error("ping", Exception("pool not open"))
        
        try:
            with self._pool.connection() as conn:
                conn.execute("SELECT 1")
        except Exception as e:
            raise map_error(e, "ping")
    
    def stats(self) -> PoolStats:
        """
        Return connection pool statistics.
        
        Returns:
            PoolStats with current pool status
            
        Example:
            stats = client.stats()
            print(f"Pool: {stats}")
            print(f"Health: {stats.health_status()}")
        """
        if self._pool is None:
            raise connection_error("get stats", Exception("pool not open"))
        
        return PoolStats(self._pool)
    
    @contextmanager
    def connection(self) -> Generator[Connection, None, None]:
        """
        Get a connection from the pool.
        
        The connection is automatically returned to the pool when done.
        Use this as a context manager.
        
        Yields:
            Database connection
            
        Raises:
            Error: If acquiring connection fails
            
        Example:
            with client.connection() as conn:
                with conn.cursor() as cur:
                    cur.execute("SELECT * FROM users")
                    users = cur.fetchall()
        """
        if self._pool is None:
            raise connection_error("acquire connection", Exception("pool not open"))
        
        try:
            with self._pool.connection() as conn:
                yield conn
        except Exception as e:
            raise map_error(e, "acquire connection")
    
    @contextmanager
    def cursor(self, **kwargs) -> Generator[Cursor, None, None]:
        """
        Get a cursor from the pool (convenience method).
        
        This acquires a connection and creates a cursor in one step.
        
        Args:
            **kwargs: Additional cursor options (e.g., row_factory=dict_row)
            
        Yields:
            Database cursor
            
        Raises:
            Error: If acquiring connection or cursor fails
            
        Example:
            with client.cursor(row_factory=dict_row) as cur:
                cur.execute("SELECT * FROM users WHERE id = %s", (user_id,))
                user = cur.fetchone()
        """
        with self.connection() as conn:
            with conn.cursor(**kwargs) as cur:
                yield cur
    
    def execute(self, query: str, params: Optional[tuple] = None) -> Cursor:
        """
        Execute a query and return the cursor (convenience method).
        
        This is useful for simple queries without needing to manage
        connection/cursor context managers.
        
        Args:
            query: SQL query to execute
            params: Query parameters
            
        Returns:
            Cursor with query results
            
        Raises:
            Error: If query execution fails
            
        Example:
            cur = client.execute("SELECT * FROM users WHERE active = %s", (True,))
            users = cur.fetchall()
        """
        with self.connection() as conn:
            return conn.execute(query, params)
    
    @property
    def pool(self) -> ConnectionPool:
        """
        Get the underlying connection pool for advanced usage.
        
        Returns:
            The psycopg ConnectionPool
            
        Raises:
            Error: If pool is not open
        """
        if self._pool is None:
            raise connection_error("get pool", Exception("pool not open"))
        
        return self._pool


def create_client(config: Config) -> Client:
    """
    Create and open a PostgreSQL client.
    
    This is a convenience function that creates and opens a client in one step.
    
    Args:
        config: Database configuration
        
    Returns:
        Opened Client instance
        
    Raises:
        Error: If client creation or opening fails
        
    Example:
        config = from_service_config(...)
        client = create_client(config)
        # Client is ready to use
    """
    client = Client(config)
    client.open()
    return client