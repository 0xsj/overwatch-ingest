# platform/pylib/scout_common/observability/metrics/histogram.py
"""Histogram metrics for tracking distributions."""

import math
import threading
from typing import Dict, List, Optional


# Default histogram buckets (in seconds)
# Suitable for HTTP request latencies: 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s
DEFAULT_BUCKETS = [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]


class Histogram:
    """Histogram tracks the distribution of values.
    
    Used for tracking latencies, request sizes, etc.
    
    Example:
        histogram = Histogram("http_request_duration_seconds", "HTTP request latency")
        histogram.observe(0.042)  # 42ms
        histogram.observe(0.156)  # 156ms
        print(f"Count: {histogram.count}, Sum: {histogram.sum}")
        print(f"Buckets: {histogram.buckets}")
    """

    def __init__(
        self,
        name: str,
        help: str,
        labels: Optional[Dict[str, str]] = None,
        buckets: Optional[List[float]] = None
    ):
        """Initialize histogram.
        
        Args:
            name: Metric name
            help: Help description
            labels: Label key-value pairs
            buckets: Histogram buckets (defaults to DEFAULT_BUCKETS)
        """
        self._name = name
        self._help = help
        self._labels = labels or {}
        
        # Sort buckets and add +Inf
        if buckets is None:
            buckets = DEFAULT_BUCKETS.copy()
        self._buckets = sorted(buckets) + [math.inf]
        self._counts = [0] * len(self._buckets)
        
        self._sum = 0.0
        self._count = 0
        self._lock = threading.RLock()

    def observe(self, value: float) -> None:
        """Record a value.
        
        Args:
            value: Value to observe
        """
        with self._lock:
            self._sum += value
            self._count += 1
            
            # Find bucket and increment
            for i, bucket in enumerate(self._buckets):
                if value <= bucket:
                    self._counts[i] += 1

    @property
    def count(self) -> int:
        """Get total number of observations."""
        with self._lock:
            return self._count

    @property
    def sum(self) -> float:
        """Get sum of all observed values."""
        with self._lock:
            return self._sum

    @property
    def buckets(self) -> Dict[float, int]:
        """Get histogram buckets with counts."""
        with self._lock:
            return dict(zip(self._buckets, self._counts))

    @property
    def name(self) -> str:
        """Get histogram name."""
        return self._name

    @property
    def help(self) -> str:
        """Get histogram help text."""
        return self._help

    @property
    def labels(self) -> Dict[str, str]:
        """Get histogram labels."""
        return self._labels.copy()

    def __repr__(self) -> str:
        """Get debug representation."""
        labels_str = ", ".join(f"{k}={v}" for k, v in self._labels.items())
        return f"Histogram({self._name}{{{labels_str}}}, count={self._count}, sum={self._sum})"