package post

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
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, in entity.CreatePostInput) (entity.Post, error) {
	sql := `
        INSERT INTO guest.posts(
            customer_id,
            venue_id,
            text,
            rating
        ) VALUES ($1, $2, $3, $4)
        RETURNING id, customer_id, venue_id, text, rating, status, created_at, updated_at
    `

	q := database.Query{
		Name: "post_repository.Create",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, in.CustomerID, in.VenueID, in.Text, in.Rating)
	if err != nil {
		return entity.Post{}, executeSQLError(err)
	}
	defer row.Close()

	var post Post
	if err := pgxscan.ScanOne(&post, row); err != nil {
		return entity.Post{}, scanRowError(err)
	}

	return post.ToEntity(), nil
}

func (r *Repository) List(ctx context.Context, in entity.ListPostsInput) (entity.ListPostsOutput, error) {
	countSQL := `SELECT COUNT(*) FROM guest.posts`
	countQ := database.Query{
		Name: "post_repository.List.Count",
		Sql:  countSQL,
	}
	var total int
	err := r.db.DB().QueryRowContext(ctx, countQ).Scan(&total)
	if err != nil {
		return entity.ListPostsOutput{}, scanRowError(err)
	}

	// Get paginated posts
	sql := `
		SELECT
		       id,
		       customer_id,
		       venue_id,
		       text,
		       rating,
		       status,
		       created_at,
		       updated_at
		FROM guest.posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	q := database.Query{
		Name: "post_repository.List",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, in.Limit, in.Offset)
	if err != nil {
		return entity.ListPostsOutput{}, executeSQLError(err)
	}
	defer rows.Close()

	var posts Posts
	if err := pgxscan.ScanAll(&posts, rows); err != nil {
		return entity.ListPostsOutput{}, scanRowsError(err)
	}

	return entity.ListPostsOutput{
		Posts: posts.ToEntities(),
		Total: total,
	}, nil
}

func (r *Repository) Get(ctx context.Context, postID string) (entity.Post, error) {
	sql := `
		SELECT
		       id,
		       customer_id,
		       venue_id,
		       text,
		       rating,
		       status,
		       created_at,
		       updated_at
		FROM guest.posts
		WHERE id = $1
	`
	q := database.Query{
		Name: "post_repository.Get",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, postID)
	if err != nil {
		return entity.Post{}, executeSQLError(err)
	}
	defer row.Close()

	var post Post
	if err := pgxscan.ScanOne(&post, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Post{}, apperror.PostNotFoundID(postID)
		}

		return entity.Post{}, scanRowError(err)
	}

	return post.ToEntity(), nil
}
