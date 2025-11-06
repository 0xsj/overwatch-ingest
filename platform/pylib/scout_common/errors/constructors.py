"""Convenient factory functions for creating common errors."""

from typing import Any

from .base import Error
from .types import ErrorType
from .codes import Code, code


# =========================================
# Generic Error Constructors
# =========================================

def not_found(resource_type: str, resource_id: str) -> Error:
    """
    Create a not found error for a resource.
    
    Args:
        resource_type: Type of resource (e.g., "user", "incident")
        resource_id: ID of the resource that wasn't found
        
    Returns:
        Error with NOT_FOUND type
        
    Example:
        >>> err = not_found("user", "123")
        >>> assert err.error_type == ErrorType.NOT_FOUND
        >>> assert err.get_detail("resource_type") == "user"
    """
    return Error(
        error_type=ErrorType.NOT_FOUND,
        code=code("RESOURCE_NOT_FOUND"),
        message=f"{resource_type} not found: {resource_id}",
        details={
            "resource_type": resource_type,
            "resource_id": resource_id,
        },
    )


def already_exists(resource_type: str, resource_id: str) -> Error:
    """
    Create a conflict error for a resource that already exists.
    
    Args:
        resource_type: Type of resource
        resource_id: ID of the resource that already exists
        
    Returns:
        Error with ALREADY_EXISTS type
    """
    return Error(
        error_type=ErrorType.ALREADY_EXISTS,
        code=code("RESOURCE_ALREADY_EXISTS"),
        message=f"{resource_type} already exists: {resource_id}",
        details={
            "resource_type": resource_type,
            "resource_id": resource_id,
        },
    )


def validation(message: str) -> Error:
    """
    Create a generic validation error.
    
    Args:
        message: Validation error message
        
    Returns:
        Error with VALIDATION type
    """
    return Error(
        error_type=ErrorType.VALIDATION,
        code=code("VALIDATION_FAILED"),
        message=message,
    )


def validation_with_field(field: str, message: str) -> Error:
    """
    Create a validation error for a specific field.
    
    Args:
        field: Name of the invalid field
        message: Validation error message
        
    Returns:
        Error with VALIDATION type and field detail
        
    Example:
        >>> err = validation_with_field("email", "invalid format")
        >>> assert err.get_detail("field") == "email"
    """
    return Error(
        error_type=ErrorType.VALIDATION,
        code=code("VALIDATION_FAILED"),
        message=f"validation failed for field '{field}': {message}",
        details={"field": field},
    )


def required_field(field: str) -> Error:
    """
    Create a validation error for a missing required field.
    
    Args:
        field: Name of the required field
        
    Returns:
        Error with VALIDATION type
    """
    return Error(
        error_type=ErrorType.VALIDATION,
        code=code("REQUIRED_FIELD_MISSING"),
        message=f"required field '{field}' is missing",
        details={"field": field},
    )


def invalid_field(field: str, reason: str) -> Error:
    """
    Create a validation error for an invalid field value.
    
    Args:
        field: Name of the invalid field
        reason: Reason why the field is invalid
        
    Returns:
        Error with VALIDATION type
    """
    return Error(
        error_type=ErrorType.VALIDATION,
        code=code("INVALID_FIELD_VALUE"),
        message=f"invalid value for field '{field}': {reason}",
        details={
            "field": field,
            "reason": reason,
        },
    )


# =========================================
# Authorization Error Constructors
# =========================================

def unauthorized(reason: str) -> Error:
    """
    Create an unauthorized error.
    
    Args:
        reason: Reason for unauthorized access
        
    Returns:
        Error with UNAUTHORIZED type
    """
    return Error(
        error_type=ErrorType.UNAUTHORIZED,
        code=code("UNAUTHORIZED"),
        message=f"unauthorized: {reason}",
        details={"reason": reason},
    )


def forbidden(resource: str, action: str) -> Error:
    """
    Create a forbidden error.
    
    Args:
        resource: Resource that was attempted to access
        action: Action that was attempted
        
    Returns:
        Error with FORBIDDEN type
        
    Example:
        >>> err = forbidden("document", "delete")
        >>> assert err.get_detail("resource") == "document"
        >>> assert err.get_detail("action") == "delete"
    """
    return Error(
        error_type=ErrorType.FORBIDDEN,
        code=code("FORBIDDEN"),
        message=f"forbidden: cannot {action} {resource}",
        details={
            "resource": resource,
            "action": action,
        },
    )


# =========================================
# Internal Error Constructors
# =========================================

def internal(message: str) -> Error:
    """
    Create a generic internal error.
    
    Args:
        message: Error message
        
    Returns:
        Error with INTERNAL type
    """
    return Error(
        error_type=ErrorType.INTERNAL,
        code=code("INTERNAL_ERROR"),
        message=message,
    )


