package middleware

import (
	"net/http"
	"strings"

	"github.com/ua-academy-projects/share-bite/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	CtxUserID           = "userId"
	CtxUserRole         = "userRole"
)

func Auth(tokenManager *jwt.TokenManager) gin.HandlerFunc {
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
		userID, role, err := tokenManager.ParseAccessToken(headerParts[1])
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
		roleVal, exist := c.Get(CtxUserRole)
		if !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error: invalid role type"})
			return
		}
		hasAccess := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {

				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied: insufficient permissions"})
			return
		}

		c.Next()
	}
}
