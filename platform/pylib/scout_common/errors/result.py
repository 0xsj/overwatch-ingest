"""Result type for explicit error handling in Scout platform."""

from dataclasses import dataclass
from typing import TypeVar, Generic, Callable, Never, cast

from .base import Error


T = TypeVar("T")
U = TypeVar("U")
E = TypeVar("E")


@dataclass(frozen=True, slots=True)
class Ok(Generic[T]):
    """
    Represents a successful result containing a value.
    
    Attributes:
        value: The successful result value
    """

    value: T

    def __repr__(self) -> str:
        return f"Ok({self.value!r})"


@dataclass(frozen=True, slots=True)
class Err(Generic[E]):
    """
    Represents a failed result containing an error.
    
    Attributes:
        error: The error value
    """

    error: E

    def __repr__(self) -> str:
        return f"Err({self.error!r})"


# Result type is a union of Ok and Err
Result = Ok[T] | Err[E]


# Type guard functions
def is_ok(result: Result[T, E]) -> bool:
    """
    Check if a result is Ok.
    
    Args:
        result: The result to check
        
    Returns:
        True if the result is Ok, False otherwise
        
    Example:
        >>> result = Ok(42)
        >>> assert is_ok(result)
    """
    return isinstance(result, Ok)


def is_err(result: Result[T, E]) -> bool:
    """
    Check if a result is Err.
    
    Args:
        result: The result to check
        
    Returns:
        True if the result is Err, False otherwise
        
    Example:
        >>> from scout_common.errors import error, ErrorType, code
        >>> result = Err(error(ErrorType.VALIDATION, code("INVALID"), "invalid"))
        >>> assert is_err(result)
    """
    return isinstance(result, Err)


# Unwrapping functions
def unwrap(result: Result[T, E]) -> T:
    """
    Extract the value from an Ok result, or raise if Err.
    
    This should only be used when you're certain the result is Ok,
    or when you want to propagate the error as an exception.
    
    Args:
        result: The result to unwrap
        
    Returns:
        The value if Ok
        
    Raises:
        ValueError: If the result is Err
        
    Example:
        >>> result = Ok(42)
        >>> assert unwrap(result) == 42
        >>> 
        >>> result = Err(error(ErrorType.INTERNAL, code("ERR"), "failed"))
        >>> unwrap(result)  # Raises ValueError
    """
    match result:
        case Ok(value):
            return value
        case Err(error):
            raise ValueError(f"Called unwrap on Err: {error}")


def unwrap_or(result: Result[T, E], default: T) -> T:
    """
    Extract the value from Ok, or return a default if Err.
    
    Args:
        result: The result to unwrap
        default: Default value to return if Err
        
    Returns:
        The value if Ok, otherwise the default
        
    Example:
        >>> result = Ok(42)
        >>> assert unwrap_or(result, 0) == 42
        >>> 
        >>> result = Err(error(ErrorType.INTERNAL, code("ERR"), "failed"))
        >>> assert unwrap_or(result, 0) == 0
    """
    match result:
        case Ok(value):
            return value
        case Err(_):
            return default


def unwrap_or_else(result: Result[T, E], default_fn: Callable[[E], T]) -> T:
    """
    Extract the value from Ok, or compute a default from the error.
    
    Args:
        result: The result to unwrap
        default_fn: Function to compute default from error
        
    Returns:
        The value if Ok, otherwise the result of default_fn(error)
        
    Example:
        >>> result = Err(error(ErrorType.INTERNAL, code("ERR"), "failed"))
        >>> value = unwrap_or_else(result, lambda e: 0)
        >>> assert value == 0
    """
    match result:
        case Ok(value):
            return value
        case Err(error):
            return default_fn(error)


