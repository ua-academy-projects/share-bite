package entity

import "time"

type Comment struct {
	ID         int64
	PostID     int64
	CustomerID string
	Text       string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentWithCustomer struct {
	Comment  Comment
	Customer Customer
}

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
	PostID    int64
	PageSize  int
	PageToken string // Курсор (например, base64)
}

type ListCommentsOutput struct {
	Comments      []CommentWithCustomer
	NextPageToken string
}
