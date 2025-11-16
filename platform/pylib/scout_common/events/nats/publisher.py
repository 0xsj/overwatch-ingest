# platform/pylib/scout_common/events/nats/publisher.py
"""NATS publisher implementation."""

import json
from typing import Any

from nats.aio.client import Client as NATS

from ..config import Config


class NATSPublisher:
    """NATS implementation of Publisher."""

    def __init__(self, config: Config):
        """Initialize NATS publisher.
        
        Args:
            config: NATS configuration
        """
        self._config = config
        self._client: NATS = NATS()
        self._connected = False

    async def connect(self) -> None:
        """Connect to NATS server."""
        if not self._connected:
            await self._client.connect(
                servers=[self._config.url],
                max_reconnect_attempts=self._config.max_reconnects,
                reconnect_time_wait=self._config.reconnect_wait,
            )
            self._connected = True

    async def publish(self, subject: str, data: bytes) -> None:
        """Publish raw bytes to a subject.
        
        Args:
            subject: The subject to publish to
            data: Raw message bytes
            
        Raises:
            Exception: If not connected or publishing fails
        """
        if not self._connected:
            await self.connect()
        
        await self._client.publish(subject, data)

    async def publish_json(self, subject: str, event: Any) -> None:
        """Publish a JSON-serialized event to a subject.
        
        Args:
            subject: The subject to publish to
            event: Event object to serialize and publish
            
        Raises:
            Exception: If serialization or publishing fails
        """
        # Handle pydantic models
        if hasattr(event, "model_dump"):
            # Pydantic v2
            data = json.dumps(event.model_dump()).encode()
        elif hasattr(event, "dict"):
            # Pydantic v1
            data = json.dumps(event.dict()).encode()
        else:
            # Regular dict or dataclass
            data = json.dumps(event).encode()
        
        await self.publish(subject, data)

    async def close(self) -> None:
        """Close the NATS connection."""
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