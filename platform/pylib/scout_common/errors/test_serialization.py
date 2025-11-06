"""Tests for error serialization."""

import json
import pytest
from hypothesis import given, strategies as st

from .serialization import (
    ErrorDict,
    to_dict,
    to_dict_verbose,
    from_dict,
    to_json,
    to_json_verbose,
    from_json,
    to_http_response,
    format_error,
    to_dict_list,
    from_dict_list,
    is_valid_error_dict,
    get_error_schema,
)
from .base import Error
from .types import ErrorType
from .codes import code
from .constructors import validation, internal_with_cause, database_error


class TestToDict:
    """Test to_dict function."""
    
    def test_simple_error(self):
        """Test serializing simple error."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test message")
        data = to_dict(err)
        
        assert data["type"] == "VALIDATION"
        assert data["code"] == "TEST"
        assert data["message"] == "test message"
        assert "details" not in data  # Empty details excluded
        assert "cause" not in data
    
    def test_error_with_details(self):
        """Test serializing error with details."""
        err = Error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test",
            details={"field": "email", "reason": "invalid"},
        )
        data = to_dict(err)
        
        assert data["details"]["field"] == "email"
        assert data["details"]["reason"] == "invalid"
    
    def test_error_without_cause_by_default(self):
        """Test that cause is excluded by default."""
        root = Error(ErrorType.DATABASE, code("DB_ERROR"), "db failed")
        wrapped = Error(ErrorType.INTERNAL, code("WRAPPER"), "wrapped", cause=root)
        
        data = to_dict(wrapped)
        
        assert "cause" not in data
    
    def test_error_with_cause_when_requested(self):
        """Test including cause when requested."""
        root = Error(ErrorType.DATABASE, code("DB_ERROR"), "db failed")
        wrapped = Error(ErrorType.INTERNAL, code("WRAPPER"), "wrapped", cause=root)
        
        data = to_dict(wrapped, include_cause=True)
        
        assert "cause" in data
        assert data["cause"]["type"] == "DATABASE"
        assert data["cause"]["code"] == "DB_ERROR"
    
    def test_error_with_nested_cause(self):
        """Test serializing error with nested cause chain."""
        root = Error(ErrorType.NETWORK, code("TIMEOUT"), "timeout")
        middle = Error(ErrorType.DATABASE, code("DB_ERROR"), "db failed", cause=root)
        top = Error(ErrorType.INTERNAL, code("WRAPPER"), "wrapped", cause=middle)
        
        data = to_dict(top, include_cause=True)
        
        assert data["cause"]["type"] == "DATABASE"
        assert data["cause"]["cause"]["type"] == "NETWORK"


class TestToDictVerbose:
    """Test to_dict_verbose function."""
    
    def test_verbose_includes_cause(self):
        """Test that verbose mode includes cause."""
        root = Error(ErrorType.DATABASE, code("DB_ERROR"), "db failed")
        wrapped = Error(ErrorType.INTERNAL, code("WRAPPER"), "wrapped", cause=root)
        
        data = to_dict_verbose(wrapped)
        
        assert "cause" in data
        assert data["cause"]["type"] == "DATABASE"


class TestFromDict:
    """Test from_dict function."""
    
    def test_simple_error(self):
        """Test deserializing simple error."""
        data = {
            "type": "VALIDATION",
            "code": "TEST_CODE",
            "message": "test message",
        }
        
        err = from_dict(data)
        
        assert err.error_type == ErrorType.VALIDATION
        assert err.code == code("TEST_CODE")
        assert err.message == "test message"
        assert err.details == {}
        assert err.cause is None
    
    def test_error_with_details(self):
        """Test deserializing error with details."""
        data = {
            "type": "VALIDATION",
            "code": "TEST",
            "message": "test",
            "details": {"field": "email"},
        }
        
        err = from_dict(data)
        
        assert err.details["field"] == "email"
    
    def test_error_with_cause(self):
        """Test deserializing error with cause."""
        data = {
            "type": "INTERNAL",
            "code": "WRAPPER",
            "message": "wrapped",
            "cause": {
                "type": "DATABASE",
                "code": "DB_ERROR",
                "message": "db failed",
            },
        }
        
        err = from_dict(data)
        
        assert err.error_type == ErrorType.INTERNAL
        assert err.cause is not None
        assert err.cause.error_type == ErrorType.DATABASE
        assert err.cause.code == code("DB_ERROR")
    
    def test_missing_required_field_raises(self):
        """Test that missing required fields raise ValueError."""
        with pytest.raises(ValueError, match="Missing required field: type"):
            from_dict({"code": "TEST", "message": "test"})
        
        with pytest.raises(ValueError, match="Missing required field: code"):
            from_dict({"type": "VALIDATION", "message": "test"})
        
        with pytest.raises(ValueError, match="Missing required field: message"):
            from_dict({"type": "VALIDATION", "code": "TEST"})
    
    def test_invalid_error_type_raises(self):
        """Test that invalid error type raises ValueError."""
        data = {
            "type": "INVALID_TYPE",
            "code": "TEST",
            "message": "test",
        }
        
        with pytest.raises(ValueError, match="Invalid error type"):
            from_dict(data)
    
    def test_invalid_details_type_raises(self):
        """Test that invalid details type raises ValueError."""
        data = {
            "type": "VALIDATION",
            "code": "TEST",
            "message": "test",
            "details": "not a dict",
        }
        
        with pytest.raises(ValueError, match="details must be a dictionary"):
            from_dict(data)


class TestRoundTrip:
    """Test round-trip serialization."""
    
    def test_simple_error_roundtrip(self):
        """Test simple error round-trip."""
        original = Error(ErrorType.VALIDATION, code("TEST"), "test message")
        
        data = to_dict(original)
        restored = from_dict(data)
        
        assert restored.error_type == original.error_type
        assert restored.code == original.code
        assert restored.message == original.message
    
    def test_error_with_details_roundtrip(self):
        """Test error with details round-trip."""
        original = Error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test",
            details={"field": "email", "reason": "invalid"},
        )
        
        data = to_dict(original)
        restored = from_dict(data)
        
        assert restored.details == original.details
    
    def test_error_with_cause_roundtrip(self):
        """Test error with cause round-trip."""
        root = Error(ErrorType.DATABASE, code("DB_ERROR"), "db failed")
        original = Error(ErrorType.INTERNAL, code("WRAPPER"), "wrapped", cause=root)
        
        data = to_dict_verbose(original)
        restored = from_dict(data)
        
        assert restored.cause is not None
        assert restored.cause.error_type == root.error_type
        assert restored.cause.code == root.code


class TestJSON:
    """Test JSON serialization."""
    
    def test_to_json(self):
        """Test to_json."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        json_str = to_json(err)
        
        # Should be valid JSON
        data = json.loads(json_str)
        assert data["type"] == "VALIDATION"
        assert data["code"] == "TEST"
    
    def test_to_json_with_indent(self):
        """Test to_json with indentation."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test")
        json_str = to_json(err, indent=2)
        
        # Should be pretty-printed
        assert "\n" in json_str
        assert "  " in json_str
    
    def test_to_json_verbose(self):
        """Test to_json_verbose includes cause."""
        root = Error(ErrorType.DATABASE, code("DB_ERROR"), "db failed")
        wrapped = Error(ErrorType.INTERNAL, code("WRAPPER"), "wrapped", cause=root)
        
        json_str = to_json_verbose(wrapped)
        data = json.loads(json_str)
        
        assert "cause" in data
    
    def test_from_json(self):
        """Test from_json."""
        json_str = '{"type": "VALIDATION", "code": "TEST", "message": "test"}'
        err = from_json(json_str)
        
        assert err.error_type == ErrorType.VALIDATION
        assert err.code == code("TEST")
    
    def test_from_json_invalid_json_raises(self):
        """Test from_json with invalid JSON raises."""
        with pytest.raises(ValueError, match="Invalid JSON"):
            from_json("not valid json")
    
    def test_from_json_non_object_raises(self):
        """Test from_json with non-object raises."""
        with pytest.raises(ValueError, match="must represent an object"):
            from_json("[1, 2, 3]")
    
    def test_json_roundtrip(self):
        """Test JSON round-trip."""
        original = Error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test",
            details={"field": "email"},
        )
        
        json_str = to_json(original)
        restored = from_json(json_str)
        
        assert restored.error_type == original.error_type
        assert restored.code == original.code
        assert restored.message == original.message
        assert restored.details == original.details


class TestHTTPResponse:
    """Test HTTP response helpers."""
    
    def test_to_http_response(self):
        """Test to_http_response."""
        err = validation("invalid input")
        error_dict, status_code = to_http_response(err)
        
        assert status_code == 400
        assert error_dict["type"] == "VALIDATION"
    
    def test_to_http_response_various_types(self):
        """Test to_http_response with various error types."""
        test_cases = [
            (ErrorType.VALIDATION, 400),
            (ErrorType.UNAUTHORIZED, 401),
            (ErrorType.FORBIDDEN, 403),
            (ErrorType.NOT_FOUND, 404),
            (ErrorType.CONFLICT, 409),
            (ErrorType.RATE_LIMIT, 429),
            (ErrorType.INTERNAL, 500),
            (ErrorType.UNAVAILABLE, 503),
            (ErrorType.TIMEOUT, 504),
        ]
        
        for error_type, expected_status in test_cases:
            err = Error(error_type, code("TEST"), "test")
            _, status_code = to_http_response(err)
            assert status_code == expected_status


class TestFormatError:
    """Test format_error function."""
    
    def test_format_simple_error(self):
        """Test formatting simple error."""
        err = Error(ErrorType.VALIDATION, code("TEST"), "test message")
        formatted = format_error(err, include_cause=False)
        
        assert "[VALIDATION:TEST] test message" in formatted
    
    def test_format_error_with_details(self):
        """Test formatting error with details."""
        err = Error(
            ErrorType.VALIDATION,
            code("TEST"),
            "test",
            details={"field": "email", "reason": "invalid"},
        )
        formatted = format_error(err, include_cause=False)
        
        assert "field: email" in formatted
        assert "reason: invalid" in formatted
    
    def test_format_error_with_cause(self):
        """Test formatting error with cause."""
        root = Error(ErrorType.DATABASE, code("DB_ERROR"), "db failed")
        wrapped = Error(ErrorType.INTERNAL, code("WRAPPER"), "wrapped", cause=root)
        
        formatted = format_error(wrapped, include_cause=True)
        
        assert "[INTERNAL:WRAPPER] wrapped" in formatted
        assert "Caused by:" in formatted
        assert "[DATABASE:DB_ERROR] db failed" in formatted


class TestCollectionFunctions:
    """Test collection helper functions."""
    
    def test_to_dict_list(self):
        """Test to_dict_list."""
        errors = [
            Error(ErrorType.VALIDATION, code("E1"), "error 1"),
            Error(ErrorType.INTERNAL, code("E2"), "error 2"),
        ]
        
        dicts = to_dict_list(errors)
        
        assert len(dicts) == 2
        assert dicts[0]["code"] == "E1"
        assert dicts[1]["code"] == "E2"
    
    def test_from_dict_list(self):
        """Test from_dict_list."""
        data = [
            {"type": "VALIDATION", "code": "E1", "message": "error 1"},
            {"type": "INTERNAL", "code": "E2", "message": "error 2"},
        ]
        
        errors = from_dict_list(data)
        
        assert len(errors) == 2
        assert errors[0].code == code("E1")
        assert errors[1].code == code("E2")
    
    def test_dict_list_roundtrip(self):
        """Test round-trip for error lists."""
        original = [
            Error(ErrorType.VALIDATION, code("E1"), "error 1"),
            Error(ErrorType.INTERNAL, code("E2"), "error 2"),
        ]
        
        dicts = to_dict_list(original)
        restored = from_dict_list(dicts)
        
        assert len(restored) == len(original)
        for orig, rest in zip(original, restored):
            assert rest.error_type == orig.error_type
            assert rest.code == orig.code


class TestValidation:
    """Test validation helpers."""
    
    def test_is_valid_error_dict_valid(self):
        """Test is_valid_error_dict with valid dict."""
        data = {
            "type": "VALIDATION",
            "code": "TEST",
            "message": "test",
        }
        
        assert is_valid_error_dict(data)
    
    def test_is_valid_error_dict_with_details(self):
        """Test is_valid_error_dict with details."""
        data = {
            "type": "VALIDATION",
            "code": "TEST",
            "message": "test",
            "details": {"field": "email"},
        }
        
        assert is_valid_error_dict(data)
    
    def test_is_valid_error_dict_with_cause(self):
        """Test is_valid_error_dict with cause."""
        data = {
            "type": "INTERNAL",
            "code": "WRAPPER",
            "message": "wrapped",
            "cause": {
                "type": "DATABASE",
                "code": "DB_ERROR",
                "message": "db failed",
            },
        }
        
        assert is_valid_error_dict(data)
    
    def test_is_valid_error_dict_missing_field(self):
        """Test is_valid_error_dict with missing field."""
        data = {"type": "VALIDATION", "message": "test"}  # Missing code
        assert not is_valid_error_dict(data)
    
    def test_is_valid_error_dict_invalid_type(self):
        """Test is_valid_error_dict with invalid type."""
        data = {
            "type": "INVALID_TYPE",
            "code": "TEST",
            "message": "test",
        }
        
        assert not is_valid_error_dict(data)
    
    def test_is_valid_error_dict_invalid_details(self):
        """Test is_valid_error_dict with invalid details."""
        data = {
            "type": "VALIDATION",
            "code": "TEST",
            "message": "test",
            "details": "not a dict",
        }
        
        assert not is_valid_error_dict(data)
    
    def test_is_valid_error_dict_not_dict(self):
        """Test is_valid_error_dict with non-dict."""
        assert not is_valid_error_dict("not a dict")
        assert not is_valid_error_dict([1, 2, 3])
        assert not is_valid_error_dict(None)


class TestGetErrorSchema:
    """Test get_error_schema function."""
    
    def test_get_error_schema_structure(self):
        """Test error schema structure."""
        schema = get_error_schema()
        
        assert schema["type"] == "object"
        assert "type" in schema["required"]
        assert "code" in schema["required"]
        assert "message" in schema["required"]
        
        assert "type" in schema["properties"]
        assert "code" in schema["properties"]
        assert "message" in schema["properties"]
        assert "details" in schema["properties"]
        assert "cause" in schema["properties"]
    
    def test_get_error_schema_type_enum(self):
        """Test that schema includes all error types."""
        schema = get_error_schema()
        
        type_enum = schema["properties"]["type"]["enum"]
        assert "VALIDATION" in type_enum
        assert "INTERNAL" in type_enum
        assert "DATABASE" in type_enum
        assert len(type_enum) == 15  # All error types


# =========================================
# Property-Based Tests
# =========================================

class TestSerializationProperties:
    """Property-based tests for serialization."""
    
    @given(
        st.sampled_from(list(ErrorType)),
        st.text(min_size=1),
        st.text(min_size=1),
    )
    def test_roundtrip_always_works(
        self,
        error_type: ErrorType,
        code_str: str,
        message: str,
    ):
        """Property: to_dict → from_dict round-trip always works."""
        original = Error(error_type, code(code_str), message)
        
        data = to_dict(original)
        restored = from_dict(data)
        
        assert restored.error_type == original.error_type
        assert restored.code == original.code
        assert restored.message == original.message


# =========================================
# Benchmarks
# =========================================

class TestSerializationBenchmarks:
    """Benchmark serialization operations."""
    
    def test_to_dict_benchmark(self, benchmark, simple_error):
        """Benchmark to_dict."""
        result = benchmark(to_dict, simple_error)
        assert "type" in result
    
    def test_to_dict_verbose_benchmark(self, benchmark, error_with_cause):
        """Benchmark to_dict_verbose with cause chain."""
        result = benchmark(to_dict_verbose, error_with_cause)
        assert "cause" in result
    
    def test_from_dict_benchmark(self, benchmark):
        """Benchmark from_dict."""
        data = {
            "type": "VALIDATION",
            "code": "TEST",
            "message": "test",
            "details": {"field": "email"},
        }
        
        result = benchmark(from_dict, data)
        assert result.error_type == ErrorType.VALIDATION
    
    def test_to_json_benchmark(self, benchmark, simple_error):
        """Benchmark to_json."""
        result = benchmark(to_json, simple_error)
        assert isinstance(result, str)
    
    def test_from_json_benchmark(self, benchmark):
        """Benchmark from_json."""
        json_str = '{"type": "VALIDATION", "code": "TEST", "message": "test"}'
        result = benchmark(from_json, json_str)
        assert result.error_type == ErrorType.VALIDATION
    
    def test_format_error_benchmark(self, benchmark, error_with_cause):
        """Benchmark format_error."""
        result = benchmark(format_error, error_with_cause)
        assert isinstance(result, str)
    
    def test_to_dict_list_benchmark(self, benchmark):
        """Benchmark to_dict_list with 100 errors."""
        errors = [
            Error(ErrorType.VALIDATION, code(f"E{i}"), f"error {i}")
            for i in range(100)
        ]
        
        result = benchmark(to_dict_list, errors)
        assert len(result) == 100