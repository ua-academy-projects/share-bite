package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader       = "Authorization"
	authorizationHeaderPrefix = "Bearer "
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header not provided"})
			return
		}

		if !strings.HasPrefix(authHeader, authorizationHeaderPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		// TODO: validation
		// token := authHeader[len(authorizationHeaderPrefix):]

		// TODO: replace with real user id
		userID := "b5461b65-7244-4cef-a8ec-c44fa10b4997"

		c.Set("userId", userID)
		c.Next()
	}
}
