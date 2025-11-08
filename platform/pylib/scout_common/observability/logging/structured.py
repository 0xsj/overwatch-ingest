"""Structured logger implementation using structlog."""

from typing import Any
import structlog
from structlog.typing import EventDict, Processor
import logging
import sys

from .logger import Level


# Custom processor to add level as string
def add_log_level(logger: Any, method_name: str, event_dict: EventDict) -> EventDict:
    """Add log level as string to event dict."""
    if method_name == "debug":
        event_dict["level"] = "debug"
    elif method_name == "info":
        event_dict["level"] = "info"
    elif method_name == "warning":
        event_dict["level"] = "warn"
    elif method_name == "error":
        event_dict["level"] = "error"
    return event_dict


def configure_structlog(level: Level = Level.INFO, development: bool = False) -> None:
    """
    Configure structlog for the application.
    
    Should be called once at application startup.
    
    Args:
        level: Minimum log level
        development: If True, use console renderer. If False, use JSON.
    """
    # Configure standard logging
    logging.basicConfig(
        format="%(message)s",
        stream=sys.stdout,
        level=level.value,
    )
    
    # Choose processors based on environment
    processors: list[Processor] = [
        structlog.contextvars.merge_contextvars,
        add_log_level,
        structlog.processors.add_log_level,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.StackInfoRenderer(),
    ]
    
    if development:
        # Development: colorized console output
        processors.extend([
            structlog.dev.ConsoleRenderer(
                colors=True,
                exception_formatter=structlog.dev.plain_traceback,
            )
        ])
    else:
        # Production: JSON output
        processors.extend([
            structlog.processors.format_exc_info,
            structlog.processors.JSONRenderer(),
        ])
    
    structlog.configure(
        processors=processors,
        wrapper_class=structlog.make_filtering_bound_logger(level.value),
        context_class=dict,
        logger_factory=structlog.PrintLoggerFactory(),
        cache_logger_on_first_use=True,
    )


class StructuredLogger:
    """
    Structured logger implementation using structlog.
    
    Provides JSON logging for production and colored console for development.
    """
    
    def __init__(self, **fields: Any) -> None:
        """
        Initialize structured logger.
        
        Args:
            **fields: Persistent fields to attach to all log messages
        """
        self._logger = structlog.get_logger()
        if fields:
            self._logger = self._logger.bind(**fields)
    
    def debug(self, msg: str, **kwargs: Any) -> None:
        """Log debug message."""
        self._logger.debug(msg, **kwargs)
    
    def info(self, msg: str, **kwargs: Any) -> None:
        """Log info message."""
        self._logger.info(msg, **kwargs)
    
    def warn(self, msg: str, **kwargs: Any) -> None:
        """Log warning message."""
        self._logger.warning(msg, **kwargs)
    
    def error(self, msg: str, **kwargs: Any) -> None:
        """Log error message."""
        self._logger.error(msg, **kwargs)
    
    def with_fields(self, **kwargs: Any) -> "StructuredLogger":
        """
        Return a new logger with fields attached.
        
        Args:
            **kwargs: Fields to attach
            
        Returns:
            New logger instance with fields bound
        """
        new_logger = StructuredLogger()
        new_logger._logger = self._logger.bind(**kwargs)
        return new_logger
    
    def with_error(self, err: Exception) -> "StructuredLogger":
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


def new(level: Level = Level.INFO, development: bool = False) -> StructuredLogger:
    """
    Create a new structured logger.
    
    Configures structlog if not already configured.
    
    Args:
        level: Minimum log level
        development: If True, use console renderer. If False, use JSON.
        
    Returns:
        New structured logger instance
        
    Example:
        >>> logger = new(Level.INFO, development=True)
        >>> logger.info("server started", port=8080)
    """
    configure_structlog(level=level, development=development)
    return StructuredLogger()