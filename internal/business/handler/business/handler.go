package business

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type handler struct {
	service businessService
}

type businessService interface {
	UpdatePost(ctx context.Context, postID int64, userID string, content string) (*entity.PostWithPhotos, error)
	DeletePost(ctx context.Context, postID int64, userID string) error

	Get(ctx context.Context, id int) (*entity.OrgUnit, error)
	List(ctx context.Context, brandId, page, limit int) ([]entity.OrgUnit, error)

	GetPosts(ctx context.Context, page, limit int) ([]entity.PostWithPhotos, error)
}

func RegisterHandlers(r *gin.RouterGroup, service businessService, parser middleware.AccessTokenParser) {
	h := &handler{
		service: service,
	}

	auth := middleware.Auth(parser)

	r.GET("/:id", h.get)
	r.GET("/:id/locations", h.list)
	r.GET("/posts", h.GetPosts)

	businessOnly := r.Group("/").
		Use(auth).
		Use(middleware.RequireRoles("business"))

	businessOnly.PUT("/posts/:id", h.UpdatePost)
	businessOnly.DELETE("/posts/:id", h.DeletePost)
}

// errorResponse is used for swagger documentation.
type errorResponse struct {
	Error string `json:"error" example:"not found"`
}
