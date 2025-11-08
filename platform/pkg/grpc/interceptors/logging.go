// platform/pkg/grpc/interceptors/logging.go
package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// UnaryLoggingInterceptor logs gRPC unary requests.
func UnaryLoggingInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		log.Info("grpc request started",
			"method", info.FullMethod,
		)

		// Call the handler
		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			log.Error("grpc request failed",
				"method", info.FullMethod,
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)
		} else {
			log.Info("grpc request completed",
				"method", info.FullMethod,
				"duration_ms", duration.Milliseconds(),
			)
		}

		return resp, err
	}
}

// StreamLoggingInterceptor logs gRPC stream requests.
func StreamLoggingInterceptor(log logger.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		log.Info("grpc stream started",
			"method", info.FullMethod,
		)

		err := handler(srv, ss)

		duration := time.Since(start)

		if err != nil {
			log.Error("grpc stream failed",
				"method", info.FullMethod,
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)
		} else {
			log.Info("grpc stream completed",
				"method", info.FullMethod,
				"duration_ms", duration.Milliseconds(),
			)
		}

		return err
	}
}