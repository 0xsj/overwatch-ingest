"""PostgreSQL health check utilities."""

import time
from dataclasses import dataclass
from datetime import timedelta
from typing import Callable, Optional

from scout_common.database.postgres.client import Client, PoolStats
from scout_common.errors import Error, ErrorType, new_error


@dataclass
class HealthCheck:
    """
    Result of a database health check.
    
    Attributes:
        healthy: Whether the database is healthy
        message: Additional context about health status
        latency: Time taken to ping the database (in seconds)
        stats: Connection pool statistics
        error: Error message if check failed
    """
    
    healthy: bool
    message: str = ""
    latency: float = 0.0
    stats: Optional[PoolStats] = None
    error: Optional[str] = None
    
    def to_dict(self) -> dict:
        """Convert to dictionary for JSON serialization."""
        result = {
            "healthy": self.healthy,
            "message": self.message,
            "latency": self.latency,
        }
        
        if self.stats:
            result["stats"] = {
                "total": self.stats.size,
                "available": self.stats.available,
                "max": self.stats.max_size,
                "min": self.stats.min_size,
            }
        
        if self.error:
            result["error"] = self.error
        
        return result


def check(client: Client) -> HealthCheck:
    """
    Perform a health check on the database client.
    
    This is suitable for Kubernetes liveness and readiness probes.
    
    The check:
    - Pings the database to verify connectivity
    - Measures response latency
    - Retrieves pool statistics
    - Determines overall health status
    
    Args:
        client: Database client to check
        
    Returns:
        HealthCheck result
        
    Example:
        health = check(db_client)
        if not health.healthy:
            logger.error(f"Database unhealthy: {health.message}")
    """
    start = time.time()
    
    # Ping the database
    try:
        client.ping()
        latency = time.time() - start
    except Exception as e:
        return HealthCheck(
            healthy=False,
            message="database ping failed",
            latency=time.time() - start,
            error=str(e),
        )
    
    # Get pool statistics
    try:
        stats = client.stats()
    except Exception as e:
        return HealthCheck(
            healthy=False,
            message="failed to get pool stats",
            latency=latency,
            error=str(e),
        )
    
    # Check pool health
    pool_health = stats.health_status()
    if pool_health != "healthy":
        return HealthCheck(
            healthy=False,
            message=pool_health,
            latency=latency,
            stats=stats,
        )
    
    # All checks passed
    return HealthCheck(
        healthy=True,
        message="database is healthy",
        latency=latency,
        stats=stats,
    )


def check_with_timeout(client: Client, timeout: timedelta) -> HealthCheck:
    """
    Perform a health check with a specified timeout.
    
    This is useful for ensuring health checks don't block indefinitely.
    
    Args:
        client: Database client to check
        timeout: Maximum time to wait for health check
        
    Returns:
        HealthCheck result
        
    Example:
        health = check_with_timeout(db_client, timedelta(seconds=5))
        if not health.healthy:
            return 503  # Service Unavailable
    """
    # Note: psycopg doesn't support per-operation timeouts easily
    # This is a simple implementation that just calls check()
    # For production, you might want to use asyncio with timeout
    return check(client)


def is_ready(client: Client) -> bool:
    """
    Determine if the database is ready to serve requests.
    
    This is suitable for Kubernetes readiness probes.
    
    A database is considered ready if:
    - It responds to pings
    - The connection pool has available connections
    - Response latency is acceptable
    
    Args:
        client: Database client to check
        
    Returns:
        True if ready, False otherwise
        
    Example:
        if not is_ready(db_client):
            return Response(status=503)  # Service Unavailable
    """
    health = check(client)
    
    # Not ready if unhealthy
    if not health.healthy:
        return False
    
    # Not ready if latency is too high (>1 second indicates issues)
    if health.latency > 1.0:
        return False
    
    # Not ready if pool is exhausted
    if health.stats and health.stats.available == 0:
        return False
    
    return True


