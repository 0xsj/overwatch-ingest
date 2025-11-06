"""Tests for error constructor functions."""

import pytest

from .constructors import (
    not_found,
    already_exists,
    validation,
    validation_with_field,
    required_field,
    invalid_field,
    unauthorized,
    forbidden,
    internal,
    internal_with_cause,
    not_implemented,
    timeout,
    unavailable,
    unavailable_with_cause,
    conflict,
    rate_limit,
    rate_limit_with_retry,
    database_error,
    database_error_with_table,
    cache_error,
    cache_error_with_key,
    network_error,
    network_error_with_url,
    event_error,
    event_error_with_subject,
    wrap,
    wrap_with_details,
)
from .types import ErrorType
from .codes import code


class TestGenericConstructors:
    """Test generic error constructors."""
    
    def test_not_found(self):
        """Test not_found constructor."""
        err = not_found("user", "123")
        
        assert err.error_type == ErrorType.NOT_FOUND
        assert err.code == code("RESOURCE_NOT_FOUND")
        assert "user" in err.message
        assert "123" in err.message
        assert err.get_detail("resource_type") == "user"
        assert err.get_detail("resource_id") == "123"
    
    def test_already_exists(self):
        """Test already_exists constructor."""
        err = already_exists("user", "john@example.com")
        
        assert err.error_type == ErrorType.ALREADY_EXISTS
        assert err.code == code("RESOURCE_ALREADY_EXISTS")
        assert "user" in err.message
        assert "john@example.com" in err.message
        assert err.get_detail("resource_type") == "user"
        assert err.get_detail("resource_id") == "john@example.com"
    
    def test_validation(self):
        """Test validation constructor."""
        err = validation("invalid input")
        
        assert err.error_type == ErrorType.VALIDATION
        assert err.code == code("VALIDATION_FAILED")
        assert err.message == "invalid input"
    
    def test_validation_with_field(self):
        """Test validation_with_field constructor."""
        err = validation_with_field("email", "invalid format")
        
        assert err.error_type == ErrorType.VALIDATION
        assert err.code == code("VALIDATION_FAILED")
        assert "email" in err.message
        assert "invalid format" in err.message
        assert err.get_detail("field") == "email"
    
    def test_required_field(self):
        """Test required_field constructor."""
        err = required_field("username")
        
        assert err.error_type == ErrorType.VALIDATION
        assert err.code == code("REQUIRED_FIELD_MISSING")
        assert "username" in err.message
        assert err.get_detail("field") == "username"
    
    def test_invalid_field(self):
        """Test invalid_field constructor."""
        err = invalid_field("age", "must be positive")
        
        assert err.error_type == ErrorType.VALIDATION
        assert err.code == code("INVALID_FIELD_VALUE")
        assert "age" in err.message
        assert err.get_detail("field") == "age"
        assert err.get_detail("reason") == "must be positive"


class TestAuthorizationConstructors:
    """Test authorization error constructors."""
    
    def test_unauthorized(self):
        """Test unauthorized constructor."""
        err = unauthorized("invalid token")
        
        assert err.error_type == ErrorType.UNAUTHORIZED
        assert err.code == code("UNAUTHORIZED")
        assert "invalid token" in err.message
        assert err.get_detail("reason") == "invalid token"
    
    def test_forbidden(self):
        """Test forbidden constructor."""
        err = forbidden("document", "delete")
        
        assert err.error_type == ErrorType.FORBIDDEN
        assert err.code == code("FORBIDDEN")
        assert "document" in err.message
        assert "delete" in err.message
        assert err.get_detail("resource") == "document"
        assert err.get_detail("action") == "delete"


