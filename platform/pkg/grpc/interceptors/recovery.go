// platform/pkg/grpc/interceptors/recovery.go
package interceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/0xsj/scout/platform/pkg/observability/logger"
)

// UnaryRecoveryInterceptor recovers from panics in unary handlers.
func UnaryRecoveryInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("grpc handler panic recovered",
					"method", info.FullMethod,
					"panic", fmt.Sprintf("%v", r),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// StreamRecoveryInterceptor recovers from panics in stream handlers.
func StreamRecoveryInterceptor(log logger.Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("grpc stream panic recovered",
					"method", info.FullMethod,
					"panic", fmt.Sprintf("%v", r),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(srv, ss)
	}
}