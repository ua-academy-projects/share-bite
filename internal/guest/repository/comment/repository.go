package comment

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	apperror "github.com/ua-academy-projects/share-bite/internal/guest/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, in entity.CreateCommentInput) (entity.Comment, error) {
	sql := `
		INSERT INTO guest.comments (post_id, customer_id, text)
		VALUES ($1, $2, $3)
		RETURNING id, post_id, customer_id, text, created_at, updated_at
	`
	q := database.Query{Name: "comment_repository.Create", Sql: sql}

	row, err := r.db.DB().QueryContext(ctx, q, in.PostID, in.CustomerID, in.Text)
	if err != nil {
		return entity.Comment{}, err
	}
	defer row.Close()

	var comment Comment
	if err := pgxscan.ScanOne(&comment, row); err != nil {
		return entity.Comment{}, err
	}

	return comment.ToEntity(), nil
}

func (r *Repository) GetByID(ctx context.Context, commentID int64) (entity.Comment, error) {
	sql := `SELECT id, post_id, customer_id, text, created_at, updated_at FROM guest.comments WHERE id = $1`
	q := database.Query{Name: "comment_repository.GetByID", Sql: sql}

	row, err := r.db.DB().QueryContext(ctx, q, commentID)
	if err != nil {
		return entity.Comment{}, err
	}
	defer row.Close()

	var comment Comment
	if err := pgxscan.ScanOne(&comment, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Comment{}, apperror.CommentNotFoundID(commentID)
		}
		return entity.Comment{}, err
	}

	return comment.ToEntity(), nil
}

func (r *Repository) Update(ctx context.Context, in entity.UpdateCommentInput) (entity.Comment, error) {
	sql := `
		UPDATE guest.comments 
		SET text = $1, updated_at = NOW() 
		WHERE id = $2 
		RETURNING id, post_id, customer_id, text, created_at, updated_at
	`
	q := database.Query{Name: "comment_repository.Update", Sql: sql}

	row, err := r.db.DB().QueryContext(ctx, q, in.Text, in.CommentID)
	if err != nil {
		return entity.Comment{}, err
	}
	defer row.Close()

	var comment Comment
	if err := pgxscan.ScanOne(&comment, row); err != nil {
		return entity.Comment{}, err
	}

	return comment.ToEntity(), nil
}

func (r *Repository) Delete(ctx context.Context, commentID int64) error {
	sql := `DELETE FROM guest.comments WHERE id = $1`
	q := database.Query{Name: "comment_repository.Delete", Sql: sql}

	_, err := r.db.DB().ExecContext(ctx, q, commentID)
	return err
}

func (r *Repository) List(ctx context.Context, in entity.ListCommentsInput) (entity.ListCommentsOutput, error) {
	countSQL := `SELECT COUNT(*) FROM guest.comments WHERE post_id = $1`
	var total int
	if err := r.db.DB().QueryRowContext(ctx, database.Query{Name: "comment.Count", Sql: countSQL}, in.PostID).Scan(&total); err != nil {
		return entity.ListCommentsOutput{}, err
	}

	sql := `
		SELECT id, post_id, customer_id, text, created_at, updated_at 
		FROM guest.comments 
		WHERE post_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.DB().QueryContext(ctx, database.Query{Name: "comment.List", Sql: sql}, in.PostID, in.Limit, in.Offset)
	if err != nil {
		return entity.ListCommentsOutput{}, err
	}
	defer rows.Close()

	var comments Comments
	if err := pgxscan.ScanAll(&comments, rows); err != nil {
		return entity.ListCommentsOutput{}, err
	}

	return entity.ListCommentsOutput{
		Comments: comments.ToEntities(),
		Total:    total,
	}, nil
}
