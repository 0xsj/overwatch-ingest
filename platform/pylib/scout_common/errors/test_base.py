"""Tests for Error base class."""

import pytest
from hypothesis import given, strategies as st

from .base import Error, error
from .types import ErrorType
from .codes import code


class TestError:
    """Test Error dataclass."""
    
    def test_error_creation(self):
        """Test creating an error."""
        err = Error(
            error_type=ErrorType.VALIDATION,
            code=code("TEST"),
            message="test message",
        )
        
        assert err.error_type == ErrorType.VALIDATION
        assert err.code == code("TEST")
        assert err.message == "test message"
        assert err.details == {}
        assert err.cause is None
    
    def test_error_with_details(self):
        """Test error with details."""
        err = Error(
            error_type=ErrorType.VALIDATION,
            code=code("TEST"),
            message="test",
            details={"field": "email"},
        )
        
        assert err.details["field"] == "email"
    
    def test_error_with_cause(self):
        """Test error with cause."""
        root = Error(ErrorType.DATABASE, code("ROOT"), "root error")
        wrapped = Error(
            ErrorType.INTERNAL,
            code("WRAPPER"),
            "wrapper",
            cause=root,
        )
        
        assert wrapped.cause == root
    
    def test_error_str(self):
        """Test error string representation."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test message")
        assert str(err) == "[VALIDATION:TEST] test message"
    
    def test_error_str_with_cause(self):
        """Test error string with cause."""
        root = Error(ErrorType.DATABASE, code("ROOT"), "root")
        wrapped = Error(ErrorType.INTERNAL, code("WRAP"), "wrapped", cause=root)
        
        s = str(wrapped)
        assert "[INTERNAL:WRAP] wrapped" in s
        assert "caused by:" in s
        assert "[DATABASE:ROOT] root" in s
    
    def test_error_repr(self):
        """Test error repr."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        r = repr(err)
        
        assert "Error(" in r
        assert "error_type=" in r
        assert "code=" in r
        assert "message=" in r
    
    def test_with_detail(self):
        """Test adding a detail."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        err2 = err.with_detail("field", "email")
        
        # Original unchanged (immutable)
        assert "field" not in err.details
        
        # New error has detail
        assert err2.details["field"] == "email"
    
    def test_with_details(self):
        """Test adding multiple details."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        err2 = err.with_details({"field": "email", "reason": "invalid"})
        
        assert err2.details["field"] == "email"
        assert err2.details["reason"] == "invalid"
        assert len(err.details) == 0  # Original unchanged
    
    def test_with_cause_method(self):
        """Test adding a cause."""
        root = Error(ErrorType.DATABASE, code("ROOT"), "root")
        err = Error(ErrorType.INTERNAL, code("TEST"), "test")
        err2 = err.with_cause(root)
        
        assert err.cause is None  # Original unchanged
        assert err2.cause == root
    
    def test_get_detail(self):
        """Test getting a detail."""
        err = Error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test",
            details={"field": "email"},
        )
        
        assert err.get_detail("field") == "email"
        assert err.get_detail("missing") is None
    
    def test_has_detail(self):
        """Test checking for detail."""
        err = Error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test",
            details={"field": "email"},
        )
        
        assert err.has_detail("field")
        assert not err.has_detail("missing")
    
    def test_has_details(self):
        """Test checking for any details."""
        err1 = Error(ErrorType.VALIDATION, code("TEST"), "test")
        err2 = Error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test",
            details={"field": "email"},
        )
        
        assert not err1.has_details()
        assert err2.has_details()
    
    def test_is_client_error(self):
        """Test is_client_error."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        assert err.is_client_error()
        
        err2 = Error(ErrorType.INTERNAL, code("TEST"), "test")
        assert not err2.is_client_error()
    
    def test_is_server_error(self):
        """Test is_server_error."""
        err = Error(ErrorType.INTERNAL, code("TEST"), "test")
        assert err.is_server_error()
        
        err2 = Error(ErrorType.VALIDATION, code("TEST"), "test")
        assert not err2.is_server_error()
    
    def test_is_retryable(self):
        """Test is_retryable."""
        err = Error(ErrorType.TIMEOUT, code("TEST"), "test")
        assert err.is_retryable()
        
        err2 = Error(ErrorType.VALIDATION, code("TEST"), "test")
        assert not err2.is_retryable()
    
    def test_http_status_code(self):
        """Test http_status_code."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        assert err.http_status_code() == 400
        
        err2 = Error(ErrorType.NOT_FOUND, code("TEST"), "test")
        assert err2.http_status_code() == 404
    
    def test_get_cause(self):
        """Test get_cause."""
        root = Error(ErrorType.DATABASE, code("ROOT"), "root")
        err = Error(ErrorType.INTERNAL, code("TEST"), "test", cause=root)
        
        assert err.get_cause() == root
        
        err2 = Error(ErrorType.VALIDATION, code("TEST"), "test")
        assert err2.get_cause() is None
    
    def test_get_root_cause(self):
        """Test get_root_cause."""
        root = Error(ErrorType.DATABASE, code("ROOT"), "root")
        middle = Error(ErrorType.NETWORK, code("MID"), "middle", cause=root)
        top = Error(ErrorType.INTERNAL, code("TOP"), "top", cause=middle)
        
        assert top.get_root_cause() == root
        assert middle.get_root_cause() == root
        assert root.get_root_cause() == root  # Returns self
    
    def test_matches_type(self):
        """Test matches with error_type."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test message")
        
        assert err.matches(error_type=ErrorType.VALIDATION)
        assert not err.matches(error_type=ErrorType.INTERNAL)
    
    def test_matches_code(self):
        """Test matches with code."""
        err = Error(ErrorType.VALIDATION, code("TEST_CODE"), "test")
        
        assert err.matches(code=code("TEST_CODE"))
        assert not err.matches(code=code("OTHER_CODE"))
    
    def test_matches_message(self):
        """Test matches with message substring."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test message with keyword")
        
        assert err.matches(message_contains="keyword")
        assert err.matches(message_contains="test message")
        assert not err.matches(message_contains="missing")
    
    def test_matches_multiple_criteria(self):
        """Test matches with multiple criteria."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test message")
        
        assert err.matches(
            error_type=ErrorType.VALIDATION,
            code=code("TEST"),
            message_contains="test",
        )
        
        assert not err.matches(
            error_type=ErrorType.VALIDATION,
            code=code("WRONG"),
        )
    
    def test_error_immutable(self):
        """Test that errors are immutable."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        
        with pytest.raises(Exception):  # FrozenInstanceError or similar
            err.message = "modified"  # type: ignore
    
    def test_error_hashable(self):
        """Test that errors are hashable (can go in sets/dicts)."""
        err1 = Error(ErrorType.VALIDATION, code("TEST"), "test")
        err2 = Error(ErrorType.VALIDATION, code("TEST"), "test")
        
        # Can be used in sets
        s = {err1, err2}
        assert isinstance(s, set)
        
        # Can be dict keys
        d = {err1: "value"}
        assert d[err1] == "value"


