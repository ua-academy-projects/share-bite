package entity

import "time"

type PostImage struct {
	ID          string
	PostID      string
	ObjectKey   string
	ContentType string
	FileSize    int64
	SortOrder   int16

	ProcessingStatus ImageProcessingStatus
	ThumbnailKey     *string
	Width            *int
	Height           *int
	ProcessedAt      *time.Time
	FailureReason    *string

	CreatedAt time.Time
}
