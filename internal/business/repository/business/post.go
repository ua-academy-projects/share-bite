package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
)

func (r *Repository) GetPostByID(ctx context.Context, postID int64) (*entity.Post, error) {
	q := database.Query{
		Name: "get_post_by_id",
		Sql: `
			SELECT id, org_id, content, created_at
		FROM business.posts
		WHERE id = $1
		`,
	}

	var post entity.Post

	err := r.db.DB().QueryRowContext(ctx, q, postID).Scan(
		&post.ID,
		&post.OrgID,
		&post.Content,
		&post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	return &post, nil
}

func (r *Repository) GetOrgIDByUserID(ctx context.Context, userID string) (int, error) {
	q := database.Query{
		Name: "get_org_by_user_id",
		Sql: `
			SELECT id
			FROM business.org_units
			WHERE org_account_id = $1
		`,
	}

	var orgID int

	err := r.db.DB().QueryRowContext(ctx, q, userID).Scan(&orgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return orgID, nil
}

func (r *Repository) GetPostPhotos(ctx context.Context, postID int64) ([]string, error) {
	q := database.Query{
		Name: "get_post_photos",
		Sql: `
		SELECT image_url
		FROM business.post_photos
		WHERE post_id = $1
		ORDER BY sort_order;
		`,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var imageUrls []string

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		imageUrls = append(imageUrls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return imageUrls, nil
}

func (r *Repository) UpdatePost(ctx context.Context, postID int64, orgID int, content string) (*entity.Post, error) {

	q := database.Query{
		Name: "update_post",
		Sql: `
		UPDATE business.posts
		SET content = $1
		WHERE id = $2 AND org_id = $3
		RETURNING id, org_id, content, created_at
	`,
	}

	var post entity.Post

	err := r.db.DB().QueryRowContext(ctx, q, content, postID, orgID).
		Scan(&post.ID, &post.OrgID, &post.Content, &post.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (r *Repository) DeletePost(ctx context.Context, id int64, orgID int) error {

	q := database.Query{
		Name: "delete_post",
		Sql: `
		DELETE FROM business.posts
		WHERE id = $1 AND org_id = $2
	`,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, id, orgID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) CheckOwnership(ctx context.Context, userID string, unitID int) error {
	checkQuery := `SELECT id
					FROM business.org_units
					WHERE id = $1
					  AND (
					  	org_account_id = $2
						OR 
						parent_id IN (
							SELECT id FROM business.org_units WHERE org_account_id = $2
						)
					);`

	q := database.Query{
		Name: "check_ownership.CheckOwnership",
		Sql:  checkQuery,
	}

	var foundID int
	var err error

	if tx, ok := ctx.Value(pg.TxKey).(pgx.Tx); ok {
		err = tx.QueryRow(ctx, checkQuery, unitID, userID).Scan(&foundID)
	} else {
		err = r.db.DB().QueryRowContext(ctx, q, unitID, userID).Scan(&foundID)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return ErrForbidden
		}
		return fmt.Errorf("execute check ownership query: %w", err)
	}
	return nil
}

func (r *Repository) CreatePost(ctx context.Context, userID string, unitID int, description string) (*entity.Post, error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("transaction not found in context")
	}
	var post entity.Post
	postQuery := `INSERT INTO business.posts (org_id, content)
					SELECT $1, $2
					WHERE EXISTS (
   							SELECT 1 
							FROM business.org_units 
							WHERE id = $1 AND (org_account_id = $3 OR parent_id IN (SELECT id FROM business.org_units WHERE org_account_id = $3))
							)
					RETURNING id, org_id, content, created_at;`

	err := tx.QueryRow(ctx, postQuery, unitID, description, userID).
		Scan(&post.ID, &post.OrgID, &post.Content, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return nil, ErrForbidden
		}
		return nil, fmt.Errorf("insert post query: %w", err)
	}
	return &post, nil
}

func (r *Repository) InsertPostImages(ctx context.Context, postID int64, URLs []string) error {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("transaction not found in context")
	}
	imagesQuery := `INSERT INTO business.post_photos (post_id, image_url, sort_order) 
					VALUES ($1, $2, $3)`
	for order, url := range URLs {
		if _, err := tx.Exec(ctx, imagesQuery, postID, url, order); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetPosts(ctx context.Context, limit, offset int) (pagination.Result[entity.Post], error) {

	params := pagination.Params{
		Table:   "business.posts",
		Columns: "id, org_id, content, created_at",
		Where:   "TRUE",
		OrderBy: "created_at DESC, id DESC",
		Args:    []any{},
		Offset:  offset,
		Limit:   limit,
	}

	return pagination.List(
		ctx,
		r.db.DB(),
		"posts",
		params,
		func(rows pgx.Rows) (entity.Post, error) {
			var p entity.Post
			err := rows.Scan(&p.ID, &p.OrgID, &p.Content, &p.CreatedAt)
			return p, err
		},
	)
}
