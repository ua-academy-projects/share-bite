package business

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/dto"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

type handler struct {
	service businessService
}

type businessService interface {
	UpdatePost(ctx context.Context, postID int64, userID string, content string) (*entity.PostWithPhotos, error)
	DeletePost(ctx context.Context, postID int64, userID string) error

	Get(ctx context.Context, id int) (*entity.OrgUnit, error)
	List(ctx context.Context, brandId, skip, limit int) (pagination.Result[entity.OrgUnit], error)
	GetPosts(ctx context.Context, skip, limit int) (pagination.Result[entity.PostWithPhotos], error)

	CreateLocation(ctx context.Context, brandID int, ownerUserID string, in dto.CreateLocationInput) (*entity.OrgUnit, error)
	UpdateLocation(ctx context.Context, locationID int, ownerUserID string, in dto.UpdateLocationInput) (*entity.OrgUnit, error)
	DeleteLocation(ctx context.Context, locationID int, ownerUserID string) error
}

func RegisterHandlers(r *gin.RouterGroup, service businessService, parser middleware.AccessTokenParser) {
	h := &handler{
		service: service,
	}

	auth := middleware.Auth(parser)

	r.GET("/org-units/:id", h.get)
	r.GET("/org-units/:id/locations", h.list)
	r.GET("/posts", h.GetPosts)

	businessOnly := r.Group("/").
		Use(auth).
		Use(middleware.RequireRoles("business"))

	businessOnly.PUT("/posts/:id", h.UpdatePost)
	businessOnly.DELETE("/posts/:id", h.DeletePost)

	businessOnly.POST("/:id/locations", h.createLocation)
	businessOnly.PATCH("/locations/:id", h.updateLocation)
	businessOnly.DELETE("/locations/:id", h.deleteLocation)
}

// errorResponse is used for swagger documentation.
type errorResponse struct {
	Error string `json:"error" example:"not found"`
}
