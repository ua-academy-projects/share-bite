package dto

import "github.com/ua-academy-projects/share-bite/internal/guest/entity"

type CreateCommentInput struct {
	PostID     int64
	CustomerID string
	Text       string
}

type UpdateCommentInput struct {
	CommentID  int64
	CustomerID string
	Text       string
}

type ListCommentsInput struct {
	PostID int64
	Limit  int
	Offset int
}

type ListCommentsOutput struct {
	Total    int
	Comments []entity.CommentWithCustomer
}
