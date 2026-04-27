package dto

import (
	"io"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
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

type ExploreNearbyInput struct {
	Lat   float64 `form:"lat" binding:"required,latitude"`
	Lon   float64 `form:"lon" binding:"required,longitude"`
	Limit int     `form:"limit" binding:"max=100"`
}

type ExploreVenueItem struct {
	VenueID int64      `json:"venue_id"`
	Posts   []PostItem `json:"posts"`
}

type PostItem struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Images    []string  `json:"images"`
}
