package httpctx

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	userIDContextKey = "userId"
)

func GetUserID(c *gin.Context) (string, error) {
	val, ok := c.Get(userIDContextKey)
	if !ok {
		return "", fmt.Errorf("userId missing in context")
	}

	userID, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("userId cast failed")
	}

	if userID == "" {
		return "", fmt.Errorf("userId is empty")
	}

	return userID, nil
}
