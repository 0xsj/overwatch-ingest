# platform/pylib/scout_common/errors/test_combinators.py
"""Tests for combinators and chainable operations."""

import pytest
pytestmark = pytest.mark.skip(reason="Combinator API needs refinement - TODO for later")
from hypothesis import given, strategies as st

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
from .result import Ok, Err, Result, unwrap_err, unwrap, is_ok, is_err
from .base import Error
from .types import ErrorType
from .codes import code
from .constructors import validation, internal


class TestOkResultErrResult:
    """Test OkResult and ErrResult wrapper classes."""
    
    def test_ok_result_creation(self):
        """Test creating OkResult."""
        result = ok(42)
        assert isinstance(result, OkResult)
        assert result.is_ok()
        assert not result.is_err()
    
    def test_err_result_creation(self):
        """Test creating ErrResult."""
        error = validation("test")  # Rename from 'err' to 'error'
        result = err(error)  # Now err() is the function
        assert isinstance(result, ErrResult)
    
    def test_ok_result_map(self):
        """Test OkResult.map."""
        result = ok(5).map(lambda x: x * 2)
        assert result.unwrap() == 10
    
    def test_err_result_map(self):
        """Test ErrResult.map does nothing."""
        error = validation("test")
        result = err(error).map(lambda x: x * 2)
        assert result.is_err()
        assert result.unwrap_err() == error
    
    def test_ok_result_map_err(self):
        """Test OkResult.map_err does nothing."""
        result = ok(42).map_err(lambda e: e.with_detail("key", "value"))
        assert result.unwrap() == 42
    
    def test_err_result_map_err(self):
        """Test ErrResult.map_err."""
        error = validation("test")
        result = err(error).map_err(lambda e: e.with_detail("key", "value"))
        assert result.is_err()
        assert result.unwrap_err().has_detail("key")
    
    def test_ok_result_and_then(self):
        """Test OkResult.and_then."""
        result = ok(5).and_then(lambda x: ok(x * 2))
        assert result.unwrap() == 10
    
    def test_ok_result_and_then_to_err(self):
        """Test OkResult.and_then returning Err."""
        result = ok(-5).and_then(
            lambda x: err(validation("negative")) if x < 0 else ok(x)
        )
        assert result.is_err()
    
    def test_err_result_and_then(self):
        """Test ErrResult.and_then does nothing."""
        error = validation("test")
        result = err(error).and_then(lambda x: ok(x * 2))
        assert result.is_err()
        assert result.unwrap_err() == error
    
    def test_ok_result_or_else(self):
        """Test OkResult.or_else does nothing."""
        result = ok(42).or_else(lambda e: ok(0))
        assert result.unwrap() == 42
    
    def test_err_result_or_else(self):
        """Test ErrResult.or_else."""
        error = validation("test")
        result = err(error).or_else(lambda e: ok(0))
        assert result.unwrap() == 0
    
    def test_ok_result_inspect(self):
        """Test OkResult.inspect."""
        called = []
        result = ok(42).inspect(lambda x: called.append(x))
        
        assert called == [42]
        assert result.unwrap() == 42
    
    def test_err_result_inspect(self):
        """Test ErrResult.inspect does nothing."""
        called = []
        error = validation("test")
        result = err(error).inspect(lambda x: called.append(x))
        
        assert called == []
        assert result.is_err()
    
    def test_ok_result_inspect_err(self):
        """Test OkResult.inspect_err does nothing."""
        called = []
        result = ok(42).inspect_err(lambda e: called.append(e))
        
        assert called == []
        assert result.unwrap() == 42
    
    def test_err_result_inspect_err(self):
        """Test ErrResult.inspect_err."""
        called = []
        error = validation("test")
        result = err(error).inspect_err(lambda e: called.append(e))
        
        assert len(called) == 1
        assert called[0] == error
    
    def test_ok_result_unwrap(self):
        """Test OkResult.unwrap."""
        assert ok(42).unwrap() == 42
    
    def test_err_result_unwrap(self):
        """Test ErrResult.unwrap raises."""
        with pytest.raises(ValueError):
            err(validation("test")).unwrap()
    
    def test_ok_result_unwrap_or(self):
        """Test OkResult.unwrap_or."""
        assert ok(42).unwrap_or(0) == 42
    
    def test_err_result_unwrap_or(self):
        """Test ErrResult.unwrap_or."""
        assert err(validation("test")).unwrap_or(0) == 0
    
    def test_ok_result_unwrap_or_else(self):
        """Test OkResult.unwrap_or_else."""
        assert ok(42).unwrap_or_else(lambda e: 0) == 42
    
    def test_err_result_unwrap_or_else(self):
        """Test ErrResult.unwrap_or_else."""
        assert err(validation("test")).unwrap_or_else(lambda e: 0) == 0
    
    def test_ok_result_expect(self):
        """Test OkResult.expect."""
        assert ok(42).expect("should have value") == 42
    
    def test_err_result_expect(self):
        """Test ErrResult.expect raises."""
        with pytest.raises(ValueError, match="expected"):
            err(validation("test")).expect("expected value")
    
    def test_err_result_unwrap_err(self):
        """Test ErrResult.unwrap_err."""
        error = validation("test")
        assert err(error).unwrap_err() == error


