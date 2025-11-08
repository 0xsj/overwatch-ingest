// platform/pkg/grpc/interceptors/client_logging.go
package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// UnaryClientLoggingInterceptor logs gRPC unary client requests.
func UnaryClientLoggingInterceptor(log logger.Logger) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()

		log.Info("grpc client request started",
			"method", method,
			"target", cc.Target(),
		)

		// Invoke the RPC
		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		if err != nil {
			log.Error("grpc client request failed",
				"method", method,
				"target", cc.Target(),
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)
		} else {
			log.Info("grpc client request completed",
				"method", method,
				"target", cc.Target(),
				"duration_ms", duration.Milliseconds(),
			)
		}

		return err
	}
}

// StreamClientLoggingInterceptor logs gRPC stream client requests.
func StreamClientLoggingInterceptor(log logger.Logger) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		start := time.Now()

		log.Info("grpc client stream started",
			"method", method,
			"target", cc.Target(),
		)

		stream, err := streamer(ctx, desc, cc, method, opts...)

		duration := time.Since(start)

		if err != nil {
			log.Error("grpc client stream failed",
				"method", method,
				"target", cc.Target(),
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)
		} else {
			log.Info("grpc client stream established",
				"method", method,
				"target", cc.Target(),
				"duration_ms", duration.Milliseconds(),
			)
		}

		return stream, err
	}
}