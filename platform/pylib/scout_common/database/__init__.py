"""Scout database utilities."""

__version__ = "0.1.0"

# Re-export postgres module for convenience
from scout_common.database import postgres

__all__ = [
    "postgres",
]