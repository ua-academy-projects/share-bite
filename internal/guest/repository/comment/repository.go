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
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Comment{}, apperror.CommentNotFoundID(in.CommentID)
		}
		return entity.Comment{}, err
	}

	return comment.ToEntity(), nil
}

func (r *Repository) Delete(ctx context.Context, commentID int64) error {
	sql := `DELETE FROM guest.comments WHERE id = $1`
	q := database.Query{Name: "comment_repository.Delete", Sql: sql}

	_, err := r.db.DB().ExecContext(ctx, q, commentID)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return apperror.CommentNotFoundID(commentID)
	}
	return err
}

func (r *Repository) List(ctx context.Context, in entity.ListCommentsInput) (entity.ListCommentsOutput, error) {
	var out entity.ListCommentsOutput

	countSql := `SELECT COUNT(*) FROM guest.comments WHERE post_id = $1`
	qCount := database.Query{Name: "comment_repository.ListCount", Sql: countSql}

	if err := r.db.DB().QueryRowContext(ctx, qCount, in.PostID).Scan(&out.Total); err != nil {
		return out, err
	}

	if out.Total == 0 {
		out.Comments = make([]entity.CommentWithCustomer, 0)
		return out, nil
	}

	sql := `
		SELECT 
			c.id, c.post_id, c.customer_id, c.text, c.created_at, c.updated_at,
			cust.id AS cust_id, cust.user_id AS cust_user_id, cust.username AS cust_username,
			cust.first_name AS cust_first_name, cust.last_name AS cust_last_name,
			cust.avatar_object_key AS cust_avatar, cust.bio AS cust_bio
		FROM guest.comments c
		JOIN guest.customers cust ON c.customer_id = cust.id
		WHERE c.post_id = $1
		ORDER BY c.id DESC
		LIMIT $2 OFFSET $3
	`
	qData := database.Query{Name: "comment_repository.List", Sql: sql}

	rows, err := r.db.DB().QueryContext(ctx, qData, in.PostID, in.Limit, in.Offset)
	if err != nil {
		return out, err
	}
	defer rows.Close()

	var rowsData []CommentRow
	if err := pgxscan.ScanAll(&rowsData, rows); err != nil {
		return out, err
	}

	for _, row := range rowsData {
		out.Comments = append(out.Comments, row.ToEntity())
	}

	return out, nil
}
