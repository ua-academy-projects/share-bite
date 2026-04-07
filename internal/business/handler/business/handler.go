package business

import (
	"context"
	"mime/multipart"

	"github.com/gin-gonic/gin"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
)

type handler struct {
	service businessService
}

type businessService interface {
	CreatePost(ctx context.Context, userID string, unitID int, description string, images []*multipart.FileHeader) (*entity.PostWithPhotos, error)
	CheckOwnership(ctx context.Context, userID string, unitID int) error
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

	r.GET("/:id", h.get)
	r.GET("/:id/locations", h.list)

	businessOnly := r.Group("/").
		Use(auth).
		Use(middleware.RequireRoles("business"))

	businessOnly.PUT("/posts/:id", h.UpdatePost)
	businessOnly.DELETE("/posts/:id", h.DeletePost)
	businessOnly.POST("/posts/:id", h.CreatePost)
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