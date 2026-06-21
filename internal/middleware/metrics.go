package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type metrics interface {
	HistogramResponseTimeObserve(path string, method string, status string, time float64)
	IncRequestCounter(path string, method string, status string)
}

func Metrics(m metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		m.IncRequestCounter(c.FullPath(), c.Request.Method, status)
		m.HistogramResponseTimeObserve(c.FullPath(), c.Request.Method, status, duration)
	}
}
