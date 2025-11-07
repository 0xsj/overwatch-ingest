// platform/pkg/config/doc.go

/*
Package config provides utilities for loading and validating configuration from environment variables.

# Overview

This package offers type-safe environment variable loading with built-in validation,
clear error messages, and support for various data types. It is designed to be used
by other platform packages and domain services to load their configuration.

# Core Concepts

The config package provides three main categories of utilities:

1. **Parser** - Convert string values to typed values (int, bool, duration, URL, etc.)
2. **Validator** - Validate typed values (required, range, choice, format, etc.)
3. **Loader** - Load from environment with parsing and validation in one step

# Basic Usage

Loading configuration values:

	// Required values - returns error if missing or invalid
	port, err := config.LoadPortRequired("GATEWAY_PORT")
	dbHost, err := config.LoadStringRequired("DATABASE_HOST")

	// Optional values with defaults
	logLevel := config.LoadStringOptional("LOG_LEVEL", "info")
	timeout, err := config.LoadDurationOptional("TIMEOUT", 30*time.Second)

	// With validation
	maxConns, err := config.LoadIntWithRange("MAX_CONNECTIONS", 1, 1000, ptr(100))
	env, err := config.LoadStringWithChoice("ENV", []string{"dev", "staging", "prod"}, ptr("dev"))

# Using Prefixes

For service-specific configuration, use prefixes to avoid collisions:

	port, err := config.LoadPortRequired(config.WithPrefix("GATEWAY_", "PORT"))
	// Loads from GATEWAY_PORT environment variable

# Error Handling

All errors use the platform errors package and include rich metadata:

	port, err := config.LoadPortRequired("PORT")
	if err != nil {
		// Error will be one of:
		// - MissingRequired: required config not found
		// - InvalidValue: value cannot be parsed
		// - OutOfRange: numeric value outside valid range
		// - InvalidChoice: value not in allowed set
		// - InvalidFormat: value doesn't match expected format

		// Errors include details for debugging
		fmt.Println(errors.GetDetail(err, "key"))    // "PORT"
		fmt.Println(errors.GetDetail(err, "value"))  // "abc"
		fmt.Println(errors.GetDetail(err, "reason")) // "not a valid integer"
	}

# Supported Types

The package supports loading the following types:

- string - Basic string values
- int, int64 - Integer numbers
- bool - Boolean flags (true/false, 1/0, yes/no, on/off)
- float64 - Floating point numbers
- time.Duration - Duration strings (5s, 10m, 1h)
- url.URL - URL strings with validation
- []string - Comma-separated lists

# Validation Helpers

The package provides standalone validation functions for custom use cases:

	// Validate after loading
	value := os.Getenv("PORT")
	port, _ := config.ParseInt(value)
	if err := config.ValidatePort("PORT", port); err != nil {
		return err
	}

	// Range validation
	if err := config.ValidateRange("WORKERS", workers, 1, 100); err != nil {
		return err
	}

	// Choice validation
	if err := config.ValidateChoice("ENV", env, []string{"dev", "prod"}); err != nil {
		return err
	}

# Best Practices

1. Load all configuration at startup - fail fast if misconfigured
2. Use service-specific prefixes to avoid environment variable collisions
3. Validate all loaded values before using them
4. Provide sensible defaults for optional configuration
5. Use LoadXxxRequired for critical configuration that must be present
6. Use LoadXxxOptional for configuration that has reasonable defaults

# Example: Service Configuration

	package config

	import (
		"time"

		"github.com/0xsj/scout/platform/pkg/config"
	)

	type Config struct {
		Port         int
		LogLevel     string
		DatabaseHost string
		DatabasePort int
		Timeout      time.Duration
	}

	func Load() (*Config, error) {
		prefix := "GATEWAY_"

		port, err := config.LoadPortRequired(config.WithPrefix(prefix, "PORT"))
		if err != nil {
			return nil, err
		}

		logLevel := config.LoadStringOptional(config.WithPrefix(prefix, "LOG_LEVEL"), "info")

		dbHost, err := config.LoadStringRequired(config.WithPrefix(prefix, "DATABASE_HOST"))
		if err != nil {
			return nil, err
		}

		dbPort, err := config.LoadPortOptional(config.WithPrefix(prefix, "DATABASE_PORT"), 5432)
		if err != nil {
			return nil, err
		}

		timeout, err := config.LoadDurationOptional(
			config.WithPrefix(prefix, "TIMEOUT"),
			30*time.Second,
		)
		if err != nil {
			return nil, err
		}

		return &Config{
			Port:         port,
			LogLevel:     logLevel,
			DatabaseHost: dbHost,
			DatabasePort: dbPort,
			Timeout:      timeout,
		}, nil
	}

For more examples and detailed documentation, see the individual function documentation.
*/
package config