package tracing

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
)

// Tracer sets the tracer
func Tracer() trace.Tracer {
	if tracer == nil {
		tracer = otel.GetTracerProvider().Tracer("shipments")
	}
	return tracer
}

// ConsoleExporter exports to console
func ConsoleExporter() (oteltrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
}

// TempoExporter exports to tempo
func TempoExporter(ctx context.Context, otelEndpoint string) (sdktrace.SpanExporter, error) {
	opt := otlptracehttp.WithInsecure()
	endpointOpt := otlptracehttp.WithEndpoint(otelEndpoint)
	return otlptracehttp.New(ctx, opt, endpointOpt)
}

// TraceProvider sets the trace provider
func TraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	if provider == nil {
		hostName, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		r, err := resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("shipments"),
				semconv.ServiceVersionKey.String("v0.1.0"),
				semconv.ServiceInstanceIDKey.String(hostName),
			))

		if err != nil {
			log.Fatal(err)
		}
		provider = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(r),
		)
	}

	return provider
}

func TraceID(ctx context.Context) trace.TraceID {
	return trace.SpanContextFromContext(ctx).TraceID()
}
