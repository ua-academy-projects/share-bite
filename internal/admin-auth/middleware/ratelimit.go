package middleware

import (
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func NewAuthRecoveryLimiter(limit int, rate time.Duration) gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  rate,
		Limit: uint(limit),
	})

	return ratelimit.RateLimiter(store, &ratelimit.Options{
		KeyFunc: func(c *gin.Context) string {
			return c.ClientIP()
		},
		ErrorHandler: rateLimitErrorHandler,
	})
}

func rateLimitErrorHandler(c *gin.Context, info ratelimit.Info) {
	retryAfterSeconds := int(time.Until(info.ResetTime).Seconds())
	if retryAfterSeconds < 0 {
		retryAfterSeconds = 0
	}

	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"message":             "too many requests",
		"retry_after_seconds": retryAfterSeconds,
	})
}
