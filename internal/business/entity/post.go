package entity

import "time"

type Post struct {
	ID        int64
	OrgID     int
	Content   string
	ImageURL  string
	CreatedAt time.Time
}
