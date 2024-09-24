package tracing

import (
	"context"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	otelLog "go.opentelemetry.io/otel/sdk/log"
	"log/slog"
	"os"
)

func LoggerProvider(ctx context.Context) (*otelLog.LoggerProvider, error) {
	slog.InfoContext(ctx, "setting up logger provider")
	logExporter, err := otlploghttp.New(ctx, otlploghttp.WithInsecure())

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

	logProvider := otelLog.NewLoggerProvider(
		otelLog.WithResource(r),
		otelLog.WithProcessor(otelLog.NewBatchProcessor(logExporter)),
	)
	global.SetLoggerProvider(logProvider)

	slog.InfoContext(ctx, "logger provider configured")

	return logProvider, nil
}
