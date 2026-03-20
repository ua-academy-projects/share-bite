package entity

import "time"

type Post struct {
	ID string

	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListPostsInput struct {
	Limit  int
	Offset int
}

type ListPostsOutput struct {
	Posts []Post
	Total int
}
