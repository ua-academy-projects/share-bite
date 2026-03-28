package comment

import (
	"context"
	"errors"
	"strconv"

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

func (r *Repository) Update(ctx context.Context, postID int64, in entity.UpdateCommentInput) (entity.Comment, error) {
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

func (r *Repository) Delete(ctx context.Context, commentID int64, postID int64) error {
	sql := `DELETE FROM guest.comments WHERE id = $1 AND post_id = $2`
	q := database.Query{Name: "comment_repository.Delete", Sql: sql}

	_, err := r.db.DB().ExecContext(ctx, q, commentID, postID)
	return err
}

func (r *Repository) List(ctx context.Context, in entity.ListCommentsInput) (entity.ListCommentsOutput, error) {
	var cursorID int64
	if in.PageToken != "" {
		parsed, err := strconv.ParseInt(in.PageToken, 10, 64)
		if err == nil {
			cursorID = parsed
		}
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
	`

	args := []any{in.PostID}

	if cursorID > 0 {
		sql += ` AND c.id < $2`
		args = append(args, cursorID)
	}

	limit := in.PageSize + 1
	sql += ` ORDER BY c.id DESC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit)

	q := database.Query{Name: "comment_repository.List", Sql: sql}
	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return entity.ListCommentsOutput{}, err
	}
	defer rows.Close()

	var rowsData []CommentRow
	if err := pgxscan.ScanAll(&rowsData, rows); err != nil {
		return entity.ListCommentsOutput{}, err
	}

	var out entity.ListCommentsOutput

	if len(rowsData) > in.PageSize {
		nextItem := rowsData[in.PageSize-1]
		out.NextPageToken = strconv.FormatInt(nextItem.ID, 10)
		rowsData = rowsData[:in.PageSize]
	}

	for _, row := range rowsData {
		out.Comments = append(out.Comments, row.ToEntity())
	}

	return out, nil
}
