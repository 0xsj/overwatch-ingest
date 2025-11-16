# platform/pylib/scout_common/observability/metrics/gauge.py
"""Gauge metrics for values that can go up and down."""

import threading
from typing import Dict, Optional


class Gauge:
    """Gauge is a metric that can go up and down.
    
    Used for tracking memory usage, queue size, active connections, etc.
    
    Example:
        gauge = Gauge("active_connections", "Number of active connections")
        gauge.inc()        # 1
        gauge.add(5)       # 6
        gauge.dec()        # 5
        gauge.set(10)      # 10
        print(gauge.value) # 10
    """

    def __init__(
        self,
        name: str,
        help: str,
        labels: Optional[Dict[str, str]] = None
    ):
        """Initialize gauge.
        
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

    def set(self, value: float) -> None:
        """Set the gauge to the given value.
        
        Args:
            value: Value to set
        """
        with self._lock:
            self._value = value

    def inc(self) -> None:
        """Increment gauge by 1."""
        self.add(1)

    def dec(self) -> None:
        """Decrement gauge by 1."""
        self.sub(1)

    def add(self, delta: float) -> None:
        """Add the given value to the gauge.
        
        Args:
            delta: Value to add
        """
        with self._lock:
            self._value += delta

    def sub(self, delta: float) -> None:
        """Subtract the given value from the gauge.
        
        Args:
            delta: Value to subtract
        """
        self.add(-delta)

    @property
    def value(self) -> float:
        """Get current gauge value."""
        with self._lock:
            return self._value

    @property
    def name(self) -> str:
        """Get gauge name."""
        return self._name

    @property
    def help(self) -> str:
        """Get gauge help text."""
        return self._help

    @property
    def labels(self) -> Dict[str, str]:
        """Get gauge labels."""
        return self._labels.copy()

    def __repr__(self) -> str:
        """Get debug representation."""
        labels_str = ", ".join(f"{k}={v}" for k, v in self._labels.items())
        return f"Gauge({self._name}{{{labels_str}}}={self._value})"