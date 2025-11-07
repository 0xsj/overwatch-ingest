package config

import (
	"fmt"
	"strings"

	"github.com/0xsj/scout/platform/pkg/errors"
)

const (
	CodeMissingRequired errors.Code = "CONFIG_MISSING_REQUIRED"
	CodeInvalidValue    errors.Code = "CONFIG_INVALID_VALUE"
	CodeInvalidFormat   errors.Code = "CONFIG_INVALID_FORMAT"
	CodeOutOfRange      errors.Code = "CONFIG_OUT_OF_RANGE"
	CodeInvalidChoice   errors.Code = "CONFIG_INVALID_CHOICE"
)

func MissingRequired(key string) *errors.Error {
	return errors.New(
		errors.ErrorTypeValidation,
		CodeMissingRequired,
		fmt.Sprintf("required configuration not found: %s", key),
	).WithDetail("key", key)
}

// InvalidValue creates an error for a configuration value that cannot be parsed or converted.
//
// Example:
//
//	err := config.InvalidValue("PORT", "abc", "not a valid integer")
//	// Error: invalid configuration value for PORT: not a valid integer (got: abc)
func InvalidValue(key, value string, reason string) *errors.Error {
	return errors.New(
		errors.ErrorTypeValidation,
		CodeInvalidValue,
		fmt.Sprintf("invalid configuration value for %s: %s (got: %s)", key, reason, value),
	).WithDetail("key", key).
		WithDetail("value", value).
		WithDetail("reason", reason)
}

// InvalidFormat creates an error for a configuration value with incorrect format.
//
// Example:
//
//	err := config.InvalidFormat("DATABASE_URL", "localhost:5432", "postgresql://host:port/db")
//	// Error: invalid format for DATABASE_URL: expected postgresql://host:port/db (got: localhost:5432)
func InvalidFormat(key, value, expectedFormat string) *errors.Error {
	return errors.New(
		errors.ErrorTypeValidation,
		CodeInvalidFormat,
		fmt.Sprintf("invalid format for %s: expected %s (got: %s)", key, expectedFormat, value),
	).WithDetail("key", key).
		WithDetail("value", value).
		WithDetail("expected_format", expectedFormat)
}

// OutOfRange creates an error for a numeric configuration value outside the valid range.
//
// Example:
//
//	err := config.OutOfRange("PORT", 70000, 1024, 65535)
//	// Error: PORT value out of range: must be between 1024 and 65535 (got: 70000)
func OutOfRange(key string, value, min, max interface{}) *errors.Error {
	return errors.New(
		errors.ErrorTypeValidation,
		CodeOutOfRange,
		fmt.Sprintf("%s value out of range: must be between %v and %v (got: %v)", key, min, max, value),
	).WithDetail("key", key).
		WithDetail("value", fmt.Sprintf("%v", value)).
		WithDetail("min", fmt.Sprintf("%v", min)).
		WithDetail("max", fmt.Sprintf("%v", max))
}

// InvalidChoice creates an error for a configuration value not in the allowed set.
//
// Example:
//
//	err := config.InvalidChoice("LOG_LEVEL", "verbose", []string{"debug", "info", "warn", "error"})
//	// Error: invalid value for LOG_LEVEL: must be one of [debug, info, warn, error] (got: verbose)
func InvalidChoice(key, value string, allowedValues []string) *errors.Error {
	allowed := strings.Join(allowedValues, ", ")
	return errors.New(
		errors.ErrorTypeValidation,
		CodeInvalidChoice,
		fmt.Sprintf("invalid value for %s: must be one of [%s] (got: %s)", key, allowed, value),
	).WithDetail("key", key).
		WithDetail("value", value).
		WithDetail("allowed_values", allowed)
}
