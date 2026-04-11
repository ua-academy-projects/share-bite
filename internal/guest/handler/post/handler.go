package post

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/storage"
)

type handler struct {
	service         postService
	customerService customerService
	storage         storage.ObjectStorage
}

type postService interface {
	Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error)
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
	storage storage.ObjectStorage,
) {
	h := &handler{
		service:         service,
		customerService: customerService,
		storage:         storage,
	}

	r.GET("/", h.list)
	r.GET("/:id", h.get)

	protected := r.Group("/").Use(authMiddleware)
	protected.POST("/", h.create)
}

type postResponse struct {
	ID         string            `json:"id"`
	CustomerID string            `json:"customerId"`
	VenueID    string            `json:"venueId"`
	Text       string            `json:"text"`
	Rating     int16             `json:"rating"`
	Status     entity.PostStatus `json:"status"`
	Images     []string          `json:"images"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

func postToResponse(post entity.Post, storage storage.ObjectStorage) postResponse {
	imageURLs := make([]string, 0, len(post.Images))

	if storage != nil {
		for _, img := range post.Images {
			imageURLs = append(imageURLs, storage.BuildURL(img.ObjectKey))
		}
	}
	return postResponse{
		ID:         post.ID,
		CustomerID: post.CustomerID,
		VenueID:    post.VenueID,
		Text:       post.Text,
		Rating:     post.Rating,
		Status:     post.Status,
		Images:     imageURLs,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}
}
