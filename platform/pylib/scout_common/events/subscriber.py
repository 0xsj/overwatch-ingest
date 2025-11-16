# platform/pylib/scout_common/events/subscriber.py
"""Subscriber protocol for event bus."""

from typing import Protocol, Callable, Awaitable
from abc import abstractmethod

from .message import Message


# Type alias for message handler
MessageHandler = Callable[[Message], Awaitable[None]]


class Subscription(Protocol):
    """Represents an active subscription."""

    @abstractmethod
    async def unsubscribe(self) -> None:
        """Stop receiving messages."""
        ...

    @property
    @abstractmethod
    def subject(self) -> str:
        """Get the subject/topic being subscribed to."""
        ...


class Subscriber(Protocol):
    """Subscriber subscribes to events from an event bus.
    
    Implementations: NATS, Kafka, RabbitMQ
    """

    @abstractmethod
    async def subscribe(
        self, subject: str, handler: MessageHandler
    ) -> Subscription:
        """Subscribe to a subject/topic with a handler.
        
        Args:
            subject: The subject/topic to subscribe to
            handler: Async function to handle incoming messages
            
        Returns:
            Subscription that can be used to unsubscribe
            
        Raises:
            Exception: If subscription fails
        """
        ...

    @abstractmethod
    async def subscribe_queue(
        self, subject: str, queue: str, handler: MessageHandler
    ) -> Subscription:
        """Subscribe to a subject/topic as part of a queue group.
        
        Multiple subscribers with the same queue group will load-balance messages.
        
        Args:
            subject: The subject/topic to subscribe to
            queue: Queue group name for load balancing
            handler: Async function to handle incoming messages
            
        Returns:
            Subscription that can be used to unsubscribe
            
        Raises:
            Exception: If subscription fails
        """
        ...

    @abstractmethod
    async def close(self) -> None:
        """Close all subscriptions and the subscriber connection."""
        ...

    async def __aenter__(self):
        """Async context manager entry."""
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        await self.close()