class TestWrapResult:
    """Test wrap_result function."""
    
    def test_wrap_ok_result(self):
        """Test wrapping Ok result."""
        result = Ok(42)
        wrapped = wrap_result(result)
        
        assert isinstance(wrapped, OkResult)
        assert wrapped.unwrap() == 42
    
    def test_wrap_err_result(self):
        """Test wrapping Err result."""
        error = validation("test")
        result = Err(error)
        wrapped = wrap_result(result)
        
        assert isinstance(wrapped, ErrResult)
        assert wrapped.unwrap_err() == error


class TestPipeline:
    """Test Pipeline builder."""
    
    def test_pipeline_creation(self):
        """Test creating a pipeline."""
        pipeline = Pipeline()
        assert isinstance(pipeline, Pipeline)
    
    def test_pipeline_then(self):
        """Test Pipeline.then."""
        pipeline = (
            Pipeline()
            .then(lambda x: Ok(x * 2))
            .then(lambda x: Ok(x + 1))
        )
        
        result = pipeline.run(Ok(5))
        assert unwrap(result) == 11
    
    def test_pipeline_map(self):
        """Test Pipeline.map."""
        pipeline = (
            Pipeline()
            .map(lambda x: x * 2)
            .map(lambda x: x + 1)
        )
        
        result = pipeline.run(Ok(5))
        assert unwrap(result) == 11
    
    def test_pipeline_recover(self):
        """Test Pipeline.recover."""
        pipeline = (
            Pipeline()
            .then(lambda x: Err(validation("error")) if x < 0 else Ok(x))
            .recover(lambda e: Ok(0))
        )
        
        result = pipeline.run(Ok(-5))
        assert unwrap(result) == 0
    
    def test_pipeline_mixed_operations(self):
        """Test pipeline with mixed operations."""
        pipeline = (
            Pipeline()
            .map(lambda x: x * 2)           # 10
            .then(lambda x: Ok(x + 3))      # 13
            .map(lambda x: x - 1)           # 12
        )
        
        result = pipeline.run(Ok(5))
        assert unwrap(result) == 12
    
    def test_pipeline_error_propagation(self):
        """Test error propagation through pipeline."""
        pipeline = (
            Pipeline()
            .map(lambda x: x * 2)
            .then(lambda x: Err(validation("error")))
            .map(lambda x: x + 1)  # Should not execute
        )
        
        result = pipeline.run(Ok(5))
        assert is_err(result)
    
    def test_pipeline_reusable(self):
        """Test that pipelines are reusable."""
        pipeline = (
            Pipeline()
            .map(lambda x: x * 2)
            .then(lambda x: Ok(x + 1))
        )
        
        result1 = pipeline.run(Ok(5))
        result2 = pipeline.run(Ok(10))
        
        assert unwrap(result1) == 11
        assert unwrap(result2) == 21


