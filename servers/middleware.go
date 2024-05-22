package servers

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// DefaultStructuredLogger logs a gin HTTP request in a json format using zerolog
func DefaultStructuredLogger() gin.HandlerFunc {
	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return StructuredLogger(jsonLogger)
}

// StructuredLogger logs a HTTP request in a specific format
func StructuredLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

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

		logger.Info(
			"incoming request",
			"client_ip", param.ClientIP,
			"method", param.Method,
			"latency", param.Latency.String(),
			"body_size", param.BodySize,
			"path", param.Path,
			"status", param.StatusCode,
			"user_agent", c.Request.UserAgent(),
		)
	}
}
