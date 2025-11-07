"""Environment variable loading utilities with validation."""

import os
from datetime import timedelta
from typing import TypeVar
from urllib.parse import ParseResult

from .errors import missing_required, invalid_value, invalid_format
from .parser import (
    parse_string,
    parse_int,
    parse_bool,
    parse_float,
    parse_duration,
    parse_url,
    parse_string_list,
)
from .validator import (
    validate_required,
    validate_range,
    validate_choice,
    validate_port,
)


T = TypeVar("T")


def load_string_required(key: str) -> str:
    """
    Load a required string from environment.
    
    Args:
        key: Environment variable name
        
    Returns:
        Parsed string value
        
    Raises:
        Error: If variable is missing or empty
    """
    value = os.getenv(key, "")
    validate_required(key, value)
    return parse_string(value)


def load_string_optional(key: str, default_value: str) -> str:
    """
    Load an optional string from environment with default.
    
    Args:
        key: Environment variable name
        default_value: Default value if not set
        
    Returns:
        Parsed string value or default
    """
    value = os.getenv(key, "")
    if not value:
        return default_value
    return parse_string(value)


def load_int_required(key: str) -> int:
    """
    Load a required integer from environment.
    
    Args:
        key: Environment variable name
        
    Returns:
        Parsed integer value
        
    Raises:
        Error: If variable is missing or not a valid integer
    """
    value = os.getenv(key, "")
    validate_required(key, value)
    
    try:
        return parse_int(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid integer")


def load_int_optional(key: str, default_value: int) -> int:
    """
    Load an optional integer from environment with default.
    
    Args:
        key: Environment variable name
        default_value: Default value if not set
        
    Returns:
        Parsed integer value or default
        
    Raises:
        Error: If value is set but not a valid integer
    """
    value = os.getenv(key, "")
    if not value:
        return default_value
    
    try:
        return parse_int(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid integer")


def load_int_with_range(
    key: str,
    min_val: int,
    max_val: int,
    default_value: int | None = None,
) -> int:
    """
    Load an integer and validate it's within range.
    
    Args:
        key: Environment variable name
        min_val: Minimum allowed value (inclusive)
        max_val: Maximum allowed value (inclusive)
        default_value: Default value if not set (optional)
        
    Returns:
        Parsed and validated integer
        
    Raises:
        Error: If value is missing, invalid, or out of range
    """
    value = os.getenv(key, "")
    
    # If not set and default provided, use default
    if not value:
        if default_value is not None:
            validate_range(key, default_value, min_val, max_val)
            return default_value
        raise missing_required(key)
    
    try:
        parsed = parse_int(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid integer")
    
    validate_range(key, parsed, min_val, max_val)
    return parsed


def load_bool_required(key: str) -> bool:
    """
    Load a required boolean from environment.
    
    Args:
        key: Environment variable name
        
    Returns:
        Parsed boolean value
        
    Raises:
        Error: If variable is missing or not a valid boolean
    """
    value = os.getenv(key, "")
    validate_required(key, value)
    
    try:
        return parse_bool(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid boolean")


def load_bool_optional(key: str, default_value: bool) -> bool:
    """
    Load an optional boolean from environment with default.
    
    Args:
        key: Environment variable name
        default_value: Default value if not set
        
    Returns:
        Parsed boolean value or default
        
    Raises:
        Error: If value is set but not a valid boolean
    """
    value = os.getenv(key, "")
    if not value:
        return default_value
    
    try:
        return parse_bool(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid boolean")


def load_float_required(key: str) -> float:
    """
    Load a required float from environment.
    
    Args:
        key: Environment variable name
        
    Returns:
        Parsed float value
        
    Raises:
        Error: If variable is missing or not a valid float
    """
    value = os.getenv(key, "")
    validate_required(key, value)
    
    try:
        return parse_float(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid float")


def load_float_optional(key: str, default_value: float) -> float:
    """
    Load an optional float from environment with default.
    
    Args:
        key: Environment variable name
        default_value: Default value if not set
        
    Returns:
        Parsed float value or default
        
    Raises:
        Error: If value is set but not a valid float
    """
    value = os.getenv(key, "")
    if not value:
        return default_value
    
    try:
        return parse_float(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid float")


def load_duration_required(key: str) -> timedelta:
    """
    Load a required duration from environment.
    
    Args:
        key: Environment variable name
        
    Returns:
        Parsed duration value
        
    Raises:
        Error: If variable is missing or not a valid duration
    """
    value = os.getenv(key, "")
    validate_required(key, value)
    
    try:
        return parse_duration(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid duration")


def load_duration_optional(key: str, default_value: timedelta) -> timedelta:
    """
    Load an optional duration from environment with default.
    
    Args:
        key: Environment variable name
        default_value: Default value if not set
        
    Returns:
        Parsed duration value or default
        
    Raises:
        Error: If value is set but not a valid duration
    """
    value = os.getenv(key, "")
    if not value:
        return default_value
    
    try:
        return parse_duration(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid duration")


def load_url_required(key: str) -> ParseResult:
    """
    Load a required URL from environment.
    
    Args:
        key: Environment variable name
        
    Returns:
        Parsed URL
        
    Raises:
        Error: If variable is missing or not a valid URL
    """
    value = os.getenv(key, "")
    validate_required(key, value)
    
    try:
        return parse_url(value)
    except ValueError:
        raise invalid_format(key, value, "valid URL with scheme")


def load_url_optional(key: str, default_value: ParseResult | None = None) -> ParseResult | None:
    """
    Load an optional URL from environment with default.
    
    Args:
        key: Environment variable name
        default_value: Default value if not set (optional)
        
    Returns:
        Parsed URL or default
        
    Raises:
        Error: If value is set but not a valid URL
    """
    value = os.getenv(key, "")
    if not value:
        return default_value
    
    try:
        return parse_url(value)
    except ValueError:
        raise invalid_format(key, value, "valid URL with scheme")


def load_string_list_required(key: str, separator: str = ",") -> list[str]:
    """
    Load a required comma-separated list from environment.
    
    Args:
        key: Environment variable name
        separator: Separator character (default: ",")
        
    Returns:
        List of parsed strings
        
    Raises:
        Error: If variable is missing or list is empty
    """
    value = os.getenv(key, "")
    validate_required(key, value)
    
    parsed = parse_string_list(value, separator)
    if not parsed:
        raise invalid_value(key, value, "list is empty after parsing")
    
    return parsed


def load_string_list_optional(
    key: str,
    default_value: list[str] | None = None,
    separator: str = ",",
) -> list[str]:
    """
    Load an optional comma-separated list from environment.
    
    Args:
        key: Environment variable name
        default_value: Default value if not set (optional)
        separator: Separator character (default: ",")
        
    Returns:
        List of parsed strings or default
    """
    if default_value is None:
        default_value = []
    
    value = os.getenv(key, "")
    if not value:
        return default_value
    
    parsed = parse_string_list(value, separator)
    if not parsed:
        return default_value
    
    return parsed


def load_string_with_choice(
    key: str,
    allowed_values: list[str],
    default_value: str | None = None,
) -> str:
    """
    Load a string and validate it's in allowed set.
    
    Args:
        key: Environment variable name
        allowed_values: List of allowed values
        default_value: Default value if not set (optional)
        
    Returns:
        Validated string value
        
    Raises:
        Error: If value is missing, or not in allowed set
    """
    value = os.getenv(key, "")
    
    # If not set and default provided, use default
    if not value:
        if default_value is not None:
            validate_choice(key, default_value, allowed_values)
            return default_value
        raise missing_required(key)
    
    parsed = parse_string(value)
    validate_choice(key, parsed, allowed_values)
    return parsed


def load_port_required(key: str) -> int:
    """
    Load a required port number from environment.
    
    Args:
        key: Environment variable name
        
    Returns:
        Validated port number
        
    Raises:
        Error: If variable is missing, invalid, or out of range
    """
    port = load_int_required(key)
    validate_port(key, port)
    return port


def load_port_optional(key: str, default_value: int) -> int:
    """
    Load an optional port number from environment with default.
    
    Args:
        key: Environment variable name
        default_value: Default port number
        
    Returns:
        Validated port number or default
        
    Raises:
        Error: If default is invalid or value is set but invalid
    """
    # Validate default
    validate_port(key, default_value)
    
    value = os.getenv(key, "")
    if not value:
        return default_value
    
    try:
        port = parse_int(value)
    except ValueError:
        raise invalid_value(key, value, "not a valid integer")
    
    validate_port(key, port)
    return port


def with_prefix(prefix: str, key: str) -> str:
    """
    Return a prefixed key for environment variables.
    
    Args:
        prefix: Prefix to prepend
        key: Base key name
        
    Returns:
        Prefixed key
        
    Example:
        >>> with_prefix("GATEWAY_", "PORT")
        'GATEWAY_PORT'
    """
    return prefix + key