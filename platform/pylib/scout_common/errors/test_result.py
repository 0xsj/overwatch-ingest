"""Tests for Result type."""

import pytest
from hypothesis import given, strategies as st

from .result import (
    Ok,
    Err,
    Result,
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
    collect,
    partition,
)
from .base import Error
from .types import ErrorType
from .codes import code


class TestOkErr:
    """Test Ok and Err types."""
    
    def test_ok_creation(self):
        """Test creating Ok."""
        result = Ok(42)
        assert result.value == 42
    
    def test_err_creation(self, simple_error):
        """Test creating Err."""
        result = Err(simple_error)
        assert result.error == simple_error
    
    def test_ok_repr(self):
        """Test Ok repr."""
        result = Ok(42)
        assert repr(result) == "Ok(42)"
    
    def test_err_repr(self, simple_error):
        """Test Err repr."""
        result = Err(simple_error)
        assert "Err(" in repr(result)
    
    def test_ok_equality(self):
        """Test Ok equality."""
        assert Ok(42) == Ok(42)
        assert Ok(42) != Ok(43)
    
    def test_err_equality(self):
        """Test Err equality."""
        err1 = Error(ErrorType.VALIDATION, code("TEST"), "test")
        err2 = Error(ErrorType.VALIDATION, code("TEST"), "test")
        err3 = Error(ErrorType.INTERNAL, code("TEST"), "test")
        
        assert Err(err1) == Err(err2)
        assert Err(err1) != Err(err3)


class TestTypeGuards:
    """Test is_ok and is_err type guards."""
    
    def test_is_ok_with_ok(self):
        """Test is_ok returns True for Ok."""
        result = Ok(42)
        assert is_ok(result)
    
    def test_is_ok_with_err(self, simple_error):
        """Test is_ok returns False for Err."""
        result = Err(simple_error)
        assert not is_ok(result)
    
    def test_is_err_with_err(self, simple_error):
        """Test is_err returns True for Err."""
        result = Err(simple_error)
        assert is_err(result)
    
    def test_is_err_with_ok(self):
        """Test is_err returns False for Ok."""
        result = Ok(42)
        assert not is_err(result)


class TestUnwrapping:
    """Test unwrapping functions."""
    
    def test_unwrap_ok(self):
        """Test unwrap with Ok."""
        result = Ok(42)
        assert unwrap(result) == 42
    
    def test_unwrap_err(self, simple_error):
        """Test unwrap with Err raises."""
        result = Err(simple_error)
        with pytest.raises(ValueError, match="Called unwrap on Err"):
            unwrap(result)
    
    def test_unwrap_or_ok(self):
        """Test unwrap_or with Ok."""
        result = Ok(42)
        assert unwrap_or(result, 0) == 42
    
    def test_unwrap_or_err(self, simple_error):
        """Test unwrap_or with Err returns default."""
        result = Err(simple_error)
        assert unwrap_or(result, 0) == 0
    
    def test_unwrap_or_else_ok(self):
        """Test unwrap_or_else with Ok."""
        result = Ok(42)
        assert unwrap_or_else(result, lambda e: 0) == 42
    
    def test_unwrap_or_else_err(self, simple_error):
        """Test unwrap_or_else with Err calls function."""
        result = Err(simple_error)
        assert unwrap_or_else(result, lambda e: 0) == 0
    
    def test_expect_ok(self):
        """Test expect with Ok."""
        result = Ok(42)
        assert expect(result, "should have value") == 42
    
    def test_expect_err(self, simple_error):
        """Test expect with Err raises with message."""
        result = Err(simple_error)
        with pytest.raises(ValueError, match="expected value"):
            expect(result, "expected value")
    
    def test_unwrap_err_err(self, simple_error):
        """Test unwrap_err with Err."""
        result = Err(simple_error)
        assert unwrap_err(result) == simple_error
    
    def test_unwrap_err_ok(self):
        """Test unwrap_err with Ok raises."""
        result = Ok(42)
        with pytest.raises(ValueError, match="Called unwrap_err on Ok"):
            unwrap_err(result)


