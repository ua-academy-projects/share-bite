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
	"github.com/ua-academy-projects/share-bite/pkg/database/pagination"
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
			box.FullPrice,
			box.DiscountPrice,
			box.ExpiresAt,
		).Scan(&id, &createdAt)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
				return 0, time.Time{}, fmt.Errorf("%s: %w", op, ErrNotFound)
			}
			return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
		}

		return id, createdAt, nil
	}

	err := r.db.DB().QueryRowContext(ctx, q,
		box.VenueID,
		box.CategoryID,
		box.Image,
		box.FullPrice,
		box.DiscountPrice,
		box.ExpiresAt,
	).Scan(&id, &createdAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, sql.ErrNoRows) {
			return 0, time.Time{}, fmt.Errorf("%s: %w", op, ErrNotFound)
		}
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, createdAt, nil
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
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	_, err := r.db.DB().ExecContext(ctx, q, boxID, code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) ListNearbyBoxes(ctx context.Context, offset, limit int, lat, lon float64, categoryID *int, orgID *int) (pagination.Result[entity.BoxWithDistance], error) {
	var args []any

	where := `boxes.expires_at > NOW()
			  AND org_units.latitude IS NOT NULL
			  AND org_units.longitude IS NOT NULL
              AND EXISTS (
                  SELECT 1 
                  FROM business.box_items bi 
                  WHERE bi.box_id = boxes.id 
                    AND bi.reserved_by_user_id IS NULL
              )`

	if categoryID != nil {
		args = append(args, *categoryID)
		where += fmt.Sprintf(" AND boxes.category_id=$%d", len(args))
	}
	if orgID != nil {
		args = append(args, *orgID)
		where += fmt.Sprintf(" AND (org_units.id=$%d OR org_units.parent_id=$%d)", len(args), len(args))
	}

	scanner := func(rows pgx.Rows) (entity.BoxWithDistance, error) {
		var item entity.BoxWithDistance

		err := rows.Scan(
			&item.Box.ID,
			&item.Box.VenueID,
			&item.Box.CategoryID,
			&item.Box.Image,
			&item.Box.FullPrice,
			&item.Box.DiscountPrice,
			&item.Box.CreatedAt,
			&item.Box.ExpiresAt,
			&item.AvailabilityCount,
			&item.Distance,
		)

		if err != nil {
			return entity.BoxWithDistance{}, err
		}

		return item, nil
	}

	dynamicColumns := fmt.Sprintf("boxes.id, boxes.venue_id, boxes.category_id, " +
			"boxes.image, boxes.price_full, boxes.price_discount, " +
			"boxes.created_at, boxes.expires_at, " +
			"(SELECT COUNT(*) FROM business.box_items bi WHERE reserved_by_user_id IS NULL AND bi.box_id=boxes.id) AS availability_count, " +
			"point(%f, %f) <@> point(org_units.longitude, org_units.latitude) AS distance", 
		lon, lat)

	p := pagination.Params{
		Table:   "business.org_units JOIN business.boxes on boxes.venue_id=org_units.id",
		Columns: dynamicColumns,
		Where:   where,
		OrderBy: "distance ASC, boxes.id ASC",
		Args:    args,
		Offset:  offset,
		Limit:   limit,
	}

	return pagination.List(ctx, r.db.DB(), "business_repository.ListNearbyBoxes", p, scanner)
}
