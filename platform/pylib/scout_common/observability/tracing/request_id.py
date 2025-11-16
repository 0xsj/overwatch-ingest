# platform/pylib/scout_common/tracing/request_id.py
"""Request ID generation and validation."""

import time
import uuid
from typing import Optional


REQUEST_ID_PREFIX = "req_"
REQUEST_ID_HEADER = "X-Request-ID"


class RequestID:
    """Represents a unique request identifier.
    
    Format: req_<timestamp>_<uuid>
    Example: req_1699564800_a1b2c3d4
    """

    def __init__(self, value: str):
        """Initialize request ID.
        
        Args:
            value: Request ID string
            
        Raises:
            ValueError: If request ID is invalid
        """
        if not value:
            raise ValueError("Request ID cannot be empty")
        if not value.startswith(REQUEST_ID_PREFIX):
            raise ValueError(f"Request ID must start with '{REQUEST_ID_PREFIX}'")
        self._value = value

    @classmethod
    def generate(cls) -> "RequestID":
        """Generate a new request ID.
        
        Returns:
            New RequestID instance
        """
        timestamp = int(time.time())
        # Take first 8 chars of UUID for brevity
        short_id = uuid.uuid4().hex[:8]
        value = f"{REQUEST_ID_PREFIX}{timestamp}_{short_id}"
        return cls(value)

    @classmethod
    def parse(cls, value: str) -> Optional["RequestID"]:
        """Parse a string into a RequestID.
        
        Args:
            value: String to parse
            
        Returns:
            RequestID if valid, None otherwise
        """
        try:
            return cls(value)
        except ValueError:
            return None

    def __str__(self) -> str:
        """Get string representation."""
        return self._value

    def __repr__(self) -> str:
        """Get debug representation."""
        return f"RequestID('{self._value}')"

    def __eq__(self, other) -> bool:
        """Check equality."""
        if isinstance(other, RequestID):
            return self._value == other._value
        return self._value == str(other)

    def __hash__(self) -> int:
        """Get hash."""
        return hash(self._value)

    @property
    def value(self) -> str:
        """Get the request ID value."""
        return self._value