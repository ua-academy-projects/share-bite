package business

import (
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
)

type handler struct {
	service businessService
}

type businessService interface {
	CreatePost(ctx context.Context, userID string, unitID int, description string, images []*multipart.FileHeader) (*entity.PostWithPhotos, error)
	CheckOwnership(ctx context.Context, userID string, unitID int) error
	UpdatePost(ctx context.Context, postID int64, userID string, content string) (*entity.PostWithPhotos, error)
	DeletePost(ctx context.Context, postID int64, userID string) error
	GetPosts(ctx context.Context, skip, limit int) (pagination.Result[entity.PostWithPhotos], error)

	Get(ctx context.Context, id int) (*entity.OrgUnit, error)
	List(ctx context.Context, brandId, skip, limit int) (pagination.Result[entity.OrgUnit], error)
	GetVenuesByIDs(ctx context.Context, ids []int) ([]entity.OrgUnit, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service businessService,
	parser middleware.AccessTokenParser,
) {
	h := &handler{
		service: service,
	}

	auth := middleware.Auth(parser)

	org_units := r.Group("/org-units")
	{
		org_units.GET("/:id", h.get)
		org_units.GET("/:id/locations", h.list)
		org_units.POST("/venues", h.getVenuesByIDs)
	}

	r.GET("/posts", h.GetPosts)

	businessOnly := r.Group("/posts").
		Use(auth).
		Use(middleware.RequireRoles("business"))
		
	{
		businessOnly.PUT("/:id", h.UpdatePost)
		businessOnly.DELETE("/:id", h.DeletePost)
		businessOnly.POST("/:id", h.CreatePost)
	}
}

// errorResponse is used for swagger documentation.
type errorResponse struct {
	Error string `json:"error" example:"not found"`
}

func getUserID(c *gin.Context) (string, bool) {
	val, exists := c.Get(middleware.CtxUserID)

	if !exists {
		return "", false
	}

	userIDStr, ok := val.(string)
	if !ok {
		return "", false
	}

	return userIDStr, true
}