"""
Configuration loading and validation utilities for Scout platform.

This package provides type-safe environment variable loading with built-in validation,
clear error messages, and support for various data types. It is designed to be used
by other platform packages and domain services to load their configuration.

Core Concepts
-------------
The config package provides three main categories of utilities:

1. **Parser** - Convert string values to typed values (int, bool, duration, URL, etc.)
2. **Validator** - Validate typed values (required, range, choice, format, etc.)
3. **Loader** - Load from environment with parsing and validation in one step

Basic Usage
-----------
Loading configuration values:

    from scout_common.config import (
        load_port_required,
        load_string_required,
        load_string_optional,
        load_duration_optional,
        load_int_with_range,
        load_string_with_choice,
    )
    from datetime import timedelta
    
    # Required values - raises Error if missing or invalid
    port = load_port_required("TOOLS_PORT")
    db_host = load_string_required("DATABASE_HOST")
    
    # Optional values with defaults
    log_level = load_string_optional("LOG_LEVEL", "info")
    timeout = load_duration_optional("TIMEOUT", timedelta(seconds=30))
    
    # With validation
    max_conns = load_int_with_range("MAX_CONNECTIONS", 1, 1000, default_value=100)
    env = load_string_with_choice("ENV", ["dev", "staging", "prod"], default_value="dev")

Using Prefixes
--------------
For service-specific configuration, use prefixes to avoid collisions:

    from scout_common.config import with_prefix, load_port_required
    
    port = load_port_required(with_prefix("TOOLS_", "PORT"))
    # Loads from TOOLS_PORT environment variable

Error Handling
--------------
All errors use the platform errors package and include rich metadata:

    from scout_common.config import load_port_required
    from scout_common.errors import Error
    
    try:
        port = load_port_required("PORT")
    except Error as e:
        # Error will be one of:
        # - missing_required: required config not found
        # - invalid_value: value cannot be parsed
        # - out_of_range: numeric value outside valid range
        # - invalid_choice: value not in allowed set
        # - invalid_format: value doesn't match expected format
        
        # Errors include details for debugging
        print(e.get_detail("key"))     # "PORT"
        print(e.get_detail("value"))   # "abc"
        print(e.get_detail("reason"))  # "not a valid integer"

Supported Types
---------------
The package supports loading the following types:

- str - Basic string values
- int - Integer numbers
- bool - Boolean flags (true/false, 1/0, yes/no, on/off)
- float - Floating point numbers
- timedelta - Duration strings (5s, 10m, 1h)
- ParseResult - URL strings with validation
- list[str] - Comma-separated lists

Validation Helpers
------------------
The package provides standalone validation functions for custom use cases:

    import os
    from scout_common.config import parse_int, validate_port, validate_range, validate_choice
    
    # Validate after loading
    value = os.getenv("PORT", "")
    port = parse_int(value)
    validate_port("PORT", port)
    
    # Range validation
    validate_range("WORKERS", workers, 1, 100)
    
    # Choice validation
    validate_choice("ENV", env, ["dev", "prod"])

Example: Service Configuration
-------------------------------

    from dataclasses import dataclass
    from datetime import timedelta
    
    from scout_common.config import (
        load_port_required,
        load_port_optional,
        load_string_required,
        load_string_optional,
        load_duration_optional,
        with_prefix,
    )
    
    
    @dataclass
    class Config:
        port: int
        log_level: str
        database_host: str
        database_port: int
        timeout: timedelta
    
    
    def load() -> Config:
        prefix = "TOOLS_"
        
        port = load_port_required(with_prefix(prefix, "PORT"))
        log_level = load_string_optional(with_prefix(prefix, "LOG_LEVEL"), "info")
        db_host = load_string_required(with_prefix(prefix, "DATABASE_HOST"))
        db_port = load_port_optional(with_prefix(prefix, "DATABASE_PORT"), 5432)
        timeout = load_duration_optional(
            with_prefix(prefix, "TIMEOUT"),
            timedelta(seconds=30),
        )
        
        return Config(
            port=port,
            log_level=log_level,
            database_host=db_host,
            database_port=db_port,
            timeout=timeout,
        )
"""

__version__ = "0.1.0"

# Error constructors
from .errors import (
    CODE_MISSING_REQUIRED,
    CODE_INVALID_VALUE,
    CODE_INVALID_FORMAT,
    CODE_OUT_OF_RANGE,
    CODE_INVALID_CHOICE,
    missing_required,
    invalid_value,
    invalid_format,
    out_of_range,
    invalid_choice,
)

# Parsers
from .parser import (
    parse_string,
    parse_int,
    parse_bool,
    parse_float,
    parse_duration,
    parse_url,
    parse_string_list,
)

# Validators
from .validator import (
    validate_required,
    validate_range,
    validate_min_max,
    validate_choice,
    validate_pattern,
    validate_url,
    validate_port,
    validate_non_zero,
    validate_positive,
    validate_non_negative,
    validate_min_length,
    validate_max_length,
)

# Loaders
from .loader import (
    load_string_required,
    load_string_optional,
    load_int_required,
    load_int_optional,
    load_int_with_range,
    load_bool_required,
    load_bool_optional,
    load_float_required,
    load_float_optional,
    load_duration_required,
    load_duration_optional,
    load_url_required,
    load_url_optional,
    load_string_list_required,
    load_string_list_optional,
    load_string_with_choice,
    load_port_required,
    load_port_optional,
    with_prefix,
)

__all__ = [
    # Version
    "__version__",
    # Error codes
    "CODE_MISSING_REQUIRED",
    "CODE_INVALID_VALUE",
    "CODE_INVALID_FORMAT",
    "CODE_OUT_OF_RANGE",
    "CODE_INVALID_CHOICE",
    # Error constructors
    "missing_required",
    "invalid_value",
    "invalid_format",
    "out_of_range",
    "invalid_choice",
    # Parsers
    "parse_string",
    "parse_int",
    "parse_bool",
    "parse_float",
    "parse_duration",
    "parse_url",
    "parse_string_list",
    # Validators
    "validate_required",
    "validate_range",
    "validate_min_max",
    "validate_choice",
    "validate_pattern",
    "validate_url",
    "validate_port",
    "validate_non_zero",
    "validate_positive",
    "validate_non_negative",
    "validate_min_length",
    "validate_max_length",
    # Loaders
    "load_string_required",
    "load_string_optional",
    "load_int_required",
    "load_int_optional",
    "load_int_with_range",
    "load_bool_required",
    "load_bool_optional",
    "load_float_required",
    "load_float_optional",
    "load_duration_required",
    "load_duration_optional",
    "load_url_required",
    "load_url_optional",
    "load_string_list_required",
    "load_string_list_optional",
    "load_string_with_choice",
    "load_port_required",
    "load_port_optional",
    "with_prefix",
]