package middleware

import (
	"net/http"
	"slices"
	"strconv"
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
		userID, role, err := parser.ParseAccessToken(headerParts[1])
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

func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get(CtxUserID)
	if !exists {
		return 0, false
	}

	userIDStr, ok := val.(string)
	if !ok {
		return 0, false
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, false
	}

	return userID, true
}
