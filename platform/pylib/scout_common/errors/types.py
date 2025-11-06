"""Error type categorization for Scout platform errors."""

from enum import Enum
from typing import Self

class ErrorType(str, Enum):
    """
    Error type categories.
    
    Types are intentionally generic to support all platform services.
    Services can use these types with specific error codes for domain errors.
    """

    # Client Errors (4xx equivalent) - caused by invalid client input or state
    # These errors typically should not be retried without changing the request.
    VALIDATION = "VALIDATION"
    NOT_FOUND = "NOT_FOUND"
    ALREADY_EXISTS = "ALREADY_EXISTS"
    UNAUTHORIZED = "UNAUTHORIZED"
    FORBIDDEN = "FORBIDDEN"
    CONFLICT = "CONFLICT"
    RATE_LIMIT = "RATE_LIMIT"

    # Server Errors (5xx equivalent) - caused by server-side issues
    # These errors may be retryable depending on the specific type.
    INTERNAL = "INTERNAL"
    UNAVAILABLE = "UNAVAILABLE"
    TIMEOUT = "TIMEOUT"
    NOT_IMPLEMENTED = "NOT_IMPLEMENTED"

    # Infrastructure Errors - platform component failures
    DATABASE = "DATABASE"
    CACHE = "CACHE"
    NETWORK = "NETWORK"
    EVENT = "EVENT"

    def is_client_error(self) -> bool:
        """
        Check if this error is caused by client input (4xx equivalent).
        
        These errors typically should not be retried without changing the request.
        """

        return self in {
            ErrorType.VALIDATION,
            ErrorType.NOT_FOUND,
            ErrorType.ALREADY_EXISTS,
            ErrorType.UNAUTHORIZED,
            ErrorType.FORBIDDEN,
            ErrorType.CONFLICT,
            ErrorType.RATE_LIMIT,
        }

    def is_server_error(self) -> bool:
        """
        Check if this error is caused by server-side issues (5xx equivalent).
        
        These errors may be retryable depending on the specific type.
        """
        return self in {
            ErrorType.INTERNAL,
            ErrorType.UNAVAILABLE,
            ErrorType.TIMEOUT,
            ErrorType.NOT_IMPLEMENTED,
            ErrorType.DATABASE,
            ErrorType.CACHE,
            ErrorType.NETWORK,
            ErrorType.EVENT,
        }

    def is_retryable(self) -> bool:
        """
        Check if errors of this type can typically be retried.
        
        Note: Even retryable errors should use exponential backoff and respect retry limits.
        """
        return self in {
            ErrorType.TIMEOUT,
            ErrorType.UNAVAILABLE,
            ErrorType.RATE_LIMIT,
            ErrorType.INTERNAL,
            ErrorType.NETWORK,
            ErrorType.CACHE,
            ErrorType.DATABASE,
            ErrorType.EVENT,
        }

    def http_status_code(self) -> int:
        """Get the recommended HTTP status code for this error type."""
        match self:
            case ErrorType.VALIDATION:
                return 400  # Bad Request
            case ErrorType.UNAUTHORIZED:
                return 401  # Unauthorized
            case ErrorType.FORBIDDEN:
                return 403  # Forbidden
            case ErrorType.NOT_FOUND:
                return 404  # Not Found
            case ErrorType.CONFLICT | ErrorType.ALREADY_EXISTS:
                return 409  # Conflict
            case ErrorType.RATE_LIMIT:
                return 429  # Too Many Requests
            case ErrorType.INTERNAL:
                return 500  # Internal Server Error
            case ErrorType.NOT_IMPLEMENTED:
                return 501  # Not Implemented
            case ErrorType.UNAVAILABLE | ErrorType.DATABASE | ErrorType.CACHE:
                return 503  # Service Unavailable
            case ErrorType.TIMEOUT | ErrorType.NETWORK:
                return 504  # Gateway Timeout
            case _:
                return 500  # Internal Server Error (safe default)

    @classmethod
    def from_string(cls, value: str) -> Self | None:
        """
        Parse an ErrorType from a string.
        
        Returns None if the string is not a valid ErrorType.
        """
        try:
            return cls(value)
        except ValueError:
            return None

    @classmethod
    def all_types(cls) -> list[Self]:
        """Get all valid ErrorTypes."""
        return list(cls)


# Type alias for convenience
ErrorTypeValue = ErrorType