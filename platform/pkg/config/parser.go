// platform/pkg/config/parser.go
package config

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ParseString trims whitespace and returns the string value.
func ParseString(value string) string {
	return strings.TrimSpace(value)
}

// ParseInt parses a string to an integer.
func ParseInt(value string) (int, error) {
	trimmed := strings.TrimSpace(value)
	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

// ParseInt64 parses a string to an int64.
func ParseInt64(value string) (int64, error) {
	trimmed := strings.TrimSpace(value)
	parsed, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

// ParseBool parses a string to a boolean.
// Accepts: "true", "false", "1", "0", "yes", "no", "on", "off" (case-insensitive).
func ParseBool(value string) (bool, error) {
	trimmed := strings.TrimSpace(strings.ToLower(value))

	switch trimmed {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, strconv.ErrSyntax
	}
}

// ParseFloat64 parses a string to a float64.
func ParseFloat64(value string) (float64, error) {
	trimmed := strings.TrimSpace(value)
	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

// ParseDuration parses a string to a time.Duration.
// Accepts formats like "5s", "10m", "1h30m".
func ParseDuration(value string) (time.Duration, error) {
	trimmed := strings.TrimSpace(value)
	duration, err := time.ParseDuration(trimmed)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

// ParseURL parses and validates a URL string.
func ParseURL(value string) (*url.URL, error) {
	trimmed := strings.TrimSpace(value)
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, err
	}

	// Validate that we have a scheme and host
	if parsed.Scheme == "" {
		return nil, &url.Error{Op: "parse", URL: trimmed, Err: strconv.ErrSyntax}
	}

	return parsed, nil
}

// ParseStringSlice parses a separated string into a slice.
// Trims whitespace from each element and filters out empty strings.
func ParseStringSlice(value string, separator string) []string {
	if value == "" {
		return []string{}
	}

	parts := strings.Split(value, separator)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
