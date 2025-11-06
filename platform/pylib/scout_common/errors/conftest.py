"""Pytest configuration and fixtures for error tests."""

import pytest
from typing import Any

from .base import Error
from .types import ErrorType
from .codes import code
from .result import Ok, Err, Result


# =========================================
# Fixtures
# =========================================

@pytest.fixture
def simple_error() -> Error:
    """Simple error for testing."""
    return Error(
        error_type=ErrorType.VALIDATION,
        code=code("TEST_CODE"),
        message="test message",
    )


@pytest.fixture
def error_with_details() -> Error:
    """Error with details."""
    return Error(
        error_type=ErrorType.VALIDATION,
        code=code("TEST_CODE"),
        message="test message",
        details={"field": "email", "reason": "invalid"},
    )


@pytest.fixture
def error_with_cause(simple_error: Error) -> Error:
    """Error with cause chain."""
    root = Error(
        error_type=ErrorType.DATABASE,
        code=code("DB_ERROR"),
        message="database failed",
    )
    return Error(
        error_type=ErrorType.INTERNAL,
        code=code("WRAPPER"),
        message="operation failed",
        cause=root,
    )


@pytest.fixture
def ok_result() -> Result[int, Error]:
    """Ok result for testing."""
    return Ok(42)


@pytest.fixture
def err_result(simple_error: Error) -> Result[int, Error]:
    """Err result for testing."""
    return Err(simple_error)