class TestTransformations:
    """Test transformation functions."""
    
    def test_map_value_ok(self):
        """Test map_value with Ok."""
        result = Ok(5)
        mapped = map_value(result, lambda x: x * 2)
        assert unwrap(mapped) == 10
    
    def test_map_value_err(self, simple_error):
        """Test map_value with Err returns unchanged."""
        result = Err(simple_error)
        mapped = map_value(result, lambda x: x * 2)
        assert is_err(mapped)
        assert unwrap_err(mapped) == simple_error
    
    def test_map_error_err(self, simple_error):
        """Test map_error with Err."""
        result = Err(simple_error)
        mapped = map_error(result, lambda e: e.with_detail("key", "value"))
        
        assert is_err(mapped)
        err = unwrap_err(mapped)
        assert err.has_detail("key")
    
    def test_map_error_ok(self):
        """Test map_error with Ok returns unchanged."""
        result = Ok(42)
        mapped = map_error(result, lambda e: e.with_detail("key", "value"))
        assert is_ok(mapped)
        assert unwrap(mapped) == 42
    
    def test_and_then_ok(self):
        """Test and_then with Ok."""
        result = Ok(5)
        
        def double(x: int) -> Result[int, Error]:
            return Ok(x * 2)
        
        chained = and_then(result, double)
        assert unwrap(chained) == 10
    
    def test_and_then_ok_to_err(self):
        """Test and_then with Ok that returns Err."""
        result = Ok(-5)
        
        def validate_positive(x: int) -> Result[int, Error]:
            if x > 0:
                return Ok(x)
            return Err(Error(ErrorType.VALIDATION, code("NEGATIVE"), "must be positive"))
        
        chained = and_then(result, validate_positive)
        assert is_err(chained)
    
    def test_and_then_err(self, simple_error):
        """Test and_then with Err returns unchanged."""
        result = Err(simple_error)
        
        def double(x: int) -> Result[int, Error]:
            return Ok(x * 2)
        
        chained = and_then(result, double)
        assert is_err(chained)
        assert unwrap_err(chained) == simple_error
    
    def test_or_else_ok(self):
        """Test or_else with Ok returns unchanged."""
        result = Ok(42)
        
        def fallback(e: Error) -> Result[int, Error]:
            return Ok(0)
        
        handled = or_else(result, fallback)
        assert unwrap(handled) == 42
    
    def test_or_else_err(self, simple_error):
        """Test or_else with Err calls handler."""
        result = Err(simple_error)
        
        def fallback(e: Error) -> Result[int, Error]:
            return Ok(0)
        
        handled = or_else(result, fallback)
        assert unwrap(handled) == 0
    
    def test_flatten_ok_ok(self):
        """Test flatten with Ok(Ok(value))."""
        inner = Ok(42)
        outer = Ok(inner)
        flattened = flatten(outer)
        assert unwrap(flattened) == 42
    
    def test_flatten_ok_err(self, simple_error):
        """Test flatten with Ok(Err(error))."""
        inner = Err(simple_error)
        outer = Ok(inner)
        flattened = flatten(outer)
        assert is_err(flattened)
        assert unwrap_err(flattened) == simple_error
    
    def test_flatten_err(self, simple_error):
        """Test flatten with Err."""
        result = Err(simple_error)
        flattened = flatten(result)
        assert is_err(flattened)
        assert unwrap_err(flattened) == simple_error


class TestInspection:
    """Test inspection functions."""
    
    def test_inspect_ok_calls_function(self):
        """Test inspect with Ok calls function."""
        result = Ok(42)
        called = []
        
        inspected = inspect(result, lambda x: called.append(x))
        
        assert called == [42]
        assert unwrap(inspected) == 42
    
    def test_inspect_err_does_not_call(self, simple_error):
        """Test inspect with Err does not call function."""
        result = Err(simple_error)
        called = []
        
        inspected = inspect(result, lambda x: called.append(x))
        
        assert called == []
        assert is_err(inspected)
    
    def test_inspect_err_calls_function(self, simple_error):
        """Test inspect_err with Err calls function."""
        result = Err(simple_error)
        called = []
        
        inspected = inspect_err(result, lambda e: called.append(e))
        
        assert len(called) == 1
        assert called[0] == simple_error
        assert is_err(inspected)
    
    def test_inspect_err_ok_does_not_call(self):
        """Test inspect_err with Ok does not call function."""
        result = Ok(42)
        called = []
        
        inspected = inspect_err(result, lambda e: called.append(e))
        
        assert called == []
        assert unwrap(inspected) == 42


class TestCollections:
    """Test collection functions."""
    
    def test_collect_all_ok(self):
        """Test collect with all Ok results."""
        results = [Ok(1), Ok(2), Ok(3)]
        collected = collect(results)
        
        assert is_ok(collected)
        assert unwrap(collected) == [1, 2, 3]
    
    def test_collect_with_err(self, simple_error):
        """Test collect with an Err returns first Err."""
        results = [Ok(1), Err(simple_error), Ok(3)]
        collected = collect(results)
        
        assert is_err(collected)
        assert unwrap_err(collected) == simple_error
    
    def test_collect_empty_list(self):
        """Test collect with empty list."""
        results = []
        collected = collect(results)
        
        assert is_ok(collected)
        assert unwrap(collected) == []
    
    def test_partition_mixed_results(self, simple_error):
        """Test partition with mixed results."""
        err1 = Error(ErrorType.VALIDATION, code("E1"), "error 1")
        err2 = Error(ErrorType.INTERNAL, code("E2"), "error 2")
        
        results = [Ok(1), Err(err1), Ok(3), Err(err2), Ok(5)]
        values, errors = partition(results)
        
        assert values == [1, 3, 5]
        assert errors == [err1, err2]
    
    def test_partition_all_ok(self):
        """Test partition with all Ok."""
        results = [Ok(1), Ok(2), Ok(3)]
        values, errors = partition(results)
        
        assert values == [1, 2, 3]
        assert errors == []
    
    def test_partition_all_err(self):
        """Test partition with all Err."""
        err1 = Error(ErrorType.VALIDATION, code("E1"), "e1")
        err2 = Error(ErrorType.INTERNAL, code("E2"), "e2")
        
        results = [Err(err1), Err(err2)]
        values, errors = partition(results)
        
        assert values == []
        assert errors == [err1, err2]


