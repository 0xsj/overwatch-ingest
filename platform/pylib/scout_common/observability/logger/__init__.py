"""
Structured logging for Scout platform.

This package provides structured logging with different implementations:

- NoopLogger - Colorized console output for local development
- StructuredLogger - Production-grade structured logging with JSON output

Basic Usage
-----------
Create a logger based on your environment:

    from scout_common.observability.logger import new, noop, Level
    
    # Development: colorized console
    log = noop()
    
    # Production: JSON structured logging
    log = new(Level.INFO, development=False)
    
    # Log messages with structured fields
    log.info("server started", port=8080, environment="production")

Structured Logging
------------------
All loggers use key-value pairs for structured logging:

    log.info("user created",
        user_id="123",
        email="user@example.com",
        tenant="acme",
    )
    
    # Output (JSON in production):
    # {"level":"info","timestamp":"2025-01-01T12:00:00Z","event":"user created","user_id":"123","email":"user@example.com","tenant":"acme"}

Log Levels
----------
The package supports four log levels:

    Level.DEBUG - Detailed debugging information
    Level.INFO  - General informational messages
    Level.WARN  - Warning messages for potentially harmful situations
    Level.ERROR - Error messages for serious problems

Set the level when creating the logger:

    log = new(Level.DEBUG)   # Show all messages
    log = new(Level.ERROR)   # Only errors

Attaching Fields
----------------
Create child loggers with persistent fields:

    # Request-scoped logger
    req_logger = log.with_fields(
        request_id="req-123",
        user_id="user-456",
    )
    
    req_logger.info("processing request")  # Includes request_id and user_id
    req_logger.info("request completed")   # Still includes those fields

Error Logging
-------------
Attach errors to log messages:

    try:
        process_request(user_id)
    except Exception as e:
        log.with_error(e).error("failed to process request", user_id=user_id)

Example: Service Logger
-----------------------

    from scout_common.observability.logger import new, Level
    
    def main():
        # Create production logger
        log = new(Level.INFO, development=False)
        
        # Add service-level fields
        log = log.with_fields(
            service="tools",
            version="1.0.0",
        )
        
        log.info("service starting", port=8083)
        
        # Create request-scoped logger
        handle_request(log, "req-123")
        
        log.info("service stopped")
    
    def handle_request(log, request_id: str):
        req_log = log.with_fields(request_id=request_id)
        
        req_log.info("handling request")
        req_log.debug("validating input")
        req_log.info("request completed", duration_ms=42)
"""

__version__ = "0.1.0"

# Core types
from .logger import Logger, Level, set_context_logger, get_context_logger

# Implementations
from .noop import NoopLogger
from .structured import StructuredLogger, new, configure_structlog

# Convenience functions
def noop(**fields) -> Logger:
    """
    Create a colorized console logger for development.
    
    Args:
        **fields: Persistent fields to attach
        
    Returns:
        Console logger instance
    """
    return NoopLogger(**fields)


__all__ = [
    # Version
    "__version__",
    # Core types
    "Logger",
    "Level",
    "set_context_logger",
    "get_context_logger",
    # Implementations
    "NoopLogger",
    "StructuredLogger",
    "new",
    "configure_structlog",
    # Convenience
    "noop",
]