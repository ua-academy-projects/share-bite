package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccessTokenParser interface {
	ParseAccessToken(token string) (string, string, error)
}

const (
	authorizationHeader = "Authorization"
	CtxUserID           = "userId"
	CtxUserRole         = "userRole"
)

func Auth(parser AccessTokenParser) gin.HandlerFunc {
	if parser == nil {
		panic("auth middleware is not configured: parser cannot be nil")
	}

	return func(c *gin.Context) {
		token := c.Query("access_token")
		if authHeader := c.GetHeader(authorizationHeader); authHeader != "" {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing token"})
			return
		}

		userID, role, err := parser.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(CtxUserID, userID)
		c.Set(CtxUserRole, role)
		c.Next()
	}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(CtxUserRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error: invalid role type"})
			return
		}

		if !slices.Contains(allowedRoles, userRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: insufficient permissions"})
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	val, exists := c.Get(CtxUserID)
	if !exists {
		return "", false
	}

	userIDStr, ok := val.(string)
	if !ok {
		return "", false
	}

	return userIDStr, true
}

func OptionalAuth(parser AccessTokenParser) gin.HandlerFunc {
	if parser == nil {
		panic("optional auth middleware is not configured: parser cannot be nil")
	}

	return func(c *gin.Context) {
		token := c.Query("access_token")
		if authHeader := c.GetHeader(authorizationHeader); authHeader != "" {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			c.Next()
			return
		}

		userID, role, err := parser.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(CtxUserID, userID)
		c.Set(CtxUserRole, role)
		c.Next()
	}
}
