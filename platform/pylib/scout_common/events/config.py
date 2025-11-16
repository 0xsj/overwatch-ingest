# platform/pylib/scout_common/events/config.py
"""Config protocol for event bus providers."""

from typing import Protocol
from abc import abstractmethod


class Config(Protocol):
    """Common configuration for event bus providers."""

    @property
    @abstractmethod
    def url(self) -> str:
        """Connection URL."""
        ...

    @property
    @abstractmethod
    def max_reconnects(self) -> int:
        """Maximum number of reconnection attempts."""
        ...

    @property
    @abstractmethod
    def reconnect_wait(self) -> float:
        """Wait time between reconnection attempts (seconds)."""
        ...