"""Combinators and chainable operations for Result types."""

from typing import TypeVar, Callable, ParamSpec, Concatenate

from .result import (
    Result,
    Ok,
    Err,
    is_ok,
    is_err,
    unwrap,
    unwrap_or,
    unwrap_or_else,
    expect,
    unwrap_err,
    map_value,
    map_error,
    and_then,
    or_else,
    flatten,
    inspect,
    inspect_err,
)
from .base import Error


T = TypeVar("T")
U = TypeVar("U")
E = TypeVar("E")
P = ParamSpec("P")


class ResultMixin:
    """
    Mixin to add chainable methods to Result types.
    
    This allows for more ergonomic method chaining:
    
    Example:
        >>> result = (
        ...     get_user(user_id)
        ...     .map(lambda u: u.email)
        ...     .and_then(send_email)
        ...     .map(lambda _: "Email sent")
        ...     .unwrap_or("Failed to send")
        ... )
    
    Note: Python doesn't support true method chaining on union types,
    so we provide helper functions that return wrapped results.
    """
    pass


# Wrapper class for ergonomic chaining
class OkResult:
    """Wrapper around Ok for method chaining."""
    
    def __init__(self, value: T):
        self._result: Result[T, E] = Ok(value)
    
    @property
    def result(self) -> Result[T, E]:
        """Get the underlying Result."""
        return self._result
    
    def map(self, func: Callable[[T], U]) -> "OkResult[U] | ErrResult[E]":
        """Transform the value inside Ok."""
        mapped = map_value(self._result, func)
        return _wrap_result(mapped)
    
    def map_err(self, func: Callable[[E], U]) -> "OkResult[T] | ErrResult[U]":
        """Transform the error (no-op for Ok)."""
        mapped = map_error(self._result, func)
        return _wrap_result(mapped)
    
    def and_then(self, func: Callable[[T], Result[U, E]]) -> "OkResult[U] | ErrResult[E]":
        """Chain another operation that can fail."""
        chained = and_then(self._result, func)
        return _wrap_result(chained)
    
    def or_else(self, func: Callable[[E], Result[T, U]]) -> "OkResult[T] | ErrResult[U]":
        """Handle errors (no-op for Ok)."""
        handled = or_else(self._result, func)
        return _wrap_result(handled)
    
    def inspect(self, func: Callable[[T], None]) -> "OkResult[T]":
        """Call a function for side effects."""
        inspect(self._result, func)
        return self
    
    def inspect_err(self, func: Callable[[E], None]) -> "OkResult[T]":
        """Call a function for side effects on error (no-op for Ok)."""
        inspect_err(self._result, func)
        return self
    
    def unwrap(self) -> T:
        """Extract the value."""
        return unwrap(self._result)
    
    def unwrap_or(self, default: T) -> T:
        """Extract the value or return default."""
        return unwrap_or(self._result, default)
    
    def unwrap_or_else(self, default_fn: Callable[[E], T]) -> T:
        """Extract the value or compute default."""
        return unwrap_or_else(self._result, default_fn)
    
    def expect(self, message: str) -> T:
        """Extract the value or raise with message."""
        return expect(self._result, message)
    
    def is_ok(self) -> bool:
        """Check if this is Ok."""
        return True
    
    def is_err(self) -> bool:
        """Check if this is Err."""
        return False


