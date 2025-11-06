"""Tests for error handling decorators."""

import pytest
import asyncio
import time
from typing import List

from .decorators import (
    safe,
    safe_async,
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
from .result import Ok, Err, Result, unwrap, is_ok, is_err, unwrap_err
from .base import Error
from .types import ErrorType
from .codes import code
from .constructors import validation, timeout as timeout_error


class TestSafeDecorator:
    """Test @safe decorator."""
    
    def test_safe_success(self):
        """Test @safe with successful execution."""
        @safe
        def divide(a: int, b: int) -> float:
            return a / b
        
        result = divide(10, 2)
        assert is_ok(result)
        assert unwrap(result) == 5.0
    
    def test_safe_exception(self):
        """Test @safe with exception."""
        @safe
        def divide(a: int, b: int) -> float:
            return a / b
        
        result = divide(10, 0)
        assert is_err(result)
        error = unwrap_err(result)
        assert error.error_type == ErrorType.INTERNAL
        assert "divide" in error.message
        assert error.has_detail("exception_type")
        assert error.get_detail("exception_type") == "ZeroDivisionError"
    
    def test_safe_preserves_metadata(self):
        """Test @safe preserves function metadata."""
        @safe
        def my_function(x: int) -> int:
            """My docstring."""
            return x * 2
        
        assert my_function.__name__ == "my_function"
        assert my_function.__doc__ == "My docstring."
    
    def test_safe_with_multiple_exceptions(self):
        """Test @safe captures different exception types."""
        @safe
        def risky(mode: str) -> int:
            if mode == "value":
                raise ValueError("value error")
            elif mode == "type":
                raise TypeError("type error")
            elif mode == "key":
                raise KeyError("key error")
            return 42
        
        for mode in ["value", "type", "key"]:
            result = risky(mode)
            assert is_err(result)
            error = unwrap_err(result)
            assert error.has_detail("exception_type")


class TestSafeAsyncDecorator:
    """Test @safe_async decorator."""
    
    @pytest.mark.asyncio
    async def test_safe_async_success(self):
        """Test @safe_async with successful execution."""
        @safe_async
        async def fetch_data(value: int) -> int:
            await asyncio.sleep(0.001)
            return value * 2
        
        result = await fetch_data(21)
        assert is_ok(result)
        assert unwrap(result) == 42
    
    @pytest.mark.asyncio
    async def test_safe_async_exception(self):
        """Test @safe_async with exception."""
        @safe_async
        async def fetch_data(value: int) -> int:
            await asyncio.sleep(0.001)
            raise ValueError("fetch failed")
        
        result = await fetch_data(21)
        assert is_err(result)
        error = unwrap_err(result)
        assert error.error_type == ErrorType.INTERNAL
        assert "fetch_data" in error.message
    
    @pytest.mark.asyncio
    async def test_safe_async_preserves_metadata(self):
        """Test @safe_async preserves function metadata."""
        @safe_async
        async def my_async_function(x: int) -> int:
            """Async docstring."""
            return x * 2
        
        assert my_async_function.__name__ == "my_async_function"
        assert my_async_function.__doc__ == "Async docstring."


class TestRetryDecorator:
    """Test @retry decorator."""
    
    def test_retry_success_first_attempt(self):
        """Test @retry when first attempt succeeds."""
        attempts = []
        
        @retry(max_attempts=3)
        def operation() -> Result[int, Error]:
            attempts.append(1)
            return Ok(42)
        
        result = operation()
        assert is_ok(result)
        assert unwrap(result) == 42
        assert len(attempts) == 1
    
    def test_retry_success_after_retries(self):
        """Test @retry when operation succeeds after retries."""
        attempts = []
        
        @retry(max_attempts=3, backoff_ms=10)
        def operation() -> Result[int, Error]:
            attempts.append(1)
            if len(attempts) < 3:
                return Err(timeout_error("operation", 1000))
            return Ok(42)
        
        result = operation()
        assert is_ok(result)
        assert unwrap(result) == 42
        assert len(attempts) == 3
    
    def test_retry_all_attempts_fail(self):
        """Test @retry when all attempts fail."""
        attempts = []
        
        @retry(max_attempts=3, backoff_ms=10)
        def operation() -> Result[int, Error]:
            attempts.append(1)
            return Err(timeout_error("operation", 1000))
        
        result = operation()
        assert is_err(result)
        assert len(attempts) == 3
        
        error = unwrap_err(result)
        assert error.has_detail("retry_attempts")
        assert error.get_detail("retry_attempts") == "3"
    
    def test_retry_non_retryable_error(self):
        """Test @retry stops on non-retryable error."""
        attempts = []
        
        @retry(max_attempts=3)
        def operation() -> Result[int, Error]:
            attempts.append(1)
            return Err(validation("invalid input"))
        
        result = operation()
        assert is_err(result)
        assert len(attempts) == 1  # No retries for validation error
    
    def test_retry_custom_predicate(self):
        """Test @retry with custom retry predicate."""
        attempts = []
        
        @retry(
            max_attempts=3,
            retry_on=lambda e: e.code == code("RETRY_ME"),
            backoff_ms=10
        )
        def operation() -> Result[int, Error]:
            attempts.append(1)
            if len(attempts) < 3:
                return Err(Error(ErrorType.INTERNAL, code("RETRY_ME"), "retry"))
            return Ok(42)
        
        result = operation()
        assert is_ok(result)
        assert len(attempts) == 3
    
    def test_retry_exponential_backoff(self):
        """Test @retry uses exponential backoff."""
        start_time = time.time()
        attempts = []
        
        @retry(max_attempts=3, backoff_ms=50, exponential=True)
        def operation() -> Result[int, Error]:
            attempts.append(time.time())
            if len(attempts) < 3:
                return Err(timeout_error("op", 1000))
            return Ok(42)
        
        result = operation()
        duration = time.time() - start_time
        
        # Exponential backoff: 50ms, 100ms = ~150ms total
        # Add some tolerance for timing variations
        assert duration >= 0.1  # At least 100ms
        assert is_ok(result)


class TestRetryAsyncDecorator:
    """Test @retry_async decorator."""
    
    @pytest.mark.asyncio
    async def test_retry_async_success_first_attempt(self):
        """Test @retry_async when first attempt succeeds."""
        attempts = []
        
        @retry_async(max_attempts=3)
        async def operation() -> Result[int, Error]:
            attempts.append(1)
            await asyncio.sleep(0.001)
            return Ok(42)
        
        result = await operation()
        assert is_ok(result)
        assert unwrap(result) == 42
        assert len(attempts) == 1
    
    @pytest.mark.asyncio
    async def test_retry_async_success_after_retries(self):
        """Test @retry_async when operation succeeds after retries."""
        attempts = []
        
        @retry_async(max_attempts=3, backoff_ms=10)
        async def operation() -> Result[int, Error]:
            attempts.append(1)
            await asyncio.sleep(0.001)
            if len(attempts) < 3:
                return Err(timeout_error("operation", 1000))
            return Ok(42)
        
        result = await operation()
        assert is_ok(result)
        assert len(attempts) == 3
    
    @pytest.mark.asyncio
    async def test_retry_async_all_attempts_fail(self):
        """Test @retry_async when all attempts fail."""
        attempts = []
        
        @retry_async(max_attempts=3, backoff_ms=10)
        async def operation() -> Result[int, Error]:
            attempts.append(1)
            return Err(timeout_error("operation", 1000))
        
        result = await operation()
        assert is_err(result)
        assert len(attempts) == 3


class TestWithTimeoutDecorator:
    """Test @with_timeout decorator."""
    
    def test_with_timeout_completes_in_time(self):
        """Test @with_timeout when operation completes in time."""
        @with_timeout(1000)  # 1 second
        def operation() -> Result[int, Error]:
            time.sleep(0.01)  # 10ms
            return Ok(42)
        
        result = operation()
        assert is_ok(result)
        assert unwrap(result) == 42
    
    def test_with_timeout_exceeds_timeout(self):
        """Test @with_timeout when operation exceeds timeout."""
        @with_timeout(50)  # 50ms
        def operation() -> Result[int, Error]:
            time.sleep(0.2)  # 200ms
            return Ok(42)
        
        result = operation()
        assert is_err(result)
        error = unwrap_err(result)
        assert error.error_type == ErrorType.TIMEOUT


class TestWithTimeoutAsyncDecorator:
    """Test @with_timeout_async decorator."""
    
    @pytest.mark.asyncio
    async def test_with_timeout_async_completes_in_time(self):
        """Test @with_timeout_async when operation completes in time."""
        @with_timeout_async(1000)  # 1 second
        async def operation() -> Result[int, Error]:
            await asyncio.sleep(0.01)  # 10ms
            return Ok(42)
        
        result = await operation()
        assert is_ok(result)
        assert unwrap(result) == 42
    
    @pytest.mark.asyncio
    async def test_with_timeout_async_exceeds_timeout(self):
        """Test @with_timeout_async when operation exceeds timeout."""
        @with_timeout_async(50)  # 50ms
        async def operation() -> Result[int, Error]:
            await asyncio.sleep(0.2)  # 200ms
            return Ok(42)
        
        result = await operation()
        assert is_err(result)
        error = unwrap_err(result)
        assert error.error_type == ErrorType.TIMEOUT


class TestValidateArgsDecorator:
    """Test @validate_args decorator."""
    
    def test_validate_args_success(self):
        """Test @validate_args with valid arguments."""
        def validate_positive(x: int) -> Result[int, Error]:
            if x > 0:
                return Ok(x)
            return Err(validation("must be positive"))
        
        @validate_args(validate_positive)
        def operation(x: int) -> Result[int, Error]:
            return Ok(x * 2)
        
        result = operation(5)
        assert is_ok(result)
        assert unwrap(result) == 10
    
    def test_validate_args_validation_fails(self):
        """Test @validate_args when validation fails."""
        def validate_positive(x: int) -> Result[int, Error]:
            if x > 0:
                return Ok(x)
            return Err(validation("must be positive"))
        
        @validate_args(validate_positive)
        def operation(x: int) -> Result[int, Error]:
            return Ok(x * 2)
        
        result = operation(-5)
        assert is_err(result)
        error = unwrap_err(result)
        assert "must be positive" in error.message
    
    def test_validate_args_multiple_validators(self):
        """Test @validate_args with multiple validators."""
        def validate_positive(x: int) -> Result[int, Error]:
            if x > 0:
                return Ok(x)
            return Err(validation("x must be positive"))
        
        def validate_non_zero(y: int) -> Result[int, Error]:
            if y != 0:
                return Ok(y)
            return Err(validation("y must be non-zero"))
        
        @validate_args(validate_positive, validate_non_zero)
        def divide(x: int, y: int) -> Result[float, Error]:
            return Ok(x / y)
        
        # Valid case
        result = divide(10, 2)
        assert is_ok(result)
        assert unwrap(result) == 5.0
        
        # Invalid x
        result = divide(-10, 2)
        assert is_err(result)
        
        # Invalid y
        result = divide(10, 0)
        assert is_err(result)


class TestWithErrorContextDecorator:
    """Test @with_error_context decorator."""
    
    def test_with_error_context_success(self):
        """Test @with_error_context with successful operation."""
        @with_error_context(
            ErrorType.DATABASE,
            "DB_OPERATION_FAILED",
            "database operation failed"
        )
        def operation() -> Result[int, Error]:
            return Ok(42)
        
        result = operation()
        assert is_ok(result)
        assert unwrap(result) == 42
    
    def test_with_error_context_wraps_error(self):
        """Test @with_error_context wraps errors."""
        @with_error_context(
            ErrorType.DATABASE,
            "DB_OPERATION_FAILED",
            "database operation failed"
        )
        def operation() -> Result[int, Error]:
            return Err(timeout_error("connection timeout", 5000))
        
        result = operation()
        assert is_err(result)
        
        error = unwrap_err(result)
        assert error.error_type == ErrorType.DATABASE
        assert error.code == code("DB_OPERATION_FAILED")
        assert error.message == "database operation failed"
        assert error.cause is not None
        assert error.cause.error_type == ErrorType.TIMEOUT


class TestWithErrorContextAsyncDecorator:
    """Test @with_error_context_async decorator."""
    
    @pytest.mark.asyncio
    async def test_with_error_context_async_wraps_error(self):
        """Test @with_error_context_async wraps errors."""
        @with_error_context_async(
            ErrorType.DATABASE,
            "DB_OPERATION_FAILED",
            "database operation failed"
        )
        async def operation() -> Result[int, Error]:
            await asyncio.sleep(0.001)
            return Err(timeout_error("connection timeout", 5000))
        
        result = await operation()
        assert is_err(result)
        
        error = unwrap_err(result)
        assert error.error_type == ErrorType.DATABASE
        assert error.cause is not None


class TestWithFallbackDecorator:
    """Test @with_fallback decorator."""
    
    def test_with_fallback_success(self):
        """Test @with_fallback with successful operation."""
        @with_fallback(0)
        def operation() -> Result[int, Error]:
            return Ok(42)
        
        result = operation()
        assert result == 42
    
    def test_with_fallback_returns_default(self):
        """Test @with_fallback returns default on error."""
        @with_fallback(0)
        def operation() -> Result[int, Error]:
            return Err(validation("error"))
        
        result = operation()
        assert result == 0
    
    def test_with_fallback_type_conversion(self):
        """Test @with_fallback converts Result to value."""
        @with_fallback("default")
        def operation() -> Result[str, Error]:
            return Ok("success")
        
        # Return type is now str, not Result
        result = operation()
        assert isinstance(result, str)
        assert result == "success"


class TestLogErrorsDecorator:
    """Test @log_errors decorator."""
    
    def test_log_errors_success_no_log(self):
        """Test @log_errors doesn't log on success."""
        logs: List[str] = []
        
        class FakeLogger:
            def error(self, msg):
                logs.append(msg)
        
        logger = FakeLogger()
        
        @log_errors(logger)
        def operation() -> Result[int, Error]:
            return Ok(42)
        
        result = operation()
        assert is_ok(result)
        assert len(logs) == 0
    
    def test_log_errors_logs_on_error(self):
        """Test @log_errors logs on error."""
        logs: List[str] = []
        
        class FakeLogger:
            def error(self, msg):
                logs.append(msg)
        
        logger = FakeLogger()
        
        @log_errors(logger)
        def operation() -> Result[int, Error]:
            return Err(validation("invalid input"))
        
        result = operation()
        assert is_err(result)
        assert len(logs) == 1
        assert "operation failed" in logs[0]
    
    def test_log_errors_returns_result_unchanged(self):
        """Test @log_errors returns result unchanged."""
        class FakeLogger:
            def error(self, msg):
                pass
        
        logger = FakeLogger()
        
        @log_errors(logger)
        def operation() -> Result[int, Error]:
            return Err(validation("invalid"))
        
        result = operation()
        assert is_err(result)
        error = unwrap_err(result)
        assert error.error_type == ErrorType.VALIDATION


# =========================================
# Integration Tests
# =========================================

class TestDecoratorIntegration:
    """Integration tests combining multiple decorators."""
    
    def test_stacked_decorators(self):
        """Test stacking multiple decorators."""
        attempts = []
        
        @with_error_context(
            ErrorType.INTERNAL,
            "OPERATION_FAILED",
            "operation failed"
        )
        @retry(max_attempts=3, backoff_ms=10)
        def operation() -> Result[int, Error]:
            attempts.append(1)
            if len(attempts) < 3:
                return Err(timeout_error("timeout", 1000))
            return Ok(42)
        
        result = operation()
        assert is_ok(result)
        assert len(attempts) == 3
    
    @pytest.mark.asyncio
    async def test_async_stacked_decorators(self):
        """Test stacking async decorators."""
        attempts = []
        
        @with_error_context_async(
            ErrorType.INTERNAL,
            "OPERATION_FAILED",
            "operation failed"
        )
        @retry_async(max_attempts=3, backoff_ms=10)
        @with_timeout_async(1000)
        async def operation() -> Result[int, Error]:
            attempts.append(1)
            await asyncio.sleep(0.001)
            if len(attempts) < 3:
                return Err(timeout_error("timeout", 1000))
            return Ok(42)
        
        result = await operation()
        assert is_ok(result)
        assert len(attempts) == 3


# =========================================
# Benchmarks
# =========================================

class TestDecoratorBenchmarks:
    """Benchmark decorator operations."""
    
    def test_safe_decorator_benchmark(self, benchmark):
        """Benchmark @safe decorator."""
        @safe
        def operation(x: int) -> int:
            return x * 2
        
        result = benchmark(operation, 21)
        assert unwrap(result) == 42
    
    def test_retry_decorator_benchmark(self, benchmark):
        """Benchmark @retry decorator (no retries)."""
        @retry(max_attempts=3)
        def operation() -> Result[int, Error]:
            return Ok(42)
        
        result = benchmark(operation)
        assert unwrap(result) == 42
    
    def test_with_error_context_benchmark(self, benchmark):
        """Benchmark @with_error_context decorator."""
        @with_error_context(
            ErrorType.INTERNAL,
            "TEST",
            "test"
        )
        def operation() -> Result[int, Error]:
            return Ok(42)
        
        result = benchmark(operation)
        assert unwrap(result) == 42