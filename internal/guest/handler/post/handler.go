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
	List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error)
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
) {
	h := &handler{
		service:         service,
		customerService: customerService,
	}

	r.GET("/", h.list)
	r.GET("/:id", h.get)

	protected := r.Group("/").Use(authMiddleware)
	protected.POST("/", h.create)
	protected.POST("/:id/like", h.like)
	protected.DELETE("/:id/like", h.unlike)
}

type postResponse struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customerId"`
	VenueID     string    `json:"venueId"`
	Text        string    `json:"text"`
	Rating      int16     `json:"rating"`
	Status      string    `json:"status"`
	LikesCount  int       `json:"likesCount"`
	IsLikedByMe bool      `json:"isLikedByMe"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func postToResponse(post entity.Post) postResponse {
	return postResponse{
		ID:          post.ID,
		CustomerID:  post.CustomerID,
		VenueID:     post.VenueID,
		Text:        post.Text,
		Rating:      post.Rating,
		Status:      post.Status,
		LikesCount:  post.LikesCount,
		IsLikedByMe: post.IsLikedByMe,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}
