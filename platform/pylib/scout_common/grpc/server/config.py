"""gRPC server configuration."""

from dataclasses import dataclass

from scout_common.observability.logger import Logger


@dataclass
class ServerConfig:
    """Configuration for gRPC server."""
    
    # Address to listen on (e.g., "[::]:50051", "0.0.0.0:50051")
    address: str = "[::]:50051"
    
    # Logger for server operations
    logger: Logger | None = None
    
    # Maximum number of concurrent workers
    max_workers: int = 10
    
    # Maximum number of concurrent connections
    max_connection_idle_sec: int = 900  # 15 minutes
    
    # Maximum time a connection can exist
    max_connection_age_sec: int = 1800  # 30 minutes
    
    # Grace period for connection closure
    max_connection_age_grace_sec: int = 5
    
    # Keepalive time
    keepalive_time_sec: int = 300  # 5 minutes
    
    # Keepalive timeout
    keepalive_timeout_sec: int = 20
    
    def with_address(self, address: str) -> "ServerConfig":
        """Set the server address."""
        self.address = address
        return self
    
    def with_logger(self, logger: Logger) -> "ServerConfig":
        """Set the logger."""
        self.logger = logger
        return self
    
    def with_max_workers(self, max_workers: int) -> "ServerConfig":
        """Set the maximum number of workers."""
        self.max_workers = max_workers
        return self


def default_config() -> ServerConfig:
    """Return a ServerConfig with sensible defaults."""
    return ServerConfig()