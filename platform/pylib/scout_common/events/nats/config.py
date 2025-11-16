# platform/pylib/scout_common/events/nats/config.py
"""NATS configuration."""

from dataclasses import dataclass


@dataclass
class NATSConfig:
    """NATS configuration implementing events.Config protocol."""

    url: str
    """NATS server URL (e.g., nats://localhost:4224)."""

    max_reconnects: int = 10
    """Maximum number of reconnection attempts."""

    reconnect_wait: float = 2.0
    """Wait time between reconnection attempts (seconds)."""