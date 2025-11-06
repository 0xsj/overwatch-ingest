"""Scout Common - Shared Python libraries for Scout platform."""

__version__ = "0.1.0"
__author__ = "SJ Lee"

# Re-export errors subpackage for convenience
from . import errors

__all__ = [
    "errors",
    "__version__",
]