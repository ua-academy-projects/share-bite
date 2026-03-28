package comment

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type handler struct {
	service         commentService
	customerService customerService
}

type commentService interface {
	Create(ctx context.Context, in entity.CreateCommentInput) (entity.Comment, error)
	Update(ctx context.Context, in entity.UpdateCommentInput) (entity.Comment, error)
	Delete(ctx context.Context, commentID int64, customerID string) error
	List(ctx context.Context, in entity.ListCommentsInput) (entity.ListCommentsOutput, error)
}

type customerService interface {
	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
}

func RegisterHandlers(r *gin.RouterGroup, service commentService, customerService customerService, authMiddleware gin.HandlerFunc) {
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
	ID         int64     `json:"id"`
	PostID     int64     `json:"postId"`
	CustomerID string    `json:"customerId"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func commentToResponse(comment entity.Comment) commentResponse {
	return commentResponse{
		ID:         comment.ID,
		PostID:     comment.PostID,
		CustomerID: comment.CustomerID,
		Text:       comment.Text,
		CreatedAt:  comment.CreatedAt,
		UpdatedAt:  comment.UpdatedAt,
	}
}
