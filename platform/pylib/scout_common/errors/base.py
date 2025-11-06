"""Core Error dataclass for Scout platform errors."""

from dataclasses import dataclass, field
from typing import Self, Optional

from .types import ErrorType
from .codes import Code


@dataclass(frozen=True, slots=True)
class Error:
    """
    Core error value with rich metadata.
    
    Error is an immutable dataclass (not an exception) that represents
    an error condition with structured information. It supports error
    chaining, contextual details, and builder-style construction.
    
    Attributes:
        error_type: Categorizes the error (validation, not_found, etc.)
        code: Specific error identifier
        message: Human-readable error description
        details: Additional context as key-value pairs
        cause: Optional underlying error that caused this error
        
    Example:
        >>> from scout_common.errors import Error, ErrorType, code
        >>> 
        >>> err = Error(
        ...     error_type=ErrorType.VALIDATION,
        ...     code=code("REQUIRED_FIELD"),
        ...     message="email is required",
        ...     details={"field": "email"}
        ... )
        >>> 
        >>> # Builder pattern
        >>> err = Error(
        ...     error_type=ErrorType.DATABASE,
        ...     code=code("QUERY_FAILED"),
        ...     message="failed to query users"
        ... ).with_detail("table", "users").with_detail("operation", "SELECT")
    """

    error_type: ErrorType
    code: Code
    message: str
    details: dict[str, str] = field(default_factory=dict)
    cause: Optional["Error"] = None

    def __str__(self) -> str:
        """
        String representation of the error.
        
        Format: [TYPE:CODE] message
        
        If there's a cause, it's appended with " | caused by: ..."
        """
        msg = f"[{self.error_type.value}:{self.code}] {self.message}"
        
        if self.cause:
            msg += f" | caused by: {self.cause}"
        
        return msg

    def __repr__(self) -> str:
        """Detailed representation for debugging."""
        parts = [
            f"error_type={self.error_type!r}",
            f"code={self.code!r}",
            f"message={self.message!r}",
        ]
        
        if self.details:
            parts.append(f"details={self.details!r}")
        
        if self.cause:
            parts.append(f"cause={self.cause!r}")
        
        return f"Error({', '.join(parts)})"

    # Builder methods - return new Error instances (immutable)

    def with_detail(self, key: str, value: str) -> Self:
        """
        Add a single detail to the error.
        
        Returns a new Error instance with the additional detail.
        
        Args:
            key: Detail key
            value: Detail value
            
        Returns:
            New Error with the added detail
            
        Example:
            >>> err = Error(ErrorType.DATABASE, code("QUERY_FAILED"), "query failed")
            >>> err = err.with_detail("table", "users")
            >>> assert err.details["table"] == "users"
        """
        new_details = {**self.details, key: value}
        return Error(
            error_type=self.error_type,
            code=self.code,
            message=self.message,
            details=new_details,
            cause=self.cause,
        )

    def with_details(self, details: dict[str, str]) -> Self:
        """
        Add multiple details to the error.
        
        Returns a new Error instance with the additional details.
        
        Args:
            details: Dictionary of details to add
            
        Returns:
            New Error with the added details
            
        Example:
            >>> err = Error(ErrorType.DATABASE, code("QUERY_FAILED"), "query failed")
            >>> err = err.with_details({"table": "users", "operation": "SELECT"})
            >>> assert len(err.details) == 2
        """
        new_details = {**self.details, **details}
        return Error(
            error_type=self.error_type,
            code=self.code,
            message=self.message,
            details=new_details,
            cause=self.cause,
        )

    def with_cause(self, cause: "Error") -> Self:
        """
        Add a cause (underlying error) to this error.
        
        Returns a new Error instance with the cause set.
        
        Args:
            cause: The underlying error that caused this error
            
        Returns:
            New Error with the cause set
            
        Example:
            >>> root = Error(ErrorType.NETWORK, code("TIMEOUT"), "connection timeout")
            >>> wrapped = Error(ErrorType.DATABASE, code("QUERY_FAILED"), "query failed")
            >>> wrapped = wrapped.with_cause(root)
            >>> assert wrapped.cause == root
        """
        return Error(
            error_type=self.error_type,
            code=self.code,
            message=self.message,
            details=self.details,
            cause=cause,
        )

    # Query methods

    def get_detail(self, key: str) -> str | None:
        """
        Get a detail value by key.
        
        Args:
            key: The detail key
            
        Returns:
            The detail value if it exists, None otherwise
        """
        return self.details.get(key)

    def has_detail(self, key: str) -> bool:
        """
        Check if a detail key exists.
        
        Args:
            key: The detail key
            
        Returns:
            True if the key exists, False otherwise
        """
        return key in self.details

    def has_details(self) -> bool:
        """
        Check if the error has any details.
        
        Returns:
            True if there are details, False otherwise
        """
        return bool(self.details)

    def is_client_error(self) -> bool:
        """
        Check if this is a client error (4xx).
        
        Returns:
            True if this is a client error, False otherwise
        """
        return self.error_type.is_client_error()

    def is_server_error(self) -> bool:
        """
        Check if this is a server error (5xx).
        
        Returns:
            True if this is a server error, False otherwise
        """
        return self.error_type.is_server_error()

    def is_retryable(self) -> bool:
        """
        Check if this error can typically be retried.
        
        Returns:
            True if the error is retryable, False otherwise
        """
        return self.error_type.is_retryable()

    def http_status_code(self) -> int:
        """
        Get the recommended HTTP status code.
        
        Returns:
            HTTP status code (400, 404, 500, etc.)
        """
        return self.error_type.http_status_code()

    def get_cause(self) -> Optional["Error"]:
        """
        Get the underlying cause of this error.
        
        Returns:
            The cause error if present, None otherwise
        """
        return self.cause

    def get_root_cause(self) -> "Error":
        """
        Get the root cause by traversing the error chain.
        
        Returns the deepest error in the cause chain.
        
        Returns:
            The root cause error (self if no cause)
        """
        current = self
        while current.cause is not None:
            current = current.cause
        return current

    def matches(
        self,
        *,
        error_type: ErrorType | None = None,
        code: Code | None = None,
        message_contains: str | None = None,
    ) -> bool:
        """
        Check if this error matches specific criteria.
        
        All non-None criteria must match for this to return True.
        
        Args:
            error_type: Expected error type (optional)
            code: Expected error code (optional)
            message_contains: Substring expected in message (optional)
            
        Returns:
            True if all specified criteria match, False otherwise
            
        Example:
            >>> err = Error(ErrorType.VALIDATION, code("REQUIRED"), "email required")
            >>> assert err.matches(error_type=ErrorType.VALIDATION)
            >>> assert err.matches(code=code("REQUIRED"))
            >>> assert err.matches(message_contains="email")
            >>> assert err.matches(error_type=ErrorType.VALIDATION, code=code("REQUIRED"))
        """
        if error_type is not None and self.error_type != error_type:
            return False

        if code is not None and self.code != code:
            return False

        if message_contains is not None and message_contains not in self.message:
            return False

        return True


# Helper function for creating errors (alternative to constructor)
def error(
    error_type: ErrorType,
    code: Code,
    message: str,
    *,
    details: dict[str, str] | None = None,
    cause: Error | None = None,
) -> Error:
    """
    Create an Error instance.
    
    This is a convenience function that provides a more concise way
    to create errors compared to the full constructor.
    
    Args:
        error_type: The error type
        code: The error code
        message: Human-readable error message
        details: Optional details dictionary
        cause: Optional underlying error
        
    Returns:
        A new Error instance
        
    Example:
        >>> from scout_common.errors import error, ErrorType, code
        >>> 
        >>> err = error(
        ...     ErrorType.VALIDATION,
        ...     code("REQUIRED_FIELD"),
        ...     "email is required",
        ...     details={"field": "email"}
        ... )
    """
    return Error(
        error_type=error_type,
        code=code,
        message=message,
        details=details or {},
        cause=cause,
    )