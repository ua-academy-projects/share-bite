package business

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

func getUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get(middleware.CtxUserID)
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

func checkBusinessRole(c *gin.Context) bool {
	val, exists := c.Get(middleware.CtxUserRole)
	if !exists {
		return false
	}

	role, ok := val.(string)
	if !ok {
		return false
	}

	return role == "business"
}