class TestErrorHelperFunction:
    """Test error() helper function."""
    
    def test_error_function(self):
        """Test error() helper."""
        err = error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test message",
        )
        
        assert err.error_type == ErrorType.VALIDATION
        assert err.code == code("TEST")
        assert err.message == "test message"
    
    def test_error_with_details(self):
        """Test error() with details."""
        err = error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test",
            details={"field": "email"},
        )
        
        assert err.details["field"] == "email"
    
    def test_error_with_cause(self):
        """Test error() with cause."""
        root = error(ErrorType.DATABASE, code("ROOT"), "root")
        wrapped = error(
            ErrorType.INTERNAL,
            code("WRAPPER"),
            "wrapped",
            cause=root,
        )
        
        assert wrapped.cause == root


# =========================================
# Property-Based Tests
# =========================================

class TestErrorProperties:
    """Property-based tests for Error."""
    
    @given(
        st.sampled_from(list(ErrorType)),
        st.text(min_size=1),
        st.text(min_size=1),
    )
    def test_error_creation_never_crashes(
        self,
        error_type: ErrorType,
        code_str: str,
        message: str,
    ):
        """Property: Error creation should never crash."""
        err = Error(error_type, code(code_str), message)
        assert isinstance(err, Error)
    
    @given(st.text(), st.text())
    def test_with_detail_preserves_immutability(self, key: str, value: str):
        """Property: with_detail should not modify original."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        err2 = err.with_detail(key, value)
        
        assert key not in err.details
        if key:  # Only check if key is non-empty
            assert key in err2.details


# =========================================
# Benchmarks
# =========================================

class TestErrorBenchmarks:
    """Benchmark Error operations."""
    
    def test_error_creation_benchmark(self, benchmark):
        """Benchmark error creation."""
        result = benchmark(
            Error,
            ErrorType.VALIDATION,
            code("TEST"),
            "test message",
        )
        assert isinstance(result, Error)
    
    def test_with_detail_benchmark(self, benchmark, simple_error):
        """Benchmark adding a detail."""
        result = benchmark(simple_error.with_detail, "key", "value")
        assert result.has_detail("key")
    
    def test_with_details_benchmark(self, benchmark, simple_error):
        """Benchmark adding multiple details."""
        details = {"key1": "value1", "key2": "value2", "key3": "value3"}
        result = benchmark(simple_error.with_details, details)
        assert len(result.details) == 3
    
    def test_with_cause_benchmark(self, benchmark, simple_error):
        """Benchmark adding a cause."""
        root = Error(ErrorType.DATABASE, code("ROOT"), "root")
        result = benchmark(simple_error.with_cause, root)
        assert result.cause == root
    
    def test_get_root_cause_benchmark(self, benchmark, error_with_cause):
        """Benchmark getting root cause."""
        result = benchmark(error_with_cause.get_root_cause)
        assert isinstance(result, Error)
    
    def test_matches_benchmark(self, benchmark, simple_error):
        """Benchmark matches."""
        result = benchmark(
            simple_error.matches,
            error_type=ErrorType.VALIDATION,
            code=code("TEST_CODE"),
        )
        assert isinstance(result, bool)
    
    def test_str_conversion_benchmark(self, benchmark, simple_error):
        """Benchmark string conversion."""
        result = benchmark(str, simple_error)
        assert isinstance(result, str)