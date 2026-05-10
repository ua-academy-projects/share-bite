package post

import (
	"context"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"time"
)

const (
	postCleanupInterval = 5 * time.Minute
)

func StartPostCleanupJob(ctx context.Context, svc postCleanupService) {
	ticker := time.NewTicker(postCleanupInterval)
	go func() {
		defer ticker.Stop()
		logger.Info(ctx, "post cleanup job started")
		runCleanup(ctx, svc)
		for {
			select {
			case <-ctx.Done():
				logger.Info(ctx, "post cleanup job stopped")
				return
			case <-ticker.C:
				runCleanup(ctx, svc)
			}
		}
	}()

}

func runCleanup(ctx context.Context, svc postCleanupService) {
	logger.Info(ctx, "cleanup expired posts started")

	if err := svc.CleanupExpiredPosts(ctx); err != nil {
		logger.ErrorKV(ctx, "cleanup expired posts failed", "error", err)
		return
	}

	logger.Info(ctx, "cleanup expired posts finished")
}
