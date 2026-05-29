package post

import (
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"time"
)

type PostImage struct {
	ID          string `db:"id"`
	PostID      string `db:"post_id"`
	ObjectKey   string `db:"object_key"`
	ContentType string `db:"content_type"`
	FileSize    int64  `db:"file_size"`
	SortOrder   int16  `db:"sort_order"`

	ProcessingStatus entity.ImageProcessingStatus `db:"processing_status"`
	ThumbnailKey     *string                      `db:"thumbnail_key"`
	Width            *int                         `db:"width"`
	Height           *int                         `db:"height"`
	ProcessedAt      *time.Time                   `db:"processed_at"`
	FailureReason    *string                      `db:"failure_reason"`

	CreatedAt time.Time `db:"created_at"`
}

func (p *PostImage) ToEntity() entity.PostImage {
	return entity.PostImage{
		ID:               p.ID,
		PostID:           p.PostID,
		ObjectKey:        p.ObjectKey,
		ContentType:      p.ContentType,
		FileSize:         p.FileSize,
		SortOrder:        p.SortOrder,
		ProcessingStatus: p.ProcessingStatus,
		ThumbnailKey:     p.ThumbnailKey,
		Width:            p.Width,
		Height:           p.Height,
		ProcessedAt:      p.ProcessedAt,
		FailureReason:    p.FailureReason,
		CreatedAt:        p.CreatedAt,
	}
}

type PostImages []PostImage

func (p PostImages) ToEntities() []entity.PostImage {
	res := make([]entity.PostImage, 0, len(p))
	for i := range p {
		res = append(res, p[i].ToEntity())
	}
	return res
}
