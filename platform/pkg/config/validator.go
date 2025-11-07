// platform/pkg/config/validator.go
package config

import (
	"fmt"
	"regexp"
)

// ValidateRequired checks if a value is non-empty.
func ValidateRequired(key, value string) error {
	if value == "" {
		return MissingRequired(key)
	}
	return nil
}

// ValidateRange checks if a numeric value is within the specified range (inclusive).
func ValidateRange[T int | int64 | float64](key string, value, min, max T) error {
	if value < min || value > max {
		return OutOfRange(key, value, min, max)
	}
	return nil
}

// ValidateMinMax checks if a numeric value satisfies min/max constraints.
// Pass nil for min or max to skip that check.
func ValidateMinMax[T int | int64 | float64](key string, value T, min, max *T) error {
	if min != nil && value < *min {
		return OutOfRange(key, value, *min, "infinity")
	}
	if max != nil && value > *max {
		return OutOfRange(key, value, "infinity", *max)
	}
	return nil
}

// ValidateChoice checks if a value is in the allowed set.
func ValidateChoice(key, value string, allowedValues []string) error {
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	return InvalidChoice(key, value, allowedValues)
}

// ValidatePattern checks if a value matches the given regex pattern.
func ValidatePattern(key, value, pattern string) error {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return InvalidFormat(key, value, fmt.Sprintf("regex pattern: %s", pattern))
	}
	if !matched {
		return InvalidFormat(key, value, fmt.Sprintf("pattern: %s", pattern))
	}
	return nil
}

// ValidateURL checks if a value is a valid URL with required scheme.
func ValidateURL(key, value string) error {
	_, err := ParseURL(value)
	if err != nil {
		return InvalidFormat(key, value, "valid URL with scheme")
	}
	return nil
}

// ValidatePort checks if a port number is in the valid range (1-65535).
func ValidatePort(key string, port int) error {
	return ValidateRange(key, port, 1, 65535)
}

// ValidateNonZero checks if a numeric value is non-zero.
func ValidateNonZero[T int | int64 | float64](key string, value T) error {
	if value == 0 {
		return InvalidValue(key, "0", "value must be non-zero")
	}
	return nil
}

// ValidatePositive checks if a numeric value is positive (> 0).
func ValidatePositive[T int | int64 | float64](key string, value T) error {
	if value <= 0 {
		return InvalidValue(key, fmt.Sprintf("%v", value), "value must be positive")
	}
	return nil
}

// ValidateNonNegative checks if a numeric value is non-negative (>= 0).
func ValidateNonNegative[T int | int64 | float64](key string, value T) error {
	if value < 0 {
		return InvalidValue(key, fmt.Sprintf("%v", value), "value must be non-negative")
	}
	return nil
}

// ValidateMinLength checks if a string meets minimum length requirement.
func ValidateMinLength(key, value string, minLength int) error {
	if len(value) < minLength {
		return InvalidValue(
			key,
			value,
			fmt.Sprintf("must be at least %d characters", minLength),
		)
	}
	return nil
}

// ValidateMaxLength checks if a string does not exceed maximum length.
func ValidateMaxLength(key, value string, maxLength int) error {
	if len(value) > maxLength {
		return InvalidValue(
			key,
			value,
			fmt.Sprintf("must be at most %d characters", maxLength),
		)
	}
	return nil
}