"""Error code type hints and utilities for Scout platform errors."""

from typing import NewType


# Code represents a specific error code.
# Codes are string identifiers that uniquely identify error conditions.
#
# Naming convention (recommended but not enforced):
#   - Use UPPER_SNAKE_CASE
#   - Format: ENTITY_CONDITION or OPERATION_FAILURE
#   - Examples: "USER_NOT_FOUND", "DATABASE_CONNECTION_FAILED"
#
# Services should define their own code constants:
#
#     # services/incidents/errors/codes.py
#     from scout_common.errors import Code
#
#     INCIDENT_NOT_FOUND: Code = Code("INCIDENT_NOT_FOUND")
#     INVALID_SEVERITY: Code = Code("INVALID_SEVERITY")
#     GEOCODING_FAILED: Code = Code("GEOCODING_FAILED")
#
Code = NewType("Code", str)


def code(value: str) -> Code:
    """
    Create a Code from a string.
    
    This is a convenience function for creating codes without the NewType constructor.
    
    Args:
        value: The error code string
        
    Returns:
        A typed Code value
        
    Example:
        >>> error_code = code("USER_NOT_FOUND")
        >>> assert error_code == "USER_NOT_FOUND"
    """
    return Code(value)


def is_empty(code_value: Code) -> bool:
    """
    Check if a code is empty.
    
    Args:
        code_value: The code to check
        
    Returns:
        True if the code is an empty string, False otherwise
        
    Example:
        >>> assert is_empty(code(""))
        >>> assert not is_empty(code("USER_NOT_FOUND"))
    """
    return code_value == ""


# Type alias for clarity in type hints
CodeValue = Code