class TestInternalConstructors:
    """Test internal error constructors."""
    
    def test_internal(self):
        """Test internal constructor."""
        err = internal("something went wrong")
        
        assert err.error_type == ErrorType.INTERNAL
        assert err.code == code("INTERNAL_ERROR")
        assert err.message == "something went wrong"
    
    def test_internal_with_cause(self):
        """Test internal_with_cause constructor."""
        root = validation("root cause")
        err = internal_with_cause("operation failed", root)
        
        assert err.error_type == ErrorType.INTERNAL
        assert err.code == code("INTERNAL_ERROR")
        assert err.message == "operation failed"
        assert err.cause == root
    
    def test_not_implemented(self):
        """Test not_implemented constructor."""
        err = not_implemented("feature X")
        
        assert err.error_type == ErrorType.NOT_IMPLEMENTED
        assert err.code == code("NOT_IMPLEMENTED")
        assert "feature X" in err.message
        assert err.get_detail("feature") == "feature X"


class TestTimeoutAvailabilityConstructors:
    """Test timeout and availability constructors."""
    
    def test_timeout(self):
        """Test timeout constructor."""
        err = timeout("database query", 5000)
        
        assert err.error_type == ErrorType.TIMEOUT
        assert err.code == code("OPERATION_TIMEOUT")
        assert "database query" in err.message
        assert "5000" in err.message
        assert err.get_detail("operation") == "database query"
        assert err.get_detail("timeout_ms") == "5000"
    
    def test_unavailable(self):
        """Test unavailable constructor."""
        err = unavailable("payment service")
        
        assert err.error_type == ErrorType.UNAVAILABLE
        assert err.code == code("SERVICE_UNAVAILABLE")
        assert "payment service" in err.message
        assert err.get_detail("service") == "payment service"
    
    def test_unavailable_with_cause(self):
        """Test unavailable_with_cause constructor."""
        root = network_error("connection refused")
        err = unavailable_with_cause("api service", root)
        
        assert err.error_type == ErrorType.UNAVAILABLE
        assert err.code == code("SERVICE_UNAVAILABLE")
        assert "api service" in err.message
        assert err.get_detail("service") == "api service"
        assert err.cause == root


class TestConflictRateLimitConstructors:
    """Test conflict and rate limit constructors."""
    
    def test_conflict(self):
        """Test conflict constructor."""
        err = conflict("user", "username already taken")
        
        assert err.error_type == ErrorType.CONFLICT
        assert err.code == code("RESOURCE_CONFLICT")
        assert "user" in err.message
        assert "username already taken" in err.message
        assert err.get_detail("resource") == "user"
        assert err.get_detail("reason") == "username already taken"
    
    def test_rate_limit(self):
        """Test rate_limit constructor."""
        err = rate_limit(100, "minute")
        
        assert err.error_type == ErrorType.RATE_LIMIT
        assert err.code == code("RATE_LIMIT_EXCEEDED")
        assert "100" in err.message
        assert "minute" in err.message
        assert err.get_detail("limit") == "100"
        assert err.get_detail("window") == "minute"
    
    def test_rate_limit_with_retry(self):
        """Test rate_limit_with_retry constructor."""
        err = rate_limit_with_retry(100, "minute", 60)
        
        assert err.error_type == ErrorType.RATE_LIMIT
        assert err.code == code("RATE_LIMIT_EXCEEDED")
        assert err.get_detail("limit") == "100"
        assert err.get_detail("window") == "minute"
        assert err.get_detail("retry_after_seconds") == "60"


