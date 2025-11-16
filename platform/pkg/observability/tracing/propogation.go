// platform/pkg/observability/tracing/propagation.go
package tracing

import (
	"context"
	"net/http"
)

// InjectHTTP injects the request ID into HTTP headers.
// Used by HTTP clients to propagate request IDs.
func InjectHTTP(header http.Header, requestID RequestID) {
	if requestID == "" {
		return
	}
	header.Set(RequestIDHeader, requestID.String())
}

// ExtractHTTP extracts the request ID from HTTP headers.
// Used by HTTP servers to receive request IDs.
func ExtractHTTP(header http.Header) RequestID {
	value := header.Get(RequestIDHeader)
	if value == "" {
		return ""
	}
	return ParseRequestID(value)
}

// HTTPMiddleware is HTTP middleware that extracts or generates a request ID
// and attaches it to the request context.
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract from header or generate new
		rid := ExtractHTTP(r.Header)
		if rid == "" {
			rid = NewRequestID()
		}

		// Attach to context
		ctx := WithRequestID(r.Context(), rid)
		r = r.WithContext(ctx)

		// Add to response headers for clients
		w.Header().Set(RequestIDHeader, rid.String())

		next.ServeHTTP(w, r)
	})
}

// InjectEventMetadata injects the request ID into event metadata.
// Used when publishing events to NATS/RabbitMQ.
func InjectEventMetadata(metadata map[string]string, requestID RequestID) {
	if requestID == "" {
		return
	}
	if metadata == nil {
		return
	}
	metadata[RequestIDHeader] = requestID.String()
}

// ExtractEventMetadata extracts the request ID from event metadata.
// Used when consuming events from NATS/RabbitMQ.
func ExtractEventMetadata(metadata map[string]string) RequestID {
	if metadata == nil {
		return ""
	}
	value, ok := metadata[RequestIDHeader]
	if !ok {
		return ""
	}
	return ParseRequestID(value)
}

// PropagateToEvent extracts the request ID from context and injects it into event metadata.
// Returns the request ID for convenience.
func PropagateToEvent(ctx context.Context, metadata map[string]string) RequestID {
	rid := GetRequestID(ctx)
	if rid == "" {
		rid = NewRequestID()
	}
	InjectEventMetadata(metadata, rid)
	return rid
}

// ExtractFromEvent extracts the request ID from event metadata and attaches it to context.
func ExtractFromEvent(ctx context.Context, metadata map[string]string) context.Context {
	rid := ExtractEventMetadata(metadata)
	if rid == "" {
		rid = NewRequestID()
	}
	return WithRequestID(ctx, rid)
}