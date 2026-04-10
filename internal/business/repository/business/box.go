package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
	"github.com/ua-academy-projects/share-bite/pkg/database/pg"
)

func (r *Repository) CreateBox(ctx context.Context, box *entity.Box) (int64, time.Time, error) {
	const op = "repository.box.CreateBox"
	q := database.Query{
		Name: "create_box",
		Sql: `
		INSERT INTO business.boxes 
        (venue_id, category_id, image, price_full, price_discount, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at
	`,
	}

	var id int64
	var createdAt time.Time

	if tx, ok := ctx.Value(pg.TxKey).(pgx.Tx); ok {
		err := tx.QueryRow(ctx, q.Sql,
			box.VenueID,
			box.CategoryID,
			box.Image,
			box.PriceFull,
			box.PriceDiscount,
			box.ExpiresAt,
		).Scan(&id, &createdAt)
		return id, createdAt, fmt.Errorf("%s: %w", op, err)
	}

	err := r.db.DB().QueryRowContext(ctx, q,
		box.VenueID,
		box.CategoryID,
		box.Image,
		box.PriceFull,
		box.PriceDiscount,
		box.ExpiresAt,
	).Scan(&id, &createdAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return 0, time.Time{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, createdAt, err
}

func (r *Repository) CreateBoxItem(ctx context.Context, boxID int64, code string) error {
	const op = "repository.box.CreateBoxItem"
	q := database.Query{
		Name: "create_box_item",
		Sql: `
		INSERT INTO business.box_items (box_id, box_code)
        VALUES ($1, $2)
	`,
	}

	if tx, ok := ctx.Value(pg.TxKey).(pgx.Tx); ok {
		_, err := tx.Exec(ctx, q.Sql, boxID, code)
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err := r.db.DB().ExecContext(ctx, q, boxID, code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
