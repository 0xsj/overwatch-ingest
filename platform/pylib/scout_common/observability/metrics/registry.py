# platform/pylib/scout_common/observability/metrics/registry.py
"""Metric registry for storing and managing metrics."""

import threading
from typing import Dict, List, Optional

from .counter import Counter
from .histogram import Histogram, DEFAULT_BUCKETS
from .gauge import Gauge


class Registry:
    """Registry stores and manages metrics.
    
    Example:
        registry = Registry()
        
        counter = registry.register_counter("requests_total", "Total requests")
        counter.inc()
        
        histogram = registry.register_histogram("request_duration", "Request latency")
        histogram.observe(0.042)
        
        gauge = registry.register_gauge("active_connections", "Active connections")
        gauge.set(10)
    """

    def __init__(self):
        """Initialize registry."""
        self._counters: Dict[str, Counter] = {}
        self._histograms: Dict[str, Histogram] = {}
        self._gauges: Dict[str, Gauge] = {}
        self._lock = threading.RLock()

    def register_counter(
        self,
        name: str,
        help: str,
        labels: Optional[Dict[str, str]] = None
    ) -> Counter:
        """Register a new counter.
        
        Args:
            name: Metric name
            help: Help description
            labels: Label key-value pairs
            
        Returns:
            Counter instance (existing or newly created)
        """
        key = self._metric_key(name, labels)
        
        with self._lock:
            if key in self._counters:
                return self._counters[key]
            
            counter = Counter(name, help, labels)
            self._counters[key] = counter
            return counter

    def register_histogram(
        self,
        name: str,
        help: str,
        labels: Optional[Dict[str, str]] = None,
        buckets: Optional[List[float]] = None
    ) -> Histogram:
        """Register a new histogram.
        
        Args:
            name: Metric name
            help: Help description
            labels: Label key-value pairs
            buckets: Histogram buckets (defaults to DEFAULT_BUCKETS)
            
        Returns:
            Histogram instance (existing or newly created)
        """
        key = self._metric_key(name, labels)
        
        with self._lock:
            if key in self._histograms:
                return self._histograms[key]
            
            histogram = Histogram(name, help, labels, buckets)
            self._histograms[key] = histogram
            return histogram

    def register_gauge(
        self,
        name: str,
        help: str,
        labels: Optional[Dict[str, str]] = None
    ) -> Gauge:
        """Register a new gauge.
        
        Args:
            name: Metric name
            help: Help description
            labels: Label key-value pairs
            
        Returns:
            Gauge instance (existing or newly created)
        """
        key = self._metric_key(name, labels)
        
        with self._lock:
            if key in self._gauges:
                return self._gauges[key]
            
            gauge = Gauge(name, help, labels)
            self._gauges[key] = gauge
            return gauge

    def get_counter(
        self,
        name: str,
        labels: Optional[Dict[str, str]] = None
    ) -> Optional[Counter]:
        """Get a counter by name and labels.
        
        Args:
            name: Metric name
            labels: Label key-value pairs
            
        Returns:
            Counter if found, None otherwise
        """
        key = self._metric_key(name, labels)
        with self._lock:
            return self._counters.get(key)

    def get_histogram(
        self,
        name: str,
        labels: Optional[Dict[str, str]] = None
    ) -> Optional[Histogram]:
        """Get a histogram by name and labels.
        
        Args:
            name: Metric name
            labels: Label key-value pairs
            
        Returns:
            Histogram if found, None otherwise
        """
        key = self._metric_key(name, labels)
        with self._lock:
            return self._histograms.get(key)

    def get_gauge(
        self,
        name: str,
        labels: Optional[Dict[str, str]] = None
    ) -> Optional[Gauge]:
        """Get a gauge by name and labels.
        
        Args:
            name: Metric name
            labels: Label key-value pairs
            
        Returns:
            Gauge if found, None otherwise
        """
        key = self._metric_key(name, labels)
        with self._lock:
            return self._gauges.get(key)

    @property
    def counters(self) -> List[Counter]:
        """Get all registered counters."""
        with self._lock:
            return list(self._counters.values())

    @property
    def histograms(self) -> List[Histogram]:
        """Get all registered histograms."""
        with self._lock:
            return list(self._histograms.values())

    @property
    def gauges(self) -> List[Gauge]:
        """Get all registered gauges."""
        with self._lock:
            return list(self._gauges.values())

    def _metric_key(self, name: str, labels: Optional[Dict[str, str]]) -> str:
        """Generate a unique key for a metric.
        
        Args:
            name: Metric name
            labels: Label key-value pairs
            
        Returns:
            Unique metric key
        """
        if not labels:
            return name
        
        # Sort labels for consistent keys
        sorted_labels = sorted(labels.items())
        label_str = ",".join(f"{k}={v}" for k, v in sorted_labels)
        return f"{name}{{{label_str}}}"


# Global default registry
default_registry = Registry()


# Convenience functions using default registry
def register_counter(
    name: str,
    help: str,
    labels: Optional[Dict[str, str]] = None
) -> Counter:
    """Register a counter in the default registry."""
    return default_registry.register_counter(name, help, labels)


def register_histogram(
    name: str,
    help: str,
    labels: Optional[Dict[str, str]] = None,
    buckets: Optional[List[float]] = None
) -> Histogram:
    """Register a histogram in the default registry."""
    return default_registry.register_histogram(name, help, labels, buckets)


def register_gauge(
    name: str,
    help: str,
    labels: Optional[Dict[str, str]] = None
) -> Gauge:
    """Register a gauge in the default registry."""
    return default_registry.register_gauge(name, help, labels)