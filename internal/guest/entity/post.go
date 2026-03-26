package entity

import "time"

type Post struct {
	ID string

	CustomerID string
	VenueID    string
	Text       string
	Rating     int16
	Status     string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePostInput struct {
	CustomerID string
	VenueID    string
	Text       string
	Rating     int16
}

type ListPostsInput struct {
	Limit  int
	Offset int
}

type ListPostsOutput struct {
	Posts []Post
	Total int
}
