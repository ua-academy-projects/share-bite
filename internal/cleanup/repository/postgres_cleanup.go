package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/cleanup/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type PostgresCleanupRepository struct {
	db database.Client
}

func NewPostgresCleanupRepository(db database.Client) CleanupRepository {
	return &PostgresCleanupRepository{db: db}
}

func (r *PostgresCleanupRepository) ExpireOldPosts(ctx context.Context, olderThan time.Time, batchSize int, dryRun bool) (int64, error) {
	if dryRun {
		count, err := r.CountExpiredPosts(ctx, olderThan)
		if err != nil {
			return 0, fmt.Errorf("failed to count expired posts: %w", err)
		}
		log.Printf("[DRY-RUN] Would delete %d expired posts older than %v", count, olderThan)
		return count, nil
	}

	totalDeleted := int64(0)

	for {
		// Atomic delete: single statement ensures status check at delete time
		q := database.Query{
			Name: "DeleteExpiredPostsAtomic",
			Sql: `
				DELETE FROM guest.posts
				WHERE id IN (
				  SELECT id FROM guest.posts
				  WHERE created_at < $1 AND status != 'archived'
				  ORDER BY created_at ASC
				  LIMIT $2
				) AND status != 'archived'
			`,
		}

		result, err := r.db.DB().ExecContext(ctx, q, olderThan, batchSize)
		if err != nil {
			return totalDeleted, fmt.Errorf("failed to delete expired posts: %w", err)
		}

		deleted := result.RowsAffected()
		if deleted == 0 {
			break
		}

		totalDeleted += deleted
		log.Printf("Deleted batch of %d posts, total: %d", deleted, totalDeleted)

		if deleted < int64(batchSize) {
			break
		}
	}

	return totalDeleted, nil
}

func (r *PostgresCleanupRepository) CountExpiredPosts(ctx context.Context, olderThan time.Time) (int64, error) {
	var count int64

	q := database.Query{
		Name: "CountExpiredPosts",
		Sql: `
			SELECT COUNT(*) FROM guest.posts
			WHERE created_at < $1 AND status != 'archived'
		`,
	}

	err := r.db.DB().QueryRowContext(ctx, q, olderThan).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count expired posts: %w", err)
	}

	return count, nil
}

func (r *PostgresCleanupRepository) GetExpiredPosts(ctx context.Context, olderThan time.Time, limit int, offset int) ([]*entity.ExpiredPost, error) {
	var posts []*entity.ExpiredPost

	q := database.Query{
		Name: "SelectExpiredPosts",
		Sql: `
			SELECT id, customer_id, created_at FROM guest.posts
			WHERE created_at < $1 AND status != 'archived'
			ORDER BY created_at ASC
			LIMIT $2 OFFSET $3
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, olderThan, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query expired posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post entity.ExpiredPost
		if err := rows.Scan(&post.ID, &post.CustomerID, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating posts: %w", err)
	}

	return posts, nil
}

func (r *PostgresCleanupRepository) DeletePostsByID(ctx context.Context, postIDs []int64) (int64, error) {
	if len(postIDs) == 0 {
		return 0, nil
	}

	q := database.Query{
		Name: "DeletePostsByID",
		Sql: `
			DELETE FROM guest.posts
			WHERE id = ANY($1)
		`,
	}

	result, err := r.db.DB().ExecContext(ctx, q, postIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete posts: %w", err)
	}

	return result.RowsAffected(), nil
}

func (r *PostgresCleanupRepository) CountPasswordResetTokens(ctx context.Context, expiredBefore time.Time) (int64, error) {
	var count int64

	q := database.Query{
		Name: "CountPasswordResetTokens",
		Sql: `
			SELECT COUNT(*) FROM auth.password_reset_tokens
			WHERE expires_at < $1 AND used_at IS NULL
		`,
	}

	err := r.db.DB().QueryRowContext(ctx, q, expiredBefore).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count expired password reset tokens: %w", err)
	}

	return count, nil
}

func (r *PostgresCleanupRepository) DeleteExpiredPasswordResetTokens(ctx context.Context, expiredBefore time.Time, batchSize int, dryRun bool) (int64, error) {
	if dryRun {
		count, err := r.CountPasswordResetTokens(ctx, expiredBefore)
		if err != nil {
			return 0, err
		}
		log.Printf("[DRY-RUN] Would delete %d expired password reset tokens before %v", count, expiredBefore)
		return count, nil
	}

	totalDeleted := int64(0)

	for {
		// Use CTE for batch deletion since PostgreSQL doesn't support DELETE ... LIMIT directly
		q := database.Query{
			Name: "DeleteExpiredPasswordResetTokens",
			Sql: `
				WITH to_delete AS (
				  SELECT ctid FROM auth.password_reset_tokens
				  WHERE expires_at < $1 AND used_at IS NULL
				  LIMIT $2
				)
				DELETE FROM auth.password_reset_tokens
				WHERE ctid IN (SELECT ctid FROM to_delete)
			`,
		}

		result, err := r.db.DB().ExecContext(ctx, q, expiredBefore, batchSize)
		if err != nil {
			return totalDeleted, fmt.Errorf("failed to delete expired password reset tokens: %w", err)
		}

		deleted := result.RowsAffected()
		totalDeleted += deleted

		if deleted == 0 {
			break
		}

		log.Printf("Deleted batch of %d password reset tokens, total: %d", deleted, totalDeleted)

		if deleted < int64(batchSize) {
			break
		}
	}

	return totalDeleted, nil
}
