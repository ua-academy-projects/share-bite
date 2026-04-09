package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

func (r *Repository) CreateComment(ctx context.Context, postID int64, authorID, content string) (*entity.Comment, error) {
	q := database.Query{
		Name: "create_comment",
		Sql: `
			INSERT INTO business.comments (post_id, author_id, content, created_at)
			VALUES ($1, $2, $3, NOW())
			RETURNING id, post_id, author_id, content, created_at
		`,
	}

	var comment entity.Comment
	err := r.db.DB().QueryRowContext(ctx, q, postID, authorID, content).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.AuthorID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	return &comment, nil
}

func (r *Repository) GetCommentByID(ctx context.Context, commentID int64) (*entity.Comment, error) {
	q := database.Query{
		Name: "get_comment_by_id",
		Sql: `
			SELECT id, post_id, author_id, content, created_at
			FROM business.comments
			WHERE id = $1
		`,
	}

	var comment entity.Comment
	err := r.db.DB().QueryRowContext(ctx, q, commentID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.AuthorID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get comment by id: %w", err) //
	}

	return &comment, nil
}

func (r *Repository) UpdateComment(ctx context.Context, commentID int64, content string) (*entity.Comment, error) {
	q := database.Query{
		Name: "update_comment",
		Sql: `
			UPDATE business.comments
			SET content = $1
			WHERE id = $2
			RETURNING id, post_id, author_id, content, created_at
		`,
	}

	var comment entity.Comment
	err := r.db.DB().QueryRowContext(ctx, q, content, commentID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.AuthorID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update comment: %w", err)
	}

	return &comment, nil
}

func (r *Repository) DeleteComment(ctx context.Context, commentID int64) error {
	q := database.Query{
		Name: "delete_comment",
		Sql: `
			DELETE FROM business.comments
			WHERE id = $1
		`,
	}

	result, err := r.db.DB().ExecContext(ctx, q, commentID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) ListCommentsByPost(ctx context.Context, postID int64, limit, offset int) ([]entity.Comment, error) {
	q := database.Query{
		Name: "list_comments_by_post",
		Sql: `
			SELECT id, post_id, author_id, content, created_at
			FROM business.comments
			WHERE post_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list comments by post: %w", err)
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.Content,
			&comment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan comment row: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate comments: %w", err)
	}

	return comments, nil
}

func (r *Repository) ListCommentsWithAuthorsByPost(ctx context.Context, postID int64, limit, offset int) ([]entity.CommentWithAuthor, error) {
	q := database.Query{
		Name: "list_comments_with_authors_by_post",
		Sql: `
			SELECT 
				c.id, c.post_id, c.content, c.created_at,
				c.author_id,
				cu.username, cu.first_name, cu.last_name, cu.avatar_object_key
			FROM business.comments c
			JOIN guest.customers cu ON c.author_id = cu.user_id
			WHERE c.post_id = $1
			ORDER BY c.created_at DESC
			LIMIT $2 OFFSET $3
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list comments with authors: %w", err)
	}
	defer rows.Close()

	var comments []entity.CommentWithAuthor
	for rows.Next() {
		var comment entity.CommentWithAuthor
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.AuthorID,
			&comment.AuthorUsername,
			&comment.AuthorFirstName,
			&comment.AuthorLastName,
			&comment.AuthorAvatarURL,
		); err != nil {
			return nil, fmt.Errorf("scan comment with author: %w", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate comments with authors: %w", err)
	}

	return comments, nil
}

func (r *Repository) CountCommentsByPost(ctx context.Context, postID int64) (int, error) {
	q := database.Query{
		Name: "count_comments_by_post",
		Sql: `
			SELECT COUNT(*) FROM business.comments
			WHERE post_id = $1
		`,
	}

	var count int
	err := r.db.DB().QueryRowContext(ctx, q, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count comments by post: %w", err)
	}

	return count, nil
}
