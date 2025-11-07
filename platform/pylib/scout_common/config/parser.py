"""Type conversion utilities for configuration values."""

from datetime import timedelta
from typing import Any
from urllib.parse import urlparse, ParseResult


def parse_string(value: str) -> str:
    """Trim whitespace and return the string value."""
    return value.strip()


def parse_int(value: str) -> int:
    """
    Parse a string to an integer.
    
    Args:
        value: String to parse
        
    Returns:
        Parsed integer
        
    Raises:
        ValueError: If value cannot be parsed as integer
    """
    trimmed = value.strip()
    return int(trimmed)


def parse_bool(value: str) -> bool:
    """
    Parse a string to a boolean.
    
    Accepts: "true", "false", "1", "0", "yes", "no", "on", "off" (case-insensitive).
    
    Args:
        value: String to parse
        
    Returns:
        Parsed boolean
        
    Raises:
        ValueError: If value is not a recognized boolean string
    """
    trimmed = value.strip().lower()
    
    if trimmed in ("true", "1", "yes", "on"):
        return True
    elif trimmed in ("false", "0", "no", "off"):
        return False
    else:
        raise ValueError(f"invalid boolean value: {value}")


def parse_float(value: str) -> float:
    """
    Parse a string to a float.
    
    Args:
        value: String to parse
        
    Returns:
        Parsed float
        
    Raises:
        ValueError: If value cannot be parsed as float
    """
    trimmed = value.strip()
    return float(trimmed)


def parse_duration(value: str) -> timedelta:
    """
    Parse a string to a timedelta.
    
    Accepts formats like "5s", "10m", "1h", "2d".
    Supports combinations like "1h30m".
    
    Args:
        value: Duration string (e.g., "5s", "10m", "1h30m")
        
    Returns:
        Parsed timedelta
        
    Raises:
        ValueError: If format is invalid
        
    Examples:
        >>> parse_duration("5s")
        timedelta(seconds=5)
        >>> parse_duration("10m")
        timedelta(minutes=10)
        >>> parse_duration("1h30m")
        timedelta(hours=1, minutes=30)
    """
    trimmed = value.strip()
    if not trimmed:
        raise ValueError("empty duration string")
    
    # Parse duration components
    total_seconds = 0.0
    current_num = ""
    
    for char in trimmed:
        if char.isdigit() or char == ".":
            current_num += char
        elif char in "smhd":
            if not current_num:
                raise ValueError(f"invalid duration format: {value}")
            
            num = float(current_num)
            
            if char == "s":
                total_seconds += num
            elif char == "m":
                total_seconds += num * 60
            elif char == "h":
                total_seconds += num * 3600
            elif char == "d":
                total_seconds += num * 86400
            
            current_num = ""
        else:
            raise ValueError(f"invalid duration format: {value}")
    
    if current_num:
        raise ValueError(f"invalid duration format: {value}")
    
    return timedelta(seconds=total_seconds)


def parse_url(value: str) -> ParseResult:
    """
    Parse and validate a URL string.
    
    Args:
        value: URL string to parse
        
    Returns:
        Parsed URL result
        
    Raises:
        ValueError: If URL is invalid or missing scheme
        
    Examples:
        >>> result = parse_url("https://example.com:8080/path")
        >>> result.scheme
        'https'
        >>> result.netloc
        'example.com:8080'
    """
    trimmed = value.strip()
    if not trimmed:
        raise ValueError("empty URL string")
    
    parsed = urlparse(trimmed)
    
    # Validate that we have a scheme
    if not parsed.scheme:
        raise ValueError(f"URL missing scheme: {value}")
    
    return parsed


def parse_string_list(value: str, separator: str = ",") -> list[str]:
    """
    Parse a separated string into a list.
    
    Trims whitespace from each element and filters out empty strings.
    
    Args:
        value: String to parse
        separator: Separator character (default: ",")
        
    Returns:
        List of trimmed strings
        
    Examples:
        >>> parse_string_list("a, b, c")
        ['a', 'b', 'c']
        >>> parse_string_list("x:y:z", separator=":")
        ['x', 'y', 'z']
    """
    if not value:
        return []
    
    parts = value.split(separator)
    result = []
    
    for part in parts:
        trimmed = part.strip()
        if trimmed:
            result.append(trimmed)
    
    return result