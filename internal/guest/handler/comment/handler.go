package comment

import (
	"context"
	"fmt"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type handler struct {
	service         commentService
	customerService customerService
}

type commentService interface {
	Create(ctx context.Context, in dto.CreateCommentInput) (entity.Comment, error)
	Update(ctx context.Context, postID int64, in dto.UpdateCommentInput) (entity.Comment, error)
	Delete(ctx context.Context, postID int64, commentID int64, customerID string) error
	List(ctx context.Context, in dto.ListCommentsInput) (dto.ListCommentsOutput, error)
}

type customerService interface {
	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	service commentService,
	customerService customerService,
	authMiddleware gin.HandlerFunc,
) {
	h := &handler{
		service:         service,
		customerService: customerService,
	}

	r.GET("/:id/comments", h.list)

	protected := r.Group("/:id/comments").Use(authMiddleware)
	protected.POST("/", h.create)
	protected.PATCH("/:comment_id", h.update)
	protected.DELETE("/:comment_id", h.delete)
}

type commentResponse struct {
	ID        int64       `json:"id"`
	PostID    int64       `json:"postId"`
	Text      string      `json:"text"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Customer  customerDTO `json:"customer"`
}

type customerDTO struct {
	ID        string  `json:"id"`
	UserName  string  `json:"userName"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	AvatarURL *string `json:"avatarURL"`
}

func commentToResponse(comment entity.Comment, customer entity.Customer) commentResponse {
	var avatarURL *string
	if customer.AvatarObjectKey != nil {
		// TODO: replace with real s3 presigned url
		url := fmt.Sprintf("https://test.com/%s", *customer.AvatarObjectKey)
		avatarURL = &url
	}

	return commentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Customer: customerDTO{
			ID:        customer.ID,
			UserName:  customer.UserName,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			AvatarURL: avatarURL,
		},
	}
}
