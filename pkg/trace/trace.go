package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	jaeger "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Enabled  bool   `yaml:"enabled"`
	Insecure bool   `yaml:"insecure"`
	Endpoint string `yaml:"endpoint"`
}

func NewTracerProvider(ctx context.Context, conf Config, version string) (func(ctx context.Context) error, trace.TracerProvider, error) {
	var tp *sdktrace.TracerProvider
	if !conf.Enabled {
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithResource(sdkresource.NewSchemaless(
				semconv.ServiceNameKey.String("ceph-api"),
				semconv.ServiceVersionKey.String(version),
			)),
		)
	} else {
		opts := []jaeger.Option{jaeger.WithEndpoint(conf.Endpoint)}
		if conf.Insecure {
			opts = append(opts, jaeger.WithInsecure())
		}
		exp, err := jaeger.New(ctx, opts...)
		if err != nil {
			return nil, nil, err
		}

		tp = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(sdkresource.NewSchemaless(
				semconv.ServiceNameKey.String("ceph-api"),
				semconv.ServiceVersionKey.String(version),
			)),
		)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, tp, nil
}
