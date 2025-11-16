# platform/pylib/scout_common/observability/metrics/__init__.py
"""Prometheus-style metrics for observability.

Provides Counter, Histogram, and Gauge metrics with a central registry.

Example usage:
    from scout_common.observability.metrics import (
        register_counter,
        register_histogram,
        register_gauge,
    )
    
    # Counter - monotonically increasing
    requests_counter = register_counter(
        "http_requests_total",
        "Total HTTP requests",
        labels={"method": "GET", "endpoint": "/api/agents"}
    )
    requests_counter.inc()
    
    # Histogram - track distributions (latencies, sizes)
    duration_histogram = register_histogram(
        "http_request_duration_seconds",
        "HTTP request latency in seconds"
    )
    duration_histogram.observe(0.042)  # 42ms
    
    # Gauge - values that go up and down
    connections_gauge = register_gauge(
        "active_connections",
        "Number of active connections"
    )
    connections_gauge.set(10)
    connections_gauge.inc()  # 11
    connections_gauge.dec()  # 10
    
    # Custom registry
    from scout_common.observability.metrics import Registry
    
    custom_registry = Registry()
    counter = custom_registry.register_counter("my_metric", "My custom metric")
"""

from .counter import Counter
from .histogram import Histogram, DEFAULT_BUCKETS
from .gauge import Gauge
from .registry import (
    Registry,
    default_registry,
    register_counter,
    register_histogram,
    register_gauge,
)

__all__ = [
    # Metric types
    "Counter",
    "Histogram",
    "Gauge",
    # Constants
    "DEFAULT_BUCKETS",
    # Registry
    "Registry",
    "default_registry",
    # Convenience functions
    "register_counter",
    "register_histogram",
    "register_gauge",
]