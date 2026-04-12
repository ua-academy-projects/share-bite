package follow

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	customerResponse "github.com/ua-academy-projects/share-bite/internal/guest/handler/customer"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
)

type handler struct {
	service         customerFollowService
	customerService customerService
}

type customerFollowService interface {
	Follow(ctx context.Context, userID, targetCustomerID string) (entity.CustomerFollow, error)
	Unfollow(ctx context.Context, userID, targetCustomerID string) error

	ListFollowing(ctx context.Context, targetCustomerID, requesterUserID string) ([]entity.Customer, error)
	ListFollowers(ctx context.Context, targetCustomerID, requesterUserID string) ([]entity.Customer, error)
}

type customerService interface {
	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
}

func RegisterHandler(
	r *gin.RouterGroup,
	service customerFollowService,
	authMiddleware gin.HandlerFunc,
) {
	h := &handler{
		service: service,
	}

	protected := r.Group("/").Use(authMiddleware)

	protected.POST("/:id/follow", h.follow)
	protected.DELETE("/:id/follow", h.unfollow)

	protected.GET("/following", h.listMyFollowing)
	r.GET("/:id/following", h.listFollowing)
	r.GET("/:id/followers", h.listFollowers)
}

type followResponse struct {
	Follow entity.CustomerFollow `json:"follow"`
}

func getFollowerUserID(c *gin.Context) (string, error) {
	userID, err := httpctx.GetUserID(c)
	if err != nil {
		return "", apperror.ErrInvalidParam
	}
	return userID, nil
}

type listCustomersResponse struct {
	Customers []customerResponse.CustomerResponse `json:"customers"`
}

func customersToResponse(customers []entity.Customer) []customerResponse.CustomerResponse {
	res := make([]customerResponse.CustomerResponse, 0, len(customers))
	for _, c := range customers {
		res = append(res, customerResponse.CustomerToResponse(c))
	}
	return res
}
