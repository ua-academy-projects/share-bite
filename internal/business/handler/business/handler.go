package business

import (
	"context"

	"github.com/gin-gonic/gin"
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
	ToggleLike(ctx context.Context, postID int64, customerID string) (bool, error)
	GetLikes(ctx context.Context, postID int64, limit, offset int) ([]entity.Like, error)
	CreateComment(ctx context.Context, postID int64, authorID, content string) (*entity.Comment, error)
	UpdateComment(ctx context.Context, commentID int64, authorID, content string) (*entity.Comment, error)
	DeleteComment(ctx context.Context, commentID int64, authorID string) error
	GetComments(ctx context.Context, postID int64, limit, offset int) ([]entity.CommentWithAuthor, error)
}

func RegisterHandlers(r *gin.RouterGroup, service businessService, parser middleware.AccessTokenParser) {
	h := &handler{
		service: service,
	}

	auth := middleware.Auth(parser)

	r.GET("/org-units/:id", h.get)
	r.GET("/org-units/:id/locations", h.list)
	r.GET("/posts", h.GetPosts)
	r.GET("/posts/:id/likes", h.GetLikes)
	r.GET("/posts/:id/comments", h.GetComments)

	businessOnly := r.Group("/").
		Use(auth).
		Use(middleware.RequireRoles("business"))

	businessOnly.PUT("/posts/:id", h.UpdatePost)
	businessOnly.DELETE("/posts/:id", h.DeletePost)

	authenticated := r.Group("/").Use(auth)

	authenticated.POST("/posts/:id/likes", h.ToggleLike)
	authenticated.POST("/posts/:id/comments", h.CreateComment)
	authenticated.PATCH("/posts/:id/comments/:comment_id", h.UpdateComment)
	authenticated.DELETE("/posts/:id/comments/:comment_id", h.DeleteComment)
}

// errorResponse is used for swagger documentation.
type errorResponse struct {
	Error string `json:"error" example:"not found"`
}
