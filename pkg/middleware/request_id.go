package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"go.uber.org/zap"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Header("X-Request-ID", reqID)

		ctx := logger.WithFields(c.Request.Context(), zap.String("request_id", reqID))
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
