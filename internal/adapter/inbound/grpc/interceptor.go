package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

type LoggingInterceptor struct {
	logger Logger
}

func NewLoggingInterceptor(logger Logger) *LoggingInterceptor {
	return &LoggingInterceptor{logger: logger}
}

func (i *LoggingInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)

		if err != nil {
			i.logger.Error("gRPC request failed",
				"method", info.FullMethod,
				"error", err.Error(),
			)
		} else {
			i.logger.Info("gRPC request completed",
				"method", info.FullMethod,
			)
		}

		return resp, err
	}
}

type RecoveryInterceptor struct {
	logger Logger
}

func NewRecoveryInterceptor(logger Logger) *RecoveryInterceptor {
	return &RecoveryInterceptor{logger: logger}
}

func (i *RecoveryInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				i.logger.Error("panic recovered",
					"method", info.FullMethod,
					"panic", r,
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
