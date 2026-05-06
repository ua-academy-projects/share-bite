package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
)

func RequireWritableAccountStatus(statusContextKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		statusVal, exists := c.Get(statusContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var status jwt.UserStatus
		switch v := statusVal.(type) {
		case jwt.UserStatus:
			status = v
		case string:
			status = jwt.UserStatus(v)
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error: invalid status type"})
			return
		}

		switch status {
		case jwt.UserStatusActive:
			c.Next()
		case jwt.UserStatusMuted:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: muted account is read-only"})
			return
		case jwt.UserStatusSuspended:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: suspended account"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: unknown account status"})
			return
		}
	}
}
