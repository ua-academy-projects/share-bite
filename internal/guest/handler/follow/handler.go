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
	Follow(ctx context.Context, customerID, targetCustomerID string) error
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

	//protected := r.Group("/").Use(authMiddleware, customerMiddleware)

	r.POST("/id/:id/follow", h.follow)
	r.DELETE("/id/:id/follow", h.unfollow)

	optional := r.Group("/").Use(optionalAuthMiddleware, customerMiddleware)

	optional.GET("/id/:id/followers", h.listFollowers)
	optional.GET("/id/:id/following", h.listFollowing)

	r.GET("/following", h.listMyFollowing)
}

func (h *handler) followersToResponse(customers []entity.Follower) []dto.FollowerResponse {
	res := make([]dto.FollowerResponse, 0, len(customers))
	for _, c := range customers {
		res = append(res, h.toResponse(c))
	}
	return res
}

func (h *handler) toResponse(f entity.Follower) dto.FollowerResponse {
	var avatarURL *string
	if f.AvatarObjectKey != nil && h.storage != nil {
		url := h.storage.BuildURL(*f.AvatarObjectKey)
		avatarURL = &url
	}

	return dto.FollowerResponse{
		ID:           f.ID,
		UserName:     f.UserName,
		AvatarURL:    avatarURL,
		IsFollowing:  f.IsFollowing,
		IsFollowedBy: f.IsFollowedBy,
		IsMutual:     f.IsMutual,
	}
}

func (h *handler) listCustomersResponse(
	followers []entity.Follower,
	nextToken string,
) dto.ListCustomersResponse {
	return dto.ListCustomersResponse{
		Customers:     h.followersToResponse(followers),
		NextPageToken: nextToken,
	}
}
