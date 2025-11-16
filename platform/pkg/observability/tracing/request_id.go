// platform/pkg/tracing/request_id.go
package tracing

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// RequestIDPrefix is the prefix for all request IDs
	RequestIDPrefix = "req_"
	
	// RequestIDHeader is the HTTP/gRPC header name for request ID
	RequestIDHeader = "X-Request-ID"
)

// RequestID represents a unique request identifier.
type RequestID string

// NewRequestID generates a new request ID.
// Format: req_<timestamp>_<uuid>
func NewRequestID() RequestID {
	timestamp := time.Now().Unix()
	id := uuid.New().String()
	// Take first 8 chars of UUID for brevity
	shortID := strings.ReplaceAll(id, "-", "")[:8]
	return RequestID(fmt.Sprintf("%s%d_%s", RequestIDPrefix, timestamp, shortID))
}

// String returns the string representation of the request ID.
func (r RequestID) String() string {
	return string(r)
}

// IsValid checks if the request ID is valid.
func (r RequestID) IsValid() bool {
	if r == "" {
		return false
	}
	return strings.HasPrefix(string(r), RequestIDPrefix)
}

// ParseRequestID parses a string into a RequestID.
// Returns empty RequestID if invalid.
func ParseRequestID(s string) RequestID {
	rid := RequestID(s)
	if !rid.IsValid() {
		return ""
	}
	return rid
}

// MustParseRequestID parses a string into a RequestID.
// Panics if invalid.
func MustParseRequestID(s string) RequestID {
	rid := ParseRequestID(s)
	if rid == "" {
		panic(fmt.Sprintf("invalid request ID: %s", s))
	}
	return rid
}