"""gRPC client protocol for Scout platform."""

from typing import Protocol, runtime_checkable, Any


@runtime_checkable
class Client(Protocol):
    """
    Client protocol defines the gRPC client interface.
    
    Implementations manage connection lifecycle and provide access to the connection.
    """
    
    def connect(self) -> None:
        """
        Establish a connection to the gRPC server.
        
        Raises:
            Exception: If connection fails
        """
        ...
    
    def close(self) -> None:
        """
        Close the connection gracefully.
        """
        ...
    
    def get_channel(self) -> Any:
        """
        Get the underlying grpc.Channel.
        
        Returns:
            The grpc.Channel instance
        """
        ...
    
    def is_connected(self) -> bool:
        """
        Check if the client is connected.
        
        Returns:
            True if connected, False otherwise
        """
        ...