# =========================================
# Property-Based Tests
# =========================================

class TestResultProperties:
    """Property-based tests for Result."""
    
    @given(st.integers())
    def test_ok_unwrap_roundtrip(self, value: int):
        """Property: Ok(x).unwrap() == x."""
        result = Ok(value)
        assert unwrap(result) == value
    
    @given(st.integers())
    def test_map_value_identity(self, value: int):
        """Property: map_value with identity returns same value."""
        result = Ok(value)
        mapped = map_value(result, lambda x: x)
        assert unwrap(mapped) == value
    
    @given(st.integers(), st.integers())
    def test_map_value_composition(self, value: int, factor: int):
        """Property: map composition."""
        result = Ok(value)
        
        # map(f).map(g) == map(g . f)
        mapped1 = map_value(map_value(result, lambda x: x + 1), lambda x: x * factor)
        mapped2 = map_value(result, lambda x: (x + 1) * factor)
        
        assert unwrap(mapped1) == unwrap(mapped2)
    
    @given(st.integers())
    def test_and_then_with_ok_always(self, value: int):
        """Property: and_then with always-Ok function."""
        result = Ok(value)
        chained = and_then(result, lambda x: Ok(x * 2))
        assert unwrap(chained) == value * 2
    
    @given(st.integers())
    def test_unwrap_or_with_ok_ignores_default(self, value: int):
        """Property: unwrap_or with Ok ignores default."""
        result = Ok(value)
        assert unwrap_or(result, 999) == value


# =========================================
# Benchmarks
# =========================================

class TestResultBenchmarks:
    """Benchmark Result operations."""
    
    def test_ok_creation_benchmark(self, benchmark):
        """Benchmark Ok creation."""
        result = benchmark(Ok, 42)
        assert is_ok(result)
    
    def test_err_creation_benchmark(self, benchmark, simple_error):
        """Benchmark Err creation."""
        result = benchmark(Err, simple_error)
        assert is_err(result)
    
    def test_unwrap_benchmark(self, benchmark, ok_result):
        """Benchmark unwrap."""
        result = benchmark(unwrap, ok_result)
        assert result == 42
    
    def test_unwrap_or_benchmark(self, benchmark, ok_result):
        """Benchmark unwrap_or."""
        result = benchmark(unwrap_or, ok_result, 0)
        assert result == 42
    
    def test_map_value_benchmark(self, benchmark, ok_result):
        """Benchmark map_value."""
        result = benchmark(map_value, ok_result, lambda x: x * 2)
        assert unwrap(result) == 84
    
    def test_and_then_benchmark(self, benchmark, ok_result):
        """Benchmark and_then."""
        result = benchmark(and_then, ok_result, lambda x: Ok(x * 2))
        assert unwrap(result) == 84
    
    def test_collect_benchmark(self, benchmark):
        """Benchmark collect with 100 Ok results."""
        results = [Ok(i) for i in range(100)]
        collected = benchmark(collect, results)
        assert len(unwrap(collected)) == 100
    
    def test_partition_benchmark(self, benchmark, simple_error):
        """Benchmark partition with mixed results."""
        results = [Ok(i) if i % 2 == 0 else Err(simple_error) for i in range(100)]
        values, errors = benchmark(partition, results)
        assert len(values) == 50
        assert len(errors) == 50


# =========================================
# Integration Tests
# =========================================

class TestResultIntegration:
    """Integration tests for Result workflows."""
    
    def test_chaining_multiple_operations(self):
        """Test chaining multiple Result operations."""
        initial = Ok(5)
        
        # Manual chaining (Python doesn't have .pipe())
        result = map_value(initial, lambda x: x * 2)  # 10
        result = and_then(result, lambda x: Ok(x + 3) if x > 0 else Err(
            Error(ErrorType.VALIDATION, code("NEG"), "negative")
        ))  # 13
        final = map_value(result, lambda x: x - 1)  # 12
        
        assert unwrap(final) == 12
    
    def test_error_propagation(self):
        """Test that errors propagate through chains."""
        err = Error(ErrorType.VALIDATION, code("ERR"), "error")
        result = Err(err)
        
        # Chain multiple operations
        result = map_value(result, lambda x: x * 2)
        result = and_then(result, lambda x: Ok(x + 1))
        result = map_value(result, lambda x: x * 3)
        
        # Should still be the original error
        assert is_err(result)
        assert unwrap_err(result) == err