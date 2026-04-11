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
	PostStatusDeleted   PostStatus = "deleted"
)

type Post struct {
	ID string

	CustomerID string
	VenueID    int64
	Text       string
	Rating     int16
	Status     PostStatus

	Images []PostImage

	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt *time.Time
}

type CreatePostInput struct {
	CustomerID string
	VenueID    int64
	Text       string
	Rating     int16

	Images []UploadImageInput
}

type UploadImageInput struct {
	File        io.Reader
	ContentType string
	FileSize    int64
}

type UpdatePostInput struct {
	ID         string
	CustomerID string

	VenueID *int64
	Text    *string
	Rating  *int16
	Status  *PostStatus
}

type ListPostsInput struct {
	Limit  int
	Offset int
}

type ListPostsOutput struct {
	Posts []Post
	Total int
}
