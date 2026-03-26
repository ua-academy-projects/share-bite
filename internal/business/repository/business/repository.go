package business

import (
	"context"

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

func (r *Repository) UpdatePost(ctx context.Context, post *entity.Post) error {

	q := database.Query{
		Name: "update_post",
		Sql: `
		UPDATE business.posts
		SET content = $1
		WHERE id = $2 AND org_id = $3
	`,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, post.Content, post.ID, post.OrgID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return biserr.ErrNotFound
	}

	return nil
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
