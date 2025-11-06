"""Tests for Code type."""

import pytest
from hypothesis import given, strategies as st

from .codes import Code, code, is_empty


class TestCode:
    """Test Code type."""
    
    def test_code_creation(self):
        """Test creating codes."""
        c = code("TEST_CODE")
        assert c == "TEST_CODE"
        assert isinstance(c, str)
    
    def test_code_equality(self):
        """Test code equality."""
        c1 = code("TEST_CODE")
        c2 = code("TEST_CODE")
        assert c1 == c2
    
    def test_code_inequality(self):
        """Test code inequality."""
        c1 = code("CODE_A")
        c2 = code("CODE_B")
        assert c1 != c2
    
    def test_is_empty_true(self):
        """Test is_empty with empty code."""
        c = code("")
        assert is_empty(c)
    
    def test_is_empty_false(self):
        """Test is_empty with non-empty code."""
        c = code("NOT_EMPTY")
        assert not is_empty(c)
    
    def test_code_string_conversion(self):
        """Test that codes work as strings."""
        c = code("TEST_CODE")
        assert str(c) == "TEST_CODE"
        assert c.upper() == "TEST_CODE"
        assert c.lower() == "test_code"
    
    def test_code_in_dict(self):
        """Test that codes can be dict keys."""
        c1 = code("KEY_A")
        c2 = code("KEY_B")
        
        d = {c1: "value_a", c2: "value_b"}
        assert d[c1] == "value_a"
        assert d[c2] == "value_b"
    
    def test_code_in_set(self):
        """Test that codes can be in sets."""
        c1 = code("CODE_A")
        c2 = code("CODE_B")
        c3 = code("CODE_A")  # Duplicate
        
        s = {c1, c2, c3}
        assert len(s) == 2  # Duplicate removed
        assert c1 in s
        assert c2 in s
    
    @given(st.text())
    def test_code_from_any_string(self, value: str):
        """Property test: code should accept any string."""
        c = code(value)
        assert isinstance(c, str)
        assert c == value


# =========================================
# Benchmarks
# =========================================

class TestCodeBenchmarks:
    """Benchmark Code operations."""
    
    def test_code_creation_benchmark(self, benchmark):
        """Benchmark code creation."""
        result = benchmark(code, "TEST_CODE")
        assert result == "TEST_CODE"
    
    def test_is_empty_benchmark(self, benchmark):
        """Benchmark is_empty check."""
        c = code("TEST_CODE")
        result = benchmark(is_empty, c)
        assert result is False
    
    def test_code_comparison_benchmark(self, benchmark):
        """Benchmark code comparison."""
        c1 = code("CODE_A")
        c2 = code("CODE_B")
        result = benchmark(lambda: c1 == c2)
        assert result is False