def expect(result: Result[T, E], message: str) -> T:
    """
    Extract the value from Ok, or raise with a custom message if Err.
    
    Args:
        result: The result to unwrap
        message: Custom error message
        
    Returns:
        The value if Ok
        
    Raises:
        ValueError: If the result is Err, with the custom message
        
    Example:
        >>> result = Ok(42)
        >>> assert expect(result, "should have value") == 42
        >>> 
        >>> result = Err(error(ErrorType.INTERNAL, code("ERR"), "failed"))
        >>> expect(result, "expected user")  # Raises with custom message
    """
    match result:
        case Ok(value):
            return value
        case Err(error):
            raise ValueError(f"{message}: {error}")


def unwrap_err(result: Result[T, E]) -> E:
    """
    Extract the error from an Err result, or raise if Ok.
    
    Useful in tests or when you know the result should be an error.
    
    Args:
        result: The result to unwrap
        
    Returns:
        The error if Err
        
    Raises:
        ValueError: If the result is Ok
        
    Example:
        >>> result = Err(error(ErrorType.VALIDATION, code("INVALID"), "invalid"))
        >>> err = unwrap_err(result)
        >>> assert err.code == code("INVALID")
    """
    match result:
        case Ok(value):
            raise ValueError(f"Called unwrap_err on Ok: {value}")
        case Err(error):
            return error


# Transformation functions (map, and_then, or_else)
def map_value(result: Result[T, E], func: Callable[[T], U]) -> Result[U, E]:
    """
    Transform the value inside an Ok result.
    
    If the result is Err, it's returned unchanged.
    
    Args:
        result: The result to transform
        func: Function to apply to the Ok value
        
    Returns:
        Ok with transformed value, or the original Err
        
    Example:
        >>> result = Ok(5)
        >>> result = map_value(result, lambda x: x * 2)
        >>> assert unwrap(result) == 10
        >>> 
        >>> result = Err(error(ErrorType.INTERNAL, code("ERR"), "failed"))
        >>> result = map_value(result, lambda x: x * 2)
        >>> assert is_err(result)
    """
    match result:
        case Ok(value):
            return Ok(func(value))
        case Err(error):
            return Err(error)


def map_error(result: Result[T, E], func: Callable[[E], U]) -> Result[T, U]:
    """
    Transform the error inside an Err result.
    
    If the result is Ok, it's returned unchanged.
    
    Args:
        result: The result to transform
        func: Function to apply to the Err error
        
    Returns:
        The original Ok, or Err with transformed error
        
    Example:
        >>> err = error(ErrorType.INTERNAL, code("ERR"), "failed")
        >>> result = Err(err)
        >>> result = map_error(result, lambda e: e.with_detail("context", "test"))
        >>> assert unwrap_err(result).has_detail("context")
    """
    match result:
        case Ok(value):
            return Ok(value)
        case Err(error):
            return Err(func(error))


def and_then(result: Result[T, E], func: Callable[[T], Result[U, E]]) -> Result[U, E]:
    """
    Chain operations that can fail.
    
    Also known as flatMap or bind in functional programming.
    If the result is Ok, applies func to the value and returns the result.
    If the result is Err, returns the error unchanged.
    
    Args:
        result: The result to chain from
        func: Function that takes the Ok value and returns a new Result
        
    Returns:
        The result of func if Ok, otherwise the original Err
        
    Example:
        >>> def validate_positive(x: int) -> Result[int, Error]:
        ...     if x > 0:
        ...         return Ok(x)
        ...     return Err(error(ErrorType.VALIDATION, code("NOT_POSITIVE"), "must be positive"))
        >>> 
        >>> result = Ok(5)
        >>> result = and_then(result, validate_positive)
        >>> assert is_ok(result)
        >>> 
        >>> result = Ok(-5)
        >>> result = and_then(result, validate_positive)
        >>> assert is_err(result)
    """
    match result:
        case Ok(value):
            return func(value)
        case Err(error):
            return Err(error)