class ErrResult:
    """Wrapper around Err for method chaining."""
    
    def __init__(self, error: E):
        self._result: Result[T, E] = Err(error)
    
    @property
    def result(self) -> Result[T, E]:
        """Get the underlying Result."""
        return self._result
    
    def map(self, func: Callable[[T], U]) -> "OkResult[U] | ErrResult[E]":
        """Transform the value (no-op for Err)."""
        mapped = map_value(self._result, func)
        return _wrap_result(mapped)
    
    def map_err(self, func: Callable[[E], U]) -> "OkResult[T] | ErrResult[U]":
        """Transform the error inside Err."""
        mapped = map_error(self._result, func)
        return _wrap_result(mapped)
    
    def and_then(self, func: Callable[[T], Result[U, E]]) -> "OkResult[U] | ErrResult[E]":
        """Chain another operation (no-op for Err)."""
        chained = and_then(self._result, func)
        return _wrap_result(chained)
    
    def or_else(self, func: Callable[[E], Result[T, U]]) -> "OkResult[T] | ErrResult[U]":
        """Handle errors by providing alternative."""
        handled = or_else(self._result, func)
        return _wrap_result(handled)
    
    def inspect(self, func: Callable[[T], None]) -> "ErrResult[E]":
        """Call a function for side effects (no-op for Err)."""
        inspect(self._result, func)
        return self
    
    def inspect_err(self, func: Callable[[E], None]) -> "ErrResult[E]":
        """Call a function for side effects on error."""
        inspect_err(self._result, func)
        return self
    
    def unwrap(self) -> T:
        """Extract the value (raises for Err)."""
        return unwrap(self._result)
    
    def unwrap_or(self, default: T) -> T:
        """Extract the value or return default."""
        return unwrap_or(self._result, default)
    
    def unwrap_or_else(self, default_fn: Callable[[E], T]) -> T:
        """Extract the value or compute default."""
        return unwrap_or_else(self._result, default_fn)
    
    def expect(self, message: str) -> T:
        """Extract the value or raise with message."""
        return expect(self._result, message)
    
    def unwrap_err(self) -> E:
        """Extract the error."""
        return unwrap_err(self._result)
    
    def is_ok(self) -> bool:
        """Check if this is Ok."""
        return False
    
    def is_err(self) -> bool:
        """Check if this is Err."""
        return True


def _wrap_result(result: Result[T, E]) -> OkResult[T] | ErrResult[E]:
    """Wrap a Result in the appropriate wrapper class."""
    match result:
        case Ok(value):
            return OkResult(value)
        case Err(error):
            return ErrResult(error)


# Convenience functions to start chains
def ok(value: T) -> OkResult[T]:
    """
    Create an OkResult for method chaining.
    
    Args:
        value: The success value
        
    Returns:
        OkResult wrapping the value
        
    Example:
        >>> result = (
        ...     ok(5)
        ...     .map(lambda x: x * 2)
        ...     .and_then(validate)
        ...     .unwrap()
        ... )
    """
    return OkResult(value)


def err(error: E) -> ErrResult[E]:
    """
    Create an ErrResult for method chaining.
    
    Args:
        error: The error value
        
    Returns:
        ErrResult wrapping the error
        
    Example:
        >>> from scout_common.errors import validation
        >>> result = (
        ...     err(validation("invalid"))
        ...     .or_else(lambda e: Ok(default_value))
        ...     .unwrap()
        ... )
    """
    return ErrResult(error)


def wrap(result: Result[T, E]) -> OkResult[T] | ErrResult[E]:
    """
    Wrap a Result for method chaining.
    
    Args:
        result: The Result to wrap
        
    Returns:
        Wrapped result with chainable methods
        
    Example:
        >>> result = get_user(user_id)
        >>> processed = (
        ...     wrap(result)
        ...     .map(lambda u: u.email)
        ...     .inspect(lambda email: print(f"Email: {email}"))
        ...     .unwrap_or("no-email@example.com")
        ... )
    """
    return _wrap_result(result)


