package middlewares

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-finance-backend/pkg/logger"
	"go.uber.org/zap"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Create a custom response writer to capture response body
		w := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// Log request details
		logger.Info("Request received",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("clientIP", c.ClientIP()),
		)

		// Process request
		c.Next()

		// Calculate resolution time
		duration := time.Since(start)

		// Ensure we capture the status
		status := c.Writer.Status()

		// If debug mode is on, log both request and response details
		if viper.GetBool("server.debug") {
			logger.Info("Request processed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.String("clientIP", c.ClientIP()),
				zap.String("response", w.body.String()),
			)
		} else {
			// In non-debug mode, just log basic response info
			logger.Info("Request processed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", status),
				zap.Duration("duration", duration),
			)
		}
	}
}
