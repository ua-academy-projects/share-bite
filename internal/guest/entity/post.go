package entity

import (
	"io"
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

type CreatePostInput struct {
	CustomerID string
	VenueID    string
	Text       string
	Rating     int16

	Images []UploadImageInput
}

type UploadImageInput struct {
	File        io.Reader
	ContentType string
	FileSize    int64
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
