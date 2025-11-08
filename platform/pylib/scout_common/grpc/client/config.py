"""gRPC client configuration."""

from dataclasses import dataclass

from scout_common.observability.logger import Logger


@dataclass
class ClientConfig:
    """Configuration for gRPC client."""
    
    # Target address (e.g., "localhost:50051", "agents:50051")
    target: str = "localhost:50051"
    
    # Logger for client operations
    logger: Logger | None = None
    
    # Connection timeout in seconds
    connect_timeout_sec: int = 10
    
    # Maximum number of retry attempts
    max_retries: int = 3
    
    # Use insecure connection (for development only)
    insecure: bool = True
    
    # Keepalive time
    keepalive_time_sec: int = 30
    
    # Keepalive timeout
    keepalive_timeout_sec: int = 10
    
    def with_target(self, target: str) -> "ClientConfig":
        """Set the target address."""
        self.target = target
        return self
    
    def with_logger(self, logger: Logger) -> "ClientConfig":
        """Set the logger."""
        self.logger = logger
        return self
    
    def with_insecure(self, insecure: bool) -> "ClientConfig":
        """Set whether to use insecure connections."""
        self.insecure = insecure
        return self


def default_config() -> ClientConfig:
    """Return a ClientConfig with sensible defaults."""
    return ClientConfig()