def internal_with_cause(message: str, cause: Error) -> Error:
    """
    Create an internal error with a cause.
    
    Args:
        message: Error message
        cause: Underlying error that caused this error
        
    Returns:
        Error with INTERNAL type and cause
    """
    return Error(
        error_type=ErrorType.INTERNAL,
        code=code("INTERNAL_ERROR"),
        message=message,
        cause=cause,
    )


def not_implemented(feature: str) -> Error:
    """
    Create a not implemented error.
    
    Args:
        feature: Name of the unimplemented feature
        
    Returns:
        Error with NOT_IMPLEMENTED type
    """
    return Error(
        error_type=ErrorType.NOT_IMPLEMENTED,
        code=code("NOT_IMPLEMENTED"),
        message=f"feature not implemented: {feature}",
        details={"feature": feature},
    )


# =========================================
# Timeout Error Constructors
# =========================================

def timeout(operation: str, duration_ms: int) -> Error:
    """
    Create a timeout error.
    
    Args:
        operation: Operation that timed out
        duration_ms: Timeout duration in milliseconds
        
    Returns:
        Error with TIMEOUT type
        
    Example:
        >>> err = timeout("database query", 5000)
        >>> assert err.get_detail("timeout_ms") == "5000"
    """
    return Error(
        error_type=ErrorType.TIMEOUT,
        code=code("OPERATION_TIMEOUT"),
        message=f"operation timed out after {duration_ms}ms: {operation}",
        details={
            "operation": operation,
            "timeout_ms": str(duration_ms),
        },
    )


# =========================================
# Service Availability Constructors
# =========================================

def unavailable(service: str) -> Error:
    """
    Create a service unavailable error.
    
    Args:
        service: Name of the unavailable service
        
    Returns:
        Error with UNAVAILABLE type
    """
    return Error(
        error_type=ErrorType.UNAVAILABLE,
        code=code("SERVICE_UNAVAILABLE"),
        message=f"service unavailable: {service}",
        details={"service": service},
    )


def unavailable_with_cause(service: str, cause: Error) -> Error:
    """
    Create a service unavailable error with a cause.
    
    Args:
        service: Name of the unavailable service
        cause: Underlying error
        
    Returns:
        Error with UNAVAILABLE type and cause
    """
    return Error(
        error_type=ErrorType.UNAVAILABLE,
        code=code("SERVICE_UNAVAILABLE"),
        message=f"service unavailable: {service}",
        details={"service": service},
        cause=cause,
    )


# =========================================
# Conflict Error Constructors
# =========================================

def conflict(resource: str, reason: str) -> Error:
    """
    Create a conflict error.
    
    Args:
        resource: Resource with the conflict
        reason: Reason for the conflict
        
    Returns:
        Error with CONFLICT type
    """
    return Error(
        error_type=ErrorType.CONFLICT,
        code=code("RESOURCE_CONFLICT"),
        message=f"conflict on {resource}: {reason}",
        details={
            "resource": resource,
            "reason": reason,
        },
    )


# =========================================
# Rate Limiting Constructors
# =========================================

def rate_limit(limit: int, window: str) -> Error:
    """
    Create a rate limit exceeded error.
    
    Args:
        limit: Request limit
        window: Time window (e.g., "minute", "hour")
        
    Returns:
        Error with RATE_LIMIT type
    """
    return Error(
        error_type=ErrorType.RATE_LIMIT,
        code=code("RATE_LIMIT_EXCEEDED"),
        message=f"rate limit exceeded: {limit} requests per {window}",
        details={
            "limit": str(limit),
            "window": window,
        },
    )


def rate_limit_with_retry(limit: int, window: str, retry_after_seconds: int) -> Error:
    """
    Create a rate limit error with retry-after information.
    
    Args:
        limit: Request limit
        window: Time window
        retry_after_seconds: Seconds to wait before retrying
        
    Returns:
        Error with RATE_LIMIT type and retry info
    """
    return Error(
        error_type=ErrorType.RATE_LIMIT,
        code=code("RATE_LIMIT_EXCEEDED"),
        message=f"rate limit exceeded: {limit} requests per {window}",
        details={
            "limit": str(limit),
            "window": window,
            "retry_after_seconds": str(retry_after_seconds),
        },
    )


# =========================================
# Infrastructure Error Constructors
# =========================================

def database_error(operation: str, cause: Error | None = None) -> Error:
    """
    Create a database error.
    
    Args:
        operation: Database operation that failed
        cause: Optional underlying error
        
    Returns:
        Error with DATABASE type
    """
    return Error(
        error_type=ErrorType.DATABASE,
        code=code("DATABASE_ERROR"),
        message=f"database operation failed: {operation}",
        details={"operation": operation},
        cause=cause,
    )


