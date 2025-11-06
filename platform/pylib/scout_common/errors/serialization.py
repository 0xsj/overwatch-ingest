"""JSON serialization and deserialization for Error types."""

from typing import Any, TypedDict

from .base import Error
from .types import ErrorType
from .codes import Code, code


class ErrorDict(TypedDict, total=False):
    """
    TypedDict for Error JSON representation.
    
    Attributes:
        type: Error type string
        code: Error code string
        message: Error message
        details: Optional details dictionary
        cause: Optional nested error dictionary
    """
    type: str
    code: str
    message: str
    details: dict[str, str]
    cause: "ErrorDict"


def to_dict(error: Error, *, include_cause: bool = False) -> ErrorDict:
    """
    Convert an Error to a dictionary.
    
    By default, the cause chain is excluded for security (prevents leaking
    internal errors). Use include_cause=True for debugging/logging.
    
    Args:
        error: The error to serialize
        include_cause: Whether to include the cause chain (default: False)
        
    Returns:
        Dictionary representation of the error
        
    Example:
        >>> from scout_common.errors import validation
        >>> err = validation("invalid input").with_detail("field", "email")
        >>> d = to_dict(err)
        >>> assert d["type"] == "VALIDATION"
        >>> assert d["details"]["field"] == "email"
    """
    result: ErrorDict = {
        "type": error.error_type.value,
        "code": error.code,
        "message": error.message,
    }
    
    # Only include details if present
    if error.details:
        result["details"] = error.details
    
    # Include cause if requested
    if include_cause and error.cause is not None:
        result["cause"] = to_dict(error.cause, include_cause=True)
    
    return result


def to_dict_verbose(error: Error) -> ErrorDict:
    """
    Convert an Error to a dictionary with full cause chain.
    
    This is a convenience function that calls to_dict with include_cause=True.
    Use this for internal logging/debugging, not for external APIs.
    
    Args:
        error: The error to serialize
        
    Returns:
        Dictionary representation with full cause chain
        
    Example:
        >>> root = internal("root cause")
        >>> wrapped = internal_with_cause("wrapper", root)
        >>> d = to_dict_verbose(wrapped)
        >>> assert "cause" in d
        >>> assert d["cause"]["message"] == "root cause"
    """
    return to_dict(error, include_cause=True)


def from_dict(data: dict[str, Any]) -> Error:
    """
    Create an Error from a dictionary.
    
    Args:
        data: Dictionary representation of an error
        
    Returns:
        Reconstructed Error instance
        
    Raises:
        ValueError: If required fields are missing or invalid
        
    Example:
        >>> data = {
        ...     "type": "VALIDATION",
        ...     "code": "REQUIRED_FIELD",
        ...     "message": "field is required",
        ...     "details": {"field": "email"}
        ... }
        >>> err = from_dict(data)
        >>> assert err.error_type == ErrorType.VALIDATION
    """
    # Validate required fields
    if "type" not in data:
        raise ValueError("Missing required field: type")
    if "code" not in data:
        raise ValueError("Missing required field: code")
    if "message" not in data:
        raise ValueError("Missing required field: message")
    
    # Parse error type
    error_type = ErrorType.from_string(data["type"])
    if error_type is None:
        raise ValueError(f"Invalid error type: {data['type']}")
    
    # Parse details (optional)
    details = data.get("details", {})
    if not isinstance(details, dict):
        raise ValueError("details must be a dictionary")
    
    # Parse cause recursively (optional)
    cause = None
    if "cause" in data and data["cause"] is not None:
        if not isinstance(data["cause"], dict):
            raise ValueError("cause must be a dictionary")
        cause = from_dict(data["cause"])
    
    return Error(
        error_type=error_type,
        code=code(data["code"]),
        message=data["message"],
        details=details,
        cause=cause,
    )


def to_json(error: Error, *, include_cause: bool = False, indent: int | None = None) -> str:
    """
    Convert an Error to a JSON string.
    
    Args:
        error: The error to serialize
        include_cause: Whether to include the cause chain (default: False)
        indent: JSON indentation level (None for compact)
        
    Returns:
        JSON string representation
        
    Example:
        >>> err = validation("invalid")
        >>> json_str = to_json(err, indent=2)
        >>> print(json_str)
        {
          "type": "VALIDATION",
          "code": "VALIDATION_FAILED",
          "message": "invalid"
        }
    """
    import json
    
    data = to_dict(error, include_cause=include_cause)
    return json.dumps(data, indent=indent)


def to_json_verbose(error: Error, *, indent: int | None = None) -> str:
    """
    Convert an Error to a JSON string with full cause chain.
    
    Args:
        error: The error to serialize
        indent: JSON indentation level (None for compact)
        
    Returns:
        JSON string representation with full cause chain
    """
    import json
    
    data = to_dict_verbose(error)
    return json.dumps(data, indent=indent)


