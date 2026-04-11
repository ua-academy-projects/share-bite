package post

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type handler struct {
	service         postService
	customerService customerService
}

type postService interface {
	Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error)
	Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error)
	List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error)
	Get(ctx context.Context, postID string) (entity.Post, error)
}

type customerService interface {
	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service postService,
	customerService customerService,
	authMiddleware gin.HandlerFunc,
) {
	h := &handler{
		service:         service,
		customerService: customerService,
	}

	r.GET("/", h.list)
	r.GET("/:id", h.get)

	protected := r.Group("/").Use(authMiddleware)
	protected.POST("/", h.create)
	protected.PATCH("/:id", h.update)
}

type postResponse struct {
	ID          string            `json:"id"`
	CustomerID  string            `json:"customerId"`
	VenueID     int64             `json:"venueId"`
	Text        string            `json:"text"`
	Rating      int16             `json:"rating"`
	Status      entity.PostStatus `json:"status"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	PublishedAt *time.Time        `json:"publishedAt,omitempty"`
}

func postToResponse(post entity.Post) postResponse {
	return postResponse{
		ID:          post.ID,
		CustomerID:  post.CustomerID,
		VenueID:     post.VenueID,
		Text:        post.Text,
		Rating:      post.Rating,
		Status:      post.Status,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
	}
}
