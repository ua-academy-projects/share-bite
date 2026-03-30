package business

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

type Repository struct {
	db database.Client
}

func New(db database.Client) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetOrgIDByUserID(ctx context.Context, userID int64) (int, error) {
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
		if err == sql.ErrNoRows {
			return 0, biserr.ErrNotFound
		}
		return 0, err
	}

	return orgID, nil
}

func (r *Repository) UpdatePost(ctx context.Context, postID int64, orgID int, content string) (*entity.Post, error) {

	q := database.Query{
		Name: "update_post",
		Sql: `
		UPDATE business.posts
		SET content = $1
		WHERE id = $2 AND org_id = $3
		RETURNING id, org_id, content, image_url, created_at
	`,
	}

	var post entity.Post

	err := r.db.DB().QueryRowContext(ctx, q, content, postID, orgID).
		Scan(&post.ID, &post.OrgID, &post.Content, &post.ImageURL, &post.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biserr.ErrNotFound
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
		return biserr.ErrNotFound
	}

	return nil
}
