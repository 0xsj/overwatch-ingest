"""gRPC server implementation."""

from concurrent import futures
from typing import Any

import grpc

from scout_common.observability.logger import Logger
from .config import ServerConfig
from ..interceptors.logging import LoggingInterceptor
from ..interceptors.recovery import RecoveryInterceptor


class GrpcServer:
    """gRPC server implementation using grpc.aio."""
    
    def __init__(self, config: ServerConfig) -> None:
        """
        Initialize gRPC server.
        
        Args:
            config: Server configuration
            
        Raises:
            ValueError: If logger is not provided
        """
        if config.logger is None:
            raise ValueError("logger is required")
        
        self._config = config
        self._logger = config.logger
        
        # Create server with interceptors
        interceptors = [
            LoggingInterceptor(self._logger),
            RecoveryInterceptor(self._logger),
        ]
        
        self._server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=config.max_workers),
            interceptors=interceptors,
            options=[
                ("grpc.max_connection_idle_ms", config.max_connection_idle_sec * 1000),
                ("grpc.max_connection_age_ms", config.max_connection_age_sec * 1000),
                ("grpc.max_connection_age_grace_ms", config.max_connection_age_grace_sec * 1000),
                ("grpc.keepalive_time_ms", config.keepalive_time_sec * 1000),
                ("grpc.keepalive_timeout_ms", config.keepalive_timeout_sec * 1000),
                ("grpc.http2.max_pings_without_data", 0),
            ],
        )
    
    def start(self) -> None:
        """Start the gRPC server."""
        self._server.add_insecure_port(self._config.address)
        
        self._logger.info("grpc server starting", address=self._config.address)
        
        self._server.start()
        self._server.wait_for_termination()
    
    def stop(self, grace_period: float | None = None) -> None:
        """
        Gracefully stop the server.
        
        Args:
            grace_period: Seconds to wait for graceful shutdown
        """
        self._logger.info("grpc server stopping")
        
        if grace_period is not None:
            self._server.stop(grace_period)
        else:
            self._server.stop(None)
        
        self._logger.info("grpc server stopped")
    
    def add_service(self, servicer: Any, add_servicer_fn: Any) -> None:
        """
        Add a gRPC service implementation.
        
        Args:
            servicer: The service implementation
            add_servicer_fn: Function to add the servicer
        """
        add_servicer_fn(servicer, self._server)
        
        # Try to get service name for logging
        service_name = getattr(servicer, '__class__.__name__', 'unknown')
        self._logger.info("grpc service registered", service=service_name)
    
    def get_server(self) -> grpc.Server:
        """Get the underlying grpc.Server."""
        return self._server


def new(config: ServerConfig) -> GrpcServer:
    """
    Create a new gRPC server.
    
    Args:
        config: Server configuration
        
    Returns:
        New gRPC server instance
    """
    return GrpcServer(config)