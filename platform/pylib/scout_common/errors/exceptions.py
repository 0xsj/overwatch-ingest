"""Exception bridge for interoperability with exception-based code."""

from typing import TypeVar, Callable, ParamSpec, Awaitable, NoReturn

from .base import Error
from .result import Result, Ok, Err, unwrap
from .types import ErrorType
from .codes import code
from .constructors import internal


T = TypeVar("T")
P = ParamSpec("P")


# =========================================
# Exception Wrapper
# =========================================

class ErrorException(Exception):
    """
    Exception wrapper for Error values.
    
    Used to convert Result-based errors into exceptions for interop
    with exception-based code.
    
    Attributes:
        error: The wrapped Error instance
        
    Example:
        >>> err = validation("invalid input")
        >>> raise ErrorException(err)
    """

    def __init__(self, error: Error):
        self.error = error
        super().__init__(str(error))

    def __repr__(self) -> str:
        return f"ErrorException({self.error!r})"


# =========================================
# Result to Exception Conversion
# =========================================

def unwrap_or_raise(result: Result[T, Error]) -> T:
    """
    Extract the value from Ok, or raise ErrorException if Err.
    
    Use this when calling Result-based code from exception-based code.
    
    Args:
        result: The result to unwrap
        
    Returns:
        The value if Ok
        
    Raises:
        ErrorException: If the result is Err
        
    Example:
        >>> result = get_user(user_id)
        >>> try:
        ...     user = unwrap_or_raise(result)
        ... except ErrorException as e:
        ...     print(f"Error: {e.error}")
    """
    match result:
        case Ok(value):
            return value
        case Err(error):
            raise ErrorException(error)


def expect_or_raise(result: Result[T, Error], message: str) -> T:
    """
    Extract the value from Ok, or raise with custom message if Err.
    
    Args:
        result: The result to unwrap
        message: Custom error message
        
    Returns:
        The value if Ok
        
    Raises:
        ErrorException: If the result is Err, with custom message prepended
        
    Example:
        >>> result = get_user(user_id)
        >>> user = expect_or_raise(result, "failed to get user")
    """
    match result:
        case Ok(value):
            return value
        case Err(error):
            enhanced_error = error.with_detail("expect_message", message)
            raise ErrorException(enhanced_error)


# =========================================
# Exception to Result Conversion
# =========================================

def catch(func: Callable[P, T]) -> Callable[P, Result[T, Error]]:
    """
    Convert exception-raising function to Result-returning.
    
    This is an alias for the @safe decorator from decorators.py,
    provided here for convenience.
    
    Args:
        func: Function that might raise exceptions
        
    Returns:
        Function that returns Result
        
    Example:
        >>> def risky_operation(x: int) -> int:
        ...     if x < 0:
        ...         raise ValueError("x must be positive")
        ...     return x * 2
        >>> 
        >>> safe_op = catch(risky_operation)
        >>> result = safe_op(-5)
        >>> assert isinstance(result, Err)
    """
    from .decorators import safe
    return safe(func)


def catch_async(func: Callable[P, Awaitable[T]]) -> Callable[P, Awaitable[Result[T, Error]]]:
    """
    Convert async exception-raising function to Result-returning.
    
    Args:
        func: Async function that might raise exceptions
        
    Returns:
        Async function that returns Result
    """
    from .decorators import safe_async
    return safe_async(func)


# =========================================
# Specific Exception Handlers
# =========================================

def from_exception(exc: Exception) -> Error:
    """
    Convert a standard exception to an Error.
    
    Attempts to infer the appropriate ErrorType based on exception type.
    
    Args:
        exc: Exception to convert
        
    Returns:
        Error representation of the exception
        
    Example:
        >>> try:
        ...     risky_operation()
        ... except ValueError as e:
        ...     error = from_exception(e)
        ...     return Err(error)
    """
    exception_type = type(exc).__name__
    message = str(exc)
    
    # Map common exception types to ErrorTypes
    error_type_map = {
        "ValueError": ErrorType.VALIDATION,
        "TypeError": ErrorType.VALIDATION,
        "KeyError": ErrorType.NOT_FOUND,
        "FileNotFoundError": ErrorType.NOT_FOUND,
        "PermissionError": ErrorType.FORBIDDEN,
        "TimeoutError": ErrorType.TIMEOUT,
        "ConnectionError": ErrorType.NETWORK,
        "ConnectionRefusedError": ErrorType.UNAVAILABLE,
        "NotImplementedError": ErrorType.NOT_IMPLEMENTED,
    }
    
    error_type = error_type_map.get(exception_type, ErrorType.INTERNAL)
    
    return Error(
        error_type=error_type,
        code=code(f"{exception_type.upper()}_EXCEPTION"),
        message=message or f"{exception_type} occurred",
        details={"exception_type": exception_type},
    )


