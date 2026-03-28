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
	Comments []Comment
	Total    int
}