class TestInfrastructureConstructors:
    """Test infrastructure error constructors."""
    
    def test_database_error(self):
        """Test database_error constructor."""
        err = database_error("SELECT")
        
        assert err.error_type == ErrorType.DATABASE
        assert err.code == code("DATABASE_ERROR")
        assert "SELECT" in err.message
        assert err.get_detail("operation") == "SELECT"
        assert err.cause is None
    
    def test_database_error_with_cause(self):
        """Test database_error with cause."""
        root = timeout("connection timeout", 5000)
        err = database_error("SELECT", cause=root)
        
        assert err.error_type == ErrorType.DATABASE
        assert err.cause == root
    
    def test_database_error_with_table(self):
        """Test database_error_with_table constructor."""
        err = database_error_with_table("INSERT", "users")
        
        assert err.error_type == ErrorType.DATABASE
        assert err.code == code("DATABASE_ERROR")
        assert "INSERT" in err.message
        assert "users" in err.message
        assert err.get_detail("operation") == "INSERT"
        assert err.get_detail("table") == "users"
    
    def test_cache_error(self):
        """Test cache_error constructor."""
        err = cache_error("GET")
        
        assert err.error_type == ErrorType.CACHE
        assert err.code == code("CACHE_ERROR")
        assert "GET" in err.message
        assert err.get_detail("operation") == "GET"
    
    def test_cache_error_with_key(self):
        """Test cache_error_with_key constructor."""
        err = cache_error_with_key("SET", "user:123")
        
        assert err.error_type == ErrorType.CACHE
        assert err.code == code("CACHE_ERROR")
        assert "SET" in err.message
        assert "user:123" in err.message
        assert err.get_detail("operation") == "SET"
        assert err.get_detail("key") == "user:123"
    
    def test_network_error(self):
        """Test network_error constructor."""
        err = network_error("HTTP GET")
        
        assert err.error_type == ErrorType.NETWORK
        assert err.code == code("NETWORK_ERROR")
        assert "HTTP GET" in err.message
        assert err.get_detail("operation") == "HTTP GET"
    
    def test_network_error_with_url(self):
        """Test network_error_with_url constructor."""
        err = network_error_with_url("HTTP GET", "https://api.example.com")
        
        assert err.error_type == ErrorType.NETWORK
        assert err.code == code("NETWORK_ERROR")
        assert "HTTP GET" in err.message
        assert "https://api.example.com" in err.message
        assert err.get_detail("operation") == "HTTP GET"
        assert err.get_detail("url") == "https://api.example.com"
    
    def test_event_error(self):
        """Test event_error constructor."""
        err = event_error("publish")
        
        assert err.error_type == ErrorType.EVENT
        assert err.code == code("EVENT_ERROR")
        assert "publish" in err.message
        assert err.get_detail("operation") == "publish"
    
    def test_event_error_with_subject(self):
        """Test event_error_with_subject constructor."""
        err = event_error_with_subject("publish", "user.created")
        
        assert err.error_type == ErrorType.EVENT
        assert err.code == code("EVENT_ERROR")
        assert "publish" in err.message
        assert "user.created" in err.message
        assert err.get_detail("operation") == "publish"
        assert err.get_detail("subject") == "user.created"


class TestWrappingConstructors:
    """Test error wrapping constructors."""
    
    def test_wrap(self):
        """Test wrap constructor."""
        root = database_error("query failed")
        wrapped = wrap(
            root,
            ErrorType.INTERNAL,
            code("OPERATION_FAILED"),
            "failed to get user"
        )
        
        assert wrapped.error_type == ErrorType.INTERNAL
        assert wrapped.code == code("OPERATION_FAILED")
        assert wrapped.message == "failed to get user"
        assert wrapped.cause == root
    
    def test_wrap_with_details(self):
        """Test wrap_with_details constructor."""
        root = network_error("connection failed")
        wrapped = wrap_with_details(
            root,
            ErrorType.UNAVAILABLE,
            code("SERVICE_DOWN"),
            "api unavailable",
            {"service": "payment-api", "attempt": "3"}
        )
        
        assert wrapped.error_type == ErrorType.UNAVAILABLE
        assert wrapped.code == code("SERVICE_DOWN")
        assert wrapped.message == "api unavailable"
        assert wrapped.cause == root
        assert wrapped.get_detail("service") == "payment-api"
        assert wrapped.get_detail("attempt") == "3"
    
    def test_wrap_chain(self):
        """Test wrapping creates proper error chain."""
        root = database_error("connection timeout")
        middle = wrap(root, ErrorType.INTERNAL, code("DB_OP_FAILED"), "db operation failed")
        top = wrap(middle, ErrorType.UNAVAILABLE, code("SERVICE_DOWN"), "service down")
        
        assert top.get_root_cause() == root
        assert top.cause == middle
        assert middle.cause == root


