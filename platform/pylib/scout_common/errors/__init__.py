# platform/pylib/scout_common/errors/__init__.py
"""
Error handling system for Scout platform.

This package provides a robust, functional error handling system with:
- Rich error types with metadata
- Result type for explicit error handling
- Comprehensive error constructors
- Chainable operations (combinators)
- JSON serialization
- Decorators for common patterns
- Exception bridge for interoperability

Quick Start:
    >>> from scout_common.errors import (
    ...     ErrorType, Error, Ok, Err, Result,
    ...     validation, not_found, internal
    ... )
    >>> 
    >>> # Create errors
    >>> err = validation("email is required").with_detail("field", "email")
    >>> 
    >>> # Use Result type
    >>> def get_user(user_id: str) -> Result[User, Error]:
    ...     if not user_id:
    ...         return Err(validation("user_id required"))
    ...     user = db.find(user_id)
    ...     if not user:
    ...         return Err(not_found("user", user_id))
    ...     return Ok(user)
    >>> 
    >>> # Pattern matching
    >>> match get_user("123"):
    ...     case Ok(user):
    ...         print(f"Found: {user}")
    ...     case Err(error):
    ...         print(f"Error: {error}")
    >>> 
    >>> # Method chaining
    >>> result = (
    ...     get_user(user_id)
    ...     .map(lambda u: u.email)
    ...     .unwrap_or("no-email@example.com")
    ... )

For detailed documentation, see: https://docs.scout.dev/errors
"""

# Version
__version__ = "0.1.0"

# Core types
from .types import ErrorType
from .codes import Code, code
from .base import Error, error

# Result type
from .result import (
    Ok,
    Err,
    Result,
    # Type guards
    is_ok,
    is_err,
    # Unwrapping
    unwrap,
    unwrap_or,
    unwrap_or_else,
    expect,
    unwrap_err,
    # Transformations
    map_value,
    map_error,
    and_then,
    or_else,
    flatten,
    # Inspection
    inspect,
    inspect_err,
    # Collections
    collect,
    partition,
)

# Constructors
from .constructors import (
    # Generic errors
    not_found,
    already_exists,
    validation,
    validation_with_field,
    required_field,
    invalid_field,
    # Authorization
    unauthorized,
    forbidden,
    # Internal
    internal,
    internal_with_cause,
    not_implemented,
    # Timeout/availability
    timeout,
    unavailable,
    unavailable_with_cause,
    # Conflict/rate limit
    conflict,
    rate_limit,
    rate_limit_with_retry,
    # Infrastructure
    database_error,
    database_error_with_table,
    cache_error,
    cache_error_with_key,
    network_error,
    network_error_with_url,
    event_error,
    event_error_with_subject,
    # Wrapping
    wrap,
    wrap_with_details,
)

# Combinators (chainable operations)
from .combinators import (
    OkResult,
    ErrResult,
    ok,
    err,
    wrap as wrap_result,
    Pipeline,
    compose,
    safe,
    safe_async,
)

# Serialization
from .serialization import (
    ErrorDict,
    to_dict,
    to_dict_verbose,
    from_dict,
    to_json,
    to_json_verbose,
    from_json,
    to_http_response,
    format_error,
    to_dict_list,
    from_dict_list,
    is_valid_error_dict,
    get_error_schema,
)

# Decorators
from .decorators import (
    # Note: safe and safe_async already imported from combinators
    retry,
    retry_async,
    with_timeout,
    with_timeout_async,
    validate_args,
    with_error_context,
    with_error_context_async,
    with_fallback,
    log_errors,
)

# Exceptions (interop)
from .exceptions import (
    ErrorException,
    unwrap_or_raise,
    expect_or_raise,
    catch,
    catch_async,
    from_exception,
    from_error_exception,
    catch_errors,
    catch_errors_async,
    raise_error,
    raise_if_err,
)

# Define public API
__all__ = [
    # Version
    "__version__",
    
    # Core types
    "ErrorType",
    "Code",
    "code",
    "Error",
    "error",
    
    # Result type
    "Ok",
    "Err",
    "Result",
    "is_ok",
    "is_err",
    "unwrap",
    "unwrap_or",
    "unwrap_or_else",
    "expect",
    "unwrap_err",
    "map_value",
    "map_error",
    "and_then",
    "or_else",
    "flatten",
    "inspect",
    "inspect_err",
    "collect",
    "partition",
    
    # Constructors
    "not_found",
    "already_exists",
    "validation",
    "validation_with_field",
    "required_field",
    "invalid_field",
    "unauthorized",
    "forbidden",
    "internal",
    "internal_with_cause",
    "not_implemented",
    "timeout",
    "unavailable",
    "unavailable_with_cause",
    "conflict",
    "rate_limit",
    "rate_limit_with_retry",
    "database_error",
    "database_error_with_table",
    "cache_error",
    "cache_error_with_key",
    "network_error",
    "network_error_with_url",
    "event_error",
    "event_error_with_subject",
    "wrap",
    "wrap_with_details",
    
    # Combinators
    "OkResult",
    "ErrResult",
    "ok",
    "err",
    "wrap_result",
    "Pipeline",
    "compose",
    "safe",
    "safe_async",
    
    # Serialization
    "ErrorDict",
    "to_dict",
    "to_dict_verbose",
    "from_dict",
    "to_json",
    "to_json_verbose",
    "from_json",
    "to_http_response",
    "format_error",
    "to_dict_list",
    "from_dict_list",
    "is_valid_error_dict",
    "get_error_schema",
    
    # Decorators
    "retry",
    "retry_async",
    "with_timeout",
    "with_timeout_async",
    "validate_args",
    "with_error_context",
    "with_error_context_async",
    "with_fallback",
    "log_errors",
    
    # Exceptions
    "ErrorException",
    "unwrap_or_raise",
    "expect_or_raise",
    "catch",
    "catch_async",
    "from_exception",
    "from_error_exception",
    "catch_errors",
    "catch_errors_async",
    "raise_error",
    "raise_if_err",
]