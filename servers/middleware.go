package servers

import (
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"time"

	"shipments/domains/tracing"

	"github.com/gin-gonic/gin"
)

// DefaultStructuredLogger logs a gin HTTP request in a json format using zerolog
func DefaultStructuredLogger() gin.HandlerFunc {
	return StructuredLogger(slog.Default())
}

// StructuredLogger logs a HTTP request in a specific format
func StructuredLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		tracer := tracing.Tracer()
		if tracer == nil {
			slog.Error("tracing not initialized")
			c.Next()
			return
		}

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		//trace.SpanFromContext(c.Request.Context())
		fmt.Printf("Span: %v\n", trace.SpanFromContext(c.Request.Context()))

		logger.Info(
			"incoming request",
			"client_ip", param.ClientIP,
			"method", param.Method,
			"latency", param.Latency.String(),
			"body_size", param.BodySize,
			"path", param.Path,
			"status", param.StatusCode,
			"trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String(), //tracing.TraceID(c.Request.Context()),
			"user_agent", c.Request.UserAgent(),
		)
	}
}

//uptrace.TraceURL(trace.SpanFromContext(c.Request.Context())))
