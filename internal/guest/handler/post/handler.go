package post

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type handler struct {
	service postService
}

type postService interface {
	List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error)
	Get(ctx context.Context, postID string) (entity.Post, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service postService,
) {
	h := &handler{
		service: service,
	}

	r.GET("/", h.list)
	r.GET("/:id", h.get)
}

type postResponse struct {
	ID string `json:"id"`

	Description string `json:"description"`
}

func postToResponse(post entity.Post) postResponse {
	return postResponse{
		ID: post.ID,

		Description: post.Description,
	}
}
