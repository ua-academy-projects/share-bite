package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type metrics interface {
	IncRequestCounter(path, method, status string)
	HistogramResponseTimeObserve(path, method, status string, time float64)
	IncActiveRequests()
	DecActiveRequests()
}

func Metrics(m metrics, ignoredPaths []string) gin.HandlerFunc {
	ignored := make(map[string]struct{}, len(ignoredPaths))
	for _, path := range ignoredPaths {
		ignored[path] = struct{}{}
	}

	return func(c *gin.Context) {
		path := c.FullPath()
		if _, skip := ignored[path]; skip {
			c.Next()
			return
		}

		m.IncActiveRequests()
		defer m.DecActiveRequests()

		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		m.IncRequestCounter(path, c.Request.Method, status)
		m.HistogramResponseTimeObserve(path, c.Request.Method, status, duration)
	}
}
