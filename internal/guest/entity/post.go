package entity

import "time"

type Post struct {
	ID string

	CustomerID string
	VenueID    string
	Text       string
	Rating     int16
	Status     string

	LikesCount  int
	IsLikedByMe bool

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
	Limit      int
	Offset     int
	CustomerID string
}

type ListPostsOutput struct {
	Posts []Post
	Total int
}
