package business

import (
	"context"
	"fmt"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

func (r *Repository) CreateLike(ctx context.Context, postID int64, customerID string) (*entity.Like, error) {
	q := database.Query{
		Name: "create_like",
		Sql: `
			INSERT INTO business.likes (post_id, customer_id, created_at)
			VALUES ($1, $2, NOW())
			RETURNING id, post_id, customer_id, created_at
		`,
	}

	var like entity.Like
	err := r.db.DB().QueryRowContext(ctx, q, postID, customerID).Scan(
		&like.ID,
		&like.PostID,
		&like.CustomerID,
		&like.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create like: %w", err)
	}

	return &like, nil
}

func (r *Repository) DeleteLike(ctx context.Context, postID int64, customerID string) error {
	q := database.Query{
		Name: "delete_like",
		Sql: `
			DELETE FROM business.likes
			WHERE post_id = $1 AND customer_id = $2
		`,
	}

	result, err := r.db.DB().ExecContext(ctx, q, postID, customerID)
	if err != nil {
		return fmt.Errorf("delete like: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) CheckUserLiked(ctx context.Context, postID int64, customerID string) (bool, error) {
	q := database.Query{
		Name: "check_user_liked",
		Sql: `
			SELECT EXISTS(
				SELECT 1 FROM business.likes
				WHERE post_id = $1 AND customer_id = $2
			)
		`,
	}

	var exists bool
	err := r.db.DB().QueryRowContext(ctx, q, postID, customerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user liked: %w", err)
	}

	return exists, nil
}

func (r *Repository) CountLikesByPost(ctx context.Context, postID int64) (int, error) {
	q := database.Query{
		Name: "count_likes_by_post",
		Sql: `
			SELECT COUNT(*) FROM business.likes
			WHERE post_id = $1
		`,
	}

	var count int
	err := r.db.DB().QueryRowContext(ctx, q, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count likes by post: %w", err)
	}

	return count, nil
}

func (r *Repository) GetLikesByPost(ctx context.Context, postID int64, limit, offset int) ([]entity.Like, error) {
	q := database.Query{
		Name: "get_likes_by_post",
		Sql: `
			SELECT id, post_id, customer_id, created_at
			FROM business.likes
			WHERE post_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get likes by post: %w", err)
	}
	defer rows.Close()

	var likes []entity.Like
	for rows.Next() {
		var like entity.Like
		if err := rows.Scan(
			&like.ID,
			&like.PostID,
			&like.CustomerID,
			&like.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan like row: %w", err)
		}
		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate likes: %w", err)
	}

	return likes, nil
}