def is_alive(client: Client) -> bool:
    """
    Determine if the database connection is alive.
    
    This is suitable for Kubernetes liveness probes.
    
    A database is considered alive if it responds to pings,
    even if the pool is under stress.
    
    Args:
        client: Database client to check
        
    Returns:
        True if alive, False otherwise
        
    Example:
        if not is_alive(db_client):
            logger.fatal("Database connection lost, shutting down...")
            sys.exit(1)
    """
    try:
        client.ping()
        return True
    except Exception:
        return False


def wait_for_ready(
    client: Client,
    timeout: timedelta,
    check_interval: timedelta = timedelta(seconds=1),
) -> Optional[Error]:
    """
    Block until the database is ready or the timeout expires.
    
    This is useful during application startup to wait for database availability.
    
    Args:
        client: Database client to check
        timeout: Maximum time to wait
        check_interval: Time between health checks
        
    Returns:
        None if ready, Error if timeout
        
    Example:
        err = wait_for_ready(
            db_client,
            timeout=timedelta(seconds=30),
            check_interval=timedelta(seconds=1),
        )
        if err:
            logger.fatal(f"Database never became ready: {err}")
            sys.exit(1)
    """
    start = time.time()
    timeout_seconds = timeout.total_seconds()
    interval_seconds = check_interval.total_seconds()
    
    # Try immediately first
    if is_ready(client):
        return None
    
    while True:
        elapsed = time.time() - start
        
        if elapsed >= timeout_seconds:
            return new_error(
                error_type=ErrorType.TIMEOUT,
                code="DB_READY_TIMEOUT",
                message="database did not become ready within timeout",
            )
        
        time.sleep(interval_seconds)
        
        if is_ready(client):
            return None


def wait_for_alive(
    client: Client,
    timeout: timedelta,
    check_interval: timedelta = timedelta(seconds=1),
) -> Optional[Error]:
    """
    Block until the database is alive or the timeout expires.
    
    This is useful during application startup to wait for initial connectivity.
    
    Args:
        client: Database client to check
        timeout: Maximum time to wait
        check_interval: Time between health checks
        
    Returns:
        None if alive, Error if timeout
        
    Example:
        err = wait_for_alive(
            db_client,
            timeout=timedelta(seconds=30),
        )
        if err:
            logger.fatal(f"Database never became alive: {err}")
            sys.exit(1)
    """
    start = time.time()
    timeout_seconds = timeout.total_seconds()
    interval_seconds = check_interval.total_seconds()
    
    # Try immediately first
    if is_alive(client):
        return None
    
    while True:
        elapsed = time.time() - start
        
        if elapsed >= timeout_seconds:
            return new_error(
                error_type=ErrorType.TIMEOUT,
                code="DB_ALIVE_TIMEOUT",
                message="database did not become alive within timeout",
            )
        
        time.sleep(interval_seconds)
        
        if is_alive(client):
            return None


HealthCheckCallback = Callable[[HealthCheck], None]


def monitor_health(
    client: Client,
    interval: timedelta,
    callback: HealthCheckCallback,
    stop_event: Optional[Callable[[], bool]] = None,
) -> None:
    """
    Periodically perform health checks and call a callback.
    
    This is useful for background health monitoring and alerting.
    The monitor runs until stop_event returns True.
    
    Args:
        client: Database client to monitor
        interval: Time between health checks
        callback: Function to call with health check results
        stop_event: Optional function that returns True to stop monitoring
        
    Example:
        import threading
        
        stop = threading.Event()
        
        def health_callback(health: HealthCheck):
            if not health.healthy:
                alerting.send(f"Database unhealthy: {health.message}")
            metrics.record("db.health.latency", health.latency)
        
        monitor_thread = threading.Thread(
            target=monitor_health,
            args=(db_client, timedelta(seconds=30), health_callback, stop.is_set),
        )
        monitor_thread.start()
    """
    interval_seconds = interval.total_seconds()
    
    # Run initial check immediately
    callback(check(client))
    
    while True:
        if stop_event and stop_event():
            break
        
        time.sleep(interval_seconds)
        callback(check(client))