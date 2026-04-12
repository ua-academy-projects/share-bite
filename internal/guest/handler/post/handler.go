package post

import (
	"context"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"

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
	Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error)
	Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error)
	Delete(ctx context.Context, postID, customerID string) error
	List(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error)
	Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error)
	Like(ctx context.Context, postID string, customerID string) error
	Unlike(ctx context.Context, postID string, customerID string) error
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
	protected.PATCH("/:id", h.update)
	protected.DELETE("/:id", h.delete)
	protected.POST("/:id/like", h.like)
	protected.DELETE("/:id/like", h.unlike)
}

type postResponse struct {
	ID          string            `json:"id"`
	CustomerID  string            `json:"customerId"`
	VenueID     int64             `json:"venueId"`
	Text        string            `json:"text"`
	Rating      int16             `json:"rating"`
	Status      entity.PostStatus `json:"status"`
	LikesCount  int               `json:"likesCount"`
	IsLikedByMe bool              `json:"isLikedByMe"`
	Images      []string          `json:"images"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	PublishedAt *time.Time        `json:"publishedAt,omitempty"`
}

func postToResponse(post entity.Post, storage storage.ObjectStorage) postResponse {
	imageURLs := make([]string, 0, len(post.Images))

	if storage != nil {
		for _, img := range post.Images {
			imageURLs = append(imageURLs, storage.BuildURL(img.ObjectKey))
		}
	}
	return postResponse{
		ID:          post.ID,
		CustomerID:  post.CustomerID,
		VenueID:     post.VenueID,
		Text:        post.Text,
		Rating:      post.Rating,
		Status:      post.Status,
		LikesCount:  post.LikesCount,
		IsLikedByMe: post.IsLikedByMe,
		Images:      imageURLs,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
	}
}
