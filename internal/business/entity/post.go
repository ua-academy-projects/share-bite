package entity

import (
	"time"
)

type Post struct {
	ID        int64
	OrgId     int
	Content   string
	CreatedAt time.Time
}

type PostPhotos struct {
	PostID    int
	ImageURL  string
	SortOrder int
}


type PostWithPhotos struct {
	ID        int64
	OrgID     int
	Content   string
	CreatedAt time.Time
	Images    []string
}
