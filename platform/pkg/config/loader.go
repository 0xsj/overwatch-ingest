// platform/pkg/config/loader.go
package config

import (
	"fmt"
	"net/url"
	"os"
	"time"
)

// LoadStringRequired loads a required string from environment.
func LoadStringRequired(key string) (string, error) {
	value := os.Getenv(key)
	if err := ValidateRequired(key, value); err != nil {
		return "", err
	}
	return ParseString(value), nil
}

// LoadStringOptional loads an optional string from environment with default.
func LoadStringOptional(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return ParseString(value)
}

// LoadIntRequired loads a required integer from environment.
func LoadIntRequired(key string) (int, error) {
	value := os.Getenv(key)
	if err := ValidateRequired(key, value); err != nil {
		return 0, err
	}
	
	parsed, err := ParseInt(value)
	if err != nil {
		return 0, InvalidValue(key, value, "not a valid integer")
	}
	
	return parsed, nil
}

// LoadIntOptional loads an optional integer from environment with default.
func LoadIntOptional(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	
	parsed, err := ParseInt(value)
	if err != nil {
		return 0, InvalidValue(key, value, "not a valid integer")
	}
	
	return parsed, nil
}

// LoadIntWithRange loads an integer and validates it's within range.
func LoadIntWithRange(key string, min, max int, defaultValue *int) (int, error) {
	value := os.Getenv(key)
	
	// If not set and default provided, use default
	if value == "" {
		if defaultValue != nil {
			if err := ValidateRange(key, *defaultValue, min, max); err != nil {
				return 0, err
			}
			return *defaultValue, nil
		}
		return 0, MissingRequired(key)
	}
	
	parsed, err := ParseInt(value)
	if err != nil {
		return 0, InvalidValue(key, value, "not a valid integer")
	}
	
	if err := ValidateRange(key, parsed, min, max); err != nil {
		return 0, err
	}
	
	return parsed, nil
}

// LoadBoolRequired loads a required boolean from environment.
func LoadBoolRequired(key string) (bool, error) {
	value := os.Getenv(key)
	if err := ValidateRequired(key, value); err != nil {
		return false, err
	}
	
	parsed, err := ParseBool(value)
	if err != nil {
		return false, InvalidValue(key, value, "not a valid boolean")
	}
	
	return parsed, nil
}

// LoadBoolOptional loads an optional boolean from environment with default.
func LoadBoolOptional(key string, defaultValue bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	
	parsed, err := ParseBool(value)
	if err != nil {
		return false, InvalidValue(key, value, "not a valid boolean")
	}
	
	return parsed, nil
}

// LoadFloat64Required loads a required float64 from environment.
func LoadFloat64Required(key string) (float64, error) {
	value := os.Getenv(key)
	if err := ValidateRequired(key, value); err != nil {
		return 0, err
	}
	
	parsed, err := ParseFloat64(value)
	if err != nil {
		return 0, InvalidValue(key, value, "not a valid float")
	}
	
	return parsed, nil
}

// LoadFloat64Optional loads an optional float64 from environment with default.
func LoadFloat64Optional(key string, defaultValue float64) (float64, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	
	parsed, err := ParseFloat64(value)
	if err != nil {
		return 0, InvalidValue(key, value, "not a valid float")
	}
	
	return parsed, nil
}

// LoadDurationRequired loads a required duration from environment.
func LoadDurationRequired(key string) (time.Duration, error) {
	value := os.Getenv(key)
	if err := ValidateRequired(key, value); err != nil {
		return 0, err
	}
	
	parsed, err := ParseDuration(value)
	if err != nil {
		return 0, InvalidValue(key, value, "not a valid duration")
	}
	
	return parsed, nil
}

// LoadDurationOptional loads an optional duration from environment with default.
func LoadDurationOptional(key string, defaultValue time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	
	parsed, err := ParseDuration(value)
	if err != nil {
		return 0, InvalidValue(key, value, "not a valid duration")
	}
	
	return parsed, nil
}

// LoadURLRequired loads a required URL from environment.
func LoadURLRequired(key string) (*url.URL, error) {
	value := os.Getenv(key)
	if err := ValidateRequired(key, value); err != nil {
		return nil, err
	}
	
	parsed, err := ParseURL(value)
	if err != nil {
		return nil, InvalidFormat(key, value, "valid URL with scheme")
	}
	
	return parsed, nil
}

// LoadURLOptional loads an optional URL from environment with default.
func LoadURLOptional(key string, defaultValue *url.URL) (*url.URL, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	
	parsed, err := ParseURL(value)
	if err != nil {
		return nil, InvalidFormat(key, value, "valid URL with scheme")
	}
	
	return parsed, nil
}

// LoadStringSliceRequired loads a required comma-separated list from environment.
func LoadStringSliceRequired(key string) ([]string, error) {
	value := os.Getenv(key)
	if err := ValidateRequired(key, value); err != nil {
		return nil, err
	}
	
	parsed := ParseStringSlice(value, ",")
	if len(parsed) == 0 {
		return nil, InvalidValue(key, value, "list is empty after parsing")
	}
	
	return parsed, nil
}

// LoadStringSliceOptional loads an optional comma-separated list from environment.
func LoadStringSliceOptional(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	
	parsed := ParseStringSlice(value, ",")
	if len(parsed) == 0 {
		return defaultValue
	}
	
	return parsed
}

// LoadStringWithChoice loads a string and validates it's in allowed set.
func LoadStringWithChoice(key string, allowedValues []string, defaultValue *string) (string, error) {
	value := os.Getenv(key)
	
	// If not set and default provided, use default
	if value == "" {
		if defaultValue != nil {
			if err := ValidateChoice(key, *defaultValue, allowedValues); err != nil {
				return "", err
			}
			return *defaultValue, nil
		}
		return "", MissingRequired(key)
	}
	
	parsed := ParseString(value)
	if err := ValidateChoice(key, parsed, allowedValues); err != nil {
		return "", err
	}
	
	return parsed, nil
}

// LoadPortRequired loads a required port number from environment.
func LoadPortRequired(key string) (int, error) {
	port, err := LoadIntRequired(key)
	if err != nil {
		return 0, err
	}
	
	if err := ValidatePort(key, port); err != nil {
		return 0, err
	}
	
	return port, nil
}

// LoadPortOptional loads an optional port number from environment with default.
func LoadPortOptional(key string, defaultValue int) (int, error) {
	// Validate default
	if err := ValidatePort(key, defaultValue); err != nil {
		return 0, fmt.Errorf("invalid default port: %w", err)
	}
	
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	
	port, err := ParseInt(value)
	if err != nil {
		return 0, InvalidValue(key, value, "not a valid integer")
	}
	
	if err := ValidatePort(key, port); err != nil {
		return 0, err
	}
	
	return port, nil
}

// WithPrefix returns a prefixed key for environment variables.
// Example: WithPrefix("GATEWAY_", "PORT") returns "GATEWAY_PORT"
func WithPrefix(prefix, key string) string {
	return prefix + key
}