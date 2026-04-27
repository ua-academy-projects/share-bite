package dto

import (
	"io"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
)

type CreatePostInput struct {
	CustomerID string
	UserID     string
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

type ListPostsInput struct {
	Limit      int
	Offset     int
	CustomerID string
}

type ListPostsOutput struct {
	Posts []entity.Post
	Total int
}
