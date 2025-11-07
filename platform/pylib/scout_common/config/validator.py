"""Validation utilities for configuration values."""

import re
from typing import Any, TypeVar


from .errors import (
    missing_required,
    invalid_value,
    invalid_format,
    out_of_range,
    invalid_choice,
)
from .parser import parse_url


T = TypeVar("T", int, float)


def validate_required(key: str, value: str) -> None:
    """
    Check if a value is non-empty.
    
    Args:
        key: Configuration key name
        value: Value to validate
        
    Raises:
        Error: If value is empty
    """
    if not value or not value.strip():
        raise missing_required(key)


def validate_range(key: str, value: T, min_val: T, max_val: T) -> None:
    """
    Check if a numeric value is within the specified range (inclusive).
    
    Args:
        key: Configuration key name
        value: Value to validate
        min_val: Minimum allowed value (inclusive)
        max_val: Maximum allowed value (inclusive)
        
    Raises:
        Error: If value is outside the range
    """
    if value < min_val or value > max_val:
        raise out_of_range(key, value, min_val, max_val)


def validate_min_max(
    key: str,
    value: T,
    min_val: T | None = None,
    max_val: T | None = None,
) -> None:
    """
    Check if a numeric value satisfies min/max constraints.
    
    Pass None for min_val or max_val to skip that check.
    
    Args:
        key: Configuration key name
        value: Value to validate
        min_val: Minimum allowed value (optional)
        max_val: Maximum allowed value (optional)
        
    Raises:
        Error: If value violates constraints
    """
    if min_val is not None and value < min_val:
        raise out_of_range(key, value, min_val, "infinity")
    if max_val is not None and value > max_val:
        raise out_of_range(key, value, "infinity", max_val)


def validate_choice(key: str, value: str, allowed_values: list[str]) -> None:
    """
    Check if a value is in the allowed set.
    
    Args:
        key: Configuration key name
        value: Value to validate
        allowed_values: List of allowed values
        
    Raises:
        Error: If value is not in allowed set
    """
    if value not in allowed_values:
        raise invalid_choice(key, value, allowed_values)


def validate_pattern(key: str, value: str, pattern: str) -> None:
    """
    Check if a value matches the given regex pattern.
    
    Args:
        key: Configuration key name
        value: Value to validate
        pattern: Regex pattern
        
    Raises:
        Error: If value doesn't match pattern
    """
    if not re.match(pattern, value):
        raise invalid_format(key, value, f"pattern: {pattern}")


def validate_url(key: str, value: str) -> None:
    """
    Check if a value is a valid URL with required scheme.
    
    Args:
        key: Configuration key name
        value: URL string to validate
        
    Raises:
        Error: If URL is invalid
    """
    try:
        parse_url(value)
    except ValueError as e:
        raise invalid_format(key, value, "valid URL with scheme") from e


def validate_port(key: str, port: int) -> None:
    """
    Check if a port number is in the valid range (1-65535).
    
    Args:
        key: Configuration key name
        port: Port number to validate
        
    Raises:
        Error: If port is outside valid range
    """
    validate_range(key, port, 1, 65535)


def validate_non_zero(key: str, value: T) -> None:
    """
    Check if a numeric value is non-zero.
    
    Args:
        key: Configuration key name
        value: Value to validate
        
    Raises:
        Error: If value is zero
    """
    if value == 0:
        raise invalid_value(key, "0", "value must be non-zero")


def validate_positive(key: str, value: T) -> None:
    """
    Check if a numeric value is positive (> 0).
    
    Args:
        key: Configuration key name
        value: Value to validate
        
    Raises:
        Error: If value is not positive
    """
    if value <= 0:
        raise invalid_value(key, str(value), "value must be positive")


def validate_non_negative(key: str, value: T) -> None:
    """
    Check if a numeric value is non-negative (>= 0).
    
    Args:
        key: Configuration key name
        value: Value to validate
        
    Raises:
        Error: If value is negative
    """
    if value < 0:
        raise invalid_value(key, str(value), "value must be non-negative")


def validate_min_length(key: str, value: str, min_length: int) -> None:
    """
    Check if a string meets minimum length requirement.
    
    Args:
        key: Configuration key name
        value: String to validate
        min_length: Minimum required length
        
    Raises:
        Error: If string is too short
    """
    if len(value) < min_length:
        raise invalid_value(
            key,
            value,
            f"must be at least {min_length} characters",
        )


def validate_max_length(key: str, value: str, max_length: int) -> None:
    """
    Check if a string does not exceed maximum length.
    
    Args:
        key: Configuration key name
        value: String to validate
        max_length: Maximum allowed length
        
    Raises:
        Error: If string is too long
    """
    if len(value) > max_length:
        raise invalid_value(
            key,
            value,
            f"must be at most {max_length} characters",
        )