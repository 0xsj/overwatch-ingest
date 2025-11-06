"""Decorators for error handling in Scout platform."""

import functools
import asyncio
from typing import TypeVar, Callable, ParamSpec, Awaitable, Any
from collections.abc import Coroutine

from .result import Result, Ok, Err
from .base import Error
from .constructors import (
    internal,
    timeout as timeout_error,
    validation,
    wrap,
)
from .types import ErrorType
from .codes import code


T = TypeVar("T")
P = ParamSpec("P")


# =========================================
# Safe Execution Decorators
# =========================================

def safe(func: Callable[P, T]) -> Callable[P, Result[T, Error]]:
    """
    Decorator that converts exception-raising functions to Result-returning.
    
    Catches all exceptions and wraps them as Error instances.
    
    Args:
        func: Function that might raise exceptions
        
    Returns:
        Function that returns Result[T, Error]
        
    Example:
        >>> @safe
        ... def parse_int(s: str) -> int:
        ...     return int(s)
        >>> 
        >>> result = parse_int("42")
        >>> assert result == Ok(42)
        >>> 
        >>> result = parse_int("invalid")
        >>> assert isinstance(result, Err)
    """
    @functools.wraps(func)
    def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
        try:
            value = func(*args, **kwargs)
            return Ok(value)
        except Exception as e:
            error = internal(f"{func.__name__} failed: {str(e)}")
            return Err(error.with_detail("exception_type", type(e).__name__))
    
    return wrapper


def safe_async(func: Callable[P, Awaitable[T]]) -> Callable[P, Awaitable[Result[T, Error]]]:
    """
    Async version of @safe decorator.
    
    Converts async exception-raising functions to Result-returning.
    
    Args:
        func: Async function that might raise exceptions
        
    Returns:
        Async function that returns Result[T, Error]
        
    Example:
        >>> @safe_async
        ... async def fetch_data(url: str) -> dict:
        ...     async with httpx.AsyncClient() as client:
        ...         response = await client.get(url)
        ...         return response.json()
        >>> 
        >>> result = await fetch_data("https://api.example.com/data")
    """
    @functools.wraps(func)
    async def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
        try:
            value = await func(*args, **kwargs)
            return Ok(value)
        except Exception as e:
            error = internal(f"{func.__name__} failed: {str(e)}")
            return Err(error.with_detail("exception_type", type(e).__name__))
    
    return wrapper


# =========================================
# Retry Decorators
# =========================================

def retry(
    max_attempts: int = 3,
    *,
    retry_on: Callable[[Error], bool] | None = None,
    backoff_ms: int = 100,
    exponential: bool = True,
) -> Callable[[Callable[P, Result[T, Error]]], Callable[P, Result[T, Error]]]:
    """
    Decorator that retries a function on failure.
    
    Args:
        max_attempts: Maximum number of attempts (default: 3)
        retry_on: Predicate to determine if error should be retried (default: retryable errors)
        backoff_ms: Initial backoff in milliseconds (default: 100)
        exponential: Use exponential backoff (default: True)
        
    Returns:
        Decorator function
        
    Example:
        >>> @retry(max_attempts=3, backoff_ms=100)
        ... def fetch_user(user_id: str) -> Result[User, Error]:
        ...     # Might fail with network error
        ...     return get_from_api(user_id)
        >>> 
        >>> # Custom retry predicate
        >>> @retry(max_attempts=5, retry_on=lambda e: e.error_type == ErrorType.TIMEOUT)
        ... def slow_operation() -> Result[Data, Error]:
        ...     return process()
    """
    import time
    
    if retry_on is None:
        retry_on = lambda e: e.is_retryable()
    
    def decorator(func: Callable[P, Result[T, Error]]) -> Callable[P, Result[T, Error]]:
        @functools.wraps(func)
        def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
            last_error: Error | None = None
            
            for attempt in range(max_attempts):
                result = func(*args, **kwargs)
                
                match result:
                    case Ok(_):
                        return result
                    case Err(error):
                        last_error = error
                        
                        # Check if we should retry
                        if not retry_on(error):
                            return result
                        
                        # Don't sleep on last attempt
                        if attempt < max_attempts - 1:
                            if exponential:
                                sleep_ms = backoff_ms * (2 ** attempt)
                            else:
                                sleep_ms = backoff_ms
                            
                            time.sleep(sleep_ms / 1000.0)
            
            # All attempts failed
            assert last_error is not None
            return Err(
                last_error.with_detail("retry_attempts", str(max_attempts))
            )
        
        return wrapper
    
    return decorator


