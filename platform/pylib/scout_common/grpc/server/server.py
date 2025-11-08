"""gRPC server protocol for Scout platform."""

from typing import Protocol, Any, runtime_checkable
from concurrent import futures


@runtime_checkable
class Server(Protocol):
    """
    Server protocol defines the gRPC server interface.
    
    Implementations handle server lifecycle (start, stop, graceful shutdown).
    """
    
    def start(self) -> None:
        """
        Start the gRPC server on the configured address.
        
        This is a blocking call that returns when the server stops.
        
        Raises:
            Exception: If server fails to start
        """
        ...
    
    def stop(self, grace_period: float | None = None) -> None:
        """
        Gracefully stop the server.
        
        Stops accepting new connections and waits for existing RPCs to complete.
        
        Args:
            grace_period: Seconds to wait for graceful shutdown (None = wait forever)
        """
        ...
    
    def add_service(self, servicer: Any, add_servicer_fn: Any) -> None:
        """
        Add a gRPC service implementation.
        
        Must be called before start().
        
        Args:
            servicer: The service implementation
            add_servicer_fn: Function to add the servicer (from generated code)
        """
        ...
    
    def get_server(self) -> Any:
        """
        Get the underlying grpc.Server for advanced use cases.
        
        Returns:
            The grpc.Server instance
        """
        ...