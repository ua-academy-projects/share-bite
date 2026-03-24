package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/pkg/middleware"
)

type handler struct {
	service customerService
}

type customerService interface {
	Create(ctx context.Context, in entity.CreateCustomer) (string, error)
	Update(ctx context.Context, in entity.UpdateCustomer) (entity.Customer, error)

	GetByUserName(ctx context.Context, userName string) (entity.Customer, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service customerService,
) {
	h := &handler{
		service: service,
	}

	// public
	r.GET("/:username", h.getByUserName)

	// protected
	protected := r.Group("/").Use(middleware.Auth())

	protected.POST("/", h.create)
	protected.PATCH("/", h.update)
	protected.GET("/", h.getMe)
}

type customerResponse struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`

	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	AvatarURL *string `json:"avatarURL"`
	Bio       *string `json:"bio"`

	CreatedAt time.Time `json:"createdAt"`
}

func customerToResponse(customer entity.Customer) customerResponse {
	var avatarURL *string
	if customer.AvatarObjectKey != nil {
		// TODO: replace with real s3 presigned url
		url := fmt.Sprintf("https://test.com/%s", *customer.AvatarObjectKey)
		avatarURL = &url
	}

	return customerResponse{
		ID:     customer.ID,
		UserID: customer.UserID,

		UserName:  customer.UserName,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,

		Bio:       customer.Bio,
		AvatarURL: avatarURL,

		CreatedAt: customer.CreatedAt,
	}

}
