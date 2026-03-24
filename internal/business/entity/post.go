package entity

import "time"

type Post struct {
	ID        int
	OrgID     int
	Content   string
	ImageURL  string
	CreatedAt time.Time
}
