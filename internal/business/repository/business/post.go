package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	biserr "github.com/ua-academy-projects/share-bite/internal/business/error"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
)

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
		Scan(&post.ID, &post.OrgId, &post.Content, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, biserr.ErrForbidden
		}
		return nil, fmt.Errorf("insert post query: %w", err)
	}
	return &post, nil
}

func (r *Repository) InsertPostImages(ctx context.Context, postID int, URLs []string) error {
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