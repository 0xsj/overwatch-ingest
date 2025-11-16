// platform/pkg/observability/tracing/metadata.go
package tracing

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// InjectRequestID injects the request ID into gRPC outgoing metadata.
// This is used by gRPC clients to propagate request IDs to downstream services.
func InjectRequestID(ctx context.Context, requestID RequestID) context.Context {
	if requestID == "" {
		return ctx
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		// Clone to avoid modifying shared metadata
		md = md.Copy()
	}

	md.Set(RequestIDHeader, requestID.String())
	return metadata.NewOutgoingContext(ctx, md)
}

// ExtractRequestID extracts the request ID from gRPC incoming metadata.
// This is used by gRPC servers to receive request IDs from upstream services.
// Returns empty RequestID if not found.
func ExtractRequestID(ctx context.Context) RequestID {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(RequestIDHeader)
	if len(values) == 0 {
		return ""
	}

	return ParseRequestID(values[0])
}

// ExtractOrGenerateRequestID extracts the request ID from gRPC metadata,
// or generates a new one if not found, and attaches it to the context.
func ExtractOrGenerateRequestID(ctx context.Context) (context.Context, RequestID) {
	rid := ExtractRequestID(ctx)
	if rid == "" {
		rid = NewRequestID()
	}
	return WithRequestID(ctx, rid), rid
}

// PropagateRequestID extracts the request ID from incoming metadata
// and injects it into outgoing metadata for downstream calls.
// This is useful in middleware/interceptors.
func PropagateRequestID(ctx context.Context) context.Context {
	rid := GetRequestID(ctx)
	if rid == "" {
		rid = ExtractRequestID(ctx)
		if rid == "" {
			rid = NewRequestID()
		}
		ctx = WithRequestID(ctx, rid)
	}
	return InjectRequestID(ctx, rid)
}