def from_error_exception(exc: ErrorException) -> Error:
    """
    Extract Error from ErrorException.
    
    Args:
        exc: ErrorException to unwrap
        
    Returns:
        The wrapped Error
        
    Example:
        >>> try:
        ...     raise ErrorException(validation("invalid"))
        ... except ErrorException as e:
        ...     error = from_error_exception(e)
    """
    return exc.error


# =========================================
# Context Manager for Exception Catching
# =========================================

class catch_errors:
    """
    Context manager that catches exceptions and converts them to Results.
    
    Usage:
        >>> with catch_errors() as result_holder:
        ...     value = risky_operation()
        ...     result_holder.set(value)
        >>> 
        >>> result = result_holder.get()
        >>> match result:
        ...     case Ok(value):
        ...         print(f"Success: {value}")
        ...     case Err(error):
        ...         print(f"Error: {error}")
    """

    def __init__(self):
        self._result: Result[T, Error] | None = None
        self._value_set = False

    def set(self, value: T) -> None:
        """Set the success value."""
        self._result = Ok(value)
        self._value_set = True

    def get(self) -> Result[T, Error]:
        """Get the result (Ok if set, Err if exception occurred)."""
        if self._result is None:
            # No value was set and no exception occurred
            return Err(internal("no value set in catch_errors context"))
        return self._result

    def __enter__(self) -> "catch_errors":
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if exc_type is not None:
            # An exception occurred
            if isinstance(exc_val, ErrorException):
                self._result = Err(from_error_exception(exc_val))
            else:
                self._result = Err(from_exception(exc_val))
            return True  # Suppress the exception
        
        if not self._value_set:
            # No exception, but no value was set either
            self._result = Err(internal("no value set in catch_errors context"))
        
        return False


class catch_errors_async:
    """
    Async context manager that catches exceptions and converts them to Results.
    
    Usage:
        >>> async with catch_errors_async() as result_holder:
        ...     value = await async_risky_operation()
        ...     result_holder.set(value)
        >>> 
        >>> result = result_holder.get()
    """

    def __init__(self):
        self._result: Result[T, Error] | None = None
        self._value_set = False

    def set(self, value: T) -> None:
        """Set the success value."""
        self._result = Ok(value)
        self._value_set = True

    def get(self) -> Result[T, Error]:
        """Get the result."""
        if self._result is None:
            return Err(internal("no value set in catch_errors_async context"))
        return self._result

    async def __aenter__(self) -> "catch_errors_async":
        return self

    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if exc_type is not None:
            if isinstance(exc_val, ErrorException):
                self._result = Err(from_error_exception(exc_val))
            else:
                self._result = Err(from_exception(exc_val))
            return True
        
        if not self._value_set:
            self._result = Err(internal("no value set in catch_errors_async context"))
        
        return False


# =========================================
# Raise Helpers
# =========================================

def raise_error(error: Error) -> NoReturn:
    """
    Raise an Error as an ErrorException.
    
    Args:
        error: Error to raise
        
    Raises:
        ErrorException: Always
        
    Example:
        >>> err = validation("invalid input")
        >>> raise_error(err)  # Raises ErrorException
    """
    raise ErrorException(error)


def raise_if_err(result: Result[T, Error]) -> T:
    """
    Return value if Ok, raise if Err.
    
    Alias for unwrap_or_raise for consistency.
    
    Args:
        result: Result to check
        
    Returns:
        Value if Ok
        
    Raises:
        ErrorException: If Err
    """
    return unwrap_or_raise(result)