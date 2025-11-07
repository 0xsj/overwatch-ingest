"""Config-specific errors for Scout platform."""

from typing import Any
from scout_common.errors import Error, ErrorType, Code, code

# Error codes for config-related errors
CODE_MISSING_REQUIRED = code("CONFIG_MISSING_REQUIRED")
CODE_INVALID_VALUE = code("CONFIG_INVALID_VALUE")
CODE_INVALID_FORMAT = code("CONFIG_INVALID_FORMAT")
CODE_OUT_OF_RANGE = code("CONFIG_OUT_OF_RANGE")
CODE_INVALID_CHOICE = code("CONFIG_INVALID_CHOICE")

def missing_required(key: str) -> Error:
    """
    Create an error for a missing required configuration value.
    
    Args:
        key: The configuration key that is missing
        
    Returns:
        Error with details about the missing configuration
        
    Example:
        >>> err = missing_required("DATABASE_HOST")
        >>> # Error: required configuration not found: DATABASE_HOST
    """
    return Error(
        error_type=ErrorType.VALIDATION,
        code=CODE_MISSING_REQUIRED,
        message=f"required configuration not found: {key}",
    ).with_detail("key", key)


def invalid_value(key: str, value: str, reason: str) -> Error:
    """
    Create an error for a configuration value that cannot be parsed or converted.
    
    Args:
        key: The configuration key
        value: The invalid value
        reason: Why the value is invalid
        
    Returns:
        Error with details about the invalid value
        
    Example:
        >>> err = invalid_value("PORT", "abc", "not a valid integer")
        >>> # Error: invalid configuration value for PORT: not a valid integer (got: abc)
    """
    return Error(
        error_type=ErrorType.VALIDATION,
        code=CODE_INVALID_VALUE,
        message=f"invalid configuration value for {key}: {reason} (got: {value})",
    ).with_detail("key", key).with_detail("value", value).with_detail("reason", reason)


def invalid_format(key: str, value: str, expected_format: str) -> Error:
    """
    Create an error for a configuration value with incorrect format.
    
    Args:
        key: The configuration key
        value: The invalid value
        expected_format: The expected format description
        
    Returns:
        Error with details about the format mismatch
        
    Example:
        >>> err = invalid_format("DATABASE_URL", "localhost:5432", "postgresql://host:port/db")
        >>> # Error: invalid format for DATABASE_URL: expected postgresql://host:port/db (got: localhost:5432)
    """
    return Error(
        error_type=ErrorType.VALIDATION,
        code=CODE_INVALID_FORMAT,
        message=f"invalid format for {key}: expected {expected_format} (got: {value})",
    ).with_detail("key", key).with_detail("value", value).with_detail("expected_format", expected_format)


def out_of_range(key: str, value: Any, min_val: Any, max_val: Any) -> Error:
    """
    Create an error for a numeric configuration value outside the valid range.
    
    Args:
        key: The configuration key
        value: The out-of-range value
        min_val: Minimum allowed value
        max_val: Maximum allowed value
        
    Returns:
        Error with details about the range violation
        
    Example:
        >>> err = out_of_range("PORT", 70000, 1024, 65535)
        >>> # Error: PORT value out of range: must be between 1024 and 65535 (got: 70000)
    """
    return Error(
        error_type=ErrorType.VALIDATION,
        code=CODE_OUT_OF_RANGE,
        message=f"{key} value out of range: must be between {min_val} and {max_val} (got: {value})",
    ).with_detail("key", key).with_detail("value", str(value)).with_detail("min", str(min_val)).with_detail("max", str(max_val))


def invalid_choice(key: str, value: str, allowed_values: list[str]) -> Error:
    """
    Create an error for a configuration value not in the allowed set.
    
    Args:
        key: The configuration key
        value: The invalid value
        allowed_values: List of allowed values
        
    Returns:
        Error with details about the invalid choice
        
    Example:
        >>> err = invalid_choice("LOG_LEVEL", "verbose", ["debug", "info", "warn", "error"])
        >>> # Error: invalid value for LOG_LEVEL: must be one of [debug, info, warn, error] (got: verbose)
    """
    allowed = ", ".join(allowed_values)
    return Error(
        error_type=ErrorType.VALIDATION,
        code=CODE_INVALID_CHOICE,
        message=f"invalid value for {key}: must be one of [{allowed}] (got: {value})",
    ).with_detail("key", key).with_detail("value", value).with_detail("allowed_values", allowed)