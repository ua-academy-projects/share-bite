package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		ctx := c.Request.Context()

		kvs := []any{
			"status", statusCode,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
			"duration", duration,
		}

		if len(c.Errors) > 0 {
			kvs = append(kvs, "errors", c.Errors.String())
		}

		// depends on the result
		switch {
		case statusCode >= 500:
			logger.ErrorKV(ctx, "http request failed", kvs...)
		case statusCode >= 400:
			logger.WarnKV(ctx, "http request failed (client error)", kvs...)
		default:
			logger.InfoKV(ctx, "http request completed", kvs...)
		}
	}
}