def or_else(result: Result[T, E], func: Callable[[E], Result[T, U]]) -> Result[T, U]:
    """
    Handle errors by providing an alternative computation.
    
    If the result is Err, applies func to the error and returns the result.
    If the result is Ok, returns the value unchanged.
    
    Args:
        result: The result to handle
        func: Function that takes the Err and returns a new Result
        
    Returns:
        The original Ok, or the result of func if Err
        
    Example:
        >>> def fallback(e: Error) -> Result[int, Error]:
        ...     return Ok(0)  # Provide default value
        >>> 
        >>> result = Err(error(ErrorType.INTERNAL, code("ERR"), "failed"))
        >>> result = or_else(result, fallback)
        >>> assert unwrap(result) == 0
    """
    match result:
        case Ok(value):
            return Ok(value)
        case Err(error):
            return func(error)


def flatten(result: Result[Result[T, E], E]) -> Result[T, E]:
    """
    Flatten a nested Result.
    
    Converts Result[Result[T, E], E] to Result[T, E].
    
    Args:
        result: The nested result to flatten
        
    Returns:
        The flattened result
        
    Example:
        >>> inner = Ok(42)
        >>> outer = Ok(inner)
        >>> flattened = flatten(outer)
        >>> assert unwrap(flattened) == 42
    """
    match result:
        case Ok(inner_result):
            return inner_result
        case Err(error):
            return Err(error)


# Inspection functions
def inspect(result: Result[T, E], func: Callable[[T], None]) -> Result[T, E]:
    """
    Call a function with the Ok value for side effects, without modifying the result.
    
    Useful for logging or debugging.
    
    Args:
        result: The result to inspect
        func: Function to call with the Ok value (for side effects only)
        
    Returns:
        The original result unchanged
        
    Example:
        >>> result = Ok(42)
        >>> result = inspect(result, lambda x: print(f"Got value: {x}"))
        >>> assert is_ok(result)
    """
    match result:
        case Ok(value):
            func(value)
            return result
        case Err(_):
            return result


def inspect_err(result: Result[T, E], func: Callable[[E], None]) -> Result[T, E]:
    """
    Call a function with the Err error for side effects, without modifying the result.
    
    Useful for logging or debugging errors.
    
    Args:
        result: The result to inspect
        func: Function to call with the Err error (for side effects only)
        
    Returns:
        The original result unchanged
        
    Example:
        >>> err = error(ErrorType.INTERNAL, code("ERR"), "failed")
        >>> result = Err(err)
        >>> result = inspect_err(result, lambda e: print(f"Error: {e}"))
        >>> assert is_err(result)
    """
    match result:
        case Ok(_):
            return result
        case Err(error):
            func(error)
            return result


# Collection functions
def collect(results: list[Result[T, E]]) -> Result[list[T], E]:
    """
    Convert a list of Results into a Result of a list.
    
    If all results are Ok, returns Ok with a list of all values.
    If any result is Err, returns the first Err encountered.
    
    Args:
        results: List of Results to collect
        
    Returns:
        Ok with list of values if all Ok, otherwise the first Err
        
    Example:
        >>> results = [Ok(1), Ok(2), Ok(3)]
        >>> collected = collect(results)
        >>> assert unwrap(collected) == [1, 2, 3]
        >>> 
        >>> results = [Ok(1), Err(error(ErrorType.INTERNAL, code("ERR"), "failed")), Ok(3)]
        >>> collected = collect(results)
        >>> assert is_err(collected)
    """
    values = []
    for result in results:
        match result:
            case Ok(value):
                values.append(value)
            case Err(error):
                return Err(error)
    return Ok(values)


def partition(results: list[Result[T, E]]) -> tuple[list[T], list[E]]:
    """
    Partition a list of Results into separate lists of values and errors.
    
    Args:
        results: List of Results to partition
        
    Returns:
        Tuple of (values, errors)
        
    Example:
        >>> results = [Ok(1), Err(error(ErrorType.INTERNAL, code("E1"), "e1")), Ok(3)]
        >>> values, errors = partition(results)
        >>> assert values == [1, 3]
        >>> assert len(errors) == 1
    """
    values = []
    errors = []
    
    for result in results:
        match result:
            case Ok(value):
                values.append(value)
            case Err(error):
                errors.append(error)
    
    return values, errors