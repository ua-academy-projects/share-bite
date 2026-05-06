package entity

import "time"

type PostImage struct {
	ID          string
	PostID      string
	ObjectKey   string
	ContentType string
	FileSize    int64
	SortOrder   int16
	CreatedAt   time.Time
}
