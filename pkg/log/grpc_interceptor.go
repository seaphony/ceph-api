package log

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func UnaryInterceptor(cfg Config) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		l := CreateLogger(cfg)
		builder := l.With().Str(grpcMethod, info.FullMethod)
		newLogger := builder.Logger()
		ctx = newLogger.WithContext(ctx)

		return handler(ctx, req)
	}
}

func StreamInterceptor(cfg Config) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		l := CreateLogger(cfg)
		builder := l.With().Str(grpcMethod, info.FullMethod)
		newLogger := builder.Logger()
		ctx := newLogger.WithContext(stream.Context())

		return handler(srv, &grpc_middleware.WrappedServerStream{
			ServerStream:   stream,
			WrappedContext: ctx,
		})
	}
}