def retry_async(
    max_attempts: int = 3,
    *,
    retry_on: Callable[[Error], bool] | None = None,
    backoff_ms: int = 100,
    exponential: bool = True,
) -> Callable[[Callable[P, Awaitable[Result[T, Error]]]], Callable[P, Awaitable[Result[T, Error]]]]:
    """
    Async version of @retry decorator.
    
    Args:
        max_attempts: Maximum number of attempts (default: 3)
        retry_on: Predicate to determine if error should be retried
        backoff_ms: Initial backoff in milliseconds (default: 100)
        exponential: Use exponential backoff (default: True)
        
    Returns:
        Async decorator function
        
    Example:
        >>> @retry_async(max_attempts=3)
        ... async def fetch_user(user_id: str) -> Result[User, Error]:
        ...     return await api_client.get_user(user_id)
    """
    if retry_on is None:
        retry_on = lambda e: e.is_retryable()
    
    def decorator(
        func: Callable[P, Awaitable[Result[T, Error]]]
    ) -> Callable[P, Awaitable[Result[T, Error]]]:
        @functools.wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
            last_error: Error | None = None
            
            for attempt in range(max_attempts):
                result = await func(*args, **kwargs)
                
                match result:
                    case Ok(_):
                        return result
                    case Err(error):
                        last_error = error
                        
                        if not retry_on(error):
                            return result
                        
                        if attempt < max_attempts - 1:
                            if exponential:
                                sleep_ms = backoff_ms * (2 ** attempt)
                            else:
                                sleep_ms = backoff_ms
                            
                            await asyncio.sleep(sleep_ms / 1000.0)
            
            assert last_error is not None
            return Err(
                last_error.with_detail("retry_attempts", str(max_attempts))
            )
        
        return wrapper
    
    return decorator


# =========================================
# Timeout Decorators
# =========================================

def with_timeout(
    timeout_ms: int,
) -> Callable[[Callable[P, Result[T, Error]]], Callable[P, Result[T, Error]]]:
    """
    Decorator that enforces a timeout on synchronous functions.
    
    Note: This uses threading and may not work for all types of operations.
    For async functions, use with_timeout_async instead.
    
    Args:
        timeout_ms: Timeout in milliseconds
        
    Returns:
        Decorator function
        
    Example:
        >>> @with_timeout(5000)  # 5 second timeout
        ... def slow_operation() -> Result[Data, Error]:
        ...     # Long-running operation
        ...     return process_data()
    """
    import concurrent.futures
    
    def decorator(func: Callable[P, Result[T, Error]]) -> Callable[P, Result[T, Error]]:
        @functools.wraps(func)
        def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
            with concurrent.futures.ThreadPoolExecutor(max_workers=1) as executor:
                future = executor.submit(func, *args, **kwargs)
                try:
                    result = future.result(timeout=timeout_ms / 1000.0)
                    return result
                except concurrent.futures.TimeoutError:
                    return Err(timeout_error(func.__name__, timeout_ms))
        
        return wrapper
    
    return decorator


def with_timeout_async(
    timeout_ms: int,
) -> Callable[[Callable[P, Awaitable[Result[T, Error]]]], Callable[P, Awaitable[Result[T, Error]]]]:
    """
    Decorator that enforces a timeout on async functions.
    
    Args:
        timeout_ms: Timeout in milliseconds
        
    Returns:
        Async decorator function
        
    Example:
        >>> @with_timeout_async(5000)
        ... async def fetch_data(url: str) -> Result[dict, Error]:
        ...     async with httpx.AsyncClient() as client:
        ...         response = await client.get(url)
        ...         return Ok(response.json())
    """
    def decorator(
        func: Callable[P, Awaitable[Result[T, Error]]]
    ) -> Callable[P, Awaitable[Result[T, Error]]]:
        @functools.wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
            try:
                result = await asyncio.wait_for(
                    func(*args, **kwargs),
                    timeout=timeout_ms / 1000.0
                )
                return result
            except asyncio.TimeoutError:
                return Err(timeout_error(func.__name__, timeout_ms))
        
        return wrapper
    
    return decorator


# =========================================
# Validation Decorators
# =========================================

def validate_args(
    *validators: Callable[[Any], Result[Any, Error]]
) -> Callable[[Callable[P, Result[T, Error]]], Callable[P, Result[T, Error]]]:
    """
    Decorator that validates function arguments before execution.
    
    Args:
        *validators: Validator functions (one per argument)
        
    Returns:
        Decorator function
        
    Example:
        >>> def validate_email(email: str) -> Result[str, Error]:
        ...     if "@" not in email:
        ...         return Err(validation("invalid email format"))
        ...     return Ok(email)
        >>> 
        >>> @validate_args(validate_email)
        ... def send_email(email: str) -> Result[None, Error]:
        ...     # Email is already validated
        ...     return Ok(None)
    """
    def decorator(func: Callable[P, Result[T, Error]]) -> Callable[P, Result[T, Error]]:
        @functools.wraps(func)
        def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
            # Validate positional arguments
            if len(args) != len(validators):
                return Err(
                    validation(
                        f"expected {len(validators)} arguments, got {len(args)}"
                    )
                )
            
            validated_args = []
            for arg, validator in zip(args, validators):
                result = validator(arg)
                match result:
                    case Ok(value):
                        validated_args.append(value)
                    case Err(error):
                        return Err(error)
            
            return func(*validated_args, **kwargs)  # type: ignore
        
        return wrapper
    
    return decorator


