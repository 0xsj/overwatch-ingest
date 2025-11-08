"""Console logger with colored output for local development."""

from typing import Any
from datetime import datetime
import sys


# ANSI color codes
class Colors:
    RESET = "\033[0m"
    RED = "\033[31m"
    GREEN = "\033[32m"
    YELLOW = "\033[33m"
    BLUE = "\033[34m"
    PURPLE = "\033[35m"
    CYAN = "\033[36m"
    GRAY = "\033[37m"
    WHITE = "\033[97m"


class NoopLogger:
    """
    Simple console logger with colored output.
    
    Great for local development - logs are human-readable with colors.
    Not a true "no-op" but a simple stdout logger without structured backends.
    """
    
    def __init__(self, **fields: Any) -> None:
        """
        Initialize console logger.
        
        Args:
            **fields: Persistent fields to attach to all log messages
        """
        self._fields = fields
    
    def debug(self, msg: str, **kwargs: Any) -> None:
        """Log debug message with cyan color."""
        self._log("DEBUG", Colors.CYAN, msg, **kwargs)
    
    def info(self, msg: str, **kwargs: Any) -> None:
        """Log info message with green color."""
        self._log("INFO ", Colors.GREEN, msg, **kwargs)
    
    def warn(self, msg: str, **kwargs: Any) -> None:
        """Log warning message with yellow color."""
        self._log("WARN ", Colors.YELLOW, msg, **kwargs)
    
    def error(self, msg: str, **kwargs: Any) -> None:
        """Log error message with red color."""
        self._log("ERROR", Colors.RED, msg, **kwargs)
    
    def with_fields(self, **kwargs: Any) -> "NoopLogger":
        """
        Return a new logger with fields attached.
        
        Args:
            **kwargs: Fields to attach
            
        Returns:
            New logger instance with combined fields
        """
        new_fields = {**self._fields, **kwargs}
        return NoopLogger(**new_fields)
    
    def with_error(self, err: Exception) -> "NoopLogger":
        """
        Return a new logger with error attached.
        
        Args:
            err: Exception to attach
            
        Returns:
            New logger instance with error fields
        """
        return self.with_fields(
            error=str(err),
            error_type=type(err).__name__,
        )
    
    def _log(self, level: str, color: str, msg: str, **kwargs: Any) -> None:
        """
        Output a colored log message to stdout.
        
        Args:
            level: Log level string
            color: ANSI color code
            msg: Log message
            **kwargs: Additional fields
        """
        # Timestamp
        timestamp = datetime.now().strftime("%H:%M:%S.%f")[:-3]
        
        # Combine persistent fields with new fields
        all_fields = {**self._fields, **kwargs}
        
        # Format fields
        fields_str = self._format_fields(all_fields)
        
        # Build output
        parts = [
            f"{Colors.GRAY}{timestamp}{Colors.RESET}",
            f"{color}{level}{Colors.RESET}",
            f"{Colors.WHITE}{msg}{Colors.RESET}",
        ]
        
        if fields_str:
            parts.append(fields_str)
        
        # Print to stdout
        print(" ".join(parts), file=sys.stdout)
        sys.stdout.flush()
    
    def _format_fields(self, fields: dict[str, Any]) -> str:
        """
        Format key-value pairs with colors.
        
        Args:
            fields: Dictionary of fields
            
        Returns:
            Formatted string with colored keys
        """
        if not fields:
            return ""
        
        parts = []
        for key, value in fields.items():
            key_str = f"{Colors.PURPLE}{key}{Colors.RESET}"
            parts.append(f"{key_str}={value}")
        
        return " ".join(parts)