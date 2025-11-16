# platform/pylib/scout_common/events/nats/subscriber.py
"""NATS subscriber implementation."""

import logging
from typing import List

from nats.aio.client import Client as NATS
from nats.aio.subscription import Subscription as NATSSubscription

from ..config import Config
from ..message import Message
from ..subscriber import MessageHandler


logger = logging.getLogger(__name__)


class NATSSubscriber:
    """NATS implementation of Subscriber."""

    def __init__(self, config: Config):
        """Initialize NATS subscriber.
        
        Args:
            config: NATS configuration
        """
        self._config = config
        self._client: NATS = NATS()
        self._connected = False
        self._subscriptions: List[NATSSubscription] = []

    async def connect(self) -> None:
        """Connect to NATS server."""
        if not self._connected:
            await self._client.connect(
                servers=[self._config.url],
                max_reconnect_attempts=self._config.max_reconnects,
                reconnect_time_wait=self._config.reconnect_wait,
            )
            self._connected = True

    async def subscribe(self, subject: str, handler: MessageHandler) -> "Subscription":
        """Subscribe to a subject.
        
        Args:
            subject: The subject to subscribe to
            handler: Async function to handle incoming messages
            
        Returns:
            Subscription that can be used to unsubscribe
        """
        if not self._connected:
            await self.connect()

        async def message_handler(msg):
            event_msg = Message(
                subject=msg.subject,
                data=msg.data,
            )
            try:
                await handler(event_msg)
            except Exception as e:
                logger.error(f"Error handling message on subject {subject}: {e}")

        sub = await self._client.subscribe(subject, cb=message_handler)
        self._subscriptions.append(sub)
        return Subscription(sub)

    async def subscribe_queue(
        self, subject: str, queue: str, handler: MessageHandler
    ) -> "Subscription":
        """Subscribe to a subject with queue group (load balancing).
        
        Args:
            subject: The subject to subscribe to
            queue: Queue group name
            handler: Async function to handle incoming messages
            
        Returns:
            Subscription that can be used to unsubscribe
        """
        if not self._connected:
            await self.connect()

        async def message_handler(msg):
            event_msg = Message(
                subject=msg.subject,
                data=msg.data,
            )
            try:
                await handler(event_msg)
            except Exception as e:
                logger.error(f"Error handling message on subject {subject}: {e}")

        sub = await self._client.subscribe(subject, queue=queue, cb=message_handler)
        self._subscriptions.append(sub)
        return Subscription(sub)

    async def close(self) -> None:
        """Close all subscriptions and the subscriber connection."""
        for sub in self._subscriptions:
            await sub.unsubscribe()
        self._subscriptions.clear()
        
        if self._connected:
            await self._client.close()
            self._connected = False

    async def __aenter__(self):
        """Async context manager entry."""
        await self.connect()
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit."""
        await self.close()


class Subscription:
    """Wrapper for NATS subscription."""

    def __init__(self, sub: NATSSubscription):
        self._sub = sub

    async def unsubscribe(self) -> None:
        """Stop receiving messages."""
        await self._sub.unsubscribe()

    @property
    def subject(self) -> str:
        """Get the subject being subscribed to."""
        return self._sub.subject