def from_json(json_str: str) -> Error:
    """
    Create an Error from a JSON string.
    
    Args:
        json_str: JSON string representation of an error
        
    Returns:
        Reconstructed Error instance
        
    Raises:
        ValueError: If the JSON is invalid or required fields are missing
        
    Example:
        >>> json_str = '{"type": "VALIDATION", "code": "INVALID", "message": "invalid input"}'
        >>> err = from_json(json_str)
        >>> assert err.error_type == ErrorType.VALIDATION
    """
    import json
    
    try:
        data = json.loads(json_str)
    except json.JSONDecodeError as e:
        raise ValueError(f"Invalid JSON: {e}") from e
    
    if not isinstance(data, dict):
        raise ValueError("JSON must represent an object")
    
    return from_dict(data)


# HTTP response helpers
def to_http_response(
    error: Error,
    *,
    include_cause: bool = False,
) -> tuple[ErrorDict, int]:
    """
    Convert an Error to an HTTP response tuple.
    
    Returns both the error dictionary and the appropriate HTTP status code.
    
    Args:
        error: The error to convert
        include_cause: Whether to include the cause chain
        
    Returns:
        Tuple of (error_dict, status_code)
        
    Example:
        >>> err = not_found("user", "123")
        >>> error_dict, status_code = to_http_response(err)
        >>> assert status_code == 404
        >>> 
        >>> # With Flask/FastAPI:
        >>> return jsonify(error_dict), status_code
    """
    return to_dict(error, include_cause=include_cause), error.http_status_code()


# Pretty printing for debugging
def format_error(error: Error, *, include_cause: bool = True, indent: int = 2) -> str:
    """
    Format an error as a human-readable string.
    
    Args:
        error: The error to format
        include_cause: Whether to include the cause chain
        indent: Indentation level for nested causes
        
    Returns:
        Formatted string representation
        
    Example:
        >>> root = database_error("query failed")
        >>> wrapped = internal_with_cause("operation failed", root)
        >>> print(format_error(wrapped))
        [INTERNAL:INTERNAL_ERROR] operation failed
          Caused by: [DATABASE:DATABASE_ERROR] database operation failed: query failed
            operation: query failed
    """
    lines = [str(error)]
    
    # Add details
    if error.details:
        for key, value in error.details.items():
            lines.append(f"  {key}: {value}")
    
    # Add cause chain
    if include_cause and error.cause is not None:
        cause_str = format_error(error.cause, include_cause=True, indent=indent)
        for line in cause_str.split("\n"):
            lines.append(f"  Caused by: {line}" if not line.startswith(" ") else f"  {line}")
    
    return "\n".join(lines)


# Batch operations
def to_dict_list(errors: list[Error], *, include_cause: bool = False) -> list[ErrorDict]:
    """
    Convert a list of errors to a list of dictionaries.
    
    Args:
        errors: List of errors to serialize
        include_cause: Whether to include cause chains
        
    Returns:
        List of error dictionaries
    """
    return [to_dict(error, include_cause=include_cause) for error in errors]


def from_dict_list(data: list[dict[str, Any]]) -> list[Error]:
    """
    Create a list of errors from a list of dictionaries.
    
    Args:
        data: List of error dictionaries
        
    Returns:
        List of reconstructed errors
    """
    return [from_dict(item) for item in data]


# Validation helpers
def is_valid_error_dict(data: Any) -> bool:
    """
    Check if a dictionary is a valid error representation.
    
    Args:
        data: Data to validate
        
    Returns:
        True if valid, False otherwise
        
    Example:
        >>> data = {"type": "VALIDATION", "code": "INVALID", "message": "invalid"}
        >>> assert is_valid_error_dict(data)
        >>> 
        >>> assert not is_valid_error_dict({"type": "INVALID_TYPE"})
    """
    if not isinstance(data, dict):
        return False
    
    # Check required fields
    if not all(key in data for key in ["type", "code", "message"]):
        return False
    
    # Validate error type
    error_type = ErrorType.from_string(data["type"])
    if error_type is None:
        return False
    
    # Validate details if present
    if "details" in data and not isinstance(data["details"], dict):
        return False
    
    # Validate cause recursively if present
    if "cause" in data and data["cause"] is not None:
        if not is_valid_error_dict(data["cause"]):
            return False
    
    return True


# Schema generation (for OpenAPI/JSON Schema)
def get_error_schema() -> dict[str, Any]:
    """
    Get JSON Schema for Error type.
    
    Useful for generating OpenAPI documentation.
    
    Returns:
        JSON Schema dictionary
        
    Example:
        >>> schema = get_error_schema()
        >>> assert schema["type"] == "object"
        >>> assert "type" in schema["required"]
    """
    return {
        "type": "object",
        "required": ["type", "code", "message"],
        "properties": {
            "type": {
                "type": "string",
                "enum": [t.value for t in ErrorType.all_types()],
                "description": "Error type category",
            },
            "code": {
                "type": "string",
                "description": "Specific error code",
            },
            "message": {
                "type": "string",
                "description": "Human-readable error message",
            },
            "details": {
                "type": "object",
                "additionalProperties": {"type": "string"},
                "description": "Additional error context",
            },
            "cause": {
                "$ref": "#/components/schemas/Error",
                "description": "Underlying error that caused this error",
            },
        },
    }