# Pipeline operator alternative (Python doesn't have |>)
class Pipeline:
    """
    Pipeline builder for composing operations.
    
    Allows building a pipeline of operations that can be applied to Results.
    
    Example:
        >>> process_user = (
        ...     Pipeline()
        ...     .then(validate_user)
        ...     .then(save_user)
        ...     .then(send_notification)
        ... )
        >>> 
        >>> result = process_user.run(Ok(user))
    """
    
    def __init__(self):
        self._operations: list[Callable[[Result], Result]] = []
    
    def then(self, operation: Callable[[T], Result[U, E]]) -> "Pipeline":
        """
        Add an operation to the pipeline.
        
        Args:
            operation: Function to apply in the pipeline
            
        Returns:
            Self for chaining
        """
        def wrapped(result: Result[T, E]) -> Result[U, E]:
            return and_then(result, operation)
        
        self._operations.append(wrapped)
        return self
    
    def map(self, func: Callable[[T], U]) -> "Pipeline":
        """
        Add a mapping operation to the pipeline.
        
        Args:
            func: Function to map over the value
            
        Returns:
            Self for chaining
        """
        def wrapped(result: Result[T, E]) -> Result[U, E]:
            return map_value(result, func)
        
        self._operations.append(wrapped)
        return self
    
    def recover(self, handler: Callable[[E], Result[T, U]]) -> "Pipeline":
        """
        Add error recovery to the pipeline.
        
        Args:
            handler: Function to handle errors
            
        Returns:
            Self for chaining
        """
        def wrapped(result: Result[T, E]) -> Result[T, U]:
            return or_else(result, handler)
        
        self._operations.append(wrapped)
        return self
    
    def run(self, initial: Result[T, E]) -> Result:
        """
        Run the pipeline on an initial Result.
        
        Args:
            initial: Starting Result
            
        Returns:
            Final Result after all operations
        """
        result = initial
        for operation in self._operations:
            result = operation(result)
        return result


# Functional composition helpers
def compose(
    *funcs: Callable[[Result[T, E]], Result[U, E]]
) -> Callable[[Result[T, E]], Result[U, E]]:
    """
    Compose multiple Result-returning functions.
    
    Creates a new function that applies each function in sequence.
    
    Args:
        *funcs: Functions to compose (applied left to right)
        
    Returns:
        Composed function
        
    Example:
        >>> process = compose(
        ...     lambda r: and_then(r, validate),
        ...     lambda r: and_then(r, save),
        ...     lambda r: map_value(r, to_dto)
        ... )
        >>> result = process(get_user(user_id))
    """
    def composed(initial: Result[T, E]) -> Result[U, E]:
        result = initial
        for func in funcs:
            result = func(result)
        return result
    
    return composed


# Utility for trying operations that might raise
def safe(func: Callable[P, T]) -> Callable[P, Result[T, Error]]:
    """
    Decorator to convert exception-raising functions to Result-returning.
    
    Wraps a function so that any exceptions are caught and converted to Err.
    
    Args:
        func: Function that might raise exceptions
        
    Returns:
        Function that returns Result instead of raising
        
    Example:
        >>> @safe
        ... def parse_int(s: str) -> int:
        ...     return int(s)
        >>> 
        >>> result = parse_int("42")  # Ok(42)
        >>> result = parse_int("abc")  # Err(...)
    """
    from .constructors import internal
    
    def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
        try:
            value = func(*args, **kwargs)
            return Ok(value)
        except Exception as e:
            return Err(internal(f"{func.__name__} failed: {str(e)}"))
    
    return wrapper


# Async support
async def safe_async(
    func: Callable[P, T]
) -> Callable[P, Result[T, Error]]:
    """
    Async version of safe decorator.
    
    Wraps an async function to return Result instead of raising.
    
    Args:
        func: Async function that might raise exceptions
        
    Returns:
        Async function that returns Result
        
    Example:
        >>> @safe_async
        ... async def fetch_user(user_id: str) -> User:
        ...     return await db.get_user(user_id)
        >>> 
        >>> result = await fetch_user("123")
    """
    from .constructors import internal
    
    async def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
        try:
            value = await func(*args, **kwargs)
            return Ok(value)
        except Exception as e:
            return Err(internal(f"{func.__name__} failed: {str(e)}"))
    
    return wrapper