# =========================================
# Context Decorators
# =========================================

def with_error_context(
    error_type: ErrorType,
    error_code: str,
    message: str,
) -> Callable[[Callable[P, Result[T, Error]]], Callable[P, Result[T, Error]]]:
    """
    Decorator that wraps errors with additional context.
    
    Args:
        error_type: Error type for the wrapper
        error_code: Error code for the wrapper
        message: Message for the wrapper
        
    Returns:
        Decorator function
        
    Example:
        >>> @with_error_context(
        ...     ErrorType.DATABASE,
        ...     "USER_FETCH_FAILED",
        ...     "failed to fetch user from database"
        ... )
        ... def get_user_from_db(user_id: str) -> Result[User, Error]:
        ...     return db.query(user_id)
        >>> 
        >>> # If db.query fails, error will be wrapped with context
    """
    def decorator(func: Callable[P, Result[T, Error]]) -> Callable[P, Result[T, Error]]:
        @functools.wraps(func)
        def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
            result = func(*args, **kwargs)
            match result:
                case Ok(_):
                    return result
                case Err(error):
                    wrapped = wrap(error, error_type, code(error_code), message)
                    return Err(wrapped)
        
        return wrapper
    
    return decorator


def with_error_context_async(
    error_type: ErrorType,
    error_code: str,
    message: str,
) -> Callable[[Callable[P, Awaitable[Result[T, Error]]]], Callable[P, Awaitable[Result[T, Error]]]]:
    """
    Async version of @with_error_context decorator.
    
    Args:
        error_type: Error type for the wrapper
        error_code: Error code for the wrapper
        message: Message for the wrapper
        
    Returns:
        Async decorator function
    """
    def decorator(
        func: Callable[P, Awaitable[Result[T, Error]]]
    ) -> Callable[P, Awaitable[Result[T, Error]]]:
        @functools.wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
            result = await func(*args, **kwargs)
            match result:
                case Ok(_):
                    return result
                case Err(error):
                    wrapped = wrap(error, error_type, code(error_code), message)
                    return Err(wrapped)
        
        return wrapper
    
    return decorator


# =========================================
# Fallback Decorators
# =========================================

def with_fallback(
    fallback_value: T,
) -> Callable[[Callable[P, Result[T, Error]]], Callable[P, T]]:
    """
    Decorator that provides a fallback value on error.
    
    Converts Result-returning function to regular function with fallback.
    
    Args:
        fallback_value: Value to return on error
        
    Returns:
        Decorator function
        
    Example:
        >>> @with_fallback(default_user)
        ... def get_user(user_id: str) -> Result[User, Error]:
        ...     return fetch_user(user_id)
        >>> 
        >>> user = get_user("123")  # Returns User (never Error)
    """
    def decorator(func: Callable[P, Result[T, Error]]) -> Callable[P, T]:
        @functools.wraps(func)
        def wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
            result = func(*args, **kwargs)
            match result:
                case Ok(value):
                    return value
                case Err(_):
                    return fallback_value
        
        return wrapper
    
    return decorator


# =========================================
# Logging Decorators
# =========================================

def log_errors(
    logger: Any = None,
    *,
    level: str = "error",
) -> Callable[[Callable[P, Result[T, Error]]], Callable[P, Result[T, Error]]]:
    """
    Decorator that logs errors without modifying the result.
    
    Args:
        logger: Logger instance (uses print if None)
        level: Log level ("error", "warning", "info")
        
    Returns:
        Decorator function
        
    Example:
        >>> import logging
        >>> logger = logging.getLogger(__name__)
        >>> 
        >>> @log_errors(logger, level="warning")
        ... def risky_operation() -> Result[Data, Error]:
        ...     return process()
    """
    def decorator(func: Callable[P, Result[T, Error]]) -> Callable[P, Result[T, Error]]:
        @functools.wraps(func)
        def wrapper(*args: P.args, **kwargs: P.kwargs) -> Result[T, Error]:
            result = func(*args, **kwargs)
            
            match result:
                case Err(error):
                    log_message = f"{func.__name__} failed: {error}"
                    
                    if logger is None:
                        print(f"[{level.upper()}] {log_message}")
                    else:
                        log_fn = getattr(logger, level, logger.error)
                        log_fn(log_message)
            
            return result
        
        return wrapper
    
    return decorator