"""gRPC client implementation."""

import time
from typing import Any

import grpc

from scout_common.observability.logger import Logger
from .config import ClientConfig


class GrpcClient:
    """gRPC client implementation using grpc."""
    
    def __init__(self, config: ClientConfig) -> None:
        """
        Initialize gRPC client.
        
        Args:
            config: Client configuration
            
        Raises:
            ValueError: If logger is not provided
        """
        if config.logger is None:
            raise ValueError("logger is required")
        
        self._config = config
        self._logger = config.logger
        self._channel: grpc.Channel | None = None
    
    def connect(self) -> None:
        """Establish a connection to the gRPC server."""
        if self._channel is not None:
            raise RuntimeError("client already connected")
        
        self._logger.info("grpc client connecting", target=self._config.target)
        
        # Channel options
        options = [
            ("grpc.keepalive_time_ms", self._config.keepalive_time_sec * 1000),
            ("grpc.keepalive_timeout_ms", self._config.keepalive_timeout_sec * 1000),
            ("grpc.keepalive_permit_without_calls", 1),
            ("grpc.http2.max_pings_without_data", 0),
        ]
        
        # Create channel
        if self._config.insecure:
            self._channel = grpc.insecure_channel(self._config.target, options=options)
        else:
            # TODO: Add secure channel with credentials
            credentials = grpc.ssl_channel_credentials()
            self._channel = grpc.secure_channel(
                self._config.target,
                credentials,
                options=options,
            )
        
        # Wait for connection to be ready (with timeout)
        try:
            grpc.channel_ready_future(self._channel).result(
                timeout=self._config.connect_timeout_sec
            )
            self._logger.info("grpc client connected", target=self._config.target)
        except grpc.FutureTimeoutError:
            self._channel.close()
            self._channel = None
            raise RuntimeError(
                f"failed to connect to {self._config.target} within {self._config.connect_timeout_sec}s"
            )
    
    def close(self) -> None:
        """Close the connection."""
        if self._channel is None:
            return
        
        self._logger.info("grpc client closing", target=self._config.target)
        
        self._channel.close()
        self._channel = None
        
        self._logger.info("grpc client closed")
    
    def get_channel(self) -> grpc.Channel:
        """
        Get the underlying grpc.Channel.
        
        Returns:
            The grpc.Channel instance
            
        Raises:
            RuntimeError: If client is not connected
        """
        if self._channel is None:
            raise RuntimeError("client not connected")
        
        return self._channel
    
    def is_connected(self) -> bool:
        """Check if the client is connected."""
        if self._channel is None:
            return False
        
        try:
            state = self._channel._channel.check_connectivity_state(False)
            return state == grpc.ChannelConnectivity.READY
        except Exception:
            return False


def new(config: ClientConfig) -> GrpcClient:
    """
    Create a new gRPC client.
    
    Args:
        config: Client configuration
        
    Returns:
        New gRPC client instance
    """
    return GrpcClient(config)