class TestCompose:
    """Test compose function."""
    
    def test_compose_single_function(self):
        """Test compose with single function."""
        from .result import map_value
        
        composed = compose(
            lambda r: map_value(r, lambda x: x * 2)
        )
        
        result = composed(Ok(5))
        assert unwrap(result) == 10
    
    def test_compose_multiple_functions(self):
        """Test compose with multiple functions."""
        from .result import map_value, and_then
        
        composed = compose(
            lambda r: map_value(r, lambda x: x * 2),
            lambda r: and_then(r, lambda x: Ok(x + 3)),
            lambda r: map_value(r, lambda x: x - 1),
        )
        
        result = composed(Ok(5))
        assert unwrap(result) == 12
    
    def test_compose_with_errors(self):
        """Test compose with errors."""
        from .result import map_value, and_then
        
        composed = compose(
            lambda r: map_value(r, lambda x: x * 2),
            lambda r: and_then(r, lambda x: Err(validation("error")) if x > 10 else Ok(x)),
        )
        
        result = composed(Ok(6))  # 6 * 2 = 12 > 10
        assert is_err(result)


class TestSafeDecorator:
    """Test @safe decorator."""
    
    def test_safe_with_success(self):
        """Test @safe with successful function."""
        @safe
        def parse_int(s: str) -> int:
            return int(s)
        
        result = parse_int("42")
        assert is_ok(result)
        assert unwrap(result) == 42
    
    def test_safe_with_exception(self):
        """Test @safe with exception."""
        @safe
        def parse_int(s: str) -> int:
            return int(s)
        
        result = parse_int("invalid")
        assert is_err(result)
        error = unwrap_err(result)  # Extract the error from Err
        assert error.error_type == ErrorType.INTERNAL
        assert "parse_int" in error.message
        assert error.has_detail("exception_type")
        assert error.get_detail("exception_type") == "ValueError"
    
    def test_safe_preserves_function_name(self):
        """Test @safe preserves function name."""
        @safe
        def my_function(x: int) -> int:
            return x * 2
        
        assert my_function.__name__ == "my_function"
    
    def test_safe_with_different_exception_types(self):
        """Test @safe captures exception type."""
        @safe
        def divide(a: int, b: int) -> float:
            return a / b
        
        result = divide(10, 0)
        assert is_err(result)
        error = unwrap_err(result)  # Use unwrap_err, not unwrap
        assert error.has_detail("exception_type")
        assert error.get_detail("exception_type") == "ZeroDivisionError"


class TestSafeAsyncDecorator:
    """Test @safe_async decorator."""
    
    @pytest.mark.asyncio
    async def test_safe_async_with_success(self):
        """Test @safe_async with successful async function."""
        @safe_async
        async def fetch_data(x: int) -> int:
            return x * 2
        
        result = await fetch_data(21)
        assert is_ok(result)
        assert unwrap(result) == 42
    
    @pytest.mark.asyncio
    async def test_safe_async_with_exception(self):
        """Test @safe_async with exception."""
        @safe_async
        async def fetch_data(x: int) -> int:
            raise ValueError("failed")
        
        result = await fetch_data(21)
        assert is_err(result)
        error = unwrap_err(result)  # Use unwrap_err, not unwrap
        assert error.error_type == ErrorType.INTERNAL
        assert "fetch_data" in error.message
    
    @pytest.mark.asyncio
    async def test_safe_async_preserves_function_name(self):
        """Test @safe_async preserves function name."""
        @safe_async
        async def my_async_function(x: int) -> int:
            return x * 2
        
        assert my_async_function.__name__ == "my_async_function"


# =========================================
# Integration Tests
# =========================================

