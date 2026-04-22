package httpctx

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

func Get[T any](c *gin.Context, key string) (T, error) {
	var zero T

	val, ok := c.Get(key)
	if !ok {
		return zero, fmt.Errorf("%s missing in context", key)
	}

	typedVal, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("%s invalid type in context", key)
	}

	return typedVal, nil
}

func GetUserID(c *gin.Context) (string, error) {
	return Get[string](c, middleware.CtxUserID)
}

func GetUserRole(c *gin.Context) (string, error) {
	return Get[string](c, middleware.CtxUserRole)
}

func GetCustomerID(c *gin.Context) (string, error) {
	return Get[string](c, middleware.CtxCustomerID)
}

func GetOptionalCustomerID(c *gin.Context) (*string, error) {
	val, ok := c.Get(middleware.CtxCustomerID)
	if !ok {
		return nil, nil
	}

	typedVal, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("%s invalid type in context", middleware.CtxCustomerID)
	}

	return &typedVal, nil
}
