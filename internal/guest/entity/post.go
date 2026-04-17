package entity

import (
	"time"
)

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

type Post struct {
	ID string

	CustomerID string
	VenueID    string
	Text       string
	Rating     int16
	Status     PostStatus

	LikesCount  int
	IsLikedByMe bool

	Images []PostImage

	CreatedAt time.Time
	UpdatedAt time.Time
}
