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
