package follow

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage"
)

type handler struct {
	service         customerFollowService
	customerService customerService
	storage         storage.ObjectStorage
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
	st storage.ObjectStorage,
) {
	h := &handler{
		service: service,
		storage: st,
	}

	protected := r.Group("/").Use(authMiddleware, customerMiddleware)

	protected.POST("/:id/follow", h.follow)
	protected.DELETE("/:id/follow", h.unfollow)

	optional := r.Group("/").Use(optionalAuthMiddleware, customerMiddleware)

	optional.GET("/:id/followers", h.listFollowers)
	optional.GET("/:id/following", h.listFollowing)

	protected.GET("/following", h.listMyFollowing)
}

func (h *handler) customersToResponse(customers []entity.Customer) []dto.CustomerResponse {
	res := make([]dto.CustomerResponse, 0, len(customers))
	for _, c := range customers {
		res = append(res, h.toResponse(c))
	}
	return res
}

func (h *handler) toResponse(customer entity.Customer) dto.CustomerResponse {
	var avatarURL *string
	if customer.AvatarObjectKey != nil && h.storage != nil {
		url := h.storage.BuildURL(*customer.AvatarObjectKey)
		avatarURL = &url
	}

	return dto.CustomerResponse{
		ID:        customer.ID,
		UserName:  customer.UserName,
		AvatarURL: avatarURL,
	}
}

func (h *handler) listCustomersResponse(
	customers []entity.Customer,
	nextToken string,
) dto.ListCustomersResponse {
	return dto.ListCustomersResponse{
		Customers:     h.customersToResponse(customers),
		NextPageToken: nextToken,
	}
}
