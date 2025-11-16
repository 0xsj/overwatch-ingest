# platform/pylib/scout_common/events/publisher.py
"""Publisher protocol for event bus."""

from typing import Protocol, Any
from abc import abstractmethod


class Publisher(Protocol):
    """Publisher publishes events to an event bus.
    
    Implementations: NATS, Kafka, RabbitMQ
    """

    @abstractmethod
    async def publish(self, subject: str, data: bytes) -> None:
        """Publish raw bytes to a subject/topic.
        
        Args:
            subject: The subject/topic to publish to
            data: Raw message bytes
            
        Raises:
            Exception: If publishing fails
        """
        ...

    @abstractmethod
    async def publish_json(self, subject: str, event: Any) -> None:
        """Publish a JSON-serialized event to a subject/topic.
        
        Args:
            subject: The subject/topic to publish to
            event: Event object to serialize and publish
            
        Raises:
            Exception: If serialization or publishing fails
        """
        ...

    @abstractmethod
    async def close(self) -> None:
        """Close the publisher connection."""
        ...

    async def __aenter__(self):
        """Async context manager entry."""
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        await self.close()