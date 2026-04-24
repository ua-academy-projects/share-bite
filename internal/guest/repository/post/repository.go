package post

import (
	"context"
	"errors"
	"fmt"
	"github.com/ua-academy-projects/share-bite/internal/guest/dto"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *Repository) Create(ctx context.Context, in dto.CreatePostInput) (entity.Post, error) {
	sql := `
        INSERT INTO guest.posts(
            customer_id,
            venue_id,
            text,
            rating
        ) VALUES ($1, $2, $3, $4)
		RETURNING id, customer_id, venue_id, text, rating, status, created_at, updated_at, published_at
    `

	q := database.Query{
		Name: "post_repository.Create",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, in.CustomerID, in.VenueID, in.Text, in.Rating)
	if err != nil {
		if translatedErr := translatePostInsertError(err, in); translatedErr != nil {
			return entity.Post{}, translatedErr
		}

		return entity.Post{}, executeSQLError(err)
	}
	defer row.Close()

	var post Post
	if err := pgxscan.ScanOne(&post, row); err != nil {
		if translatedErr := translatePostInsertError(err, in); translatedErr != nil {
			return entity.Post{}, translatedErr
		}

		return entity.Post{}, scanRowError(err)
	}

	return post.ToEntity(), nil
}

func (r *Repository) CreateImages(ctx context.Context, images []entity.PostImage) error {
	if len(images) == 0 {
		return nil
	}
	createImagesSql := `
		INSERT INTO guest.post_images(
			post_id,
			object_key,
			content_type,
			file_size,
			sort_order
		) VALUES ($1, $2, $3, $4, $5)
	`
	q := database.Query{
		Name: "post_repository.CreateImages",
		Sql:  createImagesSql,
	}

	for _, img := range images {
		postID, err := strconv.ParseInt(img.PostID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid post ID: %w", err)
		}
		_, err = r.db.DB().ExecContext(
			ctx,
			q,
			postID,
			img.ObjectKey,
			img.ContentType,
			img.FileSize,
			img.SortOrder,
		)
		if err != nil {
			return executeSQLError(err)
		}
	}
	return nil
}

func (r *Repository) List(ctx context.Context, in dto.ListPostsInput) (dto.ListPostsOutput, error) {
	countSQL := `SELECT COUNT(*) FROM guest.posts WHERE status = $1`
	countQ := database.Query{
		Name: "post_repository.List.Count",
		Sql:  countSQL,
	}
	var total int
	err := r.db.DB().QueryRowContext(ctx, countQ, entity.PostStatusPublished).Scan(&total)
	if err != nil {
		return dto.ListPostsOutput{}, scanRowError(err)
	}

	// Get paginated posts
	sql := `
		  SELECT
		        p.id,
		        p.customer_id,
		        p.venue_id,
		        p.text,
		        p.rating,
		        p.status,
		        p.created_at,
		        p.updated_at,
		        (SELECT COUNT(*) FROM guest.post_likes pl WHERE pl.post_id = p.id) AS likes_count,
		        EXISTS(SELECT 1 FROM guest.post_likes pl WHERE pl.post_id = p.id AND pl.customer_id = NULLIF($4, '')::uuid) AS is_liked_by_me
		  FROM guest.posts p
		  WHERE p.status = $1
		  ORDER BY p.created_at DESC, p.id DESC
	 	  LIMIT $2 OFFSET $3
	  `
	q := database.Query{
		Name: "post_repository.List",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, entity.PostStatusPublished, in.Limit, in.Offset, in.CustomerID)
	if err != nil {
		return dto.ListPostsOutput{}, executeSQLError(err)
	}
	defer rows.Close()

	var posts Posts
	if err := pgxscan.ScanAll(&posts, rows); err != nil {
		return dto.ListPostsOutput{}, scanRowsError(err)
	}

	return dto.ListPostsOutput{
		Posts: posts.ToEntities(),
		Total: total,
	}, nil
}

func (r *Repository) Get(ctx context.Context, postID string, reqCustomerID string) (entity.Post, error) {
	sql := `
		SELECT
		       p.id,
		       p.customer_id,
		       p.venue_id,
		       p.text,
		       p.rating,
		       p.status,
		       p.created_at,
		       p.updated_at,
		       (SELECT COUNT(*) FROM guest.post_likes pl WHERE pl.post_id = p.id) AS likes_count,
		       EXISTS(SELECT 1 FROM guest.post_likes pl WHERE pl.post_id = p.id AND pl.customer_id = NULLIF($3, '')::uuid) AS is_liked_by_me
		FROM guest.posts p
		WHERE p.id = $1
		  AND p.status = $2
	`
	q := database.Query{
		Name: "post_repository.Get",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, postID, entity.PostStatusPublished, reqCustomerID)
	if err != nil {
		return entity.Post{}, executeSQLError(err)
	}
	defer row.Close()

	var post Post
	if err := pgxscan.ScanOne(&post, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Post{}, apperror.PostNotFoundID(postID)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation {
			return entity.Post{}, apperror.PostNotFoundID(postID)
		}

		return entity.Post{}, scanRowError(err)
	}

	return post.ToEntity(), nil
}

func translatePostInsertError(err error, in dto.CreatePostInput) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	switch pgErr.Code {
	case pgerrcode.ForeignKeyViolation:
		if strings.Contains(pgErr.ConstraintName, "customer_id") {
			return apperror.CustomerNotFoundID(in.CustomerID)
		}
	case pgerrcode.CheckViolation:
		return apperror.ErrInvalidPostData
	}

	return nil
}

func (r *Repository) Like(ctx context.Context, postID string, customerID string) error {
	sql := `
        INSERT INTO guest.post_likes (post_id, customer_id) 
        VALUES ($1, $2) 
        ON CONFLICT DO NOTHING
    `
	q := database.Query{
		Name: "post_repository.Like",
		Sql:  sql,
	}

	_, err := r.db.DB().ExecContext(ctx, q, postID, customerID)
	if err != nil {
		return executeSQLError(err)
	}

	return nil
}

func (r *Repository) Unlike(ctx context.Context, postID string, customerID string) error {
	sql := `
        DELETE FROM guest.post_likes 
        WHERE post_id = $1 AND customer_id = $2
    `
	q := database.Query{
		Name: "post_repository.Unlike",
		Sql:  sql,
	}

	_, err := r.db.DB().ExecContext(ctx, q, postID, customerID)
	if err != nil {
		return executeSQLError(err)
	}

	return nil
}
