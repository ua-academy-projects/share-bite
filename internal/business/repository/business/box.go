package business

import (
	"context"
	"time"

	"github.com/ua-academy-projects/share-bite/internal/business/entity"
	"github.com/ua-academy-projects/share-bite/pkg/database"
)

func (r *Repository) CreateBox(ctx context.Context, box *entity.Box) (int64, time.Time, error) {
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

	err := r.db.DB().QueryRowContext(ctx, q,
		box.VenueID,
		box.CategoryID,
		box.Image,
		box.PriceFull,
		box.PriceDiscount,
		box.ExpiresAt,
	).Scan(&id, &createdAt)

	return id, createdAt, err
}

func (r *Repository) CreateBoxItem(ctx context.Context, boxID int64, code string) error {
	q := database.Query{
		Name: "create_box_item",
		Sql: `
		INSERT INTO business.box_items (box_id, box_code)
        VALUES ($1, $2)
	`,
	}

	_, err := r.db.DB().ExecContext(ctx, q, boxID, code)
	return err
}