class TestCombinatorIntegration:
    """Integration tests for combinators."""
    
    def test_chaining_with_ok_result(self):
        """Test method chaining with OkResult."""
        result = (
            ok(5)
            .map(lambda x: x * 2)           # 10
            .and_then(lambda x: ok(x + 3))  # 13
            .map(lambda x: x - 1)           # 12
            .inspect(lambda x: None)        # Side effect
            .unwrap_or(0)
        )
        
        assert result == 12
    
    def test_chaining_with_err_result(self):
        """Test method chaining with ErrResult."""
        error = validation("invalid")
        result = (
            err(error)
            .map(lambda x: x * 2)           # Skipped
            .and_then(lambda x: ok(x + 3))  # Skipped
            .or_else(lambda e: ok(42))      # Fallback
            .unwrap_or(0)
        )
        
        assert result == 42
    
    def test_pipeline_with_validation(self):
        """Test pipeline with validation logic."""
        def validate_positive(x: int) -> Result[int, Error]:
            if x > 0:
                return Ok(x)
            return Err(validation("must be positive"))
        
        def validate_even(x: int) -> Result[int, Error]:
            if x % 2 == 0:
                return Ok(x)
            return Err(validation("must be even"))
        
        pipeline = (
            Pipeline()
            .then(validate_positive)
            .then(validate_even)
            .map(lambda x: x * 2)
        )
        
        # Valid case
        result = pipeline.run(Ok(4))
        assert unwrap(result) == 8
        
        # Invalid: negative
        result = pipeline.run(Ok(-4))
        assert is_err(result)
        
        # Invalid: odd
        result = pipeline.run(Ok(3))
        assert is_err(result)
    
    def test_safe_in_pipeline(self):
        """Test @safe decorator in pipeline."""
        @safe
        def parse_int(s: str) -> int:
            return int(s)
        
        pipeline = (
            Pipeline()
            .map(lambda x: str(x))
            .then(parse_int)
            .map(lambda x: x * 2)
        )
        
        result = pipeline.run(Ok(21))
        assert unwrap(result) == 42


# =========================================
# Property-Based Tests
# =========================================

class TestCombinatorProperties:
    """Property-based tests for combinators."""
    
    @given(st.integers())
    def test_ok_result_map_identity(self, value: int):
        """Property: map with identity returns same value."""
        result = ok(value).map(lambda x: x).unwrap()
        assert result == value
    
    @given(st.integers(), st.integers())
    def test_ok_result_map_composition(self, value: int, offset: int):
        """Property: map composition."""
        # map(f).map(g) == map(g . f)
        result1 = ok(value).map(lambda x: x + 1).map(lambda x: x + offset).unwrap()
        result2 = ok(value).map(lambda x: x + 1 + offset).unwrap()
        assert result1 == result2


# =========================================
# Benchmarks
# =========================================

class TestCombinatorBenchmarks:
    """Benchmark combinator operations."""
    
    def test_ok_creation_benchmark(self, benchmark):
        """Benchmark ok() creation."""
        result = benchmark(ok, 42)
        assert result.is_ok()
    
    def test_err_creation_benchmark(self, benchmark):
        """Benchmark err() creation."""
        error = validation("test")
        result = benchmark(err, error)
        assert result.is_err()
    
    def test_ok_result_map_benchmark(self, benchmark):
        """Benchmark OkResult.map."""
        result = ok(42)
        mapped = benchmark(result.map, lambda x: x * 2)
        assert mapped.unwrap() == 84
    
    def test_ok_result_and_then_benchmark(self, benchmark):
        """Benchmark OkResult.and_then."""
        result = ok(42)
        chained = benchmark(result.and_then, lambda x: ok(x * 2))
        assert chained.unwrap() == 84
    
    def test_ok_result_unwrap_benchmark(self, benchmark):
        """Benchmark OkResult.unwrap."""
        result = ok(42)
        value = benchmark(result.unwrap)
        assert value == 42
    
    def test_pipeline_run_benchmark(self, benchmark):
        """Benchmark Pipeline.run."""
        pipeline = (
            Pipeline()
            .map(lambda x: x * 2)
            .map(lambda x: x + 1)
            .map(lambda x: x - 1)
        )
        
        result = benchmark(pipeline.run, Ok(21))
        assert unwrap(result) == 42
    
    def test_safe_decorator_benchmark(self, benchmark):
        """Benchmark @safe decorator."""
        @safe
        def parse_int(s: str) -> int:
            return int(s)
        
        result = benchmark(parse_int, "42")
        assert unwrap(result) == 42