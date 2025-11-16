// platform/pkg/observability/tracing/context.go
package tracing

import (
	"context"
)

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const (
	requestIDKey contextKey = "request_id"
)

// WithRequestID returns a new context with the request ID attached.
func WithRequestID(ctx context.Context, requestID RequestID) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// GetRequestID extracts the request ID from the context.
// Returns empty RequestID if not found.
func GetRequestID(ctx context.Context) RequestID {
	if ctx == nil {
		return ""
	}
	
	if rid, ok := ctx.Value(requestIDKey).(RequestID); ok {
		return rid
	}
	
	return ""
}

// GetRequestIDOrGenerate extracts the request ID from the context,
// or generates a new one if not found.
func GetRequestIDOrGenerate(ctx context.Context) RequestID {
	rid := GetRequestID(ctx)
	if rid == "" {
		rid = NewRequestID()
	}
	return rid
}

// HasRequestID checks if the context contains a request ID.
func HasRequestID(ctx context.Context) bool {
	return GetRequestID(ctx) != ""
}

// WithNewRequestID returns a new context with a newly generated request ID.
func WithNewRequestID(ctx context.Context) context.Context {
	return WithRequestID(ctx, NewRequestID())
}