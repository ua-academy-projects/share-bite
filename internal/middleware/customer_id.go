package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/guest/util/response"
)

const (
	CtxCustomerID = "customerId"
)

type CustomerProvider interface {
	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
}

func CustomerID(provider CustomerProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(CtxUserID)
		if !exists {
			c.Next()
			return
		}

		userID, ok := val.(string)
		if !ok || userID == "" {
			c.Next()
			return
		}

		ctx := c.Request.Context()

		customer, err := provider.GetByUserID(ctx, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse{Message: "customer profile not found"})
			return
		}

		c.Set(CtxCustomerID, customer.ID)
		c.Next()
	}
}
