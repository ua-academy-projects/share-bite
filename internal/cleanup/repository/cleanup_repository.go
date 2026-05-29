package repository

import (
	"context"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/cleanup/entity"
)

type CleanupRepository interface {
	ExpireOldPosts(ctx context.Context, olderThan time.Time, batchSize int, dryRun bool) (int64, error)

	CountExpiredPosts(ctx context.Context, olderThan time.Time) (int64, error)

	GetExpiredPosts(ctx context.Context, olderThan time.Time, limit int, offset int) ([]*entity.ExpiredPost, error)

	DeletePostsByID(ctx context.Context, postIDs []int64) (int64, error)

	CountPasswordResetTokens(ctx context.Context, expiredBefore time.Time) (int64, error)

	DeleteExpiredPasswordResetTokens(ctx context.Context, expiredBefore time.Time, batchSize int, dryRun bool) (int64, error)
}
