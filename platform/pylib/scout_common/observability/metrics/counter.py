# platform/pylib/scout_common/observability/metrics/counter.py
"""Counter metrics for monotonically increasing values."""

import threading
from typing import Dict, Optional


class Counter:
    """Counter is a monotonically increasing metric.
    
    Used for counting events (requests, errors, messages published, etc.)
    
    Example:
        counter = Counter("http_requests_total", "Total HTTP requests")
        counter.inc()
        counter.add(5)
        print(counter.value)  # 6
    """

    def __init__(
        self,
        name: str,
        help: str,
        labels: Optional[Dict[str, str]] = None
    ):
        """Initialize counter.
        
        Args:
            name: Metric name
            help: Help description
            labels: Label key-value pairs
        """
        self._name = name
        self._help = help
        self._labels = labels or {}
        self._value = 0.0
        self._lock = threading.RLock()

    def inc(self) -> None:
        """Increment counter by 1."""
        self.add(1)

    def add(self, delta: float) -> None:
        """Add the given value to the counter.
        
        Args:
            delta: Value to add (must be non-negative)
        """
        if delta < 0:
            # Counters can only increase
            return
        
        with self._lock:
            self._value += delta

    @property
    def value(self) -> float:
        """Get current counter value."""
        with self._lock:
            return self._value

    @property
    def name(self) -> str:
        """Get counter name."""
        return self._name

    @property
    def help(self) -> str:
        """Get counter help text."""
        return self._help

    @property
    def labels(self) -> Dict[str, str]:
        """Get counter labels."""
        return self._labels.copy()

    def __repr__(self) -> str:
        """Get debug representation."""
        labels_str = ", ".join(f"{k}={v}" for k, v in self._labels.items())
        return f"Counter({self._name}{{{labels_str}}}={self._value})"