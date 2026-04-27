package post

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ua-academy-projects/share-bite/internal/guest/dto"

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

func (r *Repository) UpdateStatus(ctx context.Context, postID, customerID string, status entity.PostStatus) error {
	sql := `UPDATE guest.posts SET status = $1 WHERE id = $2 AND customer_id = $3`
	q := database.Query{
		Name: "post_repository.UpdateStatus",
		Sql:  sql,
	}
	result, err := r.db.DB().ExecContext(ctx, q, status, postID, customerID)
	if err != nil {
		return executeSQLError(err)
	}

	if result.RowsAffected() == 0 {
		return apperror.PostNotFoundID(postID)
	}

	return nil
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

func (r *Repository) Update(ctx context.Context, in entity.UpdatePostInput) (entity.Post, error) {
	sql := `
		UPDATE guest.posts
		SET
	`
	args := pgx.NamedArgs{}
	updates := make([]string, 0, 4)

	if in.Text != nil {
		args["text"] = *in.Text
		updates = append(updates, "text=@text")
	}

	if in.VenueID != nil {
		args["venue_id"] = *in.VenueID
		updates = append(updates, "venue_id=@venue_id")
	}

	if in.Rating != nil {
		args["rating"] = *in.Rating
		updates = append(updates, "rating=@rating")
	}

	if in.Status != nil {
		args["status"] = *in.Status
		updates = append(updates, "status=@status")
	}

	if len(updates) == 0 {
		return entity.Post{}, apperror.ErrEmptyUpdate
	}

	sql += fmt.Sprintf(
		" %s WHERE id=@id AND customer_id=@customer_id RETURNING id, customer_id, venue_id, text, rating, status, created_at, updated_at, published_at",
		strings.Join(updates, ", "),
	)
	args["id"] = in.ID
	args["customer_id"] = in.CustomerID

	q := database.Query{
		Name: "post_repository.Update",
		Sql:  sql,
	}

	row, err := r.db.DB().QueryContext(ctx, q, args)
	if err != nil {
		if translatedErr := translatePostUpdateError(err, in); translatedErr != nil {
			return entity.Post{}, translatedErr
		}

		return entity.Post{}, executeSQLError(err)
	}
	defer row.Close()

	var post Post
	if err := pgxscan.ScanOne(&post, row); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Post{}, apperror.PostNotFoundID(in.ID)
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation {
			return entity.Post{}, apperror.PostNotFoundID(in.ID)
		}

		return entity.Post{}, scanRowError(err)
	}

	result := post.ToEntity()
	images, err := r.loadImagesByPostID(ctx, result.ID)
	if err != nil {
		return entity.Post{}, err
	}

	result.Images = images

	return result, nil
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
		        published_at,
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

	result := posts.ToEntities()
	postIDs := make([]string, 0, len(result))
	for i := range result {
		postIDs = append(postIDs, result[i].ID)
	}

	imagesByPostID, err := r.loadImagesByPostIDs(ctx, postIDs)
	if err != nil {
		return dto.ListPostsOutput{}, err
	}

	for i := range result {
		result[i].Images = imagesByPostID[result[i].ID]
	}

	return dto.ListPostsOutput{
		Posts: result,
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
			   published_at,
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

	result := post.ToEntity()
	images, err := r.loadImagesByPostID(ctx, result.ID)
	if err != nil {
		return entity.Post{}, err
	}

	result.Images = images

	return result, nil
}

func (r *Repository) GetByID(ctx context.Context, postID string) (entity.Post, error) {
	sql := `
		SELECT
		       id,
		       customer_id,
		       venue_id,
		       text,
		       rating,
		       status,
		       created_at,
		       updated_at,
		       published_at
		FROM guest.posts
		WHERE id = $1
	`
	q := database.Query{
		Name: "post_repository.GetByID",
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

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation {
			return entity.Post{}, apperror.PostNotFoundID(postID)
		}

		return entity.Post{}, scanRowError(err)
	}

	result := post.ToEntity()
	images, err := r.loadImagesByPostID(ctx, result.ID)
	if err != nil {
		return entity.Post{}, err
	}

	result.Images = images

	return result, nil
}

func (r *Repository) GetAuthorUserID(ctx context.Context, postID string) (string, error) {
	sql := `
		SELECT c.user_id
		FROM guest.posts p
		JOIN guest.customers c ON p.customer_id = c.id
		WHERE p.id = $1
	`
	q := database.Query{
		Name: "post_repository.GetAuthorUserID",
		Sql:  sql,
	}

	var userID string
	err := r.db.DB().QueryRowContext(ctx, q, postID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperror.PostNotFoundID(postID)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.InvalidTextRepresentation {
			return "", apperror.PostNotFoundID(postID)
		}

		return "", scanRowError(err)
	}

	return userID, nil
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

func translatePostUpdateError(err error, in entity.UpdatePostInput) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	if pgErr.Code == pgerrcode.CheckViolation {
		return apperror.ErrInvalidPostData
	}

	return nil
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
		_, err := r.db.DB().ExecContext(
			ctx,
			q,
			img.PostID,
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

func (r *Repository) DeleteImagesByPostID(ctx context.Context, postID string) error {
	sql := `
		DELETE FROM guest.post_images
		WHERE post_id = $1
	`
	q := database.Query{
		Name: "post_repository.DeleteImagesByPostID",
		Sql:  sql,
	}

	if _, err := r.db.DB().ExecContext(ctx, q, postID); err != nil {
		return executeSQLError(err)
	}

	return nil
}

func (r *Repository) loadImagesByPostID(ctx context.Context, postID string) ([]entity.PostImage, error) {
	sql := `
		SELECT
		       id,
		       post_id,
		       object_key,
		       content_type,
		       file_size,
		       sort_order,
		       created_at
		FROM guest.post_images
		WHERE post_id = $1
		ORDER BY sort_order ASC, id ASC
	`
	q := database.Query{
		Name: "post_repository.loadImagesByPostID",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, postID)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var images PostImages
	if err := pgxscan.ScanAll(&images, rows); err != nil {
		return nil, scanRowsError(err)
	}

	return images.ToEntities(), nil
}

func (r *Repository) loadImagesByPostIDs(ctx context.Context, postIDs []string) (map[string][]entity.PostImage, error) {
	if len(postIDs) == 0 {
		return map[string][]entity.PostImage{}, nil
	}

	sql := `
		SELECT
		       id,
		       post_id,
		       object_key,
		       content_type,
		       file_size,
		       sort_order,
		       created_at
		FROM guest.post_images
		WHERE post_id::text = ANY($1)
		ORDER BY post_id ASC, sort_order ASC, id ASC
	`
	q := database.Query{
		Name: "post_repository.loadImagesByPostIDs",
		Sql:  sql,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, postIDs)
	if err != nil {
		return nil, executeSQLError(err)
	}
	defer rows.Close()

	var images PostImages
	if err := pgxscan.ScanAll(&images, rows); err != nil {
		return nil, scanRowsError(err)
	}

	grouped := make(map[string][]entity.PostImage, len(postIDs))
	for _, image := range images.ToEntities() {
		grouped[image.PostID] = append(grouped[image.PostID], image)
	}

	return grouped, nil
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
