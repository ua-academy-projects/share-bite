package entity

import "time"

type Post struct {
	ID        int64
	OrgID     int
	Content   string
	CreatedAt time.Time
}

type PostWithPhotos struct {
	ID        int64
	OrgID     int
	Content   string
	CreatedAt time.Time
	Images    []string

	OrgName     string
	ProfileType string
}
