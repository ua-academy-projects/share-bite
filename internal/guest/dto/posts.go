package dto

import (
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"io"
)

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

type ListPostsInput struct {
	Limit      int
	Offset     int
	CustomerID string
}

type ListPostsOutput struct {
	Posts []entity.Post
	Total int
}
