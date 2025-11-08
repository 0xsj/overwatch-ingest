"""Logger interface and types for Scout platform."""

from typing import Protocol, Any, runtime_checkable
from enum import IntEnum
from contextvars import ContextVar

class Level(IntEnum):
    """"Log level enumeration"""

    DEBUG = 10
    INFO = 20
    WARN = 30
    ERROR = 40

    def __str__(self) -> str:
        """Return string representation of log level."""
        return self.name.lower()

    @classmethod
    def parse(cls, s: str) -> "Level":
        """
        Parse a string into a log level.
        
        Args:
            s: Level string (debug, info, warn, warning, error)
            
        Returns:
            Parsed log level, defaults to INFO if unknown
        """
        s_lower = s.lower()
        if s_lower == "debug":
            return cls.DEBUG
        elif s_lower == "info":
            return cls.INFO
        elif s_lower in ("warn", "warning"):
            return cls.WARN
        elif s_lower == "error":
            return cls.ERROR
        else:
            return cls.INFO

@runtime_checkable
class Logger(Protocol):
    """
    Logger protocol defines the logging interface for Scout platform.
    
    Provides structured logging with levels and context awareness.
    All logging methods accept key-value pairs for structured logging.
    
    Example:
        >>> logger.info("user created",
        ...     user_id="123",
        ...     email="user@example.com",
        ... )
    """
    def debug(self, msg: str, **kwargs: Any) -> None:
        """
        Log a debug-level message with optional key-value pairs.
        
        Args:
            msg: Log message
            **kwargs: Additional fields as key-value pairs
        """
        ...
    
    def info(self, msg: str, **kwargs: Any) -> None:
        """
        Log an info-level message with optional key-value pairs.
        
        Args:
            msg: Log message
            **kwargs: Additional fields as key-value pairs
        """
        ...
    
    def warn(self, msg: str, **kwargs: Any) -> None:
        """
        Log a warning-level message with optional key-value pairs.
        
        Args:
            msg: Log message
            **kwargs: Additional fields as key-value pairs
        """
        ...
    
    def error(self, msg: str, **kwargs: Any) -> None:
        """
        Log an error-level message with optional key-value pairs.
        
        Typically used for errors that should be investigated.
        
        Args:
            msg: Log message
            **kwargs: Additional fields as key-value pairs
        """
        ...
    
    def with_fields(self, **kwargs: Any) -> "Logger":
        """
        Return a new logger with the given key-value pairs attached.
        
        These fields will be included in all subsequent log entries.
        
        Args:
            **kwargs: Fields to attach as key-value pairs
            
        Returns:
            New logger instance with fields attached
            
        Example:
            >>> user_logger = logger.with_fields(user_id="123", tenant_id="abc")
            >>> user_logger.info("action performed")  # Includes user_id and tenant_id
        """
        ...
    
    def with_error(self, err: Exception) -> "Logger":
        """
        Return a new logger with an error field attached.
        
        The error will be logged with a standard "error" key.
        
        Args:
            err: Exception to attach
            
        Returns:
            New logger instance with error attached
        """
        ...


# Context variable for storing logger in async contexts
_logger_context: ContextVar[Logger | None] = ContextVar("logger", default=None)


def set_context_logger(logger: Logger) -> None:
    """
    Set the logger in the current context.
    
    This is useful for propagating loggers through async call chains.
    
    Args:
        logger: Logger instance to set
    """
    _logger_context.set(logger)


def get_context_logger() -> Logger | None:
    """
    Get the logger from the current context.
    
    Returns:
        Logger instance if set, None otherwise
    """
    return _logger_context.get()
