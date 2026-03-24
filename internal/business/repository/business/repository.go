package business

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
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

func (r *Repository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {

	q := database.Query{
		Name: "get_post_by_id",
		Sql: `
		SELECT id, org_id, content, image_url, created_at
		FROM business.posts
		WHERE id = $1
	`,
	}

	row := r.db.DB().QueryRowContext(ctx, q, id)

	post := &entity.Post{}
	err := row.Scan(
		&post.ID,
		&post.OrgID,
		&post.Content,
		&post.ImageURL,
		&post.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, biserr.ErrNotFound
		}
		return nil, err
	}

	return post, nil
}

func (r *Repository) UpdatePost(ctx context.Context, post *entity.Post) error {

	q := database.Query{
		Name: "update_post",
		Sql: `
		UPDATE business.posts
		SET content = $1
		WHERE id = $2
	`,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, post.Content, post.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return biserr.ErrNotFound
	}

	return nil
}

func (r *Repository) DeletePost(ctx context.Context, id int) error {

	q := database.Query{
		Name: "delete_post",
		Sql: `
			DELETE FROM business.posts
			WHERE id = $1
		`,
	}

	tag, err := r.db.DB().ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return biserr.ErrNotFound
	}

	return nil
}
