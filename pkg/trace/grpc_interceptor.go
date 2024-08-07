package trace

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/seaphony/ceph-api/pkg/log"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		traceID := oteltrace.SpanFromContext(ctx).
			SpanContext().
			TraceID()
		ctx = log.WithTraceID(ctx, traceID.String())

		return handler(ctx, req)
	}
}

func StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		traceID := oteltrace.SpanFromContext(stream.Context()).
			SpanContext().
			TraceID()
		ctx := log.WithTraceID(stream.Context(), traceID.String())

		return handler(srv, &grpc_middleware.WrappedServerStream{
			ServerStream:   stream,
			WrappedContext: ctx,
		})
	}
}
