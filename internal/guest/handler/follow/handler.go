package follow

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	customerResponse "github.com/ua-academy-projects/share-bite/internal/guest/handler/customer"
)

type handler struct {
	service         customerFollowService
	customerService customerService
}

type customerFollowService interface {
	Follow(ctx context.Context, customerID, targetCustomerID string) (entity.CustomerFollow, error)
	Unfollow(ctx context.Context, customerID, targetCustomerID string) error

	ListFollowers(ctx context.Context, in entity.ListFollowersInput) (entity.ListFollowersOutput, error)
	ListFollowing(ctx context.Context, in entity.ListFollowingInput) (entity.ListFollowingOutput, error)
}

type customerService interface {
	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
}

func RegisterHandler(
	r *gin.RouterGroup,
	service customerFollowService,
	authMiddleware gin.HandlerFunc,
	optionalAuthMiddleware gin.HandlerFunc,
	customerMiddleware gin.HandlerFunc,
) {
	h := &handler{
		service: service,
	}

	protected := r.Group("/").Use(authMiddleware, customerMiddleware)

	protected.POST("/:id/follow", h.follow)
	protected.DELETE("/:id/follow", h.unfollow)

	optional := r.Group("/").Use(optionalAuthMiddleware, customerMiddleware)

	optional.GET("/:id/followers", h.listFollowers)
	optional.GET("/:id/following", h.listFollowing)

	protected.GET("/following", h.listMyFollowing)
}

// errorResponse describes guest API error payload.
type errorResponse struct {
	Message string `json:"message" example:"invalid request"`
}

func customersToResponse(customers []entity.Customer) []customerResponse.CustomerResponse {
	res := make([]customerResponse.CustomerResponse, 0, len(customers))
	for _, c := range customers {
		res = append(res, customerResponse.CustomerToResponse(c))
	}
	return res
}
