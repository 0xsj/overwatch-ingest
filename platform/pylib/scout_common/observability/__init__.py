"""Observability utilities for Scout platform."""

__version__ = "0.1.0"

# Re-export logger for convenience
from . import logger

__all__ = [
    "logger",
    "__version__",
]