def database_error_with_table(
    operation: str, table: str, cause: Error | None = None
) -> Error:
    """
    Create a database error with table context.
    
    Args:
        operation: Database operation that failed
        table: Table name
        cause: Optional underlying error
        
    Returns:
        Error with DATABASE type and table detail
    """
    return Error(
        error_type=ErrorType.DATABASE,
        code=code("DATABASE_ERROR"),
        message=f"database operation failed: {operation} on table {table}",
        details={
            "operation": operation,
            "table": table,
        },
        cause=cause,
    )


def cache_error(operation: str, cause: Error | None = None) -> Error:
    """
    Create a cache error.
    
    Args:
        operation: Cache operation that failed
        cause: Optional underlying error
        
    Returns:
        Error with CACHE type
    """
    return Error(
        error_type=ErrorType.CACHE,
        code=code("CACHE_ERROR"),
        message=f"cache operation failed: {operation}",
        details={"operation": operation},
        cause=cause,
    )


def cache_error_with_key(operation: str, key: str, cause: Error | None = None) -> Error:
    """
    Create a cache error with key context.
    
    Args:
        operation: Cache operation that failed
        key: Cache key
        cause: Optional underlying error
        
    Returns:
        Error with CACHE type and key detail
    """
    return Error(
        error_type=ErrorType.CACHE,
        code=code("CACHE_ERROR"),
        message=f"cache operation failed: {operation} for key {key}",
        details={
            "operation": operation,
            "key": key,
        },
        cause=cause,
    )


def network_error(operation: str, cause: Error | None = None) -> Error:
    """
    Create a network error.
    
    Args:
        operation: Network operation that failed
        cause: Optional underlying error
        
    Returns:
        Error with NETWORK type
    """
    return Error(
        error_type=ErrorType.NETWORK,
        code=code("NETWORK_ERROR"),
        message=f"network operation failed: {operation}",
        details={"operation": operation},
        cause=cause,
    )


def network_error_with_url(operation: str, url: str, cause: Error | None = None) -> Error:
    """
    Create a network error with URL context.
    
    Args:
        operation: Network operation that failed
        url: URL that was accessed
        cause: Optional underlying error
        
    Returns:
        Error with NETWORK type and URL detail
    """
    return Error(
        error_type=ErrorType.NETWORK,
        code=code("NETWORK_ERROR"),
        message=f"network operation failed: {operation} to {url}",
        details={
            "operation": operation,
            "url": url,
        },
        cause=cause,
    )


def event_error(operation: str, cause: Error | None = None) -> Error:
    """
    Create an event bus error.
    
    Args:
        operation: Event operation that failed
        cause: Optional underlying error
        
    Returns:
        Error with EVENT type
    """
    return Error(
        error_type=ErrorType.EVENT,
        code=code("EVENT_ERROR"),
        message=f"event operation failed: {operation}",
        details={"operation": operation},
        cause=cause,
    )


def event_error_with_subject(
    operation: str, subject: str, cause: Error | None = None
) -> Error:
    """
    Create an event error with subject context.
    
    Args:
        operation: Event operation that failed
        subject: Event subject
        cause: Optional underlying error
        
    Returns:
        Error with EVENT type and subject detail
    """
    return Error(
        error_type=ErrorType.EVENT,
        code=code("EVENT_ERROR"),
        message=f"event operation failed: {operation} for subject {subject}",
        details={
            "operation": operation,
            "subject": subject,
        },
        cause=cause,
    )


# =========================================
# Error Wrapping Functions
# =========================================

def wrap(cause: Error, error_type: ErrorType, error_code: Code, message: str) -> Error:
    """
    Wrap an error with additional context.
    
    Args:
        cause: The underlying error
        error_type: Type for the wrapper error
        error_code: Code for the wrapper error
        message: Message for the wrapper error
        
    Returns:
        New error wrapping the cause
        
    Example:
        >>> root = database_error("query failed")
        >>> wrapped = wrap(root, ErrorType.INTERNAL, code("OPERATION_FAILED"), "failed to get user")
        >>> assert wrapped.cause == root
    """
    return Error(
        error_type=error_type,
        code=error_code,
        message=message,
        cause=cause,
    )


def wrap_with_details(
    cause: Error,
    error_type: ErrorType,
    error_code: Code,
    message: str,
    details: dict[str, str],
) -> Error:
    """
    Wrap an error with additional context and details.
    
    Args:
        cause: The underlying error
        error_type: Type for the wrapper error
        error_code: Code for the wrapper error
        message: Message for the wrapper error
        details: Additional details
        
    Returns:
        New error wrapping the cause with details
    """
    return Error(
        error_type=error_type,
        code=error_code,
        message=message,
        details=details,
        cause=cause,
    )