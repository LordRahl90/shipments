package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
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

// TraceProvider sets the trace provider
func TraceProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	slog.InfoContext(ctx, "setting up trace provider")
	prop := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(prop)
	slog.InfoContext(ctx, "propagator configured")

	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())
	if err != nil {
		return nil, err
	}

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	r, err := getResource(hostName)
	if err != nil {
		return nil, err
	}

	provider = sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		sdktrace.WithBatcher(traceExporter),
	)

	otel.SetTracerProvider(provider)

	return provider, nil
}

func getResource(hostName string) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"https://opentelemetry.io/schemas/1.26.0",
			attribute.String("service.name", "shipments"),
			attribute.String("environment", os.Getenv("ENVIRONMENT")),
			attribute.String("app.version", "1.0.0")),
	)
}
