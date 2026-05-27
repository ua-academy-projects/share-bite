package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/cleanup/entity"
	"github.com/ua-academy-projects/share-bite/internal/cleanup/repository"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

type CleanupService struct {
	repo repository.CleanupRepository
}

func NewCleanupService(repo repository.CleanupRepository) *CleanupService {
	return &CleanupService{repo: repo}
}

func (s *CleanupService) ExpireOldPosts(ctx context.Context, retentionPeriod time.Duration, batchSize int, dryRun bool) (*entity.CleanupResult, error) {
	result := &entity.CleanupResult{
		Name:      "expire-old-posts",
		DryRun:    dryRun,
		StartedAt: time.Now(),
	}

	cutoffTime := time.Now().Add(-retentionPeriod)

	logger.Info(ctx, fmt.Sprintf("Starting cleanup of posts older than %v (cutoff: %v)", retentionPeriod, cutoffTime))

	count, err := s.repo.CountExpiredPosts(ctx, cutoffTime)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to count expired posts: %v", err))
		result.CompletedAt = time.Now()
		return result, err
	}

	result.RecordsFound = count
	logger.Info(ctx, fmt.Sprintf("Found %d posts older than %v", count, cutoffTime))

	if count == 0 {
		logger.Info(ctx, "No posts to cleanup")
		result.CompletedAt = time.Now()
		return result, nil
	}

	deleted, err := s.repo.ExpireOldPosts(ctx, cutoffTime, batchSize, dryRun)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to delete posts: %v", err))
		result.CompletedAt = time.Now()
		return result, err
	}

	result.RecordsDeleted = deleted
	result.CompletedAt = time.Now()

	if dryRun {
		logger.Info(ctx, fmt.Sprintf("[DRY-RUN] Would delete %d posts", deleted))
	} else {
		logger.Info(ctx, fmt.Sprintf("Successfully deleted %d posts in %v", deleted, result.Duration()))
	}

	return result, nil
}

func (s *CleanupService) CleanExpiredPasswordResetTokens(ctx context.Context, tokenRetentionPeriod time.Duration, batchSize int, dryRun bool) (*entity.CleanupResult, error) {
	result := &entity.CleanupResult{
		Name:      "cleanup-expired-tokens",
		DryRun:    dryRun,
		StartedAt: time.Now(),
	}

	cutoffTime := time.Now().Add(-tokenRetentionPeriod)

	logger.Info(ctx, fmt.Sprintf("Starting cleanup of password reset tokens expired before %v", cutoffTime))

	count, err := s.repo.CountPasswordResetTokens(ctx, cutoffTime)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to count expired tokens: %v", err))
		result.CompletedAt = time.Now()
		return result, err
	}

	result.RecordsFound = count
	logger.Info(ctx, fmt.Sprintf("Found %d expired password reset tokens", count))

	if count == 0 {
		logger.Info(ctx, "No tokens to cleanup")
		result.CompletedAt = time.Now()
		return result, nil
	}

	deleted, err := s.repo.DeleteExpiredPasswordResetTokens(ctx, cutoffTime, batchSize, dryRun)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to delete tokens: %v", err))
		result.CompletedAt = time.Now()
		return result, err
	}

	result.RecordsDeleted = deleted
	result.CompletedAt = time.Now()

	if dryRun {
		logger.Info(ctx, fmt.Sprintf("[DRY-RUN] Would delete %d password reset tokens", deleted))
	} else {
		logger.Info(ctx, fmt.Sprintf("Successfully deleted %d password reset tokens in %v", deleted, result.Duration()))
	}

	return result, nil
}

func (s *CleanupService) RunAllCleanups(ctx context.Context, retentionPeriod time.Duration, batchSize int, dryRun bool) ([]*entity.CleanupResult, error) {
	var results []*entity.CleanupResult
	var aggregatedErr error

	logger.Info(ctx, "Running cleanup operations")
	postResult, err := s.ExpireOldPosts(ctx, retentionPeriod, batchSize, dryRun)
	// Always append result even if error occurred
	results = append(results, postResult)
	if err != nil {
		log.Printf("Error during post cleanup: %v", err)
		aggregatedErr = err
	}

	tokenResult, err := s.CleanExpiredPasswordResetTokens(ctx, retentionPeriod/2, batchSize, dryRun)
	// Always append result even if error occurred
	results = append(results, tokenResult)
	if err != nil {
		log.Printf("Error during token cleanup: %v", err)
		if aggregatedErr != nil {
			aggregatedErr = fmt.Errorf("%w; token cleanup error: %w", aggregatedErr, err)
		} else {
			aggregatedErr = err
		}
	}

	totalFound := int64(0)
	totalDeleted := int64(0)
	for _, result := range results {
		totalFound += result.RecordsFound
		totalDeleted += result.RecordsDeleted
	}

	if dryRun {
		logger.Info(ctx, fmt.Sprintf("[DRY-RUN] Cleanup complete. Found %d records, would delete %d records", totalFound, totalDeleted))
	} else {
		logger.Info(ctx, fmt.Sprintf("Cleanup complete. Found %d records, deleted %d records", totalFound, totalDeleted))
	}

	return results, aggregatedErr
}
