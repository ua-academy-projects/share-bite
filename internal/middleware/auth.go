package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/pkg/jwt"
)

type AccessTokenParser interface {
	ParseAccessToken(token string) (string, string, jwt.UserStatus, error)
}

const (
	authorizationHeader = "Authorization"
	CtxUserID           = "userId"
	CtxUserRole         = "userRole"
	CtxUserStatus       = "userStatus"
)

func Auth(parser AccessTokenParser) gin.HandlerFunc {
	if parser == nil {
		panic("auth middleware is not configured: parser cannot be nil")
	}

	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}
		userID, role, status, err := parser.ParseAccessToken(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(CtxUserID, userID)
		c.Set(CtxUserRole, role)
		c.Set(CtxUserStatus, status)

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
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			// no header -> skip validation below
			c.Next()
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		userID, role, status, err := parser.ParseAccessToken(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(CtxUserID, userID)
		c.Set(CtxUserRole, role)
		c.Set(CtxUserStatus, status)

		c.Next()
	}
}