# =========================================
# Benchmarks
# =========================================

class TestConstructorBenchmarks:
    """Benchmark constructor functions."""
    
    def test_not_found_benchmark(self, benchmark):
        """Benchmark not_found constructor."""
        result = benchmark(not_found, "user", "123")
        assert result.error_type == ErrorType.NOT_FOUND
    
    def test_validation_benchmark(self, benchmark):
        """Benchmark validation constructor."""
        result = benchmark(validation, "invalid input")
        assert result.error_type == ErrorType.VALIDATION
    
    def test_validation_with_field_benchmark(self, benchmark):
        """Benchmark validation_with_field constructor."""
        result = benchmark(validation_with_field, "email", "invalid")
        assert result.error_type == ErrorType.VALIDATION
    
    def test_internal_benchmark(self, benchmark):
        """Benchmark internal constructor."""
        result = benchmark(internal, "internal error")
        assert result.error_type == ErrorType.INTERNAL
    
    def test_database_error_benchmark(self, benchmark):
        """Benchmark database_error constructor."""
        result = benchmark(database_error, "SELECT")
        assert result.error_type == ErrorType.DATABASE
    
    def test_wrap_benchmark(self, benchmark):
        """Benchmark wrap constructor."""
        root = database_error("query failed")
        result = benchmark(
            wrap,
            root,
            ErrorType.INTERNAL,
            code("OP_FAILED"),
            "operation failed"
        )
        assert result.cause == root
    
    def test_wrap_with_details_benchmark(self, benchmark):
        """Benchmark wrap_with_details constructor."""
        root = network_error("timeout")
        result = benchmark(
            wrap_with_details,
            root,
            ErrorType.UNAVAILABLE,
            code("SERVICE_DOWN"),
            "service down",
            {"service": "api", "attempt": "3"}
        )
        assert result.cause == root


# =========================================
# Integration Tests
# =========================================

class TestConstructorIntegration:
    """Integration tests for constructors."""
    
    def test_realistic_error_chain(self):
        """Test creating a realistic error chain."""
        # Simulate: network timeout → database error → service unavailable
        network_err = timeout("database connection", 5000)
        db_err = database_error_with_table("SELECT", "users", cause=network_err)
        service_err = unavailable_with_cause("user service", db_err)
        
        # Verify chain
        assert service_err.error_type == ErrorType.UNAVAILABLE
        assert service_err.cause == db_err
        assert db_err.cause == network_err
        assert service_err.get_root_cause() == network_err
        
        # Verify details preserved
        assert db_err.get_detail("table") == "users"
        assert network_err.get_detail("timeout_ms") == "5000"
    
    def test_validation_error_workflow(self):
        """Test validation error construction workflow."""
        # Simulate multiple validation errors
        errors = [
            required_field("email"),
            required_field("password"),
            invalid_field("age", "must be >= 18"),
        ]
        
        assert all(e.error_type == ErrorType.VALIDATION for e in errors)
        assert all(e.is_client_error() for e in errors)
        assert all(not e.is_retryable() for e in errors)
    
    def test_infrastructure_error_workflow(self):
        """Test infrastructure error construction workflow."""
        # Simulate infrastructure issues
        errors = [
            database_error_with_table("INSERT", "events"),
            cache_error_with_key("GET", "session:abc123"),
            network_error_with_url("POST", "https://api.example.com/webhook"),
            event_error_with_subject("publish", "order.created"),
        ]
        
        assert all(e.is_server_error() for e in errors)
        assert all(e.is_retryable() for e in errors)