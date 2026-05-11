package httpctx

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

var (
	ErrMissingContext = errors.New("missing in context")
	ErrInvalidType    = errors.New("invalid type in context")
	ErrInvalidFormat  = errors.New("invalid format")
)

func Get[T any](c *gin.Context, key string) (T, error) {
	var zero T

	val, ok := c.Get(key)
	if !ok {
		return zero, fmt.Errorf("%w: %s", ErrMissingContext, key)
	}

	typedVal, ok := val.(T)
	if !ok {	
		return zero, fmt.Errorf("%w: %s", ErrInvalidType, key)
	}

	return typedVal, nil
}

func GetUserID(c *gin.Context) (string, error) {
	return Get[string](c, middleware.CtxUserID)
}

func GetUserRole(c *gin.Context) (string, error) {
	return Get[string](c, middleware.CtxUserRole)
}

func GetUserUUID(c *gin.Context) (uuid.UUID, error) {
	userIdStr, err := GetUserID(c)
	if err != nil {
		return uuid.Nil, err
	}
	parsedId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}
	return parsedId, nil
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
