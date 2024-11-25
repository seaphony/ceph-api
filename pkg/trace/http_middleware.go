package trace

import (
	"net/http"

	"github.com/clyso/ceph-api/pkg/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func HttpMiddleware(tp trace.TracerProvider, next http.Handler) http.Handler {
	return otelhttp.NewHandler(addTraceID(next), "", otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
		return r.Method + ":" + r.URL.RawPath
	}), otelhttp.WithTracerProvider(tp))
}

func addTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := trace.SpanFromContext(r.Context()).
			SpanContext().
			TraceID()
		ctx := log.WithTraceID(